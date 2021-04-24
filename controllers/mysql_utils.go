package controllers

import (
	"context"

	"github.com/workshop/mysql-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const serviceName = "some-name-mysql"

var labels = map[string]string{
	"app.kubernetes.io/component":  "mysql",
	"app.kubernetes.io/instance":   "some-name",
	"app.kubernetes.io/name":       "mysql-cluster",
	"app.kubernetes.io/managed-by": "mysql-operator",
}

func (r *CustomMysqlReconciler) CreateService(cr *v1alpha1.CustomMysql) error {
	srv := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:    labels,
			Name:      serviceName,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name:       "mysql",
					Port:       3306,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(3306),
				},
			},
			Selector: labels,
		},
	}

	currService := &corev1.Service{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      serviceName,
		Namespace: cr.Namespace,
	}, currService)

	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if errors.IsNotFound(err) {
		return r.Client.Create(context.TODO(), srv)
	}

	return nil
}

func (r *CustomMysqlReconciler) CreateSfs(cr *v1alpha1.CustomMysql) error {
	fsGroup := int64(1001)
	sfs := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &cr.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			ServiceName: serviceName,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "perconalab/percona-server:gr-test",
							Args: []string{
								"--gtid_mode=ON",
								"--enforce-gtid-consistency=ON",
								"--plugin_load_add=group_replication.so",
							},
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: "root_password",
								},
							},
							ImagePullPolicy: corev1.PullAlways,
							Ports: []corev1.ContainerPort{
								{
									Name:          "mysql",
									ContainerPort: 3306,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "datadir",
									MountPath: "/var/lib/mysql",
								},
							},
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup:            &fsGroup,
						SupplementalGroups: []int64{fsGroup},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: cr.Namespace,
						Name:      "datadir",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("2Gi"),
							},
						},
					},
				},
			},
		},
	}

	currSFS := &appsv1.StatefulSet{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      serviceName,
		Namespace: cr.Namespace,
	}, currSFS)

	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if errors.IsNotFound(err) {
		return r.Client.Create(context.TODO(), sfs)
	}

	if *currSFS.Spec.Replicas != cr.Spec.Size {
		currSFS.Spec.Replicas = &cr.Spec.Size
		return r.Client.Update(context.TODO(), currSFS)
	}

	return nil
}
