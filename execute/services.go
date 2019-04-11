package execute

import (
	"encoding/json"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

var services = func(transfer KubeTransfer, outChan chan KubeTransfer) (err error) {
	var (
		client        = client.clientset
		service       = v1.Service{}
		k8sService    *v1.Service
		deleteOptions *metaV1.DeleteOptions
	)

	if err = json.Unmarshal(transfer.HandleJson, &service); err != nil {
		goto FAIL
	} else {
		transfer.HandleJson = nil
	}

	switch transfer.Types {
	case 0:
		if _, err = client.CoreV1().Services(service.Namespace).Create(&service); err != nil {
			goto FAIL
		}
	case 1:
		if _, err = client.CoreV1().Services(service.Namespace).Update(&service); err != nil {
			goto FAIL
		}
	case 2:
		if k8sService, err = client.CoreV1().Services(service.Namespace).Get(service.Name, metaV1.GetOptions{}); err != nil {
			goto FAIL
		} else {
			if transfer.HandleJson, err = json.Marshal(k8sService); err != nil {
				goto FAIL
			}
		}
	case 3:
		if err = client.CoreV1().Secrets(service.Namespace).Delete(service.Name, deleteOptions); err != nil {
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
