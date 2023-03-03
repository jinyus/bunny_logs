package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type kv struct {
	Key   string
	Value int
}

func summarize(tracker map[string]int, name string, count int) {
	// sort by value
	var ss []kv
	for k, v := range tracker {
		ss = append(ss, kv{k, v})
	}

	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	fmt.Println("\nTop", count, name)
	for i := 0; i < min(count, len(ss)); i++ {
		fmt.Printf("%s: %d\n", ss[i].Key, ss[i].Value)
	}

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func parseLogLine(line string) (ip, referer, useragent, country string) {
	parts := strings.Split(line, "|")

	if len(parts) < 9 {
		fmt.Println("error parsing line:", line)
		return
	}

	ip = parts[5]
	referer = parts[6]
	useragent = parts[9]
	country = parts[11]
	return
}

func scanFile(filename string, ch chan string, done chan int) {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		ch <- sc.Text() // GET the line string
	}
	done <- 1
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}
}
