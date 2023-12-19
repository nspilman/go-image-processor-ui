package main

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
)

func convertJpegToPng(jpegData []byte) ([]byte, error) {
	// Decode the JPEG data
	img, err := jpeg.Decode(bytes.NewReader(jpegData))
	if err != nil {
		return nil, err
	}

	// Encode the image to PNG
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var FS embed.FS

func main() {
	component := hello("Julian")

	http.Handle("/", templ.Handler(component))
	http.HandleFunc("/output.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/dist/output.css")
	})
	http.Handle("/pioneer", templ.Handler(PioneerForm()))
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)

		buf, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("request", err)
		}

		outputFile, err := os.Create("reconstructed.jpg")
		if err != nil {
			panic(err)
		}

		img, err := jpeg.Decode(bytes.NewReader(buf))

		defer outputFile.Close()
		jpeg.Encode(outputFile, img, nil)

		// file, err := os.ReadFile("./BlueSailboats.jpg")
		// if err != nil {
		// 	panic(err)
		// }
		// defer file.Close()

		// Step 3: Decode the image.
		// img, err := jpeg.Decode(file)
		if err != nil {
			panic(err)
		}

		// pioneer, _, err := image.Decode(bytes.NewReader(buf))
		// jpgImage, err := convertJpegToPng(file)
		base64Image := base64.StdEncoding.EncodeToString(buf)
		var opts jpeg.Options
		opts.Quality = 1

		// err = jpeg.Encode(out, pioneer, &opts)
		// pioneer, err := jpeg.Decode(buf)

		Image("data:image/png;base64,"+base64Image).Render(r.Context(), w)
	})

	fmt.Println("Listening on :4000")
	http.ListenAndServe(":4000", nil)
}
