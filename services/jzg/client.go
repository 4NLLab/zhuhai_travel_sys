package jzg

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Client struct {
	baseURL         string
	distributorCode string
	accessToken     string
	httpClient      *http.Client
}

type Response struct {
	Code       string          `json:"code"`
	Msg        string          `json:"msg"`
	Data       json.RawMessage `json:"data"`
	Success    bool            `json:"success"`
	Pagination json.RawMessage `json:"pagination"`
	Total      json.RawMessage `json:"total"`
	Raw        json.RawMessage `json:"-"`
}

func NewClient(baseURL, distributorCode, accessToken string) (*Client, error) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" || distributorCode == "" || accessToken == "" {
		return nil, errors.New("九洲港接口配置不完整")
	}
	return &Client{
		baseURL:         baseURL,
		distributorCode: strings.TrimSpace(distributorCode),
		accessToken:     strings.TrimSpace(accessToken),
		httpClient: &http.Client{
			Timeout: 12 * time.Second,
		},
	}, nil
}

func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	if body == nil {
		body = map[string]interface{}{}
	}

	payload := map[string]interface{}{
		"body": body,
		"header": map[string]interface{}{
			"distributorCode":      c.distributorCode,
			"ticketMachineAccount": "",
			"timestamp":            time.Now().UnixMilli(),
			"terminalNo":           "",
			"locatePort":           "",
		},
	}
	signature, err := signPayload(payload, c.accessToken)
	if err != nil {
		return nil, err
	}
	payload["header"].(map[string]interface{})["secretKey"] = signature

	requestBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(requestBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("九洲港接口请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, fmt.Errorf("读取九洲港接口响应失败: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("九洲港接口 HTTP %d: %s", resp.StatusCode, string(respBytes))
	}

	var parsed Response
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, fmt.Errorf("解析九洲港接口响应失败: %w", err)
	}
	parsed.Raw = append(parsed.Raw, respBytes...)
	return &parsed, nil
}

func signPayload(payload map[string]interface{}, accessToken string) (string, error) {
	copyPayload := cloneWithoutSecret(payload)
	canonical, err := marshalCanonical(copyPayload)
	if err != nil {
		return "", err
	}
	sum := md5.Sum([]byte(canonical + accessToken))
	return hex.EncodeToString(sum[:]), nil
}

func VerifySignature(payload map[string]interface{}, accessToken string) bool {
	provided := ""
	for _, key := range []string{"secretKey", "secrectKey"} {
		if value, ok := payload[key]; ok {
			provided = fmt.Sprint(value)
			break
		}
	}
	if provided == "" {
		return false
	}
	expected, err := signPayload(payload, accessToken)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(strings.ToLower(provided)), []byte(expected)) == 1
}

func cloneWithoutSecret(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(v))
		for key, item := range v {
			if key == "secretKey" || key == "secrectKey" {
				continue
			}
			out[key] = cloneWithoutSecret(item)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(v))
		for i, item := range v {
			out[i] = cloneWithoutSecret(item)
		}
		return out
	default:
		return v
	}
}

func marshalCanonical(value interface{}) (string, error) {
	var buf bytes.Buffer
	if err := writeCanonical(&buf, value); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func writeCanonical(buf *bytes.Buffer, value interface{}) error {
	switch v := value.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		buf.WriteByte('{')
		for i, key := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			keyBytes, _ := json.Marshal(key)
			buf.Write(keyBytes)
			buf.WriteByte(':')
			if err := writeCanonical(buf, v[key]); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
	case []interface{}:
		buf.WriteByte('[')
		for i, item := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := writeCanonical(buf, item); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
	default:
		valueBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		buf.Write(valueBytes)
	}
	return nil
}
