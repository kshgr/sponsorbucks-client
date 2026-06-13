package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const DefaultBaseURL = "https://enbmzimfbtnfrkpevzzf.supabase.co/functions/v1"

type Client struct {
	baseURL     string
	deviceToken string
	httpClient  *http.Client
}

func New(baseURL, deviceToken string) *Client {
	return NewWithHTTPClient(baseURL, deviceToken, &http.Client{Timeout: 15 * time.Second})
}

func NewWithHTTPClient(baseURL, deviceToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		baseURL:     strings.TrimRight(baseURL, "/"),
		deviceToken: deviceToken,
		httpClient:  httpClient,
	}
}

type StartLinkRequest struct {
	DevicePublicKey string `json:"device_public_key"`
	DeviceName      string `json:"device_name"`
	ClientVersion   string `json:"client_version"`
	BuildID         string `json:"build_id"`
	BuildChannel    string `json:"build_channel"`
	OS              string `json:"os"`
	Arch            string `json:"arch"`
}

type StartLinkResponse struct {
	DeviceID  string `json:"device_id"`
	LinkCode  string `json:"link_code"`
	LinkURL   string `json:"link_url"`
	ExpiresAt string `json:"expires_at"`
}

func (c *Client) StartLink(req StartLinkRequest) (StartLinkResponse, error) {
	var out StartLinkResponse
	err := c.postJSON("/device-start-link", req, &out, nil)
	return out, err
}

type CompleteLinkRequest struct {
	DeviceID string `json:"device_id"`
	LinkCode string `json:"link_code"`
}

type CompleteLinkResponse struct {
	Status      string `json:"status"`
	DeviceToken string `json:"device_token"`
	UserID      string `json:"user_id"`
}

func (c *Client) CompleteLink(req CompleteLinkRequest) (CompleteLinkResponse, error) {
	var out CompleteLinkResponse
	err := c.postJSON("/device-complete-link", req, &out, nil)
	return out, err
}

type AdResponse struct {
	CampaignID      string `json:"campaign_id"`
	CreativeID      string `json:"creative_id"`
	Line            string `json:"line"`
	DestinationURL  string `json:"destination_url"`
	TrackingURL     string `json:"tracking_url,omitempty"`
	DisplayMS       int    `json:"display_ms"`
	RotationAllowed bool   `json:"rotation_allowed"`
	NoAd            bool   `json:"no_ad"`
}

type nextAdRequest struct {
	DeviceID       string `json:"device_id,omitempty"`
	Surface        string `json:"surface"`
	LastCreativeID string `json:"last_creative_id,omitempty"`
}

func (c *Client) Health() error {
	httpReq, err := http.NewRequest(http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("GET /health failed: %s %s", resp.Status, string(body))
	}
	return nil
}

func (c *Client) NextAd(deviceID, surface string, lastCreativeID string) (AdResponse, error) {
	var out AdResponse
	err := c.postJSON("/ads-next", nextAdRequest{DeviceID: deviceID, Surface: surface, LastCreativeID: lastCreativeID}, &out, nil)
	return out, err
}

type ClickEvent struct {
	EventType     string `json:"event_type"`
	SessionID     string `json:"session_id"`
	DeviceID      string `json:"device_id"`
	CampaignID    string `json:"campaign_id"`
	CreativeID    string `json:"creative_id"`
	CreativeHash  string `json:"creative_hash"`
	Surface       string `json:"surface"`
	ClickedAt     string `json:"clicked_at"`
	ClientVersion string `json:"client_version"`
	BuildID       string `json:"build_id"`
	BuildChannel  string `json:"build_channel"`
}

func (c *Client) PathReachable(path string) error {
	httpReq, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%s not reachable: %s", path, resp.Status)
	}
	return nil
}

func (c *Client) PostSigned(path string, body []byte, deviceID, signature, timestamp string) error {
	endpoint := c.baseURL + path
	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-SponsorBucks-Device-Id", deviceID)
	httpReq.Header.Set("X-SponsorBucks-Signature", signature)
	httpReq.Header.Set("X-SponsorBucks-Timestamp", timestamp)
	c.addAuth(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("%s failed: %s %s", path, resp.Status, string(body))
	}
	return nil
}

func (c *Client) postJSON(path string, req any, out any, headers map[string]string) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}
	c.addAuth(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("%s failed: %s %s", path, resp.Status, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) addAuth(req *http.Request) {
	if c.deviceToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.deviceToken)
	}
}
