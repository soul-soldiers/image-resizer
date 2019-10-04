package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
)

// Cloud Function entry point
func ResizeImage(w http.ResponseWriter, r *http.Request) {
	// parse the url query sting into ResizerParams
	p, err := ParseQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// fetch input image and resize
	img, err := FetchAndResizeImage(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// encode output image to jpeg buffer
	encoded, err := EncodeImageToJpg(img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set Content-Type and Content-Length headers
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(encoded.Len()))

	// write the output image to http response body
	_, err = io.Copy(w, encoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// struct containing the initial query params
type ResizerParams struct {
	url    string
	height int
	width  int
	mode   string
}

// parse/validate the url params
func ParseQuery(r *http.Request) (*ResizerParams, error) {

	var p ResizerParams
	var mode string
	query := r.URL.Query()

	baseUrl := os.Getenv("BASE_URL")
	img := query.Get("image")
	if img == "" {
		return &p, errors.New("Url Param 'image' is missing")
	}

	mode = query.Get("mode")
	if mode == "" {
		mode = "fill"
	}

	width, _ := strconv.Atoi(query.Get("x"))
	height, _ := strconv.Atoi(query.Get("y"))

	if width == 0 && height == 0 {
		return &p, errors.New("Url Param 'x' or 'y' must be set")
	}

	log.Printf("Rendering: %s/%s. Width: %d, Height: %d, Mode: %s", baseUrl, img, width, height, mode)
	p = NewResizerParams(fmt.Sprintf("%s/%s", baseUrl, img), height, width, mode)

	return &p, nil
}

// ResizerParams factory
func NewResizerParams(url string, height int, width int, mode string) ResizerParams {
	return ResizerParams{url, height, width, mode}
}

// fetch the image from provided url and resize it
func FetchAndResizeImage(p *ResizerParams) (*image.Image, error) {
	var dst image.Image

	// fetch input data
	response, err := http.Get(p.url)
	if err != nil {
		return &dst, err
	}
	// don't forget to close the response
	defer response.Body.Close()

	// decode input data to image
	src, _, err := image.Decode(response.Body)
	if err != nil {
		return &dst, err
	}

	// resize input image
	if p.mode == "resize" {
		dst = imaging.Resize(src, p.width, p.height, imaging.Lanczos)
	} else if p.mode == "fill" {
		dst = imaging.Fill(src, p.width, p.height, imaging.Center, imaging.Lanczos)
	} else {
		return &dst, errors.New(fmt.Sprintf("Invalid mode: %s", p.mode))
	}

	return &dst, nil
}

// encode image to jpeg
func EncodeImageToJpg(img *image.Image) (*bytes.Buffer, error) {
	encoded := &bytes.Buffer{}
	err := jpeg.Encode(encoded, *img, nil)
	return encoded, err
}

// server for local testing
func main() {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		port, _ = strconv.Atoi(p)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ResizeImage", ResizeImage)
	fmt.Printf("Starting local server on port: %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
