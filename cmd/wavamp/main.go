package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-audio/wav"
	audioampgenerator "github.com/keithjdanielson/audio-amplitude-generator"
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

	// Create analyzer with default configuration
	analyzer := audioampgenerator.NewDefaultAnalyzer()
	/*
	analyzer := wavamp.NewAnalyzer(wavamp.Config{
		UseLogScale:  false,
		ResolutionMs: 256,
		BufferSize:   8192,
	})
	*/

	amplitudes, err := analyzer.GenerateAmplitudeData(decoder)
	if err != nil {
		fmt.Printf("Error generating amplitude data: %v\n", err)
		return
	}
	fmt.Print("Amplitudes: ", amplitudes)
}