package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
				"Content-Type": "application/json",
			},
			Body: `{ "error": "Failed to send message" }`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{ "message": "Message sent successfully!" }`,
	}, nil
}

func main() {
	lambda.Start(handler)

	// For local testing
	// req := events.APIGatewayProxyRequest{
	// 	Body: "Hello\nWorld",
	// }

	// resp, err := handler(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// println(resp.Body)
}
