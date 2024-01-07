package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodWithoutAnnotation(t *testing.T) {
	podReq := CreateTestPod(&TestPodOptions{})
	pod, err := adminClientSet.CoreV1().Pods("default").Create(context.Background(), &podReq, metav1.CreateOptions{})
	defer adminClientSet.CoreV1().Pods("default").Delete(context.Background(), pod.Name, metav1.DeleteOptions{})

	if err != nil {
		t.Error("Could not create Pod")
	}

	assert.Equal(t, pod.Labels[LABEL], "")
}

func TestPodWithAnnotation(t *testing.T) {
	podReq := CreateTestPod(&TestPodOptions{Annotation: ANNOTATION})

	_, err := adminClientSet.CoreV1().Pods("default").Create(context.Background(), &podReq, metav1.CreateOptions{})
	if err != nil {
		t.Error("Could not create Pod")
	}

	time.Sleep(time.Second * 1)

	pod, err := adminClientSet.CoreV1().Pods("default").Get(context.Background(), "testpod", metav1.GetOptions{})
	if err != nil {
		t.Error("Couldn't fetch updated pod")
	}

	assert.Equal(t, "true", pod.Annotations[ANNOTATION])
	assert.Equal(t, "testpod", pod.Labels[LABEL])
}
