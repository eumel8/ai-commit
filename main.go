// commitmsg.go
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	defaultModel       = "gpt-4o-mini" // wähle nach Bedarf
	apiURL             = "https://api.openai.com/v1/chat/completions"
	maxDiffBytes       = 200_000
	httpTimeoutSeconds = 25
)

func main() {
	if len(os.Args) < 2 {
		fail("usage: commitmsg <path-to-message-file>")
	}
	msgPath := os.Args[1]

	diff, err := stagedDiff()
	if err != nil {
		fail("git diff failed: " + err.Error())
	}
	if len(diff) == 0 {
		os.Exit(0)
	}
	if len(diff) > maxDiffBytes {
		diff = diff[:maxDiffBytes]
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fail("OPENAI_API_KEY not set")
	}

	msg, err := callOpenAI(apiKey, defaultModel, buildPrompt(diff))
	if err != nil || strings.TrimSpace(msg) == "" {
		msg = "chore: update changes"
	}

	if err := os.WriteFile(msgPath, []byte(strings.TrimSpace(msg)+"\n"), 0644); err != nil {
		fail("write message failed: " + err.Error())
	}
}

func stagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	out, err := cmd.Output()
	return string(out), err
}

func buildPrompt(diff string) string {
	return "Create a Commit-Message following Conventional Commits.\n" +
		"English. Subject max. 72 sign, optional body with explanation, without Markdown, only plain text.\n\nDiff:\n" +
		diff
}

// -------- OpenAI --------

type oaMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type oaReq struct {
	Model    string  `json:"model"`
	Messages []oaMsg `json:"messages"`
}
type oaResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func callOpenAI(apiKey, model, prompt string) (string, error) {
	body, _ := json.Marshal(oaReq{
		Model: model,
		Messages: []oaMsg{
			{Role: "system", Content: "You write excellent Commit-Messages."},
			{Role: "user", Content: prompt},
		},
	})
	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: httpTimeoutSeconds * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", errors.New("api status: " + resp.Status + " body: " + string(b))
	}
	var parsed oaResp
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	if len(parsed.Choices) == 0 {
		return "", errors.New("no choices")
	}
	return parsed.Choices[0].Message.Content, nil
}

// -------- utils --------

func fail(msg string) {
	os.Stderr.WriteString(msg + "\n")
	os.Exit(1)
}

