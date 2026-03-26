package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	apiKey  string
	baseURL string
	fake    bool
	http    *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func New(apiKey, baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		fake:    os.Getenv("FAKE_STREAM") == "1",
		http:    &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *Client) Stream(ctx context.Context, messages []Message, model string, temperature float64, maxTokens int) (<-chan string, error) {
	if c.fake {
		return c.fakeStream(ctx), nil
	}
	return c.openaiStream(ctx, messages, model, temperature, maxTokens)
}

func (c *Client) fakeStream(ctx context.Context) <-chan string {
	ch := make(chan string)
	words := strings.Fields(`*She tilts her head, a slow smile spreading across her face as the silver light shifts between the leaves above.* "You ask a question most travelers fear to voice." *Her fingers trace idle patterns in the air, leaving faint trails of light.* "Eldoria was not always shadowed. Once the meadows ran gold at harvest, the lake sang at dawn, and the mountain paths were crowded with merchants and wanderers alike." *Her voice dips, quiet as moss.* "Then the Shadowfangs came — not all at once, but creeping, the way frost claims a field. And I have stood here, holding this one small glade, ever since."`)
	go func() {
		defer close(ch)
		for _, word := range words {
			select {
			case <-ctx.Done():
				return
			case ch <- word + " ":
				time.Sleep(40 * time.Millisecond)
			}
		}
	}()
	return ch
}

type streamRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type streamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func (c *Client) openaiStream(ctx context.Context, messages []Message, model string, temperature float64, maxTokens int) (<-chan string, error) {
	body, _ := json.Marshal(streamRequest{
		Model:       model,
		Messages:    messages,
		Stream:      true,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("upstream status %d", resp.StatusCode)
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}
			var chunk streamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}
			if len(chunk.Choices) > 0 {
				if delta := chunk.Choices[0].Delta.Content; delta != "" {
					select {
					case <-ctx.Done():
						return
					case ch <- delta:
					}
				}
			}
		}
	}()
	return ch, nil
}
