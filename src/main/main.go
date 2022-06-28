package main

import (
	"fmt"
	"os"

	"github.com/iturricf/rooftop-challenge/challenge"
)

func main() {
	if len(os.Args) < 2 || len(os.Args[1]) == 0 {
		fmt.Println("Missing email. Try adding your email as a parameter: rtchallenge your-email")
		os.Exit(1)
	}

	err := challenge.Solve(os.Args[1])
	if err != nil {
		fmt.Printf("error while trying to solve the riddle [err=%v]\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
