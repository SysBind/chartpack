package infrastructure

import (
	"fmt"
	"github.com/SysBind/chartpack/domain"
	"os"
	"os/exec"
)

type (
	Node struct {
		domain.Node
	}
)

func (node Node) copy(filename string) error {
	fmt.Printf("copying %s to %s \n", filename, node.Hostname)
	cmd := exec.Command("scp", filename, "root@"+node.Ip+":/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}

func CopyToNodes(nodes []domain.Node, file string) error {
	for _, node := range nodes {
		sshNode := Node{node}
		err := sshNode.copy(file)

		if err != nil {
			panic(err)
		}
	}
	return nil
}
