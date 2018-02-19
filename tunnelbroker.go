package tunnelbroker

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Username string
	Password string
}

type Tunnels struct {
	Tunnels []Tunnel `xml:"tunnel"`
}

type Tunnel struct {
	Id          int    `xml:"id,attr"`
	Description string `xml:"description"`
	ServerV4    string `xml:"serverv4"`
	ClientV4    string `xml:"clientv4"`
	ServerV6    string `xml:"serverv6"`
	Clientv6    string `xml:"clientv6"`
	Routed64    string `xml:"routed64"`
	Routed48    string `xml:"routed48"`
}

func (c *Client) TunnelInfo() (Tunnels, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://tunnelbroker.net/tunnelInfo.php", nil)
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := client.Do(req)
	if err != nil {
		return Tunnels{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Tunnels{}, err
	}

	tunnels := Tunnels{}
	err = xml.Unmarshal(body, &tunnels)
	if err != nil {
		return Tunnels{}, err
	}

	return tunnels, nil
}
