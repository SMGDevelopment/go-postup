package postup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type SubscriptionStatus string

const (
	SubscriptionStatusNormal = "NORMAL"
	SubscriptionStatusUnsub  = "UNSUB"
)

type SubscribeRecipientRequest struct {
	Recipient *Recipient
	List      *List
	SourceID  string
	Confirmed bool
}

type SubscribeRecipientResponse struct {
	MailingID    interface{}        `json:"mailingId"`
	RecipientID  int                `json:"recipientId"`
	ListID       int                `json:"listId"`
	Status       SubscriptionStatus `json:"status"`
	ListStatus   string             `json:"listStatus"`
	GlobalStatus string             `json:"globalStatus"`
	DateUnsub    *time.Time         `json:"dateUnsub"`
	DateJoined   time.Time          `json:"dateJoined"`
	SourceID     string             `json:"sourceId,omitempty"`
	Confirmed    bool               `json:"confirmed"`
}

// related client methods

func (pu *PostUp) SubscribeRecipientToList(ctx context.Context, srr *SubscribeRecipientRequest) (*SubscribeRecipientResponse, error) {
	var (
		r   = srr.Recipient
		l   = srr.List
		u   = pu.url("listsubscription", nil)
		err error
	)

	rid, err := r.getID(ctx, pu)
	if err != nil {
		return nil, err
	}

	lid, err := l.getID(ctx, pu)
	if err != nil {
		return nil, err
	}

	var v = struct {
		RecipientID int                `json:"recipientId"`
		ListID      int                `json:"listId"`
		Status      SubscriptionStatus `json:"status"`
		SourceID    string             `json:"sourceId,omitempty"`
		Confirmed   bool               `json:"confirmed"`
	}{
		RecipientID: rid,
		ListID:      lid,
		Status:      SubscriptionStatusNormal,
		SourceID:    srr.SourceID,
		Confirmed:   srr.Confirmed,
	}

	payload, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	req, err := pu.newRequest(ctx, "POST", u, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := pu.do(req)
	if err != nil {
		return nil, fmt.Errorf("encountered network error while subscribing a recipient to a mailing list: %w", err)
	}

	var sr SubscribeRecipientResponse
	if err := pu.decodeJSON(resp, &sr); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	return &sr, nil
}
