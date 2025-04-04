package heavytask

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	cleanup := func() {
		_ = os.Remove("test.txt")
	}

	t.Cleanup(cleanup)

	t.Run("Task completes within timeout", func(t *testing.T) {
		cleanup()

		filename, err := Run(ctx)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if filename != "test.txt" {
			t.Errorf("Expected filename 'test.txt', got: %s", filename)
		}

		_, err = os.Stat(filename)
		if os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", filename)
		}
	})

	t.Run("Task exceeds timeout with short context", func(t *testing.T) {
		cleanup()

		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
		defer cancel()

		filename, err := Run(timeoutCtx)

		if err == nil {
			t.Error("Expected timeout error, got nil")
		}

		if filename != "" {
			t.Errorf("Expected empty filename, got: %s", filename)
		}
	})

	t.Run("Context times out", func(t *testing.T) {
		cleanup()

		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
		defer cancel()

		resultCh := make(chan struct {
			filename string
			err      error
		})

		go func() {
			filename, err := Run(timeoutCtx)
			resultCh <- struct {
				filename string
				err      error
			}{filename, err}
		}()

		time.Sleep(time.Millisecond * 10)

		result := <-resultCh

		if result.err == nil {
			t.Error("Expected timeout error, got nil")
		}

		if result.filename != "" {
			t.Errorf("Expected empty filename, got: %s", result.filename)
		}
	})
}

func TestDoHeavyTaskAndWritetoFile(t *testing.T) {
	ctx := context.Background()

	cleanup := func() {
		_ = os.Remove("test.txt")
	}

	t.Cleanup(cleanup)

	t.Run("Successful file creation", func(t *testing.T) {
		cleanup()

		filename, err := doHeavyTaskAndWritetoFile(ctx)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if filename != "test.txt" {
			t.Errorf("Expected filename 'test.txt', got: %s", filename)
		}

		fileInfo, err := os.Stat(filename)
		if os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", filename)
		}

		if fileInfo.Size() == 0 {
			t.Error("Expected file to have content, but it's empty")
		}
	})

	t.Run("Context timeout", func(t *testing.T) {
		cleanup()

		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
		defer cancel()

		filename, err := doHeavyTaskAndWritetoFile(timeoutCtx)

		if err == nil {
			t.Error("Expected error due to timeout context, got nil")
		}

		if filename != "" {
			t.Errorf("Expected empty filename, got: %s", filename)
		}

		_, err = os.Stat("test.txt")
		if !os.IsNotExist(err) {
			t.Error("Expected file not to exist, but it does")
		}
	})
}
