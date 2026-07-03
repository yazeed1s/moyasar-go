package moyasar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) do(ctx context.Context, method, path string, query url.Values, in any, out any) error {
	req, err := c.newRequest(ctx, c.baseURL, method, path, query, in)
	if err != nil {
		return err
	}
	return c.send(req, out)
}

func (c *Client) doForm(ctx context.Context, method, path string, form url.Values, out any) error {
	req, err := c.newFormRequest(ctx, c.baseURL, method, path, nil, form)
	if err != nil {
		return err
	}
	return c.send(req, out)
}

func (c *Client) doTransfer(ctx context.Context, method, path string, query url.Values, out any) error {
	req, err := c.newRequest(ctx, c.transferURL, method, path, query, nil)
	if err != nil {
		return err
	}
	return c.send(req, out)
}

func (c *Client) newRequest(ctx context.Context, baseURL *url.URL, method, path string, query url.Values, in any) (*http.Request, error) {
	if baseURL == nil {
		return nil, fmt.Errorf("moyasar: base URL is not configured")
	}

	u := *baseURL
	u.Path = strings.TrimRight(baseURL.Path, "/") + "/" + strings.TrimLeft(path, "/")
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	var body io.Reader
	if in != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(in); err != nil {
			return nil, fmt.Errorf("moyasar: encode request: %w", err)
		}
		body = &buf
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.apiKey, "")
	req.Header.Set("Accept", "application/json")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) newFormRequest(ctx context.Context, baseURL *url.URL, method, path string, query url.Values, form url.Values) (*http.Request, error) {
	if baseURL == nil {
		return nil, fmt.Errorf("moyasar: base URL is not configured")
	}

	u := *baseURL
	u.Path = strings.TrimRight(baseURL.Path, "/") + "/" + strings.TrimLeft(path, "/")
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.apiKey, "")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (c *Client) send(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("moyasar: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("moyasar: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return decodeAPIError(resp.StatusCode, body)
	}
	if out == nil || len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("moyasar: decode response: %w", err)
	}
	return nil
}
