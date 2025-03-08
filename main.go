package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-audio/wav"
)

func main() {
	// Open our file
	file, err := os.Open("tmp/song.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Using wav library, create a decoder for our file.
	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		fmt.Println("Invalid WAV file")
		return
	}
	decoder.ReadInfo()

	// Run generate amplitudes function.
	amplitudes, err := GenerateAmplitudeData(decoder)
	if err != nil {
		fmt.Printf("Error generating amplitude data: %v\n", err)
		return
	}

	fmt.Print("Amplitudes: ", amplitudes)
}


