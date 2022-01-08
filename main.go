package main

import (
	"flag"
	"log"
)

func main() {
	list := flag.Bool("l", false, "List aliases")
	alias := flag.String("a", "", "Add alias")
	flag.Parse()
	if flag.NArg() <= 0 {
		log.Fatalln("User key is required")
	}
	if !*list && len(*alias) <= 0 {
		log.Fatalln("Received illegal option")
	}
	key := flag.Arg(0)
	_ = key
}
