package deployment

import (
	gramolav1alpha1 "github.com/redhat/gramola-operator/pkg/apis/gramola/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	eventsImage         = "image-registry.openshift-image-registry.svc:5000/gramola-operator-project/events-s2i"
	eventsDatabaseImage = "image-registry.openshift-image-registry.svc:5000/openshift/postgresql:10"
)

// NewEventsDeployment returns the deployment object for Gateway
func NewEventsDeployment(cr *gramolav1alpha1.AppService, name string, namespace string) *appsv1.Deployment {
	image := eventsImage
	labels := GetAppServiceLabels(cr, name)

	/*
			- name: DB_USERNAME
		              valueFrom:
		                 secretKeyRef:
		                   name: events-database-secret
		                   key: database-user
		            - name: DB_PASSWORD
		              valueFrom:
		                 secretKeyRef:
		                   name: events-database-secret
						   key: database-password
					- name: DB_NAME
		              valueFrom:
		                 secretKeyRef:
		                   name: events-database-secret
		                   key: database-name
		            - name: DB_SERVICE_NAME
		              valueFrom:
		                 configMapKeyRef:
		                   name: events-configmap
		                   key: database_service_name
		            - name: DB_SERVICE_PORT
		              valueFrom:
		                 configMapKeyRef:
		                   name: events-configmap
		                   key: database_service_port

	*/

	env := []corev1.EnvVar{
		{
			Name: "DB_DBID",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "database-name",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name + "-mysql",
					},
				},
			},
		},
		{
			Name:  "DB_HOST",
			Value: name + "-mysql",
		},
		{
			Name: "DB_PASS",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "database-password",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name + "-mysql",
					},
				},
			},
		},
		{
			Name:  "DB_PORT",
			Value: "3306",
		},
		{
			Name: "DB_USER",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "database-user",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name + "-mysql",
					},
				},
			},
		},
		{
			Name:  "NODE_ENV",
			Value: "production",
		},
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "mysql",
							Image:           image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9001,
									Protocol:      "TCP",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: int32(9001),
										},
										Scheme: corev1.URISchemeHTTP,
									},
								},
								FailureThreshold:    5,
								InitialDelaySeconds: 60,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: int32(9001),
										},
										Scheme: corev1.URISchemeHTTP,
									},
								},
								FailureThreshold:    3,
								InitialDelaySeconds: 120,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      name + "-settings",
									MountPath: "/opt/etherpad/config",
								},
							},
							Env: env,
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: name + "-settings",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: name + "-settings",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// NewEventsDatabaseDeployment returns the DB deployment for Events
func NewEventsDatabaseDeployment(cr *gramolav1alpha1.AppService, name string, namespace string) *appsv1.Deployment {
	image := eventsDatabaseImage
	labels := GetAppServiceLabels(cr, name)

	/*
			 - name: POSTGRESQL_USER
		              valueFrom:
		                secretKeyRef:
		                  name: events-database-secret
		                  key: database-user
		            - name: POSTGRESQL_PASSWORD
		              valueFrom:
		                secretKeyRef:
		                  name: events-database-secret
		                  key: database-password
		            - name: POSTGRESQL_DATABASE
		              valueFrom:
		                secretKeyRef:
		                  name: events-database-secret
						  key: database-name


						kind: Secret
						apiVersion: v1
						metadata:
						name: events-db
						namespace: gramola-operator-project
						selfLink: /api/v1/namespaces/gramola-operator-project/secrets/events-db
						uid: 4d523899-e84a-496b-b4c1-195ef14a9a2e
						resourceVersion: '1885035'
						creationTimestamp: '2020-03-14T08:41:18Z'
						labels:
							template: postgresql-persistent-template
							template.openshift.io/template-instance-owner: 63537123-ea07-4413-8d22-683bc1b74fe9
						annotations:
							template.openshift.io/expose-database_name: '{.data[''database-name'']}'
							template.openshift.io/expose-password: '{.data[''database-password'']}'
							template.openshift.io/expose-username: '{.data[''database-user'']}'
						data:
						database-name: c2FtcGxlZGI=
						database-password: WXJtMWJKQXY2SXNDUk5Jbg==
						database-user: dXNlcjBKUQ==
						type: Opaque
	*/

	env := []corev1.EnvVar{
		{
			Name: "MYSQL_USER",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "database-user",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name,
					},
				},
			},
		},
		{
			Name: "MYSQL_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "database-password",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name,
					},
				},
			},
		},
		{
			Name: "MYSQL_ROOT_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "database-root-password",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name,
					},
				},
			},
		},
		{
			Name: "MYSQL_DATABASE",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "database-name",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name,
					},
				},
			},
		},
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "mysql",
							Image:           image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          name,
									ContainerPort: 3306,
									Protocol:      "TCP",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-i",
											"-c",
											"MYSQL_PWD=\"$MYSQL_PASSWORD\" mysql -h 127.0.0.1 -u $MYSQL_USER -D $MYSQL_DATABASE -e 'SELECT 1'",
										},
									},
								},
								InitialDelaySeconds: 5,
								FailureThreshold:    10,
								TimeoutSeconds:      1,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      name + "-data",
									MountPath: "/var/lib/mysql/data",
								},
							},
							Env: env,
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: name + "-data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: name,
								},
							},
						},
					},
				},
			},
		},
	}
}
