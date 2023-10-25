package main

import (
	"fmt"
	"sort"
	"sync"
	"toySearchEngine/internal"
)

func main() {
	wg := sync.WaitGroup{}
	searchResult := internal.Search("golang programming", &wg)
	sort.Slice(searchResult[:], func(i, j int) bool {
		return searchResult[i].Accuracy > searchResult[j].Accuracy
	})
	for _, result := range searchResult {
		fmt.Printf("Path: %s Accuracy: %d \n", result.Path, result.Accuracy)
	}

}
