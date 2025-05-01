package main

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

func main() {
	pipeReader1, pipeWriter1 := io.Pipe()
	pipeReader2, pipeWriter2 := io.Pipe()
	multiWriter := io.MultiWriter(pipeWriter1, pipeWriter2)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer pipeWriter2.Close()
		defer pipeWriter1.Close()
		for x := 1; x < 5; x++ {
			io.Copy(multiWriter, strings.NewReader(fmt.Sprintf("Hello%d", x)))
		}
	}()

	wg.Add(1)
	pipe1TransformedReader, pipe1TransformedWriter := io.Pipe()
	go func() {
		defer wg.Done()
		defer pipe1TransformedWriter.Close()
		data, _ := io.ReadAll(pipeReader1)
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
		data, _ := io.ReadAll(pipeReader2)
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
