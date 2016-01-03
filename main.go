package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

var (
	emailAddress = flag.String("email", "", "Your email address")
	debug        = flag.Bool("debug", false, "Enable debug messages")
)

func main() {
	fmt.Println("Let's email!")
	flag.Parse()
	validateCommandlineFlags()
	ui := NewUI()
	ui.ListenForChanges()
	ui.Init()
	defer ui.Close()
	msgS := SetupService()
	coord := NewGmailCoordinatorWithSubscriber(msgS, ui.Sub)
	go coord.FetchMessages()
	ui.Loop()
}

func debugPrint(msg ...interface{}) {
	if *debug {
		fmt.Println(msg)
	}
}

func debugPrintf(formatStr string, data ...interface{}) {
	if *debug {
		fmt.Printf(formatStr, data)
	}
}

func validateCommandlineFlags() {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !Re.MatchString(*emailAddress) {
		fmt.Println("Oops! It looks like the email address you provided isn't valid. Please check it.")
		os.Exit(0)
	}
}
