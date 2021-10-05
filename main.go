package main

import (
	"fmt"
	"log"
	"os"
	"pngchunks/png"
)

func main() {
	file, err := os.Open("mario.png")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(file)
	defer file.Close()

	var img *png.Png
	img, err = png.NewPng(file)
	if err != nil {
		log.Fatal(err)
	}
	for {
		chunk, err := img.NextChunk()
		if err != nil {
			log.Fatal(err)
		}
		if string(chunk.Type) == "IEND" {
			break
		} else {
			fmt.Println(string(chunk.Type))
		}
	}

	//fmt.Printf("\nPNG Parameters\nWidth: %d \nHeight: %d\n",img.Parameters.Width, img.Parameters.Height)


}
