package search

import (
	"bufio"
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type Result struct {
	Phrase  string
	Line    string
	LineNum int64
	ColNum  int64
}

func FindAll(phrase, path string) (res []Result) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("error not opened file err => ", err)
		return res
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i := 0; i < len(lines); i++ {

		if strings.Contains(lines[i], phrase) {
			line := lines[i]
			r := Result{
				Phrase:  phrase,
				Line:    line,
				LineNum: int64(i + 1),
				ColNum:  int64(strings.Index(lines[i], phrase)) + 1,
			}

			res = append(res, r)
		}
	}
	return res
}

func All(root context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result)
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(root)

	for _, file := range files {
		wg.Add(1)
		go func(ctx context.Context, ch chan []Result, file string) {
			defer wg.Done()

			result := FindAll(phrase, file)

			if len(result) > 0 {
				ch <- result
			}

		}(ctx, ch, file)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	cancel()
	return ch
}


func Any(ctx context.Context, phrase string, files []string) <-chan Result {
	resultChan := make(chan Result)
	wg := sync.WaitGroup{}
	result := Result{}

	for i := 0; i < len(files); i++ {
		data, err := ioutil.ReadFile(files[i])
		if err != nil {
			log.Println("error while open file: ", err)
		}

		if strings.Contains(string(data), phrase) {
			res := FindAny(phrase, string(data))
			if (Result{}) != res {
				result = res
				break
			}
		}
	}
	log.Println("Найден: ", result)

	wg.Add(1)
	go func(ctx context.Context, ch chan<- Result) {
		defer wg.Done()
		if (Result{}) != result {
			ch <- result
		}
	}(ctx, resultChan)

	go func() {
		defer close(resultChan)
		wg.Wait()
	}()
	return resultChan
}

func FindAny(phrase, search string) (result Result) {
	for i, line := range strings.Split(search, "\n") {
		if strings.Contains(line, phrase) {
			return Result{
				Phrase:  phrase,
				Line:    line,
				LineNum: int64(i + 1),
				ColNum:  int64(strings.Index(line, phrase)) + 1,
			}
		}
	}
	return result
}
