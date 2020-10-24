dailyrun:
	go run main.go daily

weeklyrun:
	go run main.go weekly

monthlyrun:
	go run main.go monthly

build:	
	go build main.go -o update

publishbuild:
	env GOOS=linux GOARCH=amd64 go build main.go
