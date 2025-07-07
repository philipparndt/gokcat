/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/philipparndt/go-logger"
	"gokcat/cmd"
	"os"
)

func main() {
	logger.LogTo(os.Stderr)
	cmd.Execute()
}
