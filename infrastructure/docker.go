package infrastructure

import (
	"io"
	"os"
	"log"
	"github.com/helm/helm/pkg/chartutil"
	"io/ioutil"
	"github.com/SysBind/chartpack/domain"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
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
	for _, image := range chart.Images {
		_image := Image{image}
		_image.Fetch()
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


func (image Image) Fetch() {
	log.Println("Fetching image", image)
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, image.Repo + ":" + image.Tag, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)
}
