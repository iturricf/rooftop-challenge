package main

import (
	"fmt"
	"os"

	"github.com/iturricf/rooftop-challenge/challenge"
)

const login = "iturri.cf+rt@gmail.com"

func main() {
	err := challenge.SolveWithLogin(login)
	if err != nil {
		fmt.Println("error while trying to solve the riddle [err=%v]", err)
		os.Exit(1)
	}

	os.Exit(0)
}
