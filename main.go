package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/mileusna/useragent"
	path_helpers "github.com/moisespsena-go/path-helpers"
)

// log line HIT|200|1677688609214|2802|1140308|1.1.1.1|https://www.google.com/|https://www.example.com|MI|Mozilla/5.0 (iPhone; CPU iPhone OS 16_0_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1|d2af1ddaf9efe0fe3cc9739c40413|US

// process log file
// extract ip, user agent and country from each line
// extract browser and os from user agent
// summarize by country, browser and os

type StringSlice []string

func (i *StringSlice) String() string {
	return "list of strings"
}

func (i *StringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var files StringSlice
var domains StringSlice

func main() {
	flag.Var(&files, "f", "log file to process (can be repeated)")
	flag.Var(&domains, "d", "domain to exclude from referrers (can be repeated)")
	flag.Parse()

	if len(files) == 0 {
		fmt.Println("No file specified")
		os.Exit(1)
	}

	for _, file := range files {
		if !path_helpers.IsExistingRegularFile(file) {
			fmt.Println("File doesn't exists: ", file)
			os.Exit(2)
		}
	}

	cTracker := make(map[string]int)
	refTracker := make(map[string]int)
	ipTracker := make(map[string]int)

	// 2 users can have the same ip so ip + user agent will be used to track unique users
	ipUATracker := make(map[string]int)

	broswerTracker := make(map[string]int)

	broswerVersionTracker := make(map[string]int)

	osTracker := make(map[string]int)

	osVersionTracker := make(map[string]int)

	iOSVersionTracker := make(map[string]int)

	androidVersionTracker := make(map[string]int)

	windowsVersionTracker := make(map[string]int)

	ch := make(chan string, 1_000)

	doneCH := make(chan int, len(files))

	for _, file := range files {
		go scanFile(file, ch, doneCH)
	}

	count := len(files)
	go func() {
		for range doneCH {
			count--
			if count == 0 {
				close(ch)
			}
		}
	}()

	for s := range ch {
		ip, ref, ua, country := parseLogLine(s)

		cTracker[country]++
		ipTracker[ip]++

		refParsed, err := url.Parse(ref)

		if err == nil {
			host := refParsed.Host
			if !contains(domains, host) {
				refTracker[host]++
			}
		}
		ipUAKey := ip + ua

		if ipUATracker[ipUAKey] > 0 {
			continue
		}

		ipUATracker[ipUAKey]++

		userAgent := useragent.Parse(ua)

		broserVersionKey := userAgent.Name + "_" + userAgent.Version
		broswerVersionTracker[broserVersionKey]++

		osVersionKey := userAgent.OS + "_" + userAgent.OSVersion
		osVersionTracker[osVersionKey]++

		broswerTracker[userAgent.Name]++

		osTracker[userAgent.OS]++

		if userAgent.OS == "iOS" {
			iOSVersionTracker[userAgent.OSVersion]++
		}

		if userAgent.OS == "Android" {
			androidVersionTracker[userAgent.OSVersion]++
		}

		if userAgent.OS == "Windows" {
			windowsVersionTracker[userAgent.OSVersion]++
		}
	}

	summarize(cTracker, "countries", 10)
	summarize(ipTracker, "IPs", 10)
	summarize(refTracker, "Referers", 30)
	summarize(broswerTracker, "browsers", 10)
	summarize(osTracker, "OS", 10)
	summarize(broswerVersionTracker, "browser versions", 10)
	summarize(osVersionTracker, "os versions", 10)
	summarize(iOSVersionTracker, "iOS versions", 30)
	summarize(androidVersionTracker, "Android versions", 20)
	summarize(windowsVersionTracker, "Windows versions", 20)
}
