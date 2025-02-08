/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/gsoares85/code-guardian/cmd"
	"github.com/gsoares85/code-guardian/config"
)

func main() {
	config.LoadConfig()

	cmd.Execute()
}
