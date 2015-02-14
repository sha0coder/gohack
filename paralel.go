/*
	Concurrent wordlist consumers
*/

package main

import "bufio"
import "fmt"
import "os"

var ParalelMagicNumber = "[EOF1337]"

type TypeDataCallback func(i int, data string)
type TypeFinishCallback func()

type Paralel struct {
	Stop       bool
	Ch         chan string
	Routines   int
	CbOnData   TypeDataCallback
	CbOnFinish TypeFinishCallback
}

func (p *Paralel) Init(n int) {
	p.Stop = false
	p.Ch = make(chan string, 6)
	p.Routines = n
	p.CbOnData = nil
	p.CbOnFinish = nil
}

func (p *Paralel) OnData(cb TypeDataCallback) {
	p.CbOnData = cb
}

func (p *Paralel) OnFinish(cb TypeFinishCallback) {
	p.CbOnFinish = cb
}

func (p Paralel) Check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func (p *Paralel) Load(wordlist string) {
	go func() {
		file, err := os.Open(wordlist)
		p.Check(err, "Can't load the wordlist")
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			p.Ch <- scanner.Text()
		}
		p.Ch <- ParalelMagicNumber
		close(p.Ch)
	}()
}

func (p *Paralel) Start() {
	if p.CbOnData == nil {
		fmt.Println("set the data event")
		os.Exit(1)
	}
	for i := 0; i < p.Routines; i++ {
		go func(r int) {
			for w := range p.Ch {
				if p.Stop {
					return
				}

				if w == ParalelMagicNumber {
					p.Stop = true
					if p.CbOnFinish != nil {
						p.CbOnFinish()
					}
					return
				}

				p.CbOnData(i, w) //deberia pasar el campo i
			}
		}(i)
	}
}

func (p Paralel) Wait() {
	// should loop to auto-stop?
	var i int
	fmt.Printf("Scanning, press enter to interrupt.\n")
	fmt.Scanf("%d", &i)
	fmt.Printf("interrupted.")
}

func main() {
	p := new(Paralel)
	p.Init(10)                          // ten consumers
	p.Load("/etc/passwd")               // the wordlist to consume
	p.OnData(func(i int, data string) { // the consumer code
		fmt.Println(data)
	})
	p.OnFinish(func() { // finish event
		fmt.Println("finish.")
	})
	p.Start() // Start to processinf
	p.Wait()
}
