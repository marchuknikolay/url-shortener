package main

import (
	"fmt"

	"github.com/marchuknikolay/url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// init logger

	// init repository

	// init router

	// run server
}
