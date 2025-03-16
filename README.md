# lambda-telegram-writer
A lambda function for receiving and passing a text message to telegram

## Compiling
On a Linux machine, run the following commands:
```sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go
```
Zip the compiled binary:
```
zip bootstrap.zip bootstrap
```
Upload the zip file to AWS Lambda.