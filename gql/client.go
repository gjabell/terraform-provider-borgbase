package gql

import "net/http"

type Client struct {
	client *http.Client
	url    string
}

func (c *Client) Query(name string, schema interface{}, args Arguments) error {
	return Execute(
		c.client,
		c.url,
		QUERY,
		name,
		schema,
		args,
	)
}

func (c *Client) Mutation(
	name string,
	schema interface{},
	args Arguments,
) error {
	return Execute(
		c.client,
		c.url,
		MUTATION,
		name,
		schema,
		args,
	)
}

func NewClient(url, apiKey string) *Client {
	c := Client{url: url}
	if apiKey != "" {
		c.client = &http.Client{Transport: NewAuthedTransport(apiKey)}
	} else {
		c.client = http.DefaultClient
	}

	return &c
}
