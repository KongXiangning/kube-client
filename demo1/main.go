package main

import (
	"fmt"
	"k8s-client-go/common"
	"k8s.io/client-go/kubernetes"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main()  {
	var (
		clientset *kubernetes.Clientset
		podsList *core_v1.PodList
		err error
	)

	if clientset,err = common.InitClient();err != nil{
		goto  FAIL
	}

	if podsList,err = clientset.CoreV1().Pods("").List(meta_v1.ListOptions{});err != nil  {
		goto  FAIL
	}

	clientset.AppsV1beta1().Deployments().Create()

	fmt.Println(*podsList)

	return

	FAIL:
		fmt.Println(err)
		return
}
