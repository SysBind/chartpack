package infrastructure

import (
	"context"
	"fmt"
	"github.com/SysBind/chartpack/domain"
	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

func Nodes() []domain.Node {
	var retval []domain.Node

	client, err := loadClient()
	if err != nil {
		panic(err)
	}
	fmt.Println("have kube client, listing nodes")

	var nodes corev1.NodeList
	if err := client.List(context.Background(), "", &nodes); err != nil {
		panic(err)
	}
	for _, node := range nodes.Items {
		var nodeName, nodeIp string
		var isMaster bool
		for _, addr := range node.Status.Addresses {
			switch *addr.Type {
			case "InternalIP":
				nodeIp = *addr.Address
			case "Hostname":
				nodeName = *addr.Address
			}
		}
		// check if master
		for idx, _ := range node.Metadata.Labels {
			if idx == "node-role.kubernetes.io/master" {
				isMaster = true
				break
			}
		}
		retval = append(retval, domain.Node{Hostname: nodeName, Ip: nodeIp, IsMaster: isMaster})
	}

	return retval
}

func loadClient() (*k8s.Client, error) {
	data, err := ioutil.ReadFile(homeDir() + "/.kube/config")
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig: %v", err)
	}

	// Unmarshal YAML into a Kubernetes config object.
	var config k8s.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal kubeconfig: %v", err)
	}
	return k8s.NewClient(&config)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
