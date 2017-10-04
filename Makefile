
all:
	#go get code.google.com/p/go.net/html
	#go get gopkg.in/dutchcoders/goftp.v1
	go build pipper.go requests.go
	go build fauth.go requests.go
	go build smtpEnum.go requests.go
	go build params.go vulncheck.go requests.go
	#go build crawl.go requests.go
	go build massftpanon.go 

linux32:
	GOOS=linux GOARCH=386 go build pipper.go requests.go
	GOOS=linux GOARCH=386 go build fauth.go requests.go
	GOOS=linux GOARCH=386 go build smtpEnum.go requests.go
	GOOS=linux GOARCH=386 go build params.go vulncheck.go requests.go
	#GOOS=linux GOARCH=386 go build crawl.go requests.go
	GOOS=linux GOARCH=386 go build massftpanon.go 

linux64:
	GOOS=linux GOARCH=amd64 go build pipper.go requests.go
	GOOS=linux GOARCH=amd64 go build fauth.go requests.go
	GOOS=linux GOARCH=amd64 go build smtpEnum.go requests.go
	GOOS=linux GOARCH=amd64 go build params.go vulncheck.go requests.go
	#GOOS=linux GOARCH=amd64 go build crawl.go requests.go
	GOOS=linux GOARCH=amd64 go build massftpanon.go 

win32:
	GOOS=windows GOARCH=386 go build pipper.go requests.go
	GOOS=windows GOARCH=386 go build fauth.go requests.go
	GOOS=windows GOARCH=386 go build smtpEnum.go requests.go
	GOOS=windows GOARCH=386 go build crawl.go requests.go
	#GOOS=windows GOARCH=386 go build params.go vulncheck.go requests.go
	GOOS=windows GOARCH=386 go build massftpanon.go 

win64:
	GOOS=windows GOARCH=amd64 go build pipper.go requests.go
	GOOS=windows GOARCH=amd64 go build fauth.go requests.go
	GOOS=windows GOARCH=amd64 go build smtpEnum.go requests.go
	GOOS=windows GOARCH=amd64 go build params.go vulncheck.go requests.go
	#GOOS=windows GOARCH=amd64 go build crawl.go requests.go
	GOOS=windows GOARCH=amd64 go build massftpanon.go 
 
clean:
	rm -f pipper fauth smtpEnum params massftpanon  *.exe
	
uninstall:
	rm -f /usr/bin/pipper /usr/bin/fauth /usr/bin/smtpEnum /usr/bin/params /usr/bin/massftpanon

install:
	cp pipper fauth smtpEnum params massftpanon /usr/bin/


