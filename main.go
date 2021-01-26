package main

import (
	"os"

	"github.com/manifoldfinance/greyelk/cmd"

	_ "github.com/manifoldfinance/greyelk/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
