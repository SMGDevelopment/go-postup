package postup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

const SubscriptionStatusUnsub SubscriptionStatus = "UNSUB"

type UnsubscribeRecipientRequest struct {
	Recipient *Recipient
	List      *List
}

type UnsubscribeRecipientResponse struct {
	MailingID    interface{}        `json:"mailingId"`
	RecipientID  int                `json:"recipientId"`
	ListID       int                `json:"listId"`
	Status       SubscriptionStatus `json:"status"`
	ListStatus   string             `json:"listStatus"`
	GlobalStatus string             `json:"globalStatus"`
	DateUnsub    *Time              `json:"dateUnsub"`
	DateJoined   Time               `json:"dateJoined"`
	SourceID     string             `json:"sourceId,omitempty"`
	Confirmed    bool               `json:"confirmed"`
}

// related client methods

func (pu *PostUp) UnsubscribeRecipientFromList(ctx context.Context, urr *UnsubscribeRecipientRequest) (*UnsubscribeRecipientResponse, error) {
	var (
		r   = urr.Recipient
		l   = urr.List
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
	}{
		RecipientID: rid,
		ListID:      lid,
		Status:      SubscriptionStatusUnsub,
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

	var ur UnsubscribeRecipientResponse
	if err := pu.decodeJSON(resp, &ur); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	return &ur, nil
}
