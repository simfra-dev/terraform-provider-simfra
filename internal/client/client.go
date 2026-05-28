package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL    string
	adminToken string
	httpClient *http.Client
}

func New(endpoint, adminToken string, skipTLS bool) *Client {
	transport := &http.Transport{}
	if skipTLS {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	return &Client{
		baseURL:    strings.TrimRight(endpoint, "/") + "/_simfra",
		adminToken: adminToken,
		httpClient: &http.Client{Transport: transport},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	u, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("building URL: %w", err)
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.adminToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.adminToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		var errResp struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		msg := errResp.Error
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return nil, &APIError{StatusCode: resp.StatusCode, Message: msg}
	}

	return resp, nil
}

func decode[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()
	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result, nil
}

// --- Account ---

func (c *Client) CreateAccount(ctx context.Context, req CreateAccountRequest) (*AccountDetail, error) {
	resp, err := c.do(ctx, http.MethodPost, "/accounts", req)
	if err != nil {
		return nil, err
	}
	return decode[AccountDetail](resp)
}

func (c *Client) GetAccount(ctx context.Context, id string) (*AccountDetail, error) {
	resp, err := c.do(ctx, http.MethodGet, "/accounts/"+id, nil)
	if err != nil {
		return nil, err
	}
	return decode[AccountDetail](resp)
}

func (c *Client) ListAccounts(ctx context.Context) ([]AccountSummary, error) {
	resp, err := c.do(ctx, http.MethodGet, "/accounts", nil)
	if err != nil {
		return nil, err
	}
	result, err := decode[[]AccountSummary](resp)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) DeleteAccount(ctx context.Context, id string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/accounts/"+id, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// --- Port Forward ---

func (c *Client) CreatePortForward(ctx context.Context, req CreatePortForwardRequest) (*PortForwardSession, error) {
	resp, err := c.do(ctx, http.MethodPost, "/port-forwards", req)
	if err != nil {
		return nil, err
	}
	return decode[PortForwardSession](resp)
}

func (c *Client) GetPortForward(ctx context.Context, id string) (*PortForwardSession, error) {
	resp, err := c.do(ctx, http.MethodGet, "/port-forwards/"+id, nil)
	if err != nil {
		return nil, err
	}
	return decode[PortForwardSession](resp)
}

func (c *Client) ListPortForwards(ctx context.Context) ([]PortForwardSession, error) {
	resp, err := c.do(ctx, http.MethodGet, "/port-forwards", nil)
	if err != nil {
		return nil, err
	}
	result, err := decode[[]PortForwardSession](resp)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) DeletePortForward(ctx context.Context, id string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/port-forwards/"+id, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) ListPortForwardTargets(ctx context.Context) ([]PortForwardTarget, error) {
	resp, err := c.do(ctx, http.MethodGet, "/port-forwards/targets", nil)
	if err != nil {
		return nil, err
	}
	result, err := decode[[]PortForwardTarget](resp)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

// --- SSO Session ---

func (c *Client) CreateSSOSession(ctx context.Context, req CreateSSOSessionRequest) (*SSOSessionResponse, error) {
	resp, err := c.do(ctx, http.MethodPost, "/sso/sessions", req)
	if err != nil {
		return nil, err
	}
	return decode[SSOSessionResponse](resp)
}

func (c *Client) ListSSOSessions(ctx context.Context) ([]SSOSessionInfo, error) {
	resp, err := c.do(ctx, http.MethodGet, "/sso/sessions", nil)
	if err != nil {
		return nil, err
	}
	result, err := decode[[]SSOSessionInfo](resp)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) GetSSOSession(ctx context.Context, token string) (*SSOSessionInfo, error) {
	sessions, err := c.ListSSOSessions(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range sessions {
		if s.Token == token {
			return &s, nil
		}
	}
	return nil, &APIError{StatusCode: 404, Message: "SSO session not found"}
}

func (c *Client) DeleteSSOSession(ctx context.Context, token string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/sso/sessions/"+token, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// --- ACM Validation ---

func (c *Client) ListACMPendingValidations(ctx context.Context, accountID, region string) ([]ACMPendingValidation, error) {
	path := fmt.Sprintf("/acm/%s/%s/pending-validations", accountID, region)
	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	result, err := decode[[]ACMPendingValidation](resp)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) ValidateACMCertificate(ctx context.Context, accountID, region, arn string) error {
	path := fmt.Sprintf("/acm/%s/%s/validate/%s", accountID, region, url.PathEscape(arn))
	resp, err := c.do(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *Client) ValidateACMCertificateDomain(ctx context.Context, accountID, region, arn, domain string) error {
	path := fmt.Sprintf("/acm/%s/%s/validate/%s/%s", accountID, region, url.PathEscape(arn), domain)
	resp, err := c.do(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// --- Read-only data source endpoints ---

func (c *Client) Health(ctx context.Context) (bool, error) {
	resp, err := c.do(ctx, http.MethodGet, "/health", nil)
	if err != nil {
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == 503 {
			return false, nil
		}
		return false, err
	}
	resp.Body.Close()
	return true, nil
}

func (c *Client) ListServices(ctx context.Context) ([]ServiceInfo, error) {
	resp, err := c.do(ctx, http.MethodGet, "/services", nil)
	if err != nil {
		return nil, err
	}
	result, err := decode[[]ServiceInfo](resp)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) StorageSummary(ctx context.Context) (*StorageSummary, error) {
	resp, err := c.do(ctx, http.MethodGet, "/storage/summary", nil)
	if err != nil {
		return nil, err
	}
	return decode[StorageSummary](resp)
}

func (c *Client) DockerSummary(ctx context.Context) (*DockerSummary, error) {
	resp, err := c.do(ctx, http.MethodGet, "/docker/summary", nil)
	if err != nil {
		return nil, err
	}
	return decode[DockerSummary](resp)
}

func (c *Client) ListDockerContainers(ctx context.Context) ([]DockerContainer, error) {
	resp, err := c.do(ctx, http.MethodGet, "/docker/containers", nil)
	if err != nil {
		return nil, err
	}
	result, err := decode[[]DockerContainer](resp)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) ListDockerNetworks(ctx context.Context) ([]DockerNetwork, error) {
	resp, err := c.do(ctx, http.MethodGet, "/docker/networks", nil)
	if err != nil {
		return nil, err
	}
	result, err := decode[[]DockerNetwork](resp)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) CAInfo(ctx context.Context) (*CAInfo, error) {
	resp, err := c.do(ctx, http.MethodGet, "/ca/info", nil)
	if err != nil {
		return nil, err
	}
	return decode[CAInfo](resp)
}

func (c *Client) DNSPort(ctx context.Context, accountID string) (int, error) {
	resp, err := c.do(ctx, http.MethodGet, "/dns/"+accountID, nil)
	if err != nil {
		return 0, err
	}
	result, err := decode[PortResponse](resp)
	if err != nil {
		return 0, err
	}
	return result.Port, nil
}

func (c *Client) SMTPPort(ctx context.Context, accountID string) (int, error) {
	resp, err := c.do(ctx, http.MethodGet, "/smtp/"+accountID, nil)
	if err != nil {
		return 0, err
	}
	result, err := decode[PortResponse](resp)
	if err != nil {
		return 0, err
	}
	return result.Port, nil
}
