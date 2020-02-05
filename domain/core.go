package domain

import (
	"fmt"
)

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
		chart.populateImagesFromValues(chart.Values)
		exporter.Export(chart)
	}
}

func (chart *Chart) populateImagesFromValues(values Values) {
	for k, v := range values {
		if k == "image" {
			m, ok := v.(map[string]interface{})
			if !ok {
				fmt.Println("populateImagesFromValues: could not convert v to map of strings")
			} else {
				chart.Images = append(chart.Images, Image{Repo: m["repository"].(string), Tag: m["tag"].(string)})
			}
		} else if v != nil {
			// try converting to Values (i.e: map[string]interface{}), to handle nested image spec
			m, ok := v.(map[string]interface{})
			if ok {
				chart.populateImagesFromValues(m)
			}
		}

	}
}
