package infrastructure

import (
	"fmt"
	"github.com/SysBind/chartpack/domain"
	"github.com/helm/helm/pkg/chartutil"
	"io/ioutil"
	"os"
)

type (
	Image struct {
		domain.Image
	}

	Loader interface {
		Load() []domain.Chart
	}
)

type (
	LocalLoader struct {
		Path string
	}

	Exporter struct {
		Src  string
		Dest string
	}
)

func (exporter Exporter) Export(chart domain.Chart) {
	dest := fmt.Sprintf("%s/%s", exporter.Dest, chart.Name)
	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		panic(err)
	}
	for _, image := range chart.Images {
		_image := Image{image}
		_image.Fetch(dest)
	}

	fmt.Println("copy chart over to ", exporter.Dest+"/"+chart.Name)
	err = CopyDirectory(exporter.Src+"/"+chart.Name, exporter.Dest+chart.Name+"/")
	if err != nil {
		panic(err)
	}
}

func (loader LocalLoader) Load() []domain.Chart {
	var retval []domain.Chart = nil

	files, err := ioutil.ReadDir(loader.Path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			if _, err := os.Stat(loader.Path + "/" + file.Name() + "/values.yaml"); err != nil {
				continue
			}
			filename := loader.Path + "/" + file.Name() + "/values.yaml"

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
