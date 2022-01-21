package postup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type CreateRecipientRequest struct {
	ExternalID        string            `json:"externalId,omitempty"`
	Address           string            `json:"address,omitempty"`
	Channel           string            `json:"channel,omitempty"`
	Status            string            `json:"status,omitempty"`
	Comment           string            `json:"comment,omitempty"`
	Password          string            `json:"password,omitempty"`
	SourceDescription string            `json:"sourceDescription,omitempty"`
	Demographics      map[string]string `json:"-"`
}

func (crr *CreateRecipientRequest) MarshalJSON() ([]byte, error) {
	var buf struct {
		CreateRecipientRequest
		DemographicsList []string `json:"demographics,omitempty"`
	}

	var (
		dl  = make([]string, len(crr.Demographics))
		idx int
	)

	for k, v := range crr.Demographics {
		dl[idx] = k + "=" + v
		idx++
	}

	buf.CreateRecipientRequest = *crr
	buf.DemographicsList = dl

	return json.Marshal(buf)
}

// related client methods

func (pu *PostUp) RecipientCreate(ctx context.Context, crr *CreateRecipientRequest) (*Recipient, error) {
	var reqURL = pu.url("recipient", nil)

	payload, err := json.Marshal(crr)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON recipient create payload for PostUp: %w", err)
	}

	req, err := pu.newRequest(ctx, "POST", reqURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := pu.do(req)
	if err != nil {
		return nil, fmt.Errorf("encountered network error: %w", err)
	}

	var rs Recipient
	if err := pu.decodeJSON(resp, &rs); err != nil {
		return nil, fmt.Errorf("error unmarshaling json response body from postup: %w", err)
	}

	return &rs, nil
}
