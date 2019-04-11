package execute

import (
	"encoding/json"
	"k8s.io/api/apps/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

var statefulsets = func(transfer KubeTransfer, outChan chan KubeTransfer) (err error) {
	var (
		client         = client.clientset
		statefulSet    = v1beta1.StatefulSet{}
		k8sStatefulSet *v1beta1.StatefulSet
		deleteOptions  *metaV1.DeleteOptions
	)

	if err = json.Unmarshal(transfer.HandleJson, &statefulSet); err != nil {
		goto FAIL
	} else {
		transfer.HandleJson = nil
	}

	switch transfer.Types {
	case 0:
		if _, err = client.AppsV1beta1().StatefulSets(statefulSet.Namespace).Create(&statefulSet); err != nil {
			goto FAIL
		}
	case 1:
		if _, err = client.AppsV1beta1().StatefulSets(statefulSet.Namespace).Update(&statefulSet); err != nil {
			goto FAIL
		}
	case 2:
		if k8sStatefulSet, err = client.AppsV1beta1().StatefulSets(statefulSet.Namespace).Get(statefulSet.Name, metaV1.GetOptions{}); err != nil {
			goto FAIL
		} else {
			if transfer.HandleJson, err = json.Marshal(k8sStatefulSet); err != nil {
				goto FAIL
			}
		}
	case 3:
		if err = client.AppsV1beta1().StatefulSets(statefulSet.Namespace).Delete(statefulSet.Name, deleteOptions); err != nil {
			goto FAIL
		}
	}

	transfer.Types = 1
	transfer.Result = "success"
	outChan <- transfer
	return
FAIL:
	log.Println(err)
	return
}
