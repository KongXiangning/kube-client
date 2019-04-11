package execute

import (
	"encoding/json"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

var endpoints = func(transfer KubeTransfer, outChan chan KubeTransfer) (err error) {
	var (
		client        = client.clientset
		endpoint      = v1.Endpoints{}
		k8sEndpoint   *v1.Endpoints
		deleteOptions *metaV1.DeleteOptions
	)

	if err = json.Unmarshal(transfer.HandleJson, &endpoint); err != nil {
		goto FAIL
	} else {
		transfer.HandleJson = nil
	}

	switch transfer.Types {
	case 0:
		if _, err = client.CoreV1().Endpoints(endpoint.Namespace).Create(&endpoint); err != nil {
			goto FAIL
		}
	case 1:
		if _, err = client.CoreV1().Endpoints(endpoint.Namespace).Update(&endpoint); err != nil {
			goto FAIL
		}
	case 2:
		if k8sEndpoint, err = client.CoreV1().Endpoints(endpoint.Namespace).Get(endpoint.Name, metaV1.GetOptions{}); err != nil {
			goto FAIL
		} else {
			if transfer.HandleJson, err = json.Marshal(k8sEndpoint); err != nil {
				goto FAIL
			}
		}
	case 3:
		if err = client.CoreV1().Endpoints(endpoint.Namespace).Delete(endpoint.Name, deleteOptions); err != nil {
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
