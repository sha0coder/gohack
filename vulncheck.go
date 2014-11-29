package main

import "net/http"
import "io/ioutil"
import "strings"
import "fmt"

type VCheck struct {
	url      string
	post     string
	oldparam string
	normal   int
	test     int
}

type Result struct {
	words  int
	normal bool
}

func try(url string, post string) int {
	var client = &http.Client{}
	var method string = "GET"
	var err error
	var req *http.Request
	var resp *http.Response
	var html string = ""

	if post != "" {
		method = "POST"
		req, err = http.NewRequest(method, url, strings.NewReader(post))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	check(err, "err preparing request")

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1; rv:5.0.1) Gecko/20100101 Firefox/5.0.1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Accept-Encoding", "*")
	resp, err = client.Do(req)
	check(err, "Can't connect "+url)

	if err == nil && resp != nil {
		if resp.StatusCode == 200 {
			if resp.Body != nil {
				data, err := ioutil.ReadAll(resp.Body)
				check(err, "err reading data")
				html = string(data)
				resp.Body.Close()
			}
		}
	}

	return len(strings.Split(html, " "))
}

func (v *VCheck) c(newparam string) *Result {
	var url string = ""
	var post string = ""
	var w int
	url = strings.Replace(v.url, v.oldparam, newparam, -1)
	post = strings.Replace(v.post, v.oldparam, newparam, -1)

	w = try(url, post)

	if w == v.normal {
		if *verbose {
			fmt.Printf("[+] %s\n", newparam)
		}
		return &Result{normal: true, words: w}
	}

	if *verbose {
		fmt.Printf("[-] %s\n", newparam)
	}
	return &Result{normal: true, words: w}

}
