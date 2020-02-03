package main

import (
	"log"
	"os"
	"path/filepath"
)

// "Domain"
type (
	Image struct {
		repo string
		name string
		tag  string
	}

	Chart struct {
		name   string
		images []Image
	}

	Loader interface {
		Load() []Chart
	}
)

// "Interface"

func load_charts() {
}

// Infrastructure

type (
	LocalLoader struct {
		path string
	}
)

func (loader LocalLoader) Load() []Chart {
	var files []string
	err := filepath.Walk(loader.path, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		log.Println(file)
	}

	return nil
}

func main() {
	arg := os.Args[1]

	log.Println("scanning " + arg)
	loader := LocalLoader{path: arg}

	charts := loader.Load()

	for _, el := range charts {
		log.Println(el.name)
	}
}
