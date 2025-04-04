package heavytask

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"
)

type Response struct {
	Filename string
	Err      error
}

func Run(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	resch := make(chan Response)

	go func() {
		filename, err := doHeavyTaskAndWritetoFile(ctx)
		resch <- Response{Filename: filename, Err: err}
	}()

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("timedout while waiting for response")
	case resp := <-resch:
		return resp.Filename, resp.Err
	}

}

func doHeavyTaskAndWritetoFile(ctx context.Context) (string, error) {
	select {
	case <-time.After(time.Millisecond * 100):
	case <-ctx.Done():
		return "", ctx.Err()
	}

	filename := "test.txt"

	file, err := os.Create(filename)

	if err != nil {
		return "", err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	writer.WriteString("Hello, World!\n")

	writer.Flush()

	return filename, nil
}
