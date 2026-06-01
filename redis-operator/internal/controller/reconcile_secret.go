package controller

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)


func (r *BcredisReconciler) reconcilePasswordSecret(ctx context.Context, cr *bcredisv1alpha1.Bcredis) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.RedisPasswordSecret,
			Namespace: cr.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, secret, func() error {
		if secret.Data == nil {
			secret.Data = map[string][]byte{}
		}
		if len(secret.Data["password"]) == 0 {
			pw, err := generatePassword(32)
			if err != nil {
				return err
			}
			secret.Data["password"] = []byte(pw)
		}

		secret.Type = corev1.SecretTypeOpaque
		return controllerutil.SetControllerReference(cr, secret, r.Scheme)
	})
	return err
}

func generatePassword(n int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	out := make([]byte, n)
	max := big.NewInt(int64(len(chars)))

	for i := range out {
		v, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("generate password: %w", err)
		}
		out[i] = chars[v.Int64()]
	}
	return string(out), nil
}

