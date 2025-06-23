package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sb-nour/providers-endpoints/service"
)

type SlackMessage struct {
	Text        string            `json:"text,omitempty"`
	Channel     string            `json:"channel,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

type SlackAttachment struct {
	Color     string       `json:"color,omitempty"`
	Title     string       `json:"title,omitempty"`
	Text      string       `json:"text,omitempty"`
	Fields    []SlackField `json:"fields,omitempty"`
	Timestamp int64        `json:"ts,omitempty"`
}

type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func SendSlackNotification(webhookURL, message string) error {
	return SendSlackNotificationWithChannel(webhookURL, message, "")
}

func SendSlackNotificationWithChannel(webhookURL, message, channel string) error {
	if webhookURL == "" {
		log.Printf("Slack webhook URL not configured, skipping notification: %s", message)
		return nil
	}

	slackMessage := SlackMessage{
		Text:    message,
		Channel: channel,
	}

	jsonData, err := json.Marshal(slackMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack notification failed with status: %d", resp.StatusCode)
	}

	log.Printf("Slack notification sent successfully")
	return nil
}

func SendRegionsFetchErrorNotification(provider string, err error) {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		return
	}

	// Check for error-specific channel override
	var errorChannel string
	if envChannel := os.Getenv("SLACK_ERROR_CHANNEL"); envChannel != "" {
		errorChannel = fmt.Sprintf("#%s", envChannel)
	}

	message := SlackMessage{
		Channel: errorChannel,
		Attachments: []SlackAttachment{
			{
				Color: "danger",
				Title: "🚨 Provider Regions Fetch Failed",
				Text:  fmt.Sprintf("Failed to fetch regions for provider: *%s*", provider),
				Fields: []SlackField{
					{
						Title: "Provider",
						Value: provider,
						Short: true,
					},
					{
						Title: "Error",
						Value: err.Error(),
						Short: false,
					},
					{
						Title: "Timestamp",
						Value: time.Now().Format("2006-01-02 15:04:05"),
						Short: true,
					},
				},
				Timestamp: time.Now().Unix(),
			},
		},
	}

	jsonData, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		log.Printf("Failed to marshal slack message for provider %s: %v", provider, marshalErr)
		return
	}

	resp, postErr := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if postErr != nil {
		log.Printf("Failed to send slack notification for provider %s: %v", provider, postErr)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Slack notification failed for provider %s with status: %d", provider, resp.StatusCode)
		return
	}

	log.Printf("Sent regions fetch error notification for provider: %s", provider)
}

func SendRegionsChangedNotification(provider string, oldRegions, newRegions service.Regions) {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		return
	}

	// Check for changes-specific channel override
	var changesChannel string
	if envChannel := os.Getenv("SLACK_CHANGES_CHANNEL"); envChannel != "" {
		changesChannel = fmt.Sprintf("#%s", envChannel)
	}

	// Calculate changes
	storageChanges := calculateRegionChanges(oldRegions.Storage, newRegions.Storage)
	computeChanges := calculateRegionChanges(oldRegions.Compute, newRegions.Compute)

	var changeDetails []string
	if len(storageChanges) > 0 {
		changeDetails = append(changeDetails, fmt.Sprintf("*Storage regions:* %s", storageChanges))
	}
	if len(computeChanges) > 0 {
		changeDetails = append(changeDetails, fmt.Sprintf("*Compute regions:* %s", computeChanges))
	}

	changeText := "No specific changes detected"
	if len(changeDetails) > 0 {
		changeText = fmt.Sprintf("%s", changeDetails[0])
		if len(changeDetails) > 1 {
			changeText += fmt.Sprintf("\n%s", changeDetails[1])
		}
	}

	message := SlackMessage{
		Channel: changesChannel,
		Attachments: []SlackAttachment{
			{
				Color: "warning",
				Title: "🔄 Provider Regions Changed",
				Text:  fmt.Sprintf("Regions have changed for provider: *%s*", provider),
				Fields: []SlackField{
					{
						Title: "Provider",
						Value: provider,
						Short: true,
					},
					{
						Title: "Storage Regions Count",
						Value: fmt.Sprintf("Old: %d → New: %d", len(oldRegions.Storage), len(newRegions.Storage)),
						Short: true,
					},
					{
						Title: "Compute Regions Count",
						Value: fmt.Sprintf("Old: %d → New: %d", len(oldRegions.Compute), len(newRegions.Compute)),
						Short: true,
					},
					{
						Title: "Changes",
						Value: changeText,
						Short: false,
					},
					{
						Title: "Timestamp",
						Value: time.Now().Format("2006-01-02 15:04:05"),
						Short: true,
					},
				},
				Timestamp: time.Now().Unix(),
			},
		},
	}

	jsonData, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		log.Printf("Failed to marshal slack message for provider %s: %v", provider, marshalErr)
		return
	}

	resp, postErr := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if postErr != nil {
		log.Printf("Failed to send slack notification for provider %s: %v", provider, postErr)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Slack notification failed for provider %s with status: %d", provider, resp.StatusCode)
		return
	}

	log.Printf("Sent regions changed notification for provider: %s", provider)
}

func calculateRegionChanges(oldRegions, newRegions map[string]string) string {
	var changes []string

	// Find added regions
	for key, value := range newRegions {
		if _, exists := oldRegions[key]; !exists {
			changes = append(changes, fmt.Sprintf("+ %s: %s", key, value))
		}
	}

	// Find removed regions
	for key, value := range oldRegions {
		if _, exists := newRegions[key]; !exists {
			changes = append(changes, fmt.Sprintf("- %s: %s", key, value))
		}
	}

	// Find modified regions
	for key, newValue := range newRegions {
		if oldValue, exists := oldRegions[key]; exists && oldValue != newValue {
			changes = append(changes, fmt.Sprintf("~ %s: %s → %s", key, oldValue, newValue))
		}
	}

	if len(changes) == 0 {
		return "No changes detected"
	}

	if len(changes) > 5 {
		return fmt.Sprintf("%d changes detected (showing first 5): %v...", len(changes), changes[:5])
	}

	return fmt.Sprintf("%d changes: %v", len(changes), changes)
}
