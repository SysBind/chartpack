package main

import (
	"log"
	"os"
	"github.com/SysBind/chartpack/infrastructure"
	"github.com/SysBind/chartpack/domain"
)


func main() {
	arg := os.Args[1]

	log.Println("scanning " + arg)
	loader := infrastructure.LocalLoader{Path: arg}

	charts := loader.Load()

	exporter := infrastructure.Exporter{Src: arg, Dest: "/tmp"}

	domain.Package(charts, exporter)

}
