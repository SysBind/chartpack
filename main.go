package main

import (
	"./domain"
	"io/ioutil"
	"log"
	"os"
)

// "Interface"
type (
	Loader interface {
		Load() []domain.Chart
	}
)

// Infrastructure
type (
	LocalLoader struct {
		path string
	}

	Exporter struct {
		src  string
		dest string
	}
)

func (exporter Exporter) Export(chart domain.Chart) {
	log.Println("Exporting " + chart.Name)
}

func (loader LocalLoader) Load() []domain.Chart {
	var retval []domain.Chart = nil

	files, err := ioutil.ReadDir(loader.path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			retval = append(retval, domain.Chart{Name: file.Name()})
		}
	}
	return retval
}

func main() {
	arg := os.Args[1]

	log.Println("scanning " + arg)
	loader := LocalLoader{path: arg}

	charts := loader.Load()

	exporter := Exporter{src: arg, dest: "/tmp"}

	domain.Package(charts, exporter)

}
