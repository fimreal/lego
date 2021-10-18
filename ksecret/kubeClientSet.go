package ksecret

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	// "k8s.io/client-go/tools/clientcmd"
)

// 获取 k8s 配置，支持本地读取 kubeconfig 和 POD 内自动加载 RBAC 配置
func NewKubeClient() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	} else {
		log.Println("成功获取到 kubernetes 配置")
		return kubernetes.NewForConfig(config)
	}

	// log.Println("Fallthrough...尝试读取本地配置文件")
	// if Conf.Kubeconfig == "" {
	// 	log.Fatalln("配置文件中未指定 kubeconfig 位置")
	// }
	// kubeconfig := &Conf.Kubeconfig

	// config, err = clientcmd.BuildConfigFromFlags("")
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// return kubernetes.NewForConfig(config)
}
