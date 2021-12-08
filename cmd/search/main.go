package main

import (
	"context"
	"log"

	"github.com/asusg74/search/pkg/search"
)

func main() {
	files := []string{
		"C:/project/20/search/data/test1.txt",
		"C:/project/20/search/data/test2.txt",
		"C:/project/20/search/data/test3.txt",
	}
	
	ch := search.All(context.Background(), "test", files)

	s, ok := <-ch

	if !ok {
		log.Printf(" function All error => %v", ok)
	}
	for _, r := range s {
		log.Println("=======>>>>>", r)
	}
	
}
