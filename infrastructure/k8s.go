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

type (
	kubeapi struct {
		client *k8s.Client
	}
)

func (api kubeapi) LoadNodes() corev1.NodeList {
	var nodes corev1.NodeList

	if err := api.client.List(context.Background(), "", &nodes); err != nil {
		panic(err)
	}

	return nodes
}

func Nodes() []domain.Node {
	var retval []domain.Node

	client, err := loadClient()
	if err != nil {
		panic(err)
	}
	api := kubeapi{client: client}
	retval = domain.GetNodes(api)

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
