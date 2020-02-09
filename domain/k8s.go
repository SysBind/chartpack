package domain

import (
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
)

type (
	Node struct {
		Hostname string
		Ip       string
		IsMaster bool
	}

	Loader interface {
		LoadNodes() corev1.NodeList
	}
)

func FilterNodes(nodes []Node, filter func(Node) bool) []Node {
	var retval []Node

	for _, node := range nodes {
		if filter(node) {
			retval = append(retval, node)
		}
	}
	return retval
}

func GetNodes(loader Loader) []Node {
	var retval []Node
	list := loader.LoadNodes()

	for _, node := range list.Items {
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
		retval = append(retval, Node{Hostname: nodeName, Ip: nodeIp, IsMaster: isMaster})
	}

	return retval
}
