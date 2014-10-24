/*
	pipper, fast web bruteforcer.
	@sha0coder
*/

package main

import "os"
import "fmt"
import "flag"
import "bufio"
import "net/url"
import "strings"
import "net/http"
import "io/ioutil"

var proxyurl string = ""

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func checkWebserver(surl string) {
	var client = &http.Client{}
	var server string

	if proxyurl != "" {
		proxy, err := url.Parse(proxyurl)
		check(err, "Bad proxy url")
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
	}

	req, err := http.NewRequest("GET", surl, nil)
	if err != nil {
		fmt.Println("Server is not responding :/")
		os.Exit(1)
	}

	fmt.Printf("checking %s ... \n", surl)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1; rv:5.0.1) Gecko/20100101 Firefox/5.0.1")
	resp, err := client.Do(req)
	check(err, "Can't connect")

	server = resp.Header.Get("Server")
	fmt.Printf("Server: %s\nDefault response: %d\n", server, resp.StatusCode)

	req, err = http.NewRequest("OPTIONS", surl, nil)
	resp, err = client.Do(req)
	if err == nil {
		fmt.Printf("Allowed Options: %s\n", resp.Header.Get("Allow"))
	}
}

func loadWordlist(wordlist string, c chan string) {
	file, err := os.Open(wordlist)
	check(err, "Can't load the wordlist")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c <- scanner.Text()
	}
	c <- "[EOF1337]"
	close(c)
}

func trypw(surl string, post string) (string, int) {
	var resp *http.Response
	var err error
	var html = []byte{}
	var code int = 0
	var client = &http.Client{} // 1 fixed client for goroutine?
	var req *http.Request
	var method string = "GET"

	if proxyurl != "" {
		proxy, err := url.Parse(proxyurl)
		check(err, "Bad proxy url")
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
	}

	if post != "" {
		method = "POST"
	}

	req, err = http.NewRequest(method, surl, strings.NewReader(post)) // hace la resolucion dns aqui?
	if err != nil {
		fmt.Println("Server is not responding :/")
		return "", 0
	}

	//check(err, "Can't connect")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1; rv:5.0.1) Gecko/20100101 Firefox/5.0.1")
	req.Header.Set("Accept-Encoding", "*")
	resp, err = client.Do(req)

	if err == nil && resp != nil {
		code = resp.StatusCode

		if resp.Body != nil {
			html, _ = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		}
	}

	return string(html), code
}

func main() {
	var url *string = flag.String("url", "", "the url")
	var post *string = flag.String("post", "", "post variables with ## where to bruteforce")
	var wordlist *string = flag.String("dict", "", "the wordlist")
	var goroutines *int = flag.Int("go", 1, "num of concurrent goroutines")
	var hl *int = flag.Int("hl", 0, "hide lines")
	var hw *int = flag.Int("hw", 0, "hide words")
	var hwl *int = flag.Int("hwl", 0, "hide words low")
	var hwh *int = flag.Int("hwh", 0, "hide words hight")
	var hb *int = flag.Int("hb", 0, "hide bytes")
	var hc *int = flag.Int("hc", 0, "hide code")
	var proxy *string = flag.String("proxy", "", "set proxy ip:port")

	var i int
	flag.Parse()

	if *url == "" || *wordlist == "" {
		check(nil, "bad usage --help")
	}

	if *proxy != "" {
		proxyurl = "http://" + *proxy
	}

	checkWebserver(*url)

	c := make(chan string, 6)
	go loadWordlist(*wordlist, c)

	for i = 0; i < *goroutines; i++ {
		go func(url string, post string, r int, c <-chan string) {
			var html string
			var lines int
			var words int
			var bytes int
			var code int
			var u string
			var p string

			for w := range c {
				if w == "[EOF1337]" {
					fmt.Println("end.\n")
					os.Exit(1)
				}
				u = strings.Replace(url, "##", w, -1)
				p = strings.Replace(post, "##", w, -1)
				html, code = trypw(u, p)
				lines = len(strings.Split(html, "\n"))
				words = len(strings.Split(html, " "))
				bytes = len(html)
				if *hl == lines || *hw == words || *hb == bytes || *hc == code || (*hwl <= words && words <= *hwh) {
					bytes = 0
					//fmt.Printf("\033[2K%d) (%d) [%d] [%d] [%d]\t\tword: %s\r", r, code, lines, words, bytes, w)
				} else {
					fmt.Printf("(%d) [%d] [%d] [%d]\t\t%s %s\n", code, lines, words, bytes, u, p)
					//fmt.Printf("\033[32m%d) (%d) [%d] [%d] [%d]\t\tword: %s\n\033[0m", r, code, lines, words, bytes, w)
				}

			}
		}(*url, *post, i, c)
	}

	fmt.Printf("Scanning, press enter to interrupt.\n")
	fmt.Scanf("%d", &i)
	fmt.Printf("interrupted.")

}
