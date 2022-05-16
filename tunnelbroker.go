package tunnelbroker

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Client struct {
	Username string
	Password string
}

type Tunnels struct {
	Tunnels []Tunnel `xml:"tunnel"`
}

type Tunnel struct {
	Description string `xml:"description"`
	ServerV4    string `xml:"serverv4"`
	ClientV4    string `xml:"clientv4"`
	ServerV6    string `xml:"serverv6"`
	Clientv6    string `xml:"clientv6"`
	Routed64    string `xml:"routed64"`
	Routed48    string `xml:"routed48"`
	ID          int    `xml:"id,attr"`
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

func (c *Client) UpdateTunnel(tunnelId int, ipAddress string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://ipv4.tunnelbroker.net/nic/update", nil)
	req.SetBasicAuth(c.Username, c.Password)

	q := req.URL.Query()
	q.Add("hostname", strconv.Itoa(tunnelId))

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	return nil
}
