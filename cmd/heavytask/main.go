package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dliluashvili/heavytask"
)

func main() {
	ctx := context.Background()

	filename, err := heavytask.Run(ctx)

	if err != nil {
		fmt.Println("Error:", err)
	}

	time.Sleep(time.Second * 5)

	fmt.Println("Program completed successfully")

	fmt.Println("Filename:", filename)
}
