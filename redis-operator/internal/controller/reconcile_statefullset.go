package controller

import (
	"context"
	"fmt"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type stsSpecification struct {
	replicas    int32
	stsName     string
	labels      map[string]string
	spec        bcredisv1alpha1.BcredisSpec
	bcredisName string
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
	}

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      stsName,
			Namespace: bcredis.Namespace,
		},
	}
	logger.Info("Create or update of statefullset", stsName)
	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, sts, func() error {
		if err := controllerutil.SetControllerReference(bcredis, sts, r.Scheme); err != nil {
			return err
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

		sts.Spec = getSpec(s, env, volumeMounts,idx)
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
if idx == 1 {
    configFile = "/etc/redis/replica.conf"
}

	specInfo := appsv1.StatefulSetSpec{
		Replicas:    &s.replicas,
		ServiceName: s.stsName + "-headless",
		Selector: &metav1.LabelSelector{
			MatchLabels: s.labels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: s.labels,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "redis",
						Image: s.spec.RedisImage,
						// Start with a basic config; role is managed dynamically via exec
						Command: []string{"redis-server", configFile},
						Ports: []corev1.ContainerPort{
							{ContainerPort: 6379, Name: "redis"},
						},
						Env:          env,
						VolumeMounts: volumeMounts,
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: []string{"redis-cli", "ping"},
								},
							},
							InitialDelaySeconds: 5,
							PeriodSeconds:       5,
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: []string{"redis-cli", "ping"},
								},
							},
							InitialDelaySeconds: 15,
							PeriodSeconds:       10,
						},
					},
					{
						Name:    "sentinel",
						Image:   s.spec.RedisImage,
						Command: []string{"redis-sentinel", "/etc/redis/sentinel.conf"},
						Ports: []corev1.ContainerPort{
							{ContainerPort: 26379, Name: "sentinel"},
						},
						VolumeMounts: []corev1.VolumeMount{
							{Name: "sentinel-config", MountPath: "/etc/redis"},
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
					AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
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
