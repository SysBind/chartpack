package domain

import (
	"fmt"
	"strconv"
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
		fmt.Println("calling exporter.Export on ", chart.Name)
		exporter.Export(chart)
	}
}

func (chart *Chart) populateImagesFromValues(values Values) {
	for k, v := range values {
		if k == "image" {
			fmt.Println("Found image record in chart", chart.Name)
			m, ok := v.(map[string]interface{})
			if !ok {
				fmt.Println("populateImagesFromValues: could not convert v to map of strings")
			} else {
				fmt.Printf("Tag is %v\n", m["tag"])
				imageTag := m["tag"]

				var imageTagString string
				// try as string
				if imageTagString, ok = imageTag.(string); ok == false {
					// try as float64 (odd panic: interface {} is float64, not string
					var imageTagFloat float64
					if imageTagFloat, ok = imageTag.(float64); ok {
						imageTagString = strconv.FormatFloat(imageTagFloat, 'g', -1, 64)
					}

				}
				chart.Images = append(chart.Images, Image{Repo: m["repository"].(string), Tag: imageTagString})
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
