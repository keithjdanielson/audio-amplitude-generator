package main

import (
	"math"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func GenerateAmplitudeData(decoder *wav.Decoder) ([]float64, error) {
    // Configuration
    useLogScale := false           // Use dB scale
    resolutionMs := 256             // Time resolution in ms
		bufferSize := 8192
    
    // Calculate chunk size
    chunkSize := int(decoder.SampleRate) * resolutionMs / 1000
   	maxValue := math.Pow(2, float64(decoder.BitDepth-1))
    
    amplitudes := []float64{}
    buf := &audio.IntBuffer{
        Format: &audio.Format{
            NumChannels: int(decoder.NumChans),
            SampleRate:  int(decoder.SampleRate),
        },
        Data: make([]int, bufferSize),
    }
    
    chunkSum := 0.0
    sampleCount := 0
    
    for {
        n, err := decoder.PCMBuffer(buf)
        if err != nil || n == 0 {
            break
        }
        
        for i := 0; i < len(buf.Data); i += int(decoder.NumChans) {
            var frameAmplitude float64
            
            sum := 0.0
                validChannels := 0
                for ch := 0; ch < int(decoder.NumChans); ch++ {
                    idx := i + ch
                    if idx < len(buf.Data) {
                        sum += math.Abs(float64(buf.Data[idx]))
                        validChannels++
                    }
                }
                if validChannels > 0 {
                    frameAmplitude = sum / float64(validChannels)
                }
            
            // Normalize
            frameAmplitude /= maxValue
            
            // Apply log scale
            if useLogScale && frameAmplitude > 0 {
                frameAmplitude = 20 * math.Log10(frameAmplitude+0.0000001)
            }
            
            // Add to chunk
            chunkSum += frameAmplitude
            sampleCount++
            
            // If we've collected enough samples for a chunk, store average
            if sampleCount >= chunkSize {
                amplitudes = append(amplitudes, chunkSum/float64(sampleCount))
                chunkSum = 0
                sampleCount = 0
            }
        }
    }
    
    // Handle any remainder
    if sampleCount > 0 {
        amplitudes = append(amplitudes, chunkSum/float64(sampleCount))
    }
    
    return amplitudes, nil
}