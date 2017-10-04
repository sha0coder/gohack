package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/dutchcoders/goftp.v1"
)

func end(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func ftpHost(host string, verbose bool) {
	var hostport bytes.Buffer
	var ftp *goftp.FTP
	var files []string
	var err error

	hostport.WriteString(host)
	hostport.WriteString(":21")

	if ftp, err = goftp.Connect(hostport.String()); err != nil {
		if verbose {
			fmt.Printf("%s down \n", host)
		}
		return
	}
	defer ftp.Close()

	if err = ftp.Login("ftp", "ftp"); err != nil {
		fmt.Printf("%s Denied! \n", host)
		return
	}

	fmt.Printf("%s Anonimous access!!  \n", host)

	if err = ftp.Cwd("/"); err != nil {
		return
	}

	if files, err = ftp.List(""); err != nil {
		return
	}

	for _, file := range files {
		fmt.Printf(" - %s\n", file)
	}

	return
}

func getOctet(ip []string, octet int) int {
	o, err := strconv.Atoi(ip[octet-1])
	if err != nil {
		end("wrong ip")
	}
	return o
}

func nextIP(currIP string) string {
	var octets []string
	octets = strings.Split(currIP, ".")
	if len(octets) != 4 {
		end("wrong IP")
	}
	o4 := getOctet(octets, 4)
	if o4 < 255 {
		o4++
		return fmt.Sprintf("%s.%s.%s.%d", octets[0], octets[1], octets[2], o4)
	}

	o3 := getOctet(octets, 3)
	if o3 < 255 {
		o3++
		o4 = 0
		return fmt.Sprintf("%s.%s.%d.%d", octets[0], octets[1], o3, o4)
	}

	o2 := getOctet(octets, 2)
	if o2 < 255 {
		o2++
		o3 = 0
		o4 = 0
		return fmt.Sprintf("%s.%d.%d.%d", octets[0], o2, o3, o4)
	}

	o1 := getOctet(octets, 1)
	if o1 < 255 {
		o1++
		o2 = 0
		o3 = 0
		o4 = 0
		return fmt.Sprintf("%d.%d.%d.%d", o1, o2, o3, o4)
	}

	return ""
}

func main() {
	var host string
	var from *string = flag.String("from", "", "starting ip")
	var to *string = flag.String("to", "", "ending ip")
	var verbose *bool = flag.Bool("v", false, "verbose mode")
	//var goroutines *int = flag.Int("go", 1, "num of concurrent goroutines")
	flag.Parse()

	if *from == "" || *to == "" {
		end("try --help")
	}

	host = *from
	for {
		go ftpHost(host, *verbose)

		if host == *to {
			break
		}

		host = nextIP(host)

		if host == "" {
			break
		}
	}

	var i int
	fmt.Printf("Scanning, press enter to interrupt.\n")
	fmt.Scanf("%d", &i)
	fmt.Printf("interrupted.")

}
