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

func TestPostUp_GetRecipientByAddress_OK(t *testing.T) {
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

	r, err := pu.GetRecipientByAddress(ctx, "")
	is.NoErr(err)
	delete(recipient.Demographics, "_City")
	is.Equal(*r, recipient)
}

func TestPostUp_GetRecipientByAddress_NotFound(t *testing.T) {
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
		pipe2.Write([]byte("HTTP/1.1 200 OK\nContent-Type: application/json\n\n[]"))
	}()

	r, err := pu.GetRecipientByAddress(ctx, "")
	is.True(errors.Is(err, ErrRecipientNotFound))
	is.Equal(r, nil)
}

func TestPostUp_GetRecipientByAddress_Error(t *testing.T) {
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

	r, err := pu.GetRecipientByAddress(ctx, "")
	is.True(err != nil)
	is.Equal(r, nil)
}

func TestPostUp_GetRecipientByAddress_NetworkError(t *testing.T) {
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

	r, err := pu.GetRecipientByAddress(ctx, "")
	is.True(err != nil)
	is.Equal(r, nil)
}

func TestPostUp_DeleteRecipient(t *testing.T) {
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

		drr = DeleteRecipientResponse{
			Message: "MESSAGE",
			Status:  "STATUS",
		}
	)

	go func() {
		payload, err := json.Marshal(drr)
		is.NoErr(err)

		pipe2.Write(append([]byte("HTTP/1.1 200 OK\nContent-Type: application/json\n\n"), payload...))
	}()

	r, err := pu.DeleteRecipientByAddress(ctx, "")
	is.NoErr(err)
	is.Equal(*r, drr)
}

func TestPostUp_DeleteRecipient_NetworkError(t *testing.T) {
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

	r, err := pu.DeleteRecipientByAddress(ctx, "\n")
	is.True(err != nil)
	is.True(r == nil)
}

func TestPostUp_DeleteRecipient_BadJSON(t *testing.T) {
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
		pipe2.Write([]byte("HTTP/1.1 200 OK\nContent-Type: application/json\n\n."))
		pipe2.Close()
	}()

	r, err := pu.DeleteRecipientByAddress(ctx, "")
	is.True(err != nil)
	is.True(r == nil)
}
