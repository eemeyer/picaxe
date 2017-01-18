package main

import (
	"log"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	reader, err := os.Open("testimages/baby-duck.jpeg")
	//reader, err := os.Open("testimages/bordertest4.jpeg")
	//reader, err := os.Open("testimages/rotated.jpeg")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	spec := ProcessingSpec{
		Format:               ImageFormatPng,
		Trim:                 TrimModeFuzzy,
		TrimFuzzFactor:       0.5,
		Scale:                ScaleModeCover,
		ScaleWidth:           100,
		ScaleHeight:          100,
		Crop:                 CropModeCenter,
		CropWidth:            100,
		CropHeight:           100,
		NormalizeOrientation: true,
		Quality:              0.9,
	}

	writer, err := os.OpenFile("new.png", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("Error writing: ", err)
	}

	if err = ProcessImage(reader, writer, spec); err != nil {
		log.Fatal("Error processing: ", err)
	}
}
