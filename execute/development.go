package execute

import (
	"encoding/json"
	"fmt"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	inV1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	"log"
)

var deployment = func(transfer KubeTransfer)(result []byte,err error)  {
	var(
		client *kubernetes.Clientset
		deployment = v1beta1.Deployment{}
		k8sDeployment  *v1beta1.Deployment
		deleteOptions *v1.DeleteOptions
	)

	if err = json.Unmarshal(transfer.HandleJson,&deployment);err !=nil {
		goto FAIL
	}

	switch transfer.Types {
	case 0:
		if k8sDeployment,err = client.AppsV1beta1().Deployments(deployment.Namespace).Create(&deployment);err != nil{
			goto FAIL
		}
	case 1:
		if k8sDeployment,err = client.AppsV1beta1().Deployments(deployment.Namespace).Update(&deployment);err != nil{
			goto FAIL
		}
	case 2:
		k8sDeployment = nil
		if err = client.AppsV1beta1().Deployments(deployment.Namespace).Delete(deployment.Name,deleteOptions);err != nil{
			goto FAIL
		}
	}
	if k8sDeployment != nil{
		go  watchDeployment(client.AppsV1beta1().Deployments(deployment.Namespace))
	}

FAIL:
	log.Println(err)
	return
}

func watchDeployment(deployment inV1beta1.DeploymentInterface)  {
	w,_ := deployment.Watch(v1.ListOptions{})
	for{
		select {
		case v := <- w.ResultChan():
			fmt.Println(v.Type, v.Object)
		}
	}
}
