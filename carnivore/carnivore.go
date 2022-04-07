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
	"net/url"
	"os"
	"strings"
	"sync"
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

func changeUrlParam(surl, param, value string) string {
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

func addUrlParam(surl, param, value string) string {
	oUrl, err := url.Parse(surl)
	check(err, "invalid url")

	query := oUrl.Query()
	query[param] = []string{value}
	oUrl.RawQuery = query.Encode()
	return oUrl.String()
}

func getPostParmNames(post string) []string {
	var params []string

	spl := strings.Split(post, "&")
	for _, p := range spl {
		pv := strings.Split(p, "=")
		params = append(params, pv[0])
	}

	return params
}

func changePostParam(post, param, value string) string {
	var newPost string

	spl := strings.Split(post, "&")
	for _, p := range spl {
		pv := strings.Split(p, "=")

		if pv[0] == param {
			newPost += pv[0] + "=" + value + "&"
		} else {
			newPost += pv[0] + "=" + pv[1] + "&"
		}
	}

	return newPost
}

func addPostParam(post, param, value string) string {
	return post + "&" + param + "=" + value
}

func injectPayloads(param, url, post string, new bool, params, payloads []string, curls chan<- string) {

	if param == "" {

		// pentest all the url param
		fmt.Println("pentesting url params")
		for _, p := range getUrlParamNames(url) {
			for _, v := range payloads {
				u := changeUrlParam(url, p, v)
				if post == "" {
					curls <- u
				} else {
					curls <- u + "#POST#" + post
				}
			}
		}

		if post != "" {
			fmt.Println("pentesting post params")
			for _, p := range getPostParmNames(post) {
				for _, v := range payloads {
					pst := changePostParam(post, p, v)
					curls <- url + "#POST#" + pst
				}
			}
		}

	} else {

		if strings.Contains(url, param+"=") {
			fmt.Println("pentesting selected param on url")
			for _, v := range payloads {
				u := changeUrlParam(url, param, v)
				if post == "" {
					curls <- u
				} else {
					curls <- u + "#POST#" + post
				}
			}

		} else if strings.Contains(post, param+"=") {
			fmt.Println("pentest selected param on post")
			for _, v := range payloads {
				newPost := changePostParam(post, param, v)
				curls <- url + "#POST#" + newPost
			}

		} else {
			fmt.Println("param not found.")
			os.Exit(1)
		}
	}

	if new {
		fmt.Println("pentesting new possible params on url")
		for _, p := range params {
			for _, v := range payloads {
				u := addUrlParam(url, p, v)
				if post == "" {
					curls <- u
				} else {
					curls <- u + "#POST#" + post
				}
			}
		}

		if post != "" {
			fmt.Println("pentesting new possible params on post")
			for _, p := range params {
				for _, v := range payloads {
					newPost := addPostParam(post, p, v)
					curls <- url + "#POST#" + newPost
				}
			}
		}
	}

	close(curls)
}

func process(g int, c <-chan string, cres chan<- map[string]int, indicators []string, dbg bool, wg *sync.WaitGroup) {
	for url := range c {
		var html string
		var code int
		var isPost bool
		var urlPost []string

		if strings.Contains(url, "#POST#") {
			isPost = true
			urlPost = strings.Split(url, "#POST#")

			html, code, _ = R.Post(urlPost[0], urlPost[1])
			if dbg {
				fmt.Printf("[%d] sz:%d  %s %s\n", code, len(html), urlPost[0], urlPost[1])
			}
		} else {
			isPost = false
			html, code, _ = R.Get(url)
			if dbg {
				fmt.Printf("[%d] sz:%d  %s\n", code, len(html), url)
			}
		}

		for _, i := range indicators {
			if strings.Contains(html, i) {
				fmt.Println("Pattern Found !!!!")
				if isPost {
					fmt.Printf("%s\npost: %s\nindicator: %s\n", urlPost[0], urlPost[1], i)
				} else {
					fmt.Printf("%s\nindicator: %s\n", url, i)
				}
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
	post := flag.String("post", "", "post data")
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

	go injectPayloads(*param, *url, *post, *new, params, payloads, curls)

	fmt.Println("pentesting ...")

	go reduceFPs(cres)

	wg.Wait()
	close(cres)

	var i int
	for {
		fmt.Scanf("%d", &i)
	}
}
