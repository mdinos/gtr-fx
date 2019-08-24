package main

import (
	"fmt"
	"os"
	"runtime"
	"log"
	"github.com/mjibson/go-dsp/wav"
)

const (
	Threads = 8
)

func modifySingleValue(r float32) (w float32) {
	return r + 0.1
}

func modifyArray(r []float32) (w []float32) {
	c := make(chan int, Threads)
	runtime.GOMAXPROCS(Threads)
	chunk := len(r) / Threads
	for i := 0; i < Threads; i++ {
		go func(start int) {
			end := start + chunk

			if end > len(r) {
				end = len(r)
			}

			for j := start; j < end; j = j + 1 {
				r[j] = modifySingleValue(r[j])
			}
			c <- 1
		}(i * chunk)
	}
	for i := 0; i < Threads; i++ {
		<-c
	}
	return r
}

func main() {
	wavReader, err := os.Open("wav/sample.wav")
	if err != nil {
		fmt.Println(err)
		log.Panicln("Failed to open wav file.")
	}
	file, err := wav.New(wavReader)
	if err != nil {
		log.Panicln("Failed to convert wav reader to wav struct")
	}
	floats, err := file.ReadFloats(file.Samples)
	if err != nil {
		fmt.Println(err)
		log.Panicln("Failed to read floats from wav struct")
	}
	modified := modifyArray(floats)
	fmt.Println(modified[3000:3010])
}