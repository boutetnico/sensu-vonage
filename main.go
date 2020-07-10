package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nexmo-community/nexmo-go"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
)

// HandlerConfig contains the Slack handler configuration
type HandlerConfig struct {
	sensu.PluginConfig
	vonageAPIKey     string
	vonageAPISecret  string
	vonageFrom       string
	vonageRecipients string
}

var (
	config = HandlerConfig{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-vonage-handler",
			Short:    "The Sensu Go Vonage (Nexmo) handler for sending sms alerts",
			Keyspace: "sensu.io/plugins/vonage/config",
		},
	}

	slackConfigOptions = []*sensu.PluginConfigOption{
		{
			Path:      "api-key",
			Env:       "VONAGE_API_KEY",
			Argument:  "api-key",
			Shorthand: "k",
			Usage:     "The Vonage API key",
			Value:     &config.vonageAPIKey,
		},
		{
			Path:      "api-secret",
			Env:       "VONAGE_API_SECRET",
			Argument:  "api-secret",
			Shorthand: "s",
			Usage:     "The Vonage API secret",
			Value:     &config.vonageAPISecret,
		},
		{
			Path:      "from",
			Env:       "VONAGE_FROM",
			Argument:  "from",
			Shorthand: "f",
			Default:   "Sensu Go",
			Usage:     "The number/name that will send the sms",
			Value:     &config.vonageFrom,
		},
		{
			Path:      "recipients",
			Env:       "VONAGE_RECIPIENTS",
			Argument:  "recipients",
			Shorthand: "r",
			Usage:     "Comma-separated list of numbesr of recipients",
			Value:     &config.vonageRecipients,
		},
	}
)

func main() {
	goHandler := sensu.NewGoHandler(&config.PluginConfig, slackConfigOptions, checkArgs, sendMessage)
	goHandler.Execute()
}

func checkArgs(_ *corev2.Event) error {
	// Support deprecated environment variables
	if apiKey := os.Getenv("VONAGE_API_KEY"); apiKey != "" {
		config.vonageAPIKey = apiKey
	}
	if apiSecret := os.Getenv("VONAGE_API_SECRET"); apiSecret != "" {
		config.vonageAPISecret = apiSecret
	}
	if from := os.Getenv("VONAGE_FROM"); from != "" {
		config.vonageFrom = from
	}
	if recipients := os.Getenv("VONAGE_RECIPIENTS"); recipients != "" {
		config.vonageRecipients = recipients
	}

	return nil
}

func formattedEventAction(event *corev2.Event) string {
	switch event.Check.Status {
	case 0:
		return "RESOLVED"
	default:
		return "ALERT"
	}
}

func chomp(s string) string {
	return strings.Trim(strings.Trim(strings.Trim(s, "\n"), "\r"), "\r\n")
}

func eventKey(event *corev2.Event) string {
	return fmt.Sprintf("%s/%s", event.Entity.Name, event.Check.Name)
}

func eventSummary(event *corev2.Event, maxLength int) string {
	output := chomp(event.Check.Output)
	if len(event.Check.Output) > maxLength {
		output = output[0:maxLength] + "..."
	}
	return fmt.Sprintf("%s:%s", eventKey(event), output)
}

func formattedMessage(event *corev2.Event) string {
	return fmt.Sprintf("%s - %s", formattedEventAction(event), eventSummary(event, 100))
}

func messageStatus(event *corev2.Event) string {
	switch event.Check.Status {
	case 0:
		return "Resolved"
	case 2:
		return "Critical"
	default:
		return "Warning"
	}
}

func sendMessage(event *corev2.Event) error {
	auth := nexmo.NewAuthSet()
	auth.SetAPISecret(config.vonageAPIKey, config.vonageAPISecret)

	client := nexmo.NewClient(http.DefaultClient, auth)
	smsReq := nexmo.SendSMSRequest{
		From: config.vonageFrom,
		To:   config.vonageRecipients,
		Text: formattedMessage(event),
	}

	callR, _, err := client.SMS.SendSMS(smsReq)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", callR.Messages[0].Status)

	return nil
}
