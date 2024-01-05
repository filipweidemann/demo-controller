package controller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TestPodOptions struct {
	Annotation string
}

func CreateTestPod(options *TestPodOptions) corev1.Pod {
	podMeta := metav1.ObjectMeta{Name: "testpod", Namespace: "default"}
	podSpec := corev1.PodSpec{Containers: []corev1.Container{
		{
			Name:  "nginx",
			Image: "nginx",
		},
	}}

	if options.Annotation != "" {
		podMeta.Annotations = map[string]string{
			ANNOTATION: "true",
		}
	}

	podReq := corev1.Pod{ObjectMeta: podMeta, Spec: podSpec}
	return podReq
}
