package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	anthropicAPIURL     = "https://api.anthropic.com/v1/messages"
	anthropicAPIVersion = "2023-06-01"
)

// Client handles communication with Claude API
type Client struct {
	apiKey     string
	model      string
	maxTokens  int
	timeout    time.Duration
	httpClient *http.Client
	log        *zap.Logger
}

// ClientConfig holds configuration for the Claude client
type ClientConfig struct {
	APIKey         string
	Model          string
	MaxTokens      int
	TimeoutSeconds int
}

// NewClient creates a new Claude API client
func NewClient(cfg ClientConfig, log *zap.Logger) *Client {
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		apiKey:    cfg.APIKey,
		model:     cfg.Model,
		maxTokens: cfg.MaxTokens,
		timeout:   timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		log: log,
	}
}

// Chat sends a message to Claude and returns the response
func (c *Client) Chat(ctx context.Context, systemPrompt string, messages []Message) (*AIResponse, error) {
	startTime := time.Now()

	req := ChatRequest{
		Model:       c.model,
		MaxTokens:   c.maxTokens,
		System:      systemPrompt,
		Messages:    messages,
		Temperature: 0.7,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", anthropicAPIURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", anthropicAPIVersion)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.log.Error("Claude API error",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(respBody)),
		)
		return nil, fmt.Errorf("Claude API error: %s (status %d)", string(respBody), resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract text content
	var content string
	for _, block := range chatResp.Content {
		if block.Type == "text" {
			content = block.Text
			break
		}
	}

	responseTime := time.Since(startTime)

	c.log.Debug("Claude API call completed",
		zap.Duration("response_time", responseTime),
		zap.Int("input_tokens", chatResp.Usage.InputTokens),
		zap.Int("output_tokens", chatResp.Usage.OutputTokens),
	)

	return &AIResponse{
		Content:      content,
		InputTokens:  chatResp.Usage.InputTokens,
		OutputTokens: chatResp.Usage.OutputTokens,
		ResponseTime: responseTime,
	}, nil
}

// ChatWithRetry sends a message with retry logic for transient failures
func (c *Client) ChatWithRetry(ctx context.Context, systemPrompt string, messages []Message, maxRetries int) (*AIResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		resp, err := c.Chat(ctx, systemPrompt, messages)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		c.log.Warn("Claude API call failed, retrying",
			zap.Int("attempt", attempt+1),
			zap.Int("max_retries", maxRetries),
			zap.Error(err),
		)
	}

	return nil, fmt.Errorf("all retries exhausted: %w", lastErr)
}

// IsEnabled returns true if the client has a valid API key configured
func (c *Client) IsEnabled() bool {
	return c.apiKey != ""
}
