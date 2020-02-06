package infrastructure

import (
	"github.com/SysBind/chartpack/domain"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/helm/helm/pkg/chartutil"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"log"
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

func tryFetch(imageUri string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	log.Println("Fetching: ", imageUri)
	out, err := cli.ImagePull(ctx, imageUri, types.ImagePullOptions{})

	if err != nil {
		return err
	}

	defer out.Close()

	io.Copy(os.Stdout, out)

	return nil
}

func (image Image) Fetch() {
	image_uri := image.Repo + ":" + image.Tag

	err := tryFetch(image_uri)

	if err != nil {
		// attempt to add docker.io prefix
		image_uri = "docker.io/" + image.Repo + ":" + image.Tag
		log.Println("Retrying with: ", image_uri)

		err := tryFetch(image_uri)
		if err != nil {
			image_uri = "docker.io/library/" + image.Repo + ":" + image.Tag
			log.Println("Retrying with: ", image_uri)

			err := tryFetch(image_uri)

			if err != nil {
				panic(err)
			}
		}
	}
}
