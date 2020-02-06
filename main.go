package main

import (
	"github.com/SysBind/chartpack/domain"
	"github.com/SysBind/chartpack/infrastructure"
	"log"
	"os"
)

func main() {
	src := os.Args[1]
	dest := os.Args[2]

	log.Println("scanning " + src)
	loader := infrastructure.LocalLoader{Path: src}

	charts := loader.Load()

	exporter := infrastructure.Exporter{Src: src, Dest: dest}

	domain.Package(charts, exporter)

}
