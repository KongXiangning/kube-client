package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"k8s.io/api/apps/v1beta1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	inV1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	wcoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"kube-client/common"
	"log"
)

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func Int32ToBytes(i int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func main() {
	var (
		clientset *kubernetes.Clientset
		podsList  *core_v1.PodList
		k8sDeploy *v1beta1.Deployment
		is        chan int
		err       error
	)

	var i int64 = 232422
	b := Int64ToBytes(i)
	fmt.Println(b)

	if clientset, err = common.InitClient(); err != nil {
		goto FAIL
	}

	if podsList, err = clientset.CoreV1().Pods("bms-pre").List(meta_v1.ListOptions{}); err != nil {
		goto FAIL
	}

	k8sDeploy, err = clientset.AppsV1beta1().Deployments("bms-pre").Get("biz-rest", v1.GetOptions{})
	fmt.Println(k8sDeploy)

	if b, er := json.Marshal(k8sDeploy); er != nil {
		log.Fatal(er)
	} else {
		fmt.Println(string(b))
		fmt.Println("-===============================================")
	}

	go watchDeployment(clientset.AppsV1beta1().Deployments("bms-pre"), clientset.CoreV1().Pods("bms-pre"))

	fmt.Println(*podsList)

	<-is

	return

FAIL:
	fmt.Println(err)
	return
}

func watchDeployment(deployment inV1beta1.DeploymentInterface, pod wcoreV1.PodInterface) {
	var timeout int64
	timeout = 6000
	list := v1.ListOptions{TimeoutSeconds: &timeout, LabelSelector: "app=biz-provider"}

	wd, _ := deployment.Watch(list)
	wp, _ := pod.Watch(list)
	for {
		select {
		case v := <-wd.ResultChan():
			//d := v.Object.()
			//j := v1beta1.Deployment(v.Object)
			//j := v.Object.DeepCopyObject().(v1beta1.Deployment)

			if v.Object == nil {
				fmt.Println("timeout")
				goto END
			} else {
				fmt.Println(v.Type)
				d := v.Object.(*v1beta1.Deployment)
				if d.GetName() == "biz-provider" {
					fmt.Println(*d.Spec.Replicas, d.Status.Replicas, d.Status.UpdatedReplicas, d.Status.AvailableReplicas, d.Status.UnavailableReplicas)
					listCont := d.Status.Conditions
					for index, l := range listCont {
						fmt.Println(index, l)
					}
				}
			}
		case p := <-wp.ResultChan():
			if p.Object == nil {
				fmt.Println("timeout")
				goto END
			} else {
				po := p.Object.(*core_v1.Pod)
				if po.Status.Conditions != nil {
					fmt.Println(po.Status.Conditions[0])
					fmt.Println(po.Status.Conditions[1])
					fmt.Println(po.Status.Conditions[2])
				}
				if po.Status.ContainerStatuses != nil {
					fmt.Println(po.Status.ContainerStatuses[0].RestartCount)
				}
				fmt.Println(po.Status.Message)
			}
		}
	}
END:
	fmt.Println("lll")
}
