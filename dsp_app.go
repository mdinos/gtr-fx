package main

import (
	"fmt"
	"github.com/youpy/go-wav"
	"io"
	"math"
	"os"
	"time"
)

func distortion(l int, r int, boost float64, drive float64, tone float64, mix float64) (lout int, rout int, err error) {
	gain := ((boost / 100) * 100) + 1
	a := math.Sin(((drive + 1) / 101) * (math.Pi / 2))
	k := 2 * a / (1 - a)

	gainedL := float64(l) * gain
	gainedR := float64(r) * gain

	drivenL := (((1 + k) * gainedL) / (1 + k*math.Abs(gainedL))) * math.Abs(float64(l))
	drivenR := (((1 + k) * gainedR) / (1 + k*math.Abs(gainedR))) * math.Abs(float64(r))

	return int(drivenL), int(drivenR), nil
}

func readAndModify(r *wav.Reader) (w []wav.Sample, f *wav.WavFormat) {
	var outSamples []wav.Sample
	for {
		samples, err := r.ReadSamples()
		if err == io.EOF {
			break
		}

		for _, sample := range samples {
			left, right, err := distortion(r.IntValue(sample, 0), r.IntValue(sample, 1), 70.0, 80, 50.0, 0.5)
			if err != nil {
				panic("idk")
			}

			outSamples = append(outSamples, wav.Sample{
				Values: [2]int{left, right},
			})
		}
	}
	format, err := r.Format()
	if err != nil {
		fmt.Println("Error reading format.")
	}
	return outSamples, format
}

func main() {
	file, err := os.Open("wav/sample.wav")
	if err != nil {
		fmt.Println("Error opening wav file")
	}
	t0 := time.Now()
	outSamples, format := readAndModify(wav.NewReader(file))
	sampleRate := format.SampleRate
	bitsPerSample := format.BitsPerSample
	numChannels := format.NumChannels
	t1 := time.Now()
	fmt.Println(t1.Sub(t0))
	newfile, err := os.Create("out/sample.wav")
	writer := wav.NewWriter(newfile, uint32(len(outSamples)), numChannels, sampleRate, bitsPerSample)
	t2 := time.Now()
	err = writer.WriteSamples(outSamples)
	if err != nil {
		fmt.Println("Error occurred writing to file.")
	}
	t3 := time.Now()
	fmt.Println(t3.Sub(t2))
}
