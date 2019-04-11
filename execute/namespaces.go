package execute

import (
	"encoding/json"
	"errors"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

var namespaces = func(transfer KubeTransfer, outChan chan KubeTransfer) (err error) {
	var (
		client        = client.clientset
		namespace     = coreV1.Namespace{}
		k8sNamespace  *coreV1.Namespace
		deleteOptions *metaV1.DeleteOptions
	)

	if err = json.Unmarshal(transfer.HandleJson, &namespace); err != nil {
		goto FAIL
	} else {
		transfer.HandleJson = nil
	}

	switch transfer.Types {
	case 0:
		if _, err = client.CoreV1().Namespaces().Create(&namespace); err != nil {
			goto FAIL
		}
	case 1:
		err = errors.New("no types")
		goto FAIL
	case 2:
		if k8sNamespace, err = client.CoreV1().Namespaces().Get(namespace.Name, metaV1.GetOptions{}); err != nil {
			goto FAIL
		} else {
			if transfer.HandleJson, err = json.Marshal(k8sNamespace); err != nil {
				goto FAIL
			}
		}
	case 3:
		if err = client.CoreV1().Namespaces().Delete(namespace.Name, deleteOptions); err != nil {
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
