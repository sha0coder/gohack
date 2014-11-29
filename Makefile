

all:
	go build pipper.go
	go build fauth.go
	go build smtpEnum.go
	go build params.go vulncheck.go

linux32:
	GOOS=linux GOARCH=386 go build pipper.go
	GOOS=linux GOARCH=386 go build fauth.go
	GOOS=linux GOARCH=386 go build smtpEnum.go
	GOOS=linux GOARCH=386 go build params.go vulncheck.go

linux64:
	GOOS=linux GOARCH=amd64 go build pipper.go
	GOOS=linux GOARCH=amd64 go build fauth.go
	GOOS=linux GOARCH=amd64 go build smtpEnum.go
	GOOS=linux GOARCH=amd64 go build params.go vulncheck.go

win32:
	GOOS=windows GOARCH=386 go build pipper.go
	GOOS=windows GOARCH=386 go build fauth.go
	GOOS=windows GOARCH=386 go build smtpEnum.go
	GOOS=windows GOARCH=386 go build params.go vulncheck.go

win64:
	GOOS=windows GOARCH=amd64 go build pipper.go
	GOOS=windows GOARCH=amd64 go build fauth.go
	GOOS=windows GOARCH=amd64 go build smtpEnum.go
	GOOS=windows GOARCH=amd64 go build params.go vulncheck.go
 
clean:
	rm -f pipper fauth smtpEnum params *.exe
	
uninstall:
	rm -f /usr/bin/pipper /usr/bin/fauth /usr/bin/smtpEnum /usr/bin/params

install:
	cp pipper fauth smtpEnum params /usr/bin/


