package main

import (
	"github.com/filipweidemann/demo-controller/internal"
	"os"
)

func main() {
	rc := internal.Run()
	os.Exit(rc)
}
