package postup

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Recipient struct {
	ID                int               `json:"recipientId"`
	ExternalID        string            `json:"externalId"`
	ImportID          int               `json:"importId"`
	Address           string            `json:"address"`
	Channel           string            `json:"channel"`
	Status            string            `json:"status"`
	Comment           string            `json:"comment"`
	Password          string            `json:"password"`
	SourceDescription string            `json:"sourceDescription"`
	SourceSignupDate  time.Time         `json:"sourceSignupDate"`
	SignupMethod      string            `json:"signupMethod"`
	DateJoined        time.Time         `json:"dateJoined"`
	Demographics      map[string]string `json:"-"`
}

func (r *Recipient) UnmarshalJSON(data []byte) error {
	type recipient struct {
		ID                int               `json:"recipientId"`
		ExternalID        string            `json:"externalId"`
		ImportID          int               `json:"importId"`
		Address           string            `json:"address"`
		Channel           string            `json:"channel"`
		Status            string            `json:"status"`
		Comment           string            `json:"comment"`
		Password          string            `json:"password"`
		SourceDescription string            `json:"sourceDescription"`
		SourceSignupDate  time.Time         `json:"sourceSignupDate"`
		SignupMethod      string            `json:"signupMethod"`
		DateJoined        time.Time         `json:"dateJoined"`
		Demographics      map[string]string `json:"-"`
	}
	var buf struct {
		recipient
		DemographicsList []string `json:"demographics"`
	}

	if err := json.Unmarshal(data, &buf); err != nil {
		return err
	}

	var d = make(map[string]string, len(buf.Demographics))
	for _, str := range buf.DemographicsList {
		var parts = strings.Split(str, "=")
		if len(parts) != 2 {
			continue
		}

		if parts[1] == "" {
			continue
		}

		d[parts[0]] = parts[1]
	}

	buf.recipient.Demographics = d
	*r = Recipient(buf.recipient)

	return nil
}

func (r *Recipient) MarshalJSON() ([]byte, error) {
	var buf struct {
		Recipient
		DemographicsList []string `json:"demographics"`
	}

	var (
		dl  = make([]string, len(r.Demographics))
		idx int
	)

	for k, v := range r.Demographics {
		dl[idx] = k + "=" + v
		idx++
	}

	buf.Recipient = *r
	buf.DemographicsList = dl

	return json.Marshal(buf)
}

// related client methods

func (r *Recipient) getID(ctx context.Context, pu *PostUp) (int, error) {
	if r.ID != 0 {
		return r.ID, nil
	}

	var err error
	switch {
	case r.Address != "":
		r, err = pu.GetRecipientByAddress(ctx, r.Address)
	case r.ExternalID != "":
		r, err = pu.GetRecipientByExternalID(ctx, r.ExternalID)
	default:
		err = ErrRecipientMissingIdentifier
	}
	if err != nil {
		return 0, err
	}

	return r.ID, nil
}

func (pu *PostUp) GetRecipientByAddress(ctx context.Context, addr string) (*Recipient, error) {
	var u = pu.url("recipient", url.Values{"address": []string{addr}})

	return pu.getRecipient(ctx, u)
}

func (pu *PostUp) GetRecipientByExternalID(ctx context.Context, id string) (*Recipient, error) {
	var u = pu.url("recipient", url.Values{"externalId": []string{id}})

	return pu.getRecipient(ctx, u)
}

func (pu *PostUp) GetRecipientByRecipientID(ctx context.Context, id int) (*Recipient, error) {
	var u = pu.url(fmt.Sprintf("recipient/%d", id), nil)

	return pu.getRecipient(ctx, u)
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

type DeleteRecipientResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (pu *PostUp) DeleteRecipientByAddress(ctx context.Context, addr string) (*DeleteRecipientResponse, error) {
	var u = pu.url("recipient/privacy", url.Values{
		"address": []string{addr},
		"scope":   []string{"forget"},
	})

	req, err := pu.newRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := pu.do(req)
	if err != nil {
		return nil, fmt.Errorf("encountered network error: %w", err)
	}

	var drr DeleteRecipientResponse
	if err := pu.decodeJSON(resp, &drr); err != nil {
		return nil, err
	}

	return &drr, nil
}
