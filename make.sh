### For Cellularhacker
export GOPATH="$HOME/Projects/sms-sorter/"
GOOS=linux GOARCH=amd64 go build -v sms-sorter

go build src/sms-sorter/main.go