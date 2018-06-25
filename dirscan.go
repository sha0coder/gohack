/*
	dirscan, fast web bruteforcer.
	@sha0coder
*/

package main

import "os"
import "fmt"
import "flag"
import "bufio"
import "strings"
import "strconv"

var R *Requests

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func checkWebserver(surl string) {
	_, code, resp := R.Get(surl)
	R.QuitOnFail(code, "Can't connect")

	fmt.Printf("Server: %s\nDefault response: %d\n", resp.Header.Get("Server"), resp.StatusCode)

	_, _, resp = R.Options(surl)
	fmt.Printf("Allowed Options: %s\n", resp.Header.Get("Allow"))
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
	fmt.Println("Wordlist completed.")
}

func main() {
	var url *string = flag.String("url", "", "the url")
	var post *string = flag.String("post", "", "post variables with ## where to bruteforce")
	var wordlist *string = flag.String("dict", "", "the wordlist")
	var num *int = flag.Int("num", 0, "numeric sequence")
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

	if *url == "" || (*wordlist == "" && *num == 0) {
		fmt.Println("num:%d\n", *num)
		check(nil, "bad usage --help")
	}

	R = NewRequests()

	if *proxy != "" {
		R.SetProxy("http://" + *proxy)
	}

	checkWebserver(*url)

	c := make(chan string, 6)

	if *wordlist != "" {
		go loadWordlist(*wordlist, c)
	}

	if *num > 0 {
		go func() {
			for n := 0; n < *num; n++ {
				c <- strconv.Itoa(n)
			}
			c <- "[EOF1337]"
			close(c)
		}()
	}

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
				html, code, _ = R.GetOrPost(u, p)
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
