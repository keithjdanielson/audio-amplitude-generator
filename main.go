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

	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		fmt.Println("Invalid WAV file")
		return
	}

	fmt.Println("WAV Header Information:")
	fmt.Printf("Audio Format: %d (1 = PCM)\n", decoder.WavAudioFormat)
	fmt.Printf("Number of Channels: %d\n", decoder.NumChans)
	fmt.Printf("Sample Rate: %d Hz\n", decoder.SampleRate)
	fmt.Printf("Bits Per Sample: %d-bit\n", decoder.BitDepth)
	fmt.Printf("Average Bytes Per Second: %d\n", decoder.AvgBytesPerSec)
}
