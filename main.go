package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Jamess-Lucass/interpreter-go/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! Please start typing commands\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
