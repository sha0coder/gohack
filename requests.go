/*
	http/https connections
	@sha0coder
*/

package main

import "fmt"
import "os"
import "net/url"
import "strings"
import "net/http"
import "io/ioutil"
import "crypto/tls"

type Requests struct {
	TlsCfg    *tls.Config
	Transport *http.Transport
	UserAgent string
	ProxyUrl  string
	User      string
	Passw     string
}

func NewRequests() *Requests {
	h := new(Requests)
	h.UserAgent = "Mozilla/5.0 (Windows NT 5.1; rv:5.0.1) Gecko/20100101 Firefox/5.0.1"
	h.ProxyUrl = ""
	h.User = ""
	h.Passw = ""
	h.TlsCfg = &tls.Config{InsecureSkipVerify: true}
	h.Transport = &http.Transport{
		TLSClientConfig: h.TlsCfg,
	}

	return h
}

func (h *Requests) SetProxy(proxyurl string) {
	proxy, err := url.Parse(proxyurl)
	if err != nil {
		fmt.Println("bad proxy url")
		return
	}
	h.Transport.Proxy = http.ProxyURL(proxy)
	h.ProxyUrl = proxyurl
}

func (h *Requests) Get(url string) (string, int, *http.Response) {
	return h.Launch("GET", url, "")
}

func (h *Requests) Post(url string, post string) (string, int, *http.Response) {
	return h.Launch("POST", url, post)
}

func (h *Requests) GetOrPost(url string, post string) (string, int, *http.Response) {
	if post == "" {
		return h.Launch("GET", url, "")
	} else {
		return h.Launch("POST", url, post)
	}
}

func (h *Requests) Options(url string) (string, int, *http.Response) {
	return h.Launch("OPTIONS", url, "")
}

func (h *Requests) QuitOnFail(code int, msg string) {
	if code == 0 {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func (r *Requests) SetBasicAuth(user string, passw string) {
	// This call is not thread-safe
	r.User = user
	r.Passw = passw
}

func (h *Requests) LaunchNoRead(method string, url string, post string) *http.Response {
	req, err := http.NewRequest(method, url, strings.NewReader(post)) // hace la resolucion dns aqui?

	if err != nil {
		fmt.Println("Server is not responding :/")
		return nil
	}

	if h.User != "" {
		req.SetBasicAuth(h.User, h.Passw)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("User-Agent", h.UserAgent)
	req.Header.Set("Accept-Encoding", "*")

	client := &http.Client{Transport: h.Transport}

	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return nil
	}

	if resp.Body == nil {
		return nil
	}

	return resp
}

func (h *Requests) Launch(method string, url string, post string) (string, int, *http.Response) {
	req, err := http.NewRequest(method, url, strings.NewReader(post)) // hace la resolucion dns aqui?

	if err != nil {
		fmt.Println("Server is not responding :/")
		return "", 0, nil
	}

	if h.User != "" {
		req.SetBasicAuth(h.User, h.Passw)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("User-Agent", h.UserAgent)
	req.Header.Set("Accept-Encoding", "*")

	client := &http.Client{Transport: h.Transport}

	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return "", 0, nil
	}

	code := resp.StatusCode
	if resp.Body == nil {
		return "", code, resp
	}

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", code, resp
	}
	resp.Body.Close()

	return string(html), code, resp
}

/*
func test() {
	h := NewHTTP()
	html, code := h.Get("https://google.com")
	fmt.Println(code)
	fmt.Println(html)
}*/
