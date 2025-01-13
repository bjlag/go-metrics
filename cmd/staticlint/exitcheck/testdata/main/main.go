package main

import (
	"os"
)

func exitCheckFunc() {
	defer os.Exit(1) // want "os.Exit in main package"

	os.Exit(1) // want "os.Exit in main package"
}
