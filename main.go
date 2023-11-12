package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"log"
	"net/http"

	"github.com/qeesung/image2ascii/convert"
)

func main() {
	flag.Parse()

	pokemon := flag.Arg(0)

	if pokemon == "" {
		log.Fatal("Please provide a pokemon name")
	}

	api := "https://pokeapi.co/api/v2/pokemon"
	endpoint := fmt.Sprintf("%s/%s", api, pokemon)

	resp, err := http.Get(endpoint)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Pokemon not found")
	}

	defer resp.Body.Close()

	data := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		log.Fatal(err)
	}

	image_url := data["sprites"].(map[string]interface{})["front_default"]

	image_request, err := http.Get(image_url.(string))

	if err != nil {
		log.Fatal(err)
	}

	defer image_request.Body.Close()

	buffer := make([]byte, 1024*1024*10)
	size, err := image_request.Body.Read(buffer)

	if err != nil || size == 0 {
		log.Fatal(err)
	}

	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 70
	convertOptions.FixedHeight = 30

	// convert buffer to png
	img, _, err := image.Decode(bytes.NewReader(buffer))

	if err != nil {
		log.Fatal(err)
	}

	converter := convert.NewImageConverter()
	asciiString := converter.Image2ASCIIString(img, &convertOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(asciiString)
}
