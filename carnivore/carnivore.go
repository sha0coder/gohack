/*
	Carnivore url pentest
		- test new parameters not found in url
		- allow tests only spcecific parameter
		- many payloads
		- expert system recognize errors
		- gauss to false positive reduction
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"net/url"
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

func getUrlParamNames(surl string) []string {
	var params []string
	oUrl, err := url.Parse(surl)
	check(err, "invalid url")

	for k, _ := range oUrl.Query() {
		params = append(params, k)
	}

	return params
}

func changeUrlParam(surl string, param string, value string) string {
	oUrl, err := url.Parse(surl)
	check(err, "invalid url")

	query := oUrl.Query()
	for k, _ := range query {
		if k == param {
			query[k][0] = value
			break
		}
	}

	oUrl.RawQuery = query.Encode()
	return oUrl.String()
}

func addUrlParam(surl string, param string, value string) string {
	oUrl, err := url.Parse(surl)
	check(err, "invalid url")

	query := oUrl.Query()
	query[param] = []string{value}
	oUrl.RawQuery = query.Encode()
	return oUrl.String()
}



func injectPayloads(param, url string, new bool, params, payloads []string, curls chan<- string) {

	if param == "" {

		for _, p := range getUrlParamNames(url) {
			for _, v := range payloads {
				u := changeUrlParam(url, p, v)
				curls <- u
			}
		}

	} else {

		for _, v := range payloads {
			u := changeUrlParam(url, param, v)
			curls <- u
		}
	}

	if new {
		for _, p := range params {
			for _, v := range payloads {
				u := addUrlParam(url, p, v)
				curls <- u
			}
		}
	}

	close(curls)
}

func process(g int, c <-chan string, cres chan<- map[string]int, indicators []string, dbg bool, wg *sync.WaitGroup) {
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
		res := make(map[string]int)
		res[url] = len(strings.Split(html, "\n"))
		cres <- res
	}
	wg.Done()
}

func reduceFPs(cres <-chan map[string]int) {
	results := make(map[string]int)
	mean := 0.0
	elems := 0
	sum := 0.0

	for r := range cres {
		for k, v := range r {
			results[k] = v
			mean += float64(v)
			elems += 1
		}
	}

	mean = mean / float64(elems)

	for _, v := range results {
		sum += math.Pow(float64(v)-mean, 2)
	}

	err := math.Sqrt(sum / float64(elems))

	fmt.Printf("gauss mean:%f err:%f\n", mean, err)
	fmt.Println(" --- results ---")
	for url, lines := range results {
		if float64(lines) < mean-err || float64(lines) > mean+err {
			fmt.Printf("[%d] %s\n", lines, url)
		}
	}
	fmt.Println("done.")
	os.Exit(1)
}

func main() {
	goroutines := flag.Int("go", 1, "number of goroutines")
	url := flag.String("url", "", "target url")
	dbg := flag.Bool("dbg", false, "debug mode")
	new := flag.Bool("new", false, "try guessing new parameters, slow.")
	param := flag.String("p", "", "choose parameter to test otherwise test all.")
	flag.Parse()

	if *url == "" {
		fmt.Println("select an url -url or -help")
		os.Exit(1)
	}

	fmt.Println("loading ...")
	curls := make(chan string, 6)
	cres := make(chan map[string]int, 6)

	indicators := loadWordlist("indicators.txt")
	params := loadWordlist("params.txt")
	payloads := loadWordlist("payloads.txt")
	R = NewRequests()
	var wg sync.WaitGroup
	wg.Add(*goroutines)

	if !strings.Contains(*url, "?") {
		*url = *url + "?"
	}

	for i := 0; i < *goroutines; i++ {
		go process(i, curls, cres, indicators, *dbg, &wg)
	}

	go injectPayloads(*param, *url, *new, params, payloads, curls)

	fmt.Println("pentesting ...")

	go reduceFPs(cres)

	wg.Wait()
	close(cres)

	var i int
	for {
		fmt.Scanf("%d", &i)
	}
}
