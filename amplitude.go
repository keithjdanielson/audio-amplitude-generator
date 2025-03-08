package audioampgenerator

import (
	"math"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// Config holds options for amplitude generation
type Config struct {
	UseLogScale   bool // Use dB scale
	ResolutionMs  int  // Time resolution in milliseconds
	BufferSize    int  // Buffer size for processing
}

// Default configuration
func DefaultConfig() Config {
	return Config{
		UseLogScale:   false,
		ResolutionMs:  256,
		BufferSize:    8192,
	}
}

//  extracts amplitude data from audio files
type Analyzer struct {
	config Config
}

// creates a new Analyzer (with the configuration)
func NewAnalyzer(config Config) *Analyzer {
	return &Analyzer{
		config: config,
	}
}

// creates a new Analyzer with default config
func NewDefaultAnalyzer() *Analyzer {
	return NewAnalyzer(DefaultConfig())
}

// extracts amplitude data from a WAV decoder
func (a *Analyzer) GenerateAmplitudeData(decoder *wav.Decoder) ([]float64, error) {
	// Calculate chunk size
	chunkSize := int(decoder.SampleRate) * a.config.ResolutionMs / 1000
	maxValue := math.Pow(2, float64(decoder.BitDepth-1))
	
	amplitudes := []float64{}
	buf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: int(decoder.NumChans),
			SampleRate:  int(decoder.SampleRate),
		},
		Data: make([]int, a.config.BufferSize),
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
			if a.config.UseLogScale && frameAmplitude > 0 {
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