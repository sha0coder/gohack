package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var R *Requests

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func loadWordlist(wordlist string) []string {
	var list []string

	file, err := os.Open(wordlist)
	check(err, "Can't load the wordlist")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	return list
}

func main() {
	goroutines := flag.Int("go", 1, "number of goroutines")
	url := flag.String("url", "", "target url")
	dbg := flag.Bool("dbg", false, "debug mode")
	param := flag.String("p", "", "parameter to test")
	flag.Parse()

	if *url == "" {
		fmt.Println("select an url -url or -help")
		os.Exit(1)
	}

	fmt.Println("loading ...")
	c := make(chan string, 6)
	indicators := loadWordlist("indicators.txt")
	params := loadWordlist("params.txt")
	payloads := loadWordlist("payloads.txt")
	R = NewRequests()

	if !strings.Contains(*url, "?") {
		*url = *url + "?"
	}

	for i := 0; i < *goroutines; i++ {
		go func(g int, c <-chan string, indicators []string, dbg bool) {
			for url := range c {
				html, code, _ := R.Get(url)
				if dbg {
					fmt.Printf("[%d] sz:%d  %s\n", code, len(html), url)
				}
				for _, i := range indicators {
					if strings.Contains(html, i) {
						fmt.Printf("%s  indicator: %s\n", url, i)
					}
				}
			}
		}(i, c, indicators, *dbg)
	}

	if *param == "" {
		for _, p := range params {
			for _, v := range payloads {
				c <- *url + "&" + p + "=" + v
			}
		}
	} else {
		parts := strings.Split(*url, *param+"=")
	
		if len(parts) < 2 {
			*url += "&"+*param+"=##"
		} else {
			value := strings.Split(parts[1], "&")
			*url = strings.ReplaceAll(*url, value[0], "##")
		}

		for _, v := range payloads {
			u := strings.ReplaceAll(*url, "##", v)
			c <- u
		}
	}


	var i int
	fmt.Printf("Scanning, press enter to interrupt.\n")
	fmt.Scanf("%d", &i)
	fmt.Printf("interrupted.")
}
