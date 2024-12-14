package libs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func CallOpenAI(apiKey, prompt string) (OpenAIResponse, error) {
	url := "https://api.openai.com/v1/chat/completions"

	payload := map[string]interface{}{
		"model": "gpt-4o-2024-11-20", // Use GPT-4o model
		"messages": []map[string]string{
			{"role": "system", "content": "You are an AI that helps suggest relevant article categories based on content."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  500,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return OpenAIResponse{}, fmt.Errorf("OpenAI API error: %s", body)
	}

	var apiResponse OpenAIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return apiResponse, nil
}
