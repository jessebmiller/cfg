package main

import (
	"app/cfg"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println(cfg.Get("firstKey", "default One"))
	fmt.Println(cfg.Find("PATH"))
	fmt.Println(cfg.Get("anotherKey", "default Two"))
	fmt.Println(cfg.Get("firstKey", "different default"))
	http.ListenAndServe(":8000", nil)
}
