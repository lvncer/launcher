package main

import (
	"fmt"
	"os"

	"launcher/internal/launcher"
)

func main() {
	if err := launcher.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
