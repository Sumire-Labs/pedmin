// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

// Package deepl provides a client for the DeepL translation API.
package deepl

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const translateCacheTTL = 10 * time.Minute

type cacheEntry struct {
	result    *TranslateResult
	expiresAt time.Time
}

// TranslateClient calls the DeepL API Free.
type TranslateClient struct {
	http   *http.Client
	apiKey string
	mu     sync.Mutex
	cache  map[string]cacheEntry
}

func NewTranslateClient(apiKey string, timeout time.Duration) *TranslateClient {
	return &TranslateClient{
		http:   &http.Client{Timeout: timeout},
		apiKey: apiKey,
		cache:  make(map[string]cacheEntry),
	}
}

func (c *TranslateClient) cacheKey(text, targetLang string) string {
	h := sha256.Sum256([]byte(text + "\x00" + targetLang))
	return hex.EncodeToString(h[:16])
}

func (c *TranslateClient) getCached(key string) (*TranslateResult, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.cache[key]
	if !ok || time.Now().After(e.expiresAt) {
		delete(c.cache, key)
		return nil, false
	}
	return e.result, true
}

func (c *TranslateClient) setCache(key string, result *TranslateResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{result: result, expiresAt: time.Now().Add(translateCacheTTL)}
}

// IsAvailable reports whether the client has a configured API key.
func (c *TranslateClient) IsAvailable() bool {
	return c.apiKey != ""
}

// TranslateResult holds a single translation response.
type TranslateResult struct {
	TranslatedText   string
	DetectedLanguage string // uppercase, e.g. "EN", "JA"
}

type deeplResponse struct {
	Translations []struct {
		Text               string `json:"text"`
		DetectedSourceLang string `json:"detected_source_language"`
	} `json:"translations"`
}

func (c *TranslateClient) Translate(ctx context.Context, text, targetLang string) (*TranslateResult, error) {
	key := c.cacheKey(text, targetLang)
	if cached, ok := c.getCached(key); ok {
		return cached, nil
	}

	// DeepL Free API uses api-free.deepl.com
	endpoint := "https://api-free.deepl.com/v2/translate"

	form := url.Values{}
	form.Set("text", text)
	form.Set("target_lang", strings.ToUpper(targetLang))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DeepL API returned %d: %s", resp.StatusCode, string(body))
	}

	var dResp deeplResponse
	if err := json.Unmarshal(body, &dResp); err != nil {
		return nil, fmt.Errorf("invalid DeepL response: %w", err)
	}

	if len(dResp.Translations) == 0 {
		return nil, fmt.Errorf("no translations returned")
	}

	t := dResp.Translations[0]
	result := &TranslateResult{
		TranslatedText:   t.Text,
		DetectedLanguage: strings.ToUpper(t.DetectedSourceLang),
	}
	c.setCache(key, result)
	return result, nil
}
