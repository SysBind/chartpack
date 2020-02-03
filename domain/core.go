package domain

// "Domain"
type (
	Image struct {
		Repo string
		Name string
		Tag  string
	}

	Chart struct {
		Name   string
		Images []Image
	}

	Exporter interface {
		Export(chart Chart)
	}
)

func Package(charts []Chart, exporter Exporter) {
	for _, chart := range charts {
		exporter.Export(chart)
	}
}
