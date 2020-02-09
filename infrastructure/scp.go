package infrastructure

import (
	"fmt"
	"github.com/SysBind/chartpack/domain"
	"os/exec"
	"sync"
)

type (
	Node struct {
		domain.Node
	}
)

func (node Node) copy(filename string, wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Printf("copying %s to %s \n", filename, node.Hostname)
	cmd := exec.Command("scp", filename, "root@"+node.Ip+":/")
	err := cmd.Run()

	return err
}

func CopyToNodes(nodes []domain.Node, file string) error {
	var wg sync.WaitGroup
	for _, node := range nodes {
		sshNode := Node{node}
		wg.Add(1)
		go sshNode.copy(file, &wg)

	}
	fmt.Println("CopyToNodes: waiting for scp goroutins to finish..")
	wg.Wait()
	return nil
}
