package main

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
)

func generateBase64Image(img image.Image) (string, error) {

	// Encode to PNG
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return "", err
	}

	// Convert to base64
	base64Image := base64.StdEncoding.EncodeToString(buffer.Bytes())

	return "data:image/png;base64," + base64Image, nil
}

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

	http.Handle("/", templ.Handler(PageWrapper("IAMABANANA", button("hello", "WHATS UP", "/pioneer"))))
	http.HandleFunc("/output.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/dist/output.css")
	})
	http.Handle("/pioneer", templ.Handler(PioneerForm()))
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			log.Fatal("Error retrieving the file: ", err)
			return
		}
		defer file.Close()

		fmt.Printf("Uploaded File: %+v\n", fileHeader.Filename)
		fmt.Printf("File Size: %+v\n", fileHeader.Size)
		fmt.Printf("MIME Header: %+v\n", fileHeader.Header)

		outputFile, err := os.Create("reconstructed.jpg")
		if err != nil {
			panic(err)
		}

		img, err := jpeg.Decode(file)

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
		base64Image, err := generateBase64Image(img)
		var opts jpeg.Options
		opts.Quality = 1

		// err = jpeg.Encode(out, pioneer, &opts)
		// pioneer, err := jpeg.Decode(buf)

		Image(base64Image).Render(r.Context(), w)
	})

	fmt.Println("Listening on :4000")
	http.ListenAndServe(":4000", nil)
}
