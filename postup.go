package postup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	ErrRecipientNotFound = errors.New("recipient not found")
	ErrListNotFound      = errors.New("list not found")

	ErrRecipientMissingIdentifier = errors.New("recipient has no address, recipientId or externalId")
	ErrListMissingIdentifier      = errors.New("recipient has no listId or externalId")
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "api.postup.com",
}

type PostUp struct {
	client  *http.Client
	baseURL url.URL

	username, password string
}

func NewPostUp(user, passwd string, c *http.Client) *PostUp {
	var pu = PostUp{
		username: user,
		password: passwd,
		client:   c,
		baseURL:  baseURL,
	}

	if pu.client == nil {
		pu.client = http.DefaultClient
	}

	return &pu
}

func (pu *PostUp) url(endpoint string, query url.Values) string {
	if pu.baseURL.Host == "" {
		pu.baseURL = baseURL
	}

	var reqURL = pu.baseURL
	reqURL.Path = "api/" + endpoint
	if query != nil {
		reqURL.RawQuery = query.Encode()
	}
	return reqURL.String()
}

func (pu *PostUp) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	r.SetBasicAuth(pu.username, pu.password)

	return r, nil
}

func (pu *PostUp) decodeJSON(r *http.Response, v interface{}) error {
	defer r.Body.Close()

	if code := r.StatusCode; code != http.StatusOK {
		v = nil
		return fmt.Errorf("received non-200 status: %d - %s", code, r.Status)
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func (pu *PostUp) do(r *http.Request) (*http.Response, error) {
	if pu.client == nil {
		pu.client = http.DefaultClient
	}

	return pu.client.Do(r)
}

func (pu *PostUp) getRecipient(ctx context.Context, url string) (*Recipient, error) {
	req, err := pu.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := pu.do(req)
	if err != nil {
		return nil, fmt.Errorf("encountered network error: %w", err)
	}

	var rs []*Recipient
	if err := pu.decodeJSON(resp, &rs); err != nil {
		return nil, err
	}

	if 0 < len(rs) {
		return rs[0], nil
	}

	return nil, ErrRecipientNotFound
}
