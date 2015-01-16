package main

import "strings"
import "fmt"

type VCheck struct {
	url      string
	post     string
	oldparam string
	normal   int
	test     int
	R        *Requests
}

type Result struct {
	words  int
	normal bool
}

func NewVCheck() *VCheck {
	vc := new(VCheck)
	vc.R = NewRequests()
	return vc
}

func (v *VCheck) c(newparam string) *Result {
	var url string = ""
	var post string = ""
	var w int
	url = strings.Replace(v.url, v.oldparam, newparam, -1)
	post = strings.Replace(v.post, v.oldparam, newparam, -1)

	html, code, _ := v.R.GetOrPost(url, post)
	v.R.QuitOnFail(code, "Can't connect")
	w = len(strings.Split(html, " "))

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
