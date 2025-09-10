// main.go
package main

import (
	"github.com/soyomarvaldezg/neuron-cli/internal/cmd" // Adjusted import path
)

func main() {
	// All the magic happens in the cmd package.
	// This keeps our main function clean and simple.
	cmd.Execute()
}
