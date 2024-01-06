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
	"reflect"
	"runtime"
	"strings"

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

	type FuncOneInput func(imgs []image.Image) image.Image
	functionStore := make(map[string]FuncOneInput)

	functionStore[getFunctionName(BlendImages)] = BlendImages
	functionStore[getFunctionName(ReplaceHue)] = ReplaceHue
	functionStore[getFunctionName(ReplaceBrightPixels)] = ReplaceBrightPixels
	functionStore[getFunctionName(QuantizeImage)] = QuantizeImage
	functionStore[getFunctionName(QuantizeLight)] = QuantizeLight
	var keys []string
	for key := range functionStore {
		keys = append(keys, "/"+strings.TrimPrefix(key, "main."))
	}

	http.Handle("/", templ.Handler(PageWrapper("Image processing", ImageForm("/upload", keys, 2))))
	http.HandleFunc("/output.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/dist/output.css")
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)

		file, fileHeader, err := r.FormFile("file1")
		if err != nil {
			log.Fatal("Error retrieving the file: ", err)
			return
		}
		defer file.Close()

		file2, _, err2 := r.FormFile("file2")

		if err2 != nil {
			log.Fatal("Error retrieving the file: ", err)
			return
		}
		defer file2.Close()

		if err != nil {
			panic(err)
		}

		fmt.Printf("Uploaded File: %+v\n", fileHeader.Filename)
		fmt.Printf("File Size: %+v\n", fileHeader.Size)
		fmt.Printf("MIME Header: %+v\n", fileHeader.Header)

		img, err := jpeg.Decode(file)

		if err != nil {
			panic(err)
		}

		img2, err := jpeg.Decode(file2)

		if err != nil {
			panic(err)
		}

		// secondImage, err := LoadImageFromFile("./outer-space.jpeg")

		// base64Image, err := generateBase64Image((ReplaceHue(secondImage, InvertLight(img))))
		base64Image, err := generateBase64Image(BlendImages([]image.Image{img, img2}))
		var opts jpeg.Options
		opts.Quality = 1

		Image(base64Image).Render(r.Context(), w)
	})

	http.HandleFunc("/ReplaceBrightPixels", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			PageWrapper("Replace Bright Pixels", ImageForm("/ReplaceBrightPixels", keys, 2)).Render(r.Context(), w)
		case "POST":
			handleUpload(w, r, ReplaceBrightPixels(getTwoImages(r)))
		default:
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ReplaceHue", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			PageWrapper("Replace Hue", ImageForm("/ReplaceHue", keys, 2)).Render(r.Context(), w)
		case "POST":
			handleUpload(w, r, ReplaceHue(getTwoImages((r))))
		default:
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/BlendImages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			PageWrapper("Replace Hue", ImageForm("/BlendImages", keys, 2)).Render(r.Context(), w)
		case "POST":
			handleUpload(w, r, (BlendImages(getTwoImages(r))))
		default:
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/QuantizeImage", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			PageWrapper("Quantize Image", ImageForm("/QuantizeImage", keys, 1)).Render(r.Context(), w)
		case "POST":
			handleUpload(w, r, (QuantizeImage(getOneImage(r))))
		default:
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/QuantizeLight", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			PageWrapper("Quantize Light", ImageForm("/QuantizeLight", keys, 1)).Render(r.Context(), w)
		case "POST":
			handleUpload(w, r, (QuantizeLight(getOneImage(r))))
		default:
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Listening on :4000")
	http.ListenAndServe(":4000", nil)
}

func parseImageFromRequest(r *http.Request, key string) image.Image {
	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		log.Fatal("Error retrieving the file: ", err)
		panic(err)
	}
	defer file.Close()
	img, err := jpeg.Decode(file)

	fmt.Printf("Uploaded File: %+v\n", fileHeader.Filename)
	fmt.Printf("File Size: %+v\n", fileHeader.Size)
	fmt.Printf("MIME Header: %+v\n", fileHeader.Header)
	return img
}

func getOneImage(r *http.Request) []image.Image {
	r.ParseMultipartForm(0)
	img := parseImageFromRequest(r, "file1")
	return []image.Image{img}
}
func getTwoImages(r *http.Request) []image.Image {
	r.ParseMultipartForm(0)
	img := parseImageFromRequest(r, "file1")
	img2 := parseImageFromRequest(r, "file2")
	return []image.Image{img, img2}
}

func handleUpload(w http.ResponseWriter, r *http.Request, img image.Image) {
	base64Image, err := generateBase64Image(img)
	var opts jpeg.Options
	opts.Quality = 1

	if err != nil {
		panic(err)
	}

	Image(base64Image).Render(r.Context(), w)
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
