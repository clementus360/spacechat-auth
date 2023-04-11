package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type MessageService interface {
	SendMessage(number string, message string) error
}

type TwilioService struct {
	SID       string
	Token     string
	Number    string
	TwilioCli *twilio.RestClient
}

func NewTwilioService() *TwilioService {
	sid := os.Getenv("TWILIO_SID")
	token := os.Getenv("TWILIO_TOKEN")
	number := os.Getenv("TWILIO_NUMBER")

	TwilioCli := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: sid,
		Password: token,
	})

	return &TwilioService{
		SID:       sid,
		Token:     token,
		Number:    number,
		TwilioCli: TwilioCli,
	}
}

func (ts TwilioService) SendMessage(number string, message string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(number)
	params.SetFrom(ts.Number)
	params.SetBody(message)

	resp, err := ts.TwilioCli.Api.CreateMessage(params)

	if err != nil {
		fmt.Println("Failed to send message")
		return err
	}

	respBytes, _ := json.Marshal(*resp)
	fmt.Println("Message sent: ", string(respBytes))

	return nil

}
