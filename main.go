//
// Multi-threaded mandelbrot microservice.
// Author: Bob Kozdemba (bkozdemba@gmail.com)
//
// Testing: ab -n 2000 -c8 http://localhost:8080/mandelbrot
//
// Ideas borrowed from from https://gist.githubusercontent.com/santiaago/d6d681d14c5f3b3f5d69/raw/b9cae0c3d0e10cb91e8a47a0d4e8420fb3a05c31/main.go
package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"html/template"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/jpeg"
	"log"
	"math/cmplx"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// var randNumber = uint8(0)
var root = flag.String("root", ".", "file system path")

func main() {
	go http.HandleFunc("/mandelbrot", fracHandler)
	// http.Handle("/", http.FileServer(http.Dir(*root)))
	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// plotHandler - Draw the fractal image.
func fracHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UTC().UnixNano())
	// Update the global variable used to control contrast.
	//randNumber = uint8(rand.Intn(15))
	//log.Println("fracHandler running, randNumber = ", randNumber)
	width := 1024
	height := 1024
	m := image.NewRGBA(image.Rect(0, 0, width, height))

	log.Println("fracHandler: createImage started.")
	start := time.Now()
	mandel := createImage(width, height)
	t := time.Now()
	elapsed := t.Sub(start)
	log.Println("fracHandler: createImage finished, elasped time = ", elapsed)

	draw.Draw(m, m.Bounds(), mandel, image.ZP, draw.Src)

	var img image.Image = m
	writeImageWithTemplate(w, &img)
}

var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head><style TYPE="text/css"> h1 { font-size: xx-large; font-family: sans-serif } </style> </head>
<body><h1></h1><img src="data:image/jpg;base64,{{.Image}}"></body>`

// Writeimagewithtemplate encodes an image 'img' in jpeg format and writes it into ResponseWriter using a template.
func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
}

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 61.
//!+

// Mandelbrot - Calculates and returns an image of the Mandelbrot fractal.
func createImage(width int, height int) image.Image {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
	)

	var c uint8
	c = uint8(rand.Intn(15))
	log.Println("createImage: contrast = ", c)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/float64(height)*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/float64(width)*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, mandelbrot(z, c))
		}
	}
	// png.Encode(os.Stdout, img) // NOTE: ignoring errors
	return img
}

// mandelbrot - Compute and return the pixel.
//              Implement color LUT using a go map type - key is based on 'n'?

func mandelbrot(z complex128, contrast uint8) color.Color {
	const iterations = 200
	var v complex128

	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			// return color.Gray{255 - contrast*n}
			return palette.Plan9[255-contrast*n]
		}
	}
	return color.Black
}

//!-

// Some other interesting functions:

func acos(z complex128) color.Color {
	v := cmplx.Acos(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}

func sqrt(z complex128) color.Color {
	v := cmplx.Sqrt(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{128, blue, red}
}

// f(x) = x^4 - 1
//
// z' = z - f(z)/f'(z)
//    = z - (z^4 - 1) / (4 * z^3)
//    = z - (z - 1/z^3) / 4
func newton(z complex128) color.Color {
	const iterations = 37
	const contrast = 7
	for i := uint8(0); i < iterations; i++ {
		z -= (z - 1/(z*z*z)) / 4
		if cmplx.Abs(z*z*z*z-1) < 1e-6 {
			return color.Gray{255 - contrast*i}
		}
	}
	return color.Black
}
