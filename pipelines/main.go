package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Record struct {
	Weight float64
	Height float64
}

func main() {
	parseInput := make(chan []string)
	convertInput := make(chan Record)
	convertOutput := make(chan Record)
	encodeOutput := make(chan []byte)
	done := make(chan struct{})

	go pipeLineStage(parseInput, convertInput, parse)
	go pipeLineStage(convertInput, convertOutput, convert)
	go pipeLineStage(convertOutput, encodeOutput, encode)

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	data := csv.NewReader(file)

	go func() {
		for data := range encodeOutput {
			fmt.Println(string(data))
		}
		close(done)
	}()

	_, err = data.Read()
	if err != nil {
		log.Fatalln(err)
	}

	for {
		row, err := data.Read()
		if err != nil {
			if err == io.EOF {
				close(parseInput)
				break
			}
			panic(err)
		}

		parseInput <- row
	}

	<-done
}

func pipeLineStage[I any, O any](input <-chan I, output chan<- O, stageTransform func(I) O) {
	defer close(output)

	for i := range input {
		output <- stageTransform(i)
	}
}

func parse(in []string) Record {
	rec, err := newRecord(in)
	if err != nil {
		panic(err)
	}

	return rec
}

func convert(in Record) Record {
	in.Weight = in.Weight * 0.454
	in.Height = in.Height * 2.54

	return in
}

func encode(input Record) []byte {
	data, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	return data
}

func newRecord(in []string) (Record, error) {
	var rec Record
	var err error
	rec.Weight, err = strconv.ParseFloat(in[0], 64)
	if err != nil {
		return Record{}, err
	}

	rec.Height, err = strconv.ParseFloat(in[0], 64)
	if err != nil {
		return Record{}, err
	}

	return rec, nil
}
