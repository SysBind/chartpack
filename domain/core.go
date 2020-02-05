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
		chart.populateImagesFromValues()
		exporter.Export(chart)
	}
}

func (chart *Chart) populateImagesFromValues() {
	var flat_tag, flat_image string // for cases where image specification is flat

	var popFunc func(values Values)
	
	popFunc = func(values Values) {
		for k, v := range values {
			if k == "image" {
				m, ok := v.(map[string]interface{})
				if ok {
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
				} else {
					fmt.Println("WARNING: populateImagesFromValues:  value is not to map of strings, key:", k, "trying flat parsing")
					var ok bool
					flat_image, ok = v.(string)
					if ok && (flat_tag != "") {					
						chart.Images = append(chart.Images, Image{Repo: flat_image, Tag: flat_tag})
						flat_tag = ""
						flat_image = ""
					}
				}
			} else if k == "imageTag" { // flat image spec
				var ok bool
				flat_tag, ok = v.(string)
				if ok && (flat_image != "") {				
					chart.Images = append(chart.Images, Image{Repo: flat_image, Tag: flat_tag})
					flat_tag = ""
					flat_image = ""
				}
			} else if v != nil {
				// try converting to Values (i.e: map[string]interface{}), to handle nested image spec
				m, ok := v.(map[string]interface{})
				if ok {
					// fmt.Println("recursing into ", k)
					popFunc(m)
				}
			}

		}
	}

	popFunc(chart.Values)
}
