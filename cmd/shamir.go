package main

import (
	"fmt"
	"os"

	"shamir/pkg/cmd"
	"shamir/pkg/utils/log"
)

func main() {
	command, err := cmd.NewMainCommand()
	handleError(err)
	err = command.Execute()
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		log.Error(err)
		fmt.Println(err)
		os.Exit(1)
	}
}
