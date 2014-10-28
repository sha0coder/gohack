/*
	pipper, params checker
	@sha0coder
*/

package main

import "os"
import "fmt"
import "flag"
import "net/http"
import "io/ioutil"
import "strings"

var verbose *bool

func die(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func check(err error, msg string) {
	if err != nil {
		die(msg)
	}
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

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1; rv:5.0.1) Gecko/20100101 Firefox/5.0.1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Set("Accept-Encoding", "*")
	resp, err = client.Do(req)
	check(err, "Can't connect")

	if err == nil && resp != nil {
		if resp.StatusCode == 200 {
			if resp.Body != nil {
				data, _ := ioutil.ReadAll(resp.Body)
				html = string(data)
				resp.Body.Close()
			}
		}
	}

	return len(strings.Split(html, " "))
}

func getParams(url string, post string) []string {
	var params []string

	spl := strings.Split(url, "?")
	if len(spl) > 1 {
		for _, p := range strings.Split(spl[1], "&") {
			params = append(params, p)
		}
	}

	if post != "" {
		for _, p := range strings.Split(post, "&") {
			params = append(params, p)
		}
	}

	return params
}

func tryVal(url string, post string, oldparam string, newparam string, normal int) bool {
	var w int
	url = strings.Replace(url, oldparam, newparam, -1)
	post = strings.Replace(post, oldparam, newparam, -1)

	w = try(url, post)

	if w == normal {
		if *verbose {
			fmt.Printf("[+] %s\n", newparam)
		}
		return true
	}

	if *verbose {
		fmt.Printf("[-] %s\n", newparam)
	}
	return false

}

func checkParam(url string, post string, param string, normal int) {
	var name string
	var value string
	var dynamic bool = true
	var p = strings.Split(param, "=")
	name = p[0]
	value = p[1]

	fmt.Printf("\nChecking param %s ...\n", name)

	// peticion normal
	tryVal(url, post, param, param, normal)

	// si la respuesta no varia con un parámetro vaico ni con un parametro 69 (que sirve tanto numerico como string)
	if tryVal(url, post, param, name+"=69", normal) && tryVal(url, post, param, name+"=", normal) {
		fmt.Printf("%s not dynamic!\n\n", name)
		dynamic = false
	}

	// es dinamico, el contenido no es fijo
	fmt.Printf("%s is dynamic!\n\n", name)

	// el transversal se chekea sobre prams dinamicos,
	// porque el objetivo es coneguir un output fijo en
	// parametros dinamicos
	if dynamic && tryVal(url, post, param, name+"=./"+value, normal) {
		fmt.Printf("%s potential traversal directory\n", name)
	}

	// se busca mantener el mismo resultado en un param dinamico
	// mediante un parametro equivalente xx  x''x o x'+'x
	if dynamic && tryVal(url, post, param, name+"=''"+value, normal) {
		fmt.Printf("%s potential SQL injection\n", name)
	}

	// se busca provocar un cambio de contenido,
	// en parametros estaticos
	if !dynamic && !tryVal(url, post, param, name+"='"+value, normal) {
		fmt.Printf("%s potential SQL injection\n", name)
	}

}

func main() {
	var test [3]int
	var url *string = flag.String("url", "", "the url")
	var post *string = flag.String("post", "", "post data")
	verbose = flag.Bool("v", false, "verbose")
	flag.Parse()

	if *url == "" {
		die("try --help")
	}

	if !strings.Contains(*post, "=") {
		if !strings.Contains(*url, "?") || !strings.Contains(*url, "=") {
			die("No params found")
		}
	}

	test[0] = try(*url, *post)
	test[1] = try(*url, *post)
	test[2] = try(*url, *post)

	if test[0] != test[1] || test[1] != test[2] {
		die("Non stable")
	}

	fmt.Printf("Words: %d\n", test[0])

	for _, param := range getParams(*url, *post) {
		checkParam(*url, *post, param, test[0])
	}

}
