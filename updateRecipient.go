package postup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type UpdateRecipientRequest struct {
	Address           string            `json:"address"`
	ExternalID        string            `json:"externalId"`
	Channel           string            `json:"channel"`
	Status            string            `json:"status"`
	Comment           string            `json:"comment"`
	Password          string            `json:"password"`
	SourceDescription string            `json:"sourceDescription"`
	Demographics      map[string]string `json:"-"`
}

func (urr *UpdateRecipientRequest) MarshalJSON() ([]byte, error) {
	var buf struct {
		UpdateRecipientRequest
		DemographicsList []string `json:"demographics"`
	}

	var (
		dl  = make([]string, len(urr.Demographics))
		idx int
	)

	for k, v := range urr.Demographics {
		dl[idx] = k + "=" + v
		idx++
	}

	buf.UpdateRecipientRequest = *urr
	buf.DemographicsList = dl

	return json.Marshal(buf)
}

// related client methods

func (pu *PostUp) RecipientUpdate(ctx context.Context, id int, urr *UpdateRecipientRequest) (*Recipient, error) {
	var reqURL = pu.url("recipient", nil)

	payload, err := json.Marshal(urr)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON recipient update payload for postup: %w", err)
	}

	request, err := pu.newRequest(ctx, "PUT", reqURL+fmt.Sprintf("/%d", id), bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")

	resp, err := pu.do(request)
	if err != nil {
		return nil, fmt.Errorf("encountered network error while updating user with postup: %w", err)
	}

	var rs Recipient
	if err := pu.decodeJSON(resp, &rs); err != nil {
		return nil, fmt.Errorf("error unmarshaling json response body from postup: %w", err)
	}

	return &rs, nil
}
