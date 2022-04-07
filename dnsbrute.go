/*
	subdomain 0-4 chars bruteforce for discovering subdomains
*/

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func queryNS(subdomain, nameserver string, port int) ([]string, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, nameserver+":"+strconv.Itoa(port))
		},
	}
	return r.LookupHost(context.Background(), subdomain)
}

func process(c <-chan string, domain, nameserver string, port int) {
	for dns := range c {
		var addr []string
		var err error

		subd := dns + "." + domain

		if nameserver == "" {
			addr, err = net.LookupHost(subd)
		} else {
			addr, err = queryNS(subd, nameserver, port)
		}

		if err == nil {
			fmt.Printf("%s %s\n", subd, addr)
		} else {
			fmt.Printf("%s        \r", subd)
			f := bufio.NewWriter(os.Stdout)
			f.Write([]byte(subd + "                          "))
			f.Write([]byte("\r"))
			f.Flush()

		}

	}
}

func main() {
	domain := flag.String("dom", "", "domain name ie: -dom test.com")
	gorountines := flag.Int("go", 2, "number of concurrent goroutines")
	nameserver := flag.String("ns", "", "specify the nameserver to query")
	port := flag.Int("p", 53, "name-server port")
	flag.Parse()

	if *domain == "" {
		fmt.Println("use -dom domain.com or -h")
		os.Exit(1)
	}

	ch := make(chan string, 6)

	for i := 0; i < *gorountines; i++ {
		go process(ch, *domain, *nameserver, *port)
	}

	for a := int('a'); a <= int('z'); a++ {
		ch <- string(a)

		for b := int('a'); b <= int('z'); b++ {
			ch <- string(a) + string(b)

			for c := int('a'); c <= int('z'); c++ {
				ch <- string(a) + string(b) + string(c)

				for d := int('a'); d <= int('z'); d++ {
					ch <- string(a) + string(b) + string(c) + string(d)
				}
			}
		}
	}
	close(ch)

	var i int
	for {
		fmt.Scanf("%d", &i)
	}

}
