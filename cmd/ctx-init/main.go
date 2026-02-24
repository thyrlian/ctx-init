package main

import (
	"fmt"
	"os"

	"github.com/thyrlian/ctx-init/internal/cli"
)

func main() {
	opts, err := cli.ParseAndValidate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("------------")
	fmt.Println("=== 🛠️ CTX-INIT 🛠️ ===\n")
	fmt.Printf("Manifest: %s\n", opts.ManifestPath)
	fmt.Printf("Preset:   %s\n", opts.Preset)
	fmt.Println("------------")
}
