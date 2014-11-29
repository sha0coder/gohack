/*
	pipper, params checker
	@sha0coder
*/

package main

import "os"
import "fmt"
import "flag"

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

func checkParam(url string, post string, param string, normal int) {
	var name string = ""
	var value string = ""
	var dynamic bool = true
	var p = strings.Split(param, "=")

	if len(p) >= 1 {
		name = p[0]
	}
	if len(p) >= 2 {
		value = p[1]
	}

	if *verbose {
		fmt.Printf("\nChecking param %s ...\n", name)
	}

	v := &VCheck{
		url:      url,
		post:     post,
		oldparam: param,
		normal:   normal,
	}

	// peticion normal
	//tryVal(url, post, param, param, normal)

	// si la respuesta no varia con un parámetro vaico ni con un parametro 69 (que sirve tanto numerico como string)
	if v.c(name+"=69").normal && v.c(name+"=").normal {
		if *verbose {
			fmt.Printf("%s not dynamic!\n\n", name)
		}
		dynamic = false
	}

	// es dinamico, el contenido no es fijo
	if *verbose {
		fmt.Printf("%s is dynamic!\n\n", name)
	}

	// el transversal se chekea sobre prams dinamicos,
	// porque el objetivo es coneguir un output fijo en
	// parametros dinamicos
	if dynamic && v.c(name+"=./"+value).normal {
		fmt.Printf("/!\\ %s %s potential traversal directory\n", url, name)
	}

	// se busca mantener el mismo resultado en un param dinamico
	// mediante un parametro equivalente xx  x''x o x'+'x
	if dynamic && v.c(name+"=''"+value).normal {
		fmt.Printf("/!\\ %s %s potential SQL injection\n", url, name)
	}

	// se busca provocar un cambio de contenido,
	// en parametros estaticos
	if !dynamic && !v.c(name+"='"+value).normal {
		fmt.Printf("/!\\ %s %s potential SQL injection\n", url, name)
	}

	// un patron univoco de SQLi es que las comillas pares den un resultado y las impares otro diferente
	// también se prueba la concatenación
	// ' == '''  && '' == '''' && '' == '+' && '' != '
	if dynamic {
		if v.c("'").words == v.c("'''").words && v.c("''").words == v.c("''''").words && v.c("''").words == v.c("'+'").words && v.c("'").words != v.c("''").words {
			fmt.Println("/!\\ %s %s high probability of SQL injection!\n", url, name)
		}
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

	if *verbose {
		fmt.Printf("Words: %d\n", test[0])
	}

	for _, param := range getParams(*url, *post) {
		checkParam(*url, *post, param, test[0])
	}
}

/*
	R+D Logica

	in (url)

	SQLI ->	(n("'") == n("'''") && n("''") == n("''''")  && n("'") != n("''"))
	SQLI ->	(n("'") == n("'''") && n("''") == n("''''")  && n("'") != n("''"))
	XSS  -> ()




	SQLI -> ' == '''  && '' == '''' && '' == '+' && '' != '
	XSS  ->

*/
