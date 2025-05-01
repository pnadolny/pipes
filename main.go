package main

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

func main() {
	pipe1reader, pipe1Writer := io.Pipe()
	pipe2reader, pipe2Writer := io.Pipe()
	multiWriter := io.MultiWriter(pipe1Writer, pipe2Writer)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer pipe2Writer.Close()
		defer pipe1Writer.Close()
		for x := 1; x < 5; x++ {
			io.Copy(multiWriter, strings.NewReader("Hello world test stream"))
		}
	}()

	wg.Add(1)
	pipe1TransformedReader, pipe1TransformedWriter := io.Pipe()
	go func() {
		defer wg.Done()
		defer pipe1TransformedWriter.Close()
		data, _ := io.ReadAll(pipe1reader)
		pipe1TransformedWriter.Write([]byte(strings.ToUpper(string(data))))
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, _ := io.ReadAll(pipe1TransformedReader)
		fmt.Printf("Pipe 1 transformed: %s\n", string(data))
	}()

	wg.Add(1)
	pipe2TransformedReader, pipe2TransformedWriter := io.Pipe()
	go func() {
		defer wg.Done()
		defer pipe2TransformedWriter.Close()
		data, _ := io.ReadAll(pipe2reader)
		pipe2TransformedWriter.Write([]byte("PREFIX: " + string(data)))
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, _ := io.ReadAll(pipe2TransformedReader)
		fmt.Printf("Pipe 2 transformed: %s\n", string(data))
	}()
	wg.Wait()

}
