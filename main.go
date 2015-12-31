package main

import (
	"flag"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"os"
	"regexp"
	"strings"
)

var (
	emailAddress = flag.String("email", "", "Your email address")
	debug        = flag.Bool("debug", true, "Enable debug messages")
)

func main() {
	fmt.Println("Let's email!")
	flag.Parse()
	validateCommandlineFlags()
	msgS := SetupService()
	resp, err2 := msgS.List(*emailAddress).Do()
	if err2 != nil {
		panic(err2)
	}
	messages := resp.Messages
	ids := []string{}
	for _, msg := range messages {
		ids = append(ids, msg.Id)
		fmt.Println("Message subject ", msg.Header)
	}

	for _, id := range ids {
		result, err3 := msgS.Get(*emailAddress, id).Format("full").Do()
		if err3 != nil {
			panic(err3)
		}
		e, err4 := NewEmailMessage(result)
		if err4 != nil {
			panic(err4)
		}
		e.Print()
	}

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

func printParts(part *gmail.MessagePart, indentation int) {
	for i, p := range part.Parts {
		fmt.Printf("%s", strings.Repeat("\t", indentation))
		fmt.Printf("Part %d has type %s with %d subparts.\n", i, p.MimeType, len(p.Parts))
		decodedBody := strings.Replace(strings.Replace(decode(p.Body.Data), "\n", "", -1), "\r", "", -1)
		fmt.Printf("%s", strings.Repeat("\t", indentation))
		fmt.Printf("Text of length %d, a snippet: %s\n", len(decodedBody), decodedBody[:len(decodedBody)%50])
		printParts(p, indentation+1)
	}
}
