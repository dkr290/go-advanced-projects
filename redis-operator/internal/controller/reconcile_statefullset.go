package controller

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type stsSpecification struct {
	replicas             int32
	stsName              string
	labels               map[string]string
	spec                 bcredisv1alpha1.BcredisSpec
	bcredisName          string
	currentMasterService string
}

// reconcileStatefulSet creates/updates a StatefulSet for a Redis instance.
func (r *BcredisReconciler) reconcileStatefulSet(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
	idx int, logger logr.Logger,
) error {
	stsName := fmt.Sprintf("%s-redis-%d", bcredis.Name, idx)
	labels := map[string]string{
		"app":      bcredis.Name + "-redis",
		"instance": fmt.Sprintf("%d", idx),
		"bcredis":  bcredis.Name,
	}

	s := stsSpecification{
		stsName:     stsName,
		replicas:    int32(1),
		labels:      labels,
		spec:        spec,
		bcredisName: bcredis.Name,
		currentMasterService: func() string {
			if bcredis.Status.CurrentMasterService != "" {
				return bcredis.Status.CurrentMasterService
			}
			return fmt.Sprintf("%s-redis-0", bcredis.Name)
		}(),
	}
	volumeMounts := []corev1.VolumeMount{
		{Name: "redis-data", MountPath: "/data"},
		{Name: "redis-config", MountPath: "/etc/redis"},
	}

	env := []corev1.EnvVar{
		{Name: "REDIS_INSTANCE_INDEX", Value: fmt.Sprintf("%d", idx)},
	}
	if spec.RedisPasswordSecret != "" {
		env = append(env, corev1.EnvVar{
			Name: "REDIS_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: spec.RedisPasswordSecret,
					},
					Key: "password",
				},
			},
		})
	}

	desiredSpec := getSpec(s, env, volumeMounts, idx)

	// Compute hash of desired spec
	specBytes, _ := json.Marshal(desiredSpec)
	hash := fmt.Sprintf("%x", sha256.Sum256(specBytes))
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      stsName,
			Namespace: bcredis.Namespace,
		},
	}
	logger.Info("Create or update of statefullset", "statefullset", stsName)
	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, sts, func() error {
		if err := controllerutil.SetControllerReference(bcredis, sts, r.Scheme); err != nil {
			return err
		}
		// Only update spec if hash changed
		currentHash := ""
		if sts.Annotations != nil {
			currentHash = sts.Annotations["bcredis/spec-hash"]
		}
		if currentHash != hash {
			sts.Spec = desiredSpec
			if sts.Annotations == nil {
				sts.Annotations = map[string]string{}
			}
			sts.Annotations["bcredis/spec-hash"] = hash
		}
		return nil
	})
	if err == nil {
		logger.Info("StatefulSet successfully reconciled", "operation", result)
	}
	return err
}

func getSpec(
	s stsSpecification,
	env []corev1.EnvVar,
	volumeMounts []corev1.VolumeMount,
	idx int,
) appsv1.StatefulSetSpec {
	storageClass := s.spec.StorageClassName
	configFile := "/etc/redis/master.conf"
	if idx > 0 {
		configFile = "/etc/redis/replica.conf"
	}
	redisCommand := []string{"redis-server", configFile}
	redisArgs := []string(nil)

	redisProbeCmd := []string{"redis-cli", "ping"}

	if s.spec.RedisPasswordSecret != "" {
		redisCommand = []string{"sh", "-c"}
		redisArgs = []string{
			fmt.Sprintf(
				`exec redis-server %s --requirepass "$REDIS_PASSWORD" --masterauth "$REDIS_PASSWORD"`,
				configFile,
			),
		}
		redisProbeCmd = []string{
			"sh",
			"-c",
			`if [ -n "$REDIS_PASSWORD" ]; then redis-cli -a "$REDIS_PASSWORD" ping; else redis-cli ping; fi`,
		}
	}

	sentinelArgs := []string{
		`if [ ! -f /data/sentinel.conf ]; then cp /etc/redis-config/sentinel.conf /data/sentinel.conf; fi && \
			if [ -n "$REDIS_PASSWORD" ]; then \
				if grep -q '^sentinel auth-pass mymaster ' /data/sentinel.conf; then \
					sed -i "s#^sentinel auth-pass mymaster .*#sentinel auth-pass mymaster $REDIS_PASSWORD#" /data/sentinel.conf; \
				else \
					echo "sentinel auth-pass mymaster $REDIS_PASSWORD" >> /data/sentinel.conf; \
				fi; \
			fi && \
    exec redis-server /data/sentinel.conf --sentinel`,
	}
	specInfo := appsv1.StatefulSetSpec{
		Replicas:    &s.replicas,
		ServiceName: fmt.Sprintf("%s-redis-headless", s.bcredisName),
		Selector: &metav1.LabelSelector{
			MatchLabels: s.labels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: s.labels,
			},
			Spec: corev1.PodSpec{
				Affinity: &corev1.Affinity{
					PodAntiAffinity: &corev1.PodAntiAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
							{
								LabelSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"bcredis": s.bcredisName,
									},
								},
								TopologyKey: "kubernetes.io/hostname",
							},
						},
					},
				},

				Containers: []corev1.Container{
					{
						Name:  "redis",
						Image: s.spec.RedisImage,
						// Start with a basic config; role is managed dynamically via exec
						Command: redisCommand,
						Args:    redisArgs,
						Ports: []corev1.ContainerPort{
							{ContainerPort: 6379, Name: "redis"},
						},
						Env:          env,
						VolumeMounts: volumeMounts,
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: redisProbeCmd,
								},
							},
							InitialDelaySeconds: 5,
							PeriodSeconds:       5,
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: redisProbeCmd,
								},
							},
							InitialDelaySeconds: 15,
							PeriodSeconds:       10,
						},
					},
					{
						Name:    "sentinel",
						Image:   s.spec.RedisImage,
						Command: []string{"sh", "-c"},
						Args:    sentinelArgs,
						Env:     env,
						Ports: []corev1.ContainerPort{
							{ContainerPort: 26379, Name: "sentinel"},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "sentinel-config",
								MountPath: "/etc/redis-config",
								ReadOnly:  true,
							},
							{Name: "redis-data", MountPath: "/data"},
						},

						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: []string{"redis-cli", "-p", "26379", "ping"},
								},
							},
							InitialDelaySeconds: 5,
							PeriodSeconds:       10,
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "redis-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: s.bcredisName + "-redis-config",
								},
							},
						},
					},
					{
						Name: "sentinel-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: s.bcredisName + "-sentinel-config",
								},
							},
						},
					},
				},
			},
		},
		VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "redis-data",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					StorageClassName: &storageClass,
					Resources: corev1.VolumeResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: s.spec.StorageSize,
						},
					},
				},
			},
		},
	}
	return specInfo
}
