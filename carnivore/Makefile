all:
	CGO_ENABLED=0 go build carnivore.go requests.go
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build carnivore.go requests.go
	