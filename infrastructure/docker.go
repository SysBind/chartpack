package infrastructure

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"strings"
)

func tryFetch(imageUri, imagePrefix string, dest string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	log.Println("fetching: ", imagePrefix+imageUri)
	out, err := cli.ImagePull(ctx, imagePrefix+imageUri, types.ImagePullOptions{})
	if err == nil {
		defer out.Close()
		b := make([]byte, 8)
		for {
			_, err := out.Read(b)
			// read and do nothing with it, seems to make docker client happier
			if err == io.EOF {
				break
			}
		}

		log.Println("Saving to ", dest)
		destination, err := os.Create(dest)
		if err != nil {
			panic(err)
		}
		defer destination.Close()

		img_reader, err := cli.ImageSave(ctx, []string{imageUri})
		if err != nil {
			panic(err)
		}
		defer img_reader.Close()

		written, err := io.Copy(destination, img_reader)
		if err != nil {
			panic(err)
		}
		log.Printf("written %d bytes", written)
	}
	return err
}

func (image Image) Fetch(dest string) {
	image_uri := image.Repo + ":" + image.Tag

	dest = dest + "/" + strings.Replace(image.Repo, "/", "_", -1) + ".tar"

	err := tryFetch(image_uri, "", dest)

	if err != nil {
		// attempt to add docker.io prefix
		err := tryFetch(image_uri, "docker.io/", dest)
		if err != nil {
			err := tryFetch(image_uri, "docker.io/library/", dest)
			if err != nil {
				panic(err)
			} else {
				log.Println("pulled ", image_uri)
			}

		} else {
			log.Println("pulled ", image_uri)
		}
	} else {
		log.Println("pulled ", image_uri)
	}
}
