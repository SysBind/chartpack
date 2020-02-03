package domain

/*import (
	"reflect"
	"fmt"
)*/

// "Domain"
type (
	// as also defined at github.com/helm/helm/pkg/chartutil/values.go
	Values map[string]interface{}

	Image struct {
		Repo string
		Name string
		Tag  string
	}

	Chart struct {
		Name   string
		Values Values
		Images []Image
	}

	Exporter interface {
		Export(chart Chart)
	}
)

func Package(charts []Chart, exporter Exporter) {
	for _, chart := range charts {
		// chart.populateImagesFromValues()
		exporter.Export(chart)
	}
}

/*
func (chart *Chart) populateImagesFromValues() {
	for k, v := range chart.Values {
		if k == "image" {
			if v.Kind() == reflect.Map {
				for s, b := range v {
					fmt.Printf("%s: book=%v\n", s, b)
				}
				//chart.Images = append(chart.Images, Image{Repo: v["repository"], Name: v["repository"], Tag: v["tag"]})
			}
		}
	}
}
*/
