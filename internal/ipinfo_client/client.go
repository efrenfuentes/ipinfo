package ipinfo_client

import (
	"encoding/json"
	"net/http"
)


// Client is a client for the ipinfo.io API.
type Client struct {
	host       string
	httpClient *http.Client
	accessToken     string
}

// Details contains details about an IP address.
type Details struct {
	IP          string `json:"ip"`
	City        string `json:"city"`
	Region      string `json:"region"`
	Country     string `json:"country"`
	Loc         string `json:"loc"`
	Org         string `json:"org"`
	Postal      string `json:"postal"`
	Timezone    string `json:"timezone"`
}

// NewClient returns a new Client.
func NewClient() *Client {
	return &Client{
		host:       "ipinfo.io",
		httpClient: http.DefaultClient,
	}
}

// SetAccessToken sets the access token to use for requests.
func (c *Client) SetAccessToken(token string) {
	c.accessToken = token
}

// GetDetails returns details for the given IP address.
func (c *Client) GetDetails(ip string) (*Details, error) {
	// Build the request URL.
	url := "https://" + c.host + "/" + ip
	if c.accessToken != "" {
		url += "?token=" + c.accessToken
	}

	// Make the request.
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response.
	var details Details
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}