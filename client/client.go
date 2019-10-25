package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/mapleque/vpsagent/core"
)

type Client struct {
	host  string
	token string
	v     bool
}

func New(host, token string, v bool) *Client {
	return &Client{
		host,
		token,
		v,
	}
}

func (c *Client) Do(body []byte) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	ts := fmt.Sprintf("%d", time.Now().Unix())
	sign := core.MakeSignature(c.token, ts, body)

	r, err := http.NewRequest(
		"POST",
		c.host,
		bytes.NewReader(body),
	)
	if err != nil {
		panic(err)
	}

	if c.v {
		fmt.Printf("--------\n")
		fmt.Printf("POST %s HTTP/1.1\n", c.host)
		fmt.Printf("Content-Type: text/plain\n")
		fmt.Printf("Timestamp: %s\n", ts)
		fmt.Printf("Signature: %s\n", sign)

		fmt.Printf("\n%s\n", body)
		fmt.Printf("--------\n")
	}

	r.Header.Set("Content-Type", "text/plain")
	r.Header.Set("Timestamp", ts)
	r.Header.Set("Signature", sign)

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n%s\n", resp.Status, ret)
	os.Exit(0)
}
