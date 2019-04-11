package execute

import (
	"encoding/json"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

var secrets = func(transfer KubeTransfer, outChan chan KubeTransfer) (err error) {
	var (
		client        = client.clientset
		secret        = v1.Secret{}
		k8sSecret     *v1.Secret
		deleteOptions *metaV1.DeleteOptions
	)

	if err = json.Unmarshal(transfer.HandleJson, &secret); err != nil {
		goto FAIL
	} else {
		transfer.HandleJson = nil
	}

	switch transfer.Types {
	case 0:
		if _, err = client.CoreV1().Secrets(secret.Namespace).Create(&secret); err != nil {
			goto FAIL
		}
	case 1:
		if _, err = client.CoreV1().Secrets(secret.Namespace).Update(&secret); err != nil {
			goto FAIL
		}
	case 2:
		if k8sSecret, err = client.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metaV1.GetOptions{}); err != nil {
			goto FAIL
		} else {
			if transfer.HandleJson, err = json.Marshal(k8sSecret); err != nil {
				goto FAIL
			}
		}
	case 3:
		if err = client.CoreV1().Secrets(secret.Namespace).Delete(secret.Name, deleteOptions); err != nil {
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
