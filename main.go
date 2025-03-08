package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Open our file
	file, err := os.Open("tmp/song.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// First, we need to read the header of the file.
	headerBuffer := make([]byte, 44)
	n1, err := file.Read(headerBuffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d bytes: %s\n", n1, string(headerBuffer[:n1]))
}
