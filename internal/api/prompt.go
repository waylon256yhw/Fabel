package api

import (
	"strings"

	"fabel/llm"
)

func buildPrompt(conv *ConversationDetail, userContent string) []llm.Message {
	replacer := strings.NewReplacer("{{char}}", conv.Character.Name, "{{user}}", "User")

	system := replacer.Replace(conv.Preset.SystemPrompt)
	if conv.Character.Description != "" {
		system += "\n\n[Character: " + replacer.Replace(conv.Character.Description) + "]"
	}
	if conv.Character.Personality != "" {
		system += "\n[Personality: " + replacer.Replace(conv.Character.Personality) + "]"
	}
	if conv.Character.Scenario != "" {
		system += "\n[Scenario: " + replacer.Replace(conv.Character.Scenario) + "]"
	}

	messages := []llm.Message{{Role: "system", Content: system}}
	for _, msg := range conv.Messages {
		messages = append(messages, llm.Message{Role: string(msg.Role), Content: msg.Content})
	}
	messages = append(messages, llm.Message{Role: "user", Content: userContent})
	return messages
}
