package main

import (
	"app/cfg"
	"fmt"
	"log"
)

// an example of using cfg
func main() {
	greeter, _ := cfg.Find("greeter")
	_, err := cfg.Valid()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf(
		"%s from %s",
		cfg.Get("greeting", "Hello!"),
		greeter,
	))
}
