/*
	Enumerate users via SMTP service
	@sha0coder

*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
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

func main() {
	var smtp *string = flag.String("smtp", "", "ip:port of the smtp server")
	var wlfile *string = flag.String("wl", "", "users wordlist")
	var host *string = flag.String("host", "", "users wordlist")
	var verbose *bool = flag.Bool("v", false, "verbose")

	flag.Parse()

	c := make(chan string, 6)
	go loadWordlist(*wlfile, c)

	conn, err := net.Dial("tcp", *smtp)
	check(err, "cant connect")

	read := bufio.NewReader(conn)

	str, _ := read.ReadString('\n')
	fmt.Println(str)
	if *verbose {
		fmt.Printf("helo %s\n", *host)
	}
	fmt.Fprintf(conn, "helo %s\n", *host)
	str, _ = read.ReadString('\n')
	if *verbose {
		fmt.Println(str)
		fmt.Println("mail from: <test@clouuds.com>")
	}

	//fmt.Fprintf(conn, "mail from: <test@clouuds.com>\n")
	//read.ReadString('\n')

	for w := range c {
		if w == "[EOF1337]" {
			fmt.Println("end.\n")
			os.Exit(1)
		}

		fmt.Fprintf(conn, "VRFY %s\n", w)
		str, _ := read.ReadString('\n')
		if *verbose {
			fmt.Println(str)
		}
		if strings.Contains(str, "Sender ok") || strings.Contains(str, "Recipient ok") || strings.Contains(str, "Ok") {
			fmt.Printf("=> %s\n", w)
		} else if !*verbose {
			fmt.Printf("\r%s          ", w)
		}

		time.Sleep(10 * time.Second)
	}

	fmt.Println("end.")

}
