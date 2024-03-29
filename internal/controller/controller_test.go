package controller_test

import (
	"context"
	"testing"
	"time"

	"github.com/filipweidemann/demo-controller/internal/controller"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	T_Duration = time.Second * 5
	T_Timeout  = time.Second * 5
	T_Interval = time.Millisecond * 500
)

func TestPodWithoutAnnotation(t *testing.T) {
	podReq := CreateTestPod(&TestPodOptions{})
	defer k8sClientSet.CoreV1().Pods("default").Delete(context.Background(), podReq.Name, metav1.DeleteOptions{})

	pod, err := k8sClientSet.CoreV1().Pods("default").Create(context.Background(), &podReq, metav1.CreateOptions{})

	if err != nil {
		t.Error("Could not create Pod")
	}

	Expect(pod.Labels[controller.LABEL]).Should(Equal(""))
}

func TestPodWithAnnotation(t *testing.T) {
	podReq := CreateTestPod(&TestPodOptions{SetAnnotation: true})
	defer k8sClientSet.CoreV1().Pods("default").Delete(context.Background(), podReq.Name, metav1.DeleteOptions{})

	_, err := k8sClientSet.CoreV1().Pods("default").Create(context.Background(), &podReq, metav1.CreateOptions{})
	if err != nil {
		t.Error("Could not create Pod")
	}

	Eventually(func() bool {
		pod, err := k8sClientSet.CoreV1().Pods("default").Get(context.Background(), "testpod", metav1.GetOptions{})
		if err != nil {
			t.Error("Couldn't fetch updated pod")
		}
		return pod.Labels[controller.LABEL] == "testpod"
	}, T_Timeout, T_Interval).Should(BeTrue())

}
