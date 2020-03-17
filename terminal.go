package main

import (
	"log"
	"os"
	"os/exec"
)

func runTerminal(term *exec.Cmd) error {
	term.Stdin = os.Stdin
	return term.Run()
}

func initCooked() func() *exec.Cmd {
	return func() *exec.Cmd {
		return exec.Command("stty", "cbreak", "-echo")
	}
}

func initCBreak() func() *exec.Cmd {
	return func() *exec.Cmd {
		return exec.Command("stty", "-cbreak", "echo")
	}
}

func recoverFatal(msg string, err error) {
	if err != nil {
		log.Fatalln("Error activating cbreak mode:", err)
	}
}

//Initialise - init game
func Initialise() {
	cbTerm := initCooked()()
	err := runTerminal(cbTerm)
	recoverFatal("Error activating cbreak mode:", err)
}

//Cleanup - clean game
func Cleanup() {
	cookedTerm := initCBreak()()
	err := runTerminal(cookedTerm)
	recoverFatal("Error activating cooked mode:", err)
}
