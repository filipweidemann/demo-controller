package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodWithoutAnnotation(t *testing.T) {
	podReq := CreateTestPod(&TestPodOptions{})
	pod, err := adminClientSet.CoreV1().Pods("default").Create(context.Background(), &podReq, metav1.CreateOptions{})
	if err != nil {
		t.Error("Could not create Pod")
	}

	assert.Equal(t, pod.Labels[ANNOTATION], "")
}

func TestPodWithAnnotation(t *testing.T) {
	podReq := CreateTestPod(&TestPodOptions{Annotation: ANNOTATION})

	pod, err := adminClientSet.CoreV1().Pods("default").Create(context.Background(), &podReq, metav1.CreateOptions{})
	if err != nil {
		t.Error("Could not create Pod")
	}

	assert.Equal(t, pod.Labels[ANNOTATION], "nginx")
}
