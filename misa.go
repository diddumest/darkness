package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/thecsw/darkness/emilia"
	"github.com/thecsw/darkness/yunyun"
	"github.com/thecsw/gana"
)

// misaCommandFunc will support many different tools that darkness can support,
// such as creating gallery previews, etc. WIP.
func misaCommandFunc() {
	misaCmd := darknessFlagset(misaCommand)

	buildGalleryPreviews := misaCmd.Bool("gallery-previews", false, "build gallery previews")
	removeGalleryPreviews := misaCmd.Bool("no-gallery-previews", false, "delete gallery previews")
	dryRun := misaCmd.Bool("dry-run", false, "skip writing files (but do the reading)")

	options := getEmiliaOptions(misaCmd)
	options.Dev = true
	emilia.InitDarkness(options)

	if *buildGalleryPreviews {
		buildGalleryFiles(*dryRun)
		return
	}
	if *removeGalleryPreviews {
		removeGalleryFiles()
		return
	}

	fmt.Println("I don't know what you want me to do, see -help")
}

const (
	galleryPreviewImageSize = 500
	galleryPreviewImageBlur = 20
)

// buildGalleryFiles finds all the gallery entries and build a resized
// preview version of it.
func buildGalleryFiles(dryRun bool) {
	galleryFiles := getGalleryFiles()
	for i, galleryFile := range galleryFiles {
		fmt.Printf("[%d/%d] ", i+1, len(galleryFiles))
		newFile := emilia.GalleryPreview(galleryFile)
		if info, err := os.Stat(string(newFile)); info != nil && !os.IsNotExist(err) {
			fmt.Printf("%s already exists -- skipping\n", emilia.RelPathToWorkdir(newFile))
			continue
		}
		fmt.Printf("Building %s... ", emilia.RelPathToWorkdir(newFile))

		// Retrieve image contents reader:
		// - For local files, it's a reader of the file.
		// - For remote files, it's a reader of the response body.
		imgReader, err := galleryItemToReader(galleryFile)
		if err != nil {
			fmt.Println("gallery item to reader: " + err.Error())
			continue
		}

		// Encode preview image into a buffer.
		buf := new(bytes.Buffer)
		if err := blurImageForPreview(imgReader, buf); err != nil {
			fmt.Println("gallery reader to writer: " + err.Error())
			continue
		}

		// Read the encoded image buffer into a byte slice.
		img, err := ioutil.ReadAll(buf)
		if err != nil {
			fmt.Println("reading gallery buffer: " + err.Error())
			continue
		}

		// Don't save the file if it's in dry run mode.
		if !dryRun {
			// Write the final preview image file.
			if err := os.WriteFile(string(newFile), img, os.ModePerm); err != nil {
				fmt.Println("failed to save image: " + err.Error())
				continue
			}
		}

		fmt.Println("done")
	}
}

// galleryItemToReader takes in a gallery item and returns an `io.ReadCloser`
// for the image's contents.
func galleryItemToReader(item *emilia.GalleryItem) (io.ReadCloser, error) {
	// If it's a local file, simply open the os file.
	if !item.IsExternal {
		file := emilia.JoinWorkdir(yunyun.JoinRelativePaths(item.Path, item.Item))
		return os.Open(string(file))
	}
	// If it's a remote file, then run a get request and return
	// the body reader.
	resp, err := http.Get(string(item.Item))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't retrieve external gallery item ("+string(item.Item)+")")
	}
	return resp.Body, nil
}

// blurImageForPreview decodes image from `source`, makes a preview out of it,
// and finally encodes it into `target`.
func blurImageForPreview(source io.ReadCloser, target io.Writer) error {
	img, err := imaging.Decode(source, imaging.AutoOrientation(true))
	defer source.Close()
	if err != nil {
		return errors.Wrap(err, "couldn't read the image reader")
	}
	// Resize the image to save up on storage.
	img = imaging.Resize(img, galleryPreviewImageSize, 0, imaging.Lanczos)
	blurred := imaging.Blur(img, galleryPreviewImageBlur)
	return imaging.Encode(target, blurred, imaging.JPEG)
}

// removeGalleryFiles removes all generate gallery previews.
func removeGalleryFiles() {
	for _, galleryFile := range getGalleryFiles() {
		newFile := emilia.GalleryPreview(galleryFile)
		if err := os.Remove(string(newFile)); err != nil && !os.IsNotExist(err) {
			fmt.Println("Couldn't delete", newFile, "| reason:", err.Error())
		}
	}
}

// getGalleryFiles returns a slice of all gallery images represented as `emilia.GalleryItem`.
func getGalleryFiles() []*emilia.GalleryItem {
	inputFilenames := make(chan yunyun.FullPathFile, customChannelCapacity)
	pages := gana.GenericWorkers(gana.GenericWorkers(inputFilenames,
		func(v yunyun.FullPathFile) gana.Tuple[yunyun.FullPathFile, string] {
			data, err := ioutil.ReadFile(filepath.Clean(string(v)))
			if err != nil {
				fmt.Printf("Failed to open %s: %s\n", v, err.Error())
			}
			return gana.NewTuple(v, string(data))
		}, 1, customChannelCapacity), func(v gana.Tuple[yunyun.FullPathFile, string]) *yunyun.Page {
		return emilia.ParserBuilder.BuildParser(emilia.PackRef(v.UnpackRef())).Parse()
	}, customNumWorkers, customChannelCapacity)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go emilia.FindFilesByExt(inputFilenames, emilia.Config.Project.Input, wg)

	// Launch a second discovery for gallery files.
	galleryFiles := make([]*emilia.GalleryItem, 0, 32)
	go func(wg *sync.WaitGroup) {
		for page := range pages {
			for _, gc := range page.Contents.Galleries() {
				for _, item := range gc.List {
					galleryFiles = append(galleryFiles, emilia.NewGalleryItem(page, gc, item))
				}
			}
			wg.Done()
		}
		wg.Done()
	}(wg)

	wg.Wait()
	return galleryFiles
}
