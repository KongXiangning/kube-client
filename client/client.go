package client

import (
	"k8s.io/client-go/kubernetes"
	"kube-client/execute"
)

type Kube struct {
	clientset   *kubernetes.Clientset
	typeFunMaps map[string]func(execute.KubeTransfer, chan execute.KubeTransfer) error
}
