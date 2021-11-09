package postup

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/matryer/is"
)

func TestPostUp_RecipientCreate_OK(t *testing.T) {
	var (
		is = is.New(t)

		ctx = context.Background()

		pipe1, pipe2 = net.Pipe()
		pu           = PostUp{
			baseURL: url.URL{
				Scheme: "http",
				Host:   "kraft.singles",
			},
			client: &http.Client{
				Transport: &http.Transport{
					DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
						return pipe1, nil
					},
				},
			},
		}

		recipient = Recipient{
			ID:         3,
			Address:    "test@test.com",
			ExternalID: "test@test.com",
			Demographics: map[string]string{
				"FirstName": "FIRST_NAME",
				"LastName":  "LAST_NAME",
				"_City":     "",
			},
		}
	)

	go func() {
		payload, err := json.Marshal([]interface{}{&recipient})
		is.NoErr(err)

		pipe2.Write(append([]byte("HTTP/1.1 200 OK\nContent-Type: application/json\n\n"), payload...))
	}()

	r, err := pu.RecipientCreate(ctx, &CreateRecipientRequest{
		Address:    "test@test.com",
		ExternalID: "test@test.com",
		Demographics: map[string]string{
			"FirstName": "FIRST_NAME",
			"LastName":  "LAST_NAME",
			"_City":     "",
		},
	})
	is.NoErr(err)
	delete(recipient.Demographics, "_City")
	is.Equal(*r, recipient)
}

func TestPostUp_RecipientCreate_Error(t *testing.T) {
	var (
		is = is.New(t)

		ctx = context.Background()

		pipe1, pipe2 = net.Pipe()
		pu           = PostUp{
			baseURL: url.URL{
				Scheme: "http",
				Host:   "kraft.singles",
			},
			client: &http.Client{
				Transport: &http.Transport{
					DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
						return pipe1, nil
					},
				},
			},
		}
	)

	go func() {
		pipe2.Write([]byte("HTTP/1.1 500 OK\nContent-Type: application/json\n\n"))
	}()

	r, err := pu.RecipientCreate(ctx, &CreateRecipientRequest{
		Address:    "test@test.com",
		ExternalID: "test@test.com",
		Demographics: map[string]string{
			"FirstName": "FIRST_NAME",
			"LastName":  "LAST_NAME",
			"_City":     "",
		},
	})
	is.True(err != nil)
	is.Equal(r, nil)
}

func TestPostUp_RecipientCreate_NetworkError(t *testing.T) {
	var (
		is = is.New(t)

		ctx = context.Background()

		pu = PostUp{
			baseURL: url.URL{
				Scheme: "http",
				Host:   "kraft.singles",
			},
			client: &http.Client{
				Transport: &http.Transport{
					DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
						return nil, errors.New("ERROR")
					},
				},
			},
		}
	)

	r, err := pu.RecipientCreate(ctx, &CreateRecipientRequest{
		Address:    "test@test.com",
		ExternalID: "test@test.com",
		Demographics: map[string]string{
			"FirstName": "FIRST_NAME",
			"LastName":  "LAST_NAME",
			"_City":     "",
		},
	})
	is.True(err != nil)
	is.Equal(r, nil)
}
