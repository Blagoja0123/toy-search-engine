package main

import (
	"fmt"
	"toyBrowser/internal"
)

func main() {
	searchResult := internal.Search("golang programming")

	for _, result := range searchResult {
		fmt.Printf("Path: %s Accuracy: %d \n", result.Path, result.Accuracy)
	}

}
