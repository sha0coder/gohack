package main

import "fmt"
import "flag"
import "strings"
import "os"
import "code.google.com/p/go.net/html"

type Crawl struct {
	R       *Requests
	Crawled []string
	Hosts   []string
	Pending chan string
	DoStop  bool
}

func NewCrawl() *Crawl {
	c := new(Crawl)
	c.R = NewRequests()
	c.Pending = make(chan string)
	c.DoStop = false
	return c
}

func (c *Crawl) AddHost(host string) {
	c.Hosts = append(c.Hosts, host)
}

func (c *Crawl) Fix(url string) string {
	return url
}

func (c *Crawl) IsCrawled(url string) bool {
	for _, u := range c.Crawled {
		if url == u {
			return true
		}
	}
	return false
}

func (c *Crawl) IsAllowed(host string) bool {
	for _, h := range c.Hosts {
		if host == h {
			return true
		}
	}
	return false
}

func (c *Crawl) Queue(url string) {
	host := strings.Split(url, "/")
	if len(host) < 3 {
		//fmt.Println("bad url: " + url)
		return
	}

	if !c.IsAllowed(host[2]) {
		//fmt.Printf("out scope: " + url)
		return
	}

	if !c.IsCrawled(url) {
		//fmt.Println("queued: " + url)
		c.Pending <- url // crash 
	}
}

func (c *Crawl) Stop() {
	c.DoStop = false
}

func (c *Crawl) Process(url string) {
	if c.IsCrawled(url) {
		return
	}

	fmt.Println("scanning " + url)

	resp := c.R.LaunchNoRead("GET", url, "")
	defer resp.Body.Close()

	page := html.NewTokenizer(resp.Body)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			c.Crawled = append(c.Crawled, url)
			return
		}
		token := page.Token()

		//if tokenType == html.StartTagToken { //&& token.DataAtom.String() == "a" {
		for _, attr := range token.Attr {
			if attr.Key == "href" || attr.Key == "action" || attr.Key == "src" {
				go c.Queue(attr.Val)
			}
		}
		//}
	}
}

func (c *Crawl) Worker() {
	for {
		for url := range c.Pending {
			if c.DoStop {
				return
			}
			c.Process(url)
		}
	}
}

func main() {
	var url *string = flag.String("url", "", "the url to start crawling")
	var th *int = flag.Int("go", 5, "number of concurrent goroutines")
	flag.Parse()

	upart := strings.Split(*url, "/")

	if len(upart) < 3 {
		fmt.Println("bad url")
		os.Exit(1)
	}

	c := NewCrawl()
	c.AddHost(upart[2])
	for i := 0; i < *th; i++ {
		go c.Worker()
	}
	c.Queue(*url)

	fmt.Println("Press enter to stop crawling")
	var i int
	fmt.Scan(&i)
	c.Stop()

	for _, u := range c.Crawled {
		fmt.Println(u)
	}
}

