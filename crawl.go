package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	colly "github.com/gocolly/colly/v2"
)

func onHtml(e *colly.HTMLElement) {
	link := e.Attr("href")
	e.Request.Visit(link)
	src := e.Attr("src")
	e.Request.Visit(src)
	act := e.Attr("action")
	e.Request.Visit(act)
}

func onRequest(r *colly.Request) {
	fmt.Println(r.URL)
}

func main() {
	url := flag.String("url", "", "url to crawl")
	flag.Parse()

	if *url == "" {
		fmt.Println("select a domain to crawl -url or -h")
		os.Exit(1)
	}

	purl := strings.Split(*url, "/")
	if len(purl) < 3 {
		fmt.Println("bad url")
		os.Exit(1)
	}

	dom := purl[2]

	c := colly.NewCollector(
		colly.AllowedDomains(dom),
		colly.CacheDir("./crawler_cache"),
	)

	c.OnHTML("a[href]", onHtml)
	c.OnRequest(onRequest)
	c.Visit(*url)
}
