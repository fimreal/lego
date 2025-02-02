package ksecret

import (
	"context"
	"strings"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func UpdateSecret(clientset *kubernetes.Clientset, s *Secret) (err error) {

	// create secret struct
	metadata := metav1.ObjectMeta{
		Name:        s.SecretName,
		Namespace:   s.SecretNamespace,
		Annotations: map[string]string{},
	}
	metadata.Annotations["replicator.v1.mittwald.de/replicate-to"] = ".*"
	// metadata.Name = s.SecretName
	// metadata.Namespace = s.SecretNamespace

	newSecret := v1.Secret{
		TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		// ObjectMeta: metav1.ObjectMeta{Name: s.SecretName, Namespace: s.SecretNamespace},
		ObjectMeta: metadata,
		Data:       map[string][]byte{},
		StringData: map[string]string{},
		Type:       "kubernetes.io/tls",
	}
	newSecret.Data["tls.key"] = s.Key
	newSecret.Data["tls.crt"] = s.Crt

	// update secret
	_, err = clientset.CoreV1().Secrets(s.SecretNamespace).Update(context.TODO(), &newSecret, metav1.UpdateOptions{})
	if err == nil {
		return
	}
	log.Warnf("Update secret[%s/%s]: %s", s.SecretNamespace, s.SecretName, err.Error())

	// create secret
	_, err = clientset.CoreV1().Secrets(s.SecretNamespace).Create(context.TODO(), &newSecret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return
}

type Secret struct {
	SecretName      string
	SecretNamespace string
	Crt             []byte
	Key             []byte
}

func DeployToSecret(secretName *string, cert *certificate.Resource) (err error) {
	clientset, err := NewKubeClient()
	if err != nil {
		return err
	}
	sName := strings.Split(*secretName, "/")
	s := Secret{
		SecretName:      sName[1],
		SecretNamespace: sName[0],
		Crt:             cert.Certificate,
		Key:             cert.PrivateKey,
	}
	return UpdateSecret(clientset, &s)
}
