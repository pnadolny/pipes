package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

func main() {
	// Pipe creates a synchronous in-memory pipe
	pipe1reader, pipe1Writer := io.Pipe()
	pipe2reader, pipe2Writer := io.Pipe()
	//MultiWriter creates a writer that duplicates its writes to all the provided writers,
	multiWriter := io.MultiWriter(pipe1Writer, pipe2Writer)

	// A WaitGroup waits for a collection of goroutines to finish
	var wg sync.WaitGroup

	// Writer goroutine: writes input to MultiWriter
	wg.Add(1)
	go func() {
		defer wg.Done()

		// MultiWriter writes to both for me, so close them when done.
		defer pipe2Writer.Close()
		defer pipe1Writer.Close()

		input := "Hello world test stream"
		_, err := io.Copy(multiWriter, strings.NewReader(input))
		if err != nil {
			log.Println(err)
		}
		input2 := "Hi again."
		_, err = io.Copy(multiWriter, strings.NewReader(input2))
		if err != nil {
			log.Println(err)
		}

	}()

	wg.Add(1)
	// Transformer goroutine 1: reads from pipe1, transforms to uppercase, writes to new pipe
	pipe1TransformedReader, pipe1TransformedWriter := io.Pipe()
	go func() {
		defer wg.Done()
		defer pipe1TransformedWriter.Close()
		data, err := io.ReadAll(pipe1reader)
		if err != nil {
			log.Println(err)
		}
		// Transform: Convert to uppercase
		transformed := strings.ToUpper(string(data))
		_, err = pipe1TransformedWriter.Write([]byte(transformed))
		if err != nil {
			fmt.Printf("Pipe 1 transform write error: %v\n", err)
			return
		}
	}()

	wg.Add(1)

	pipe2TransformedReader, pipe2TransformedWriter := io.Pipe()
	go func() {
		defer wg.Done()
		defer pipe2TransformedWriter.Close()
		data, err := io.ReadAll(pipe2reader)
		if err != nil {
			log.Println(err)
		}
		// Transform: Add prefix
		transformed := "PREFIX: " + string(data)
		_, err = pipe2TransformedWriter.Write([]byte(transformed))
		if err != nil {
			fmt.Printf("Pipe 1 transform write error: %v\n", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, err := io.ReadAll(pipe1TransformedReader)
		if err != nil {
			fmt.Printf("Pipe 1 transformed read error: %v\n", err)
			return
		}
		fmt.Printf("Pipe 1 transformed: %s\n", string(data))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, err := io.ReadAll(pipe2TransformedReader)
		if err != nil {
			fmt.Printf("System: Pipe 2 transformed read error: %v\n", err)
			return
		}
		fmt.Printf("Pipe 2 transformed: %s\n", string(data))
	}()
	wg.Wait()

}
