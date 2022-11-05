package main

import (
	"fmt"
	"os"

	cmd2 "shamir/pkg/cmd"
	"shamir/pkg/utils/log"
)

func main() {
	cmd := cmd2.NewCommand()
	err := cmd.Execute()
	if err != nil {
		log.Error("run command failed:%v", err)
		fmt.Println(err)
		os.Exit(1)
	}
}
