package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

type TelegramMessageRequest struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

var telegramWriter *TelegramWriter

var environment string

type TelegramWriter struct {
	token  string
	chatId string
}

func (tw *TelegramWriter) sendMessage(msg string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tw.token)

	body := &TelegramMessageRequest{tw.chatId, msg}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewReader(jsonBody))

	return err
}

func init() {
	godotenv.Load(".env")

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatId := os.Getenv("TELEGRAM_CHAT_ID")
	if botToken == "" || chatId == "" {
		log.Fatal("Missing environment variables!")
	}

	environment = os.Getenv("ENVIRONMENT")

	telegramWriter = &TelegramWriter{
		token:  botToken,
		chatId: chatId,
	}
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := telegramWriter.sendMessage(req.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type",
			},
			Body: `{ "error": "Failed to send message" }`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type",
		},
		Body: `{ "message": "Message sent successfully!" }`,
	}, nil
}

func devToLambdaHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	lambdaReq := events.APIGatewayProxyRequest{
		Body: string(body),
	}

	res, err := handler(lambdaReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong!"))
		return
	}

	w.WriteHeader(res.StatusCode)
	w.Write([]byte(res.Body))
}

func main() {

	if environment == "dev" {
		http.HandleFunc("/telegramWriter", devToLambdaHandler)
		http.ListenAndServe(":8080", nil)
	} else {
		lambda.Start(handler)
	}
}
