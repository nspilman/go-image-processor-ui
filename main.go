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

	// type FuncOneInput func(img1 image.Image) image.Image
	type FuncTwoInputs func(img1 image.Image, img2 image.Image) image.Image
	functionStoreTwoInputs := make(map[string]FuncTwoInputs)
	// functionStoreOneInput := make(map[string]FuncOneInput)

	functionStoreTwoInputs[getFunctionName(BlendImages)] = BlendImages
	functionStoreTwoInputs[getFunctionName(ReplaceHue)] = ReplaceHue
	functionStoreTwoInputs[getFunctionName(ReplaceBrightPixels)] = ReplaceBrightPixels
	// functionStoreOneInput[getFunctionName(ReplaceBrightPixels)] = ReplaceBrightPixels
	var keys []string
	for key := range functionStoreTwoInputs {
		keys = append(keys, "/"+strings.TrimPrefix(key, "main."))
	}

	http.Handle("/", templ.Handler(PageWrapper("Image processing", ImageForm("/upload", keys))))
	http.HandleFunc("/output.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/dist/output.css")
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)

		file, fileHeader, err := r.FormFile("file")
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
		base64Image, err := generateBase64Image(BlendImages(img, img2))
		var opts jpeg.Options
		opts.Quality = 1

		Image(base64Image).Render(r.Context(), w)
	})

	http.HandleFunc("/ReplaceBrightPixels", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			PageWrapper("Replace Bright Pixels", ImageForm("/ReplaceBrightPixels", keys)).Render(r.Context(), w)
		case "POST":
			handleUpload(w, r, functionStoreTwoInputs["main.ReplaceBrightPixels"])
		default:
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ReplaceHue", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			PageWrapper("Replace Hue", ImageForm("/ReplaceHue", keys)).Render(r.Context(), w)
		case "POST":
			handleUpload(w, r, functionStoreTwoInputs["main.ReplaceHue"])
		default:
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/BlendImages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			PageWrapper("Replace Hue", ImageForm("/BlendImages", keys)).Render(r.Context(), w)
		case "POST":
			handleUpload(w, r, functionStoreTwoInputs["main.BlendImages"])
		default:
			http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
		}
	})

	for key := range functionStoreTwoInputs {
		fmt.Println("/" + key)
		http.HandleFunc("/"+key, func(w http.ResponseWriter, r *http.Request) {
			handleUpload(w, r, functionStoreTwoInputs[key])
		})
	}

	fmt.Println("Listening on :4000")
	http.ListenAndServe(":4000", nil)
}

func handleUpload(w http.ResponseWriter, r *http.Request, mutationFunction func(img1 image.Image, img2 image.Image) image.Image) {
	r.ParseMultipartForm(0)

	file, fileHeader, err := r.FormFile("file")
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
	base64Image, err := generateBase64Image(mutationFunction(img, img2))
	var opts jpeg.Options
	opts.Quality = 1

	Image(base64Image).Render(r.Context(), w)
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
