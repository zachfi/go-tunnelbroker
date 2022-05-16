package tunnelbroker

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	Username string
	Password string
}

type Tunnels struct {
	Tunnels []Tunnel `xml:"tunnel"`
}

type Tunnel struct {
	ClientV4    string `xml:"clientv4"`
	ClientV6    string `xml:"clientv6"`
	Description string `xml:"description"`
	ID          string `xml:"id,attr"`
	Routed48    string `xml:"routed48"`
	Routed64    string `xml:"routed64"`
	ServerV4    string `xml:"serverv4"`
	ServerV6    string `xml:"serverv6"`
}

func NewClient(username, password *string) (*Client, error) {
	if username == nil || password == nil {
		return nil, fmt.Errorf("username and password must not be nil")
	}

	if *username == "" || *password == "" {
		return nil, fmt.Errorf("username and password must not be empty")
	}

	return &Client{
		Username: *username,
		Password: *password,
	}, nil
}

func (c *Client) TunnelInfo() (*Tunnels, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://tunnelbroker.net/tunnelInfo.php", nil)

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return &Tunnels{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Tunnels{}, err
	}

	tunnels := &Tunnels{}
	err = xml.Unmarshal(body, tunnels)
	if err != nil {
		return &Tunnels{}, err
	}

	return tunnels, nil
}

func (c *Client) UpdateTunnel(tunnelId string, ipAddress string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://ipv4.tunnelbroker.net/nic/update", nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)

	q := req.URL.Query()
	q.Add("hostname", tunnelId)
	q.Add("myip", ipAddress)

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) GetTunnel(id string) (Tunnel, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://tunnelbroker.net/tunnelInfo.php", nil)
	if err != nil {
		return Tunnel{}, errors.Wrap(err, "failed to create new request")
	}

	req.SetBasicAuth(c.Username, c.Password)

	q := req.URL.Query()
	q.Add("tid", id)

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return Tunnel{}, errors.Wrap(err, "request call failed")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Tunnel{}, errors.Wrap(err, "failed to read response body")
	}

	switch string(body) {
	case "Invalid username or password.":
		return Tunnel{}, errInvalidCredentials
	}

	tunnels := &Tunnels{}
	err = xml.Unmarshal(body, tunnels)
	if err != nil {
		return Tunnel{}, errors.Wrap(err, fmt.Sprintf("failed to unmarshal XML response: %s", body))
	}

	return tunnels.Tunnels[0], nil
}
