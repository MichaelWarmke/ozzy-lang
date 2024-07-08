package main

import (
	"fmt"
	"os"
	"os/user"
	"ozzy/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! this is the ozzy programming language!\n", user.Username)
	fmt.Printf("Feel Free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
