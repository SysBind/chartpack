package main

import (
	"github.com/SysBind/chartpack/domain"
	"github.com/helm/helm/pkg/chartutil"
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
	for _, image := range chart.Images {
		log.Println("Image: ", image)
	}
}

func (loader LocalLoader) Load() []domain.Chart {
	var retval []domain.Chart = nil

	files, err := ioutil.ReadDir(loader.path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			if _, err := os.Stat(loader.path + "/" + file.Name() + "/values.yaml"); err != nil {
				continue
			}
			filename := loader.path + "/" + file.Name() + "/values.yaml"

			source, err := ioutil.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			var vals chartutil.Values
			vals, err = chartutil.ReadValues(source)
			if err != nil {
				panic(err)
			}
			retval = append(retval, domain.Chart{Name: file.Name(), Values: domain.Values(vals)})
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
