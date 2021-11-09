package postup

import (
	"context"
	"fmt"
	"net/url"
)

type List struct {
	ID              int    `json:"listId"`
	Title           string `json:"title"`
	FriendlyTitle   string `json:"friendlyTitle"`
	Description     string `json:"description"`
	Populated       bool   `json:"populated"`
	PublicSignup    bool   `json:"publicSignup"`
	GlobalUnsub     bool   `json:"globalUnsub"`
	Query           string `json:"query"`
	CategoryID      int    `json:"categoryId"`
	BlockDomains    string `json:"blockDomains"`
	SeedListID      int    `json:"seedListId"`
	CreateTime      Time   `json:"createTime"`
	Creator         string `json:"creator"`
	ExternalID      string `json:"externalID"`
	Custom1         string `json:"custom1"`
	Channel         string `json:"channel"`
	CountRecips     bool   `json:"countRecips"`
	BrandIDs        []int  `json:"brandIds"`
	ListCount       int    `json:"listCount"`
	TestMessageList bool   `json:"testMessageList"`
}

func (l *List) getID(ctx context.Context, pu *PostUp) (int, error) {
	if l.ID != 0 {
		return l.ID, nil
	}

	var eid = l.ExternalID
	if eid == "" {
		return 0, ErrListMissingIdentifier
	}

	var err error
	if l, err = pu.GetListByExternalID(ctx, eid); err != nil {
		return 0, err
	}

	return l.ID, nil
}

// related client methods

func (pu *PostUp) GetLists(ctx context.Context) ([]*List, error) {
	var reqURL = pu.url("list", nil)

	req, err := pu.newRequest(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := pu.do(req)
	if err != nil {
		return nil, fmt.Errorf("encountered network error while updating user with postup: %w", err)
	}

	var lists []*List
	if err := pu.decodeJSON(resp, &lists); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	return lists, err
}

func (pu *PostUp) GetListsByBrandID(ctx context.Context, id int) ([]*List, error) {
	var reqURL = pu.url("list", url.Values{"brandId": []string{fmt.Sprintf("%d", id)}})

	req, err := pu.newRequest(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := pu.do(req)
	if err != nil {
		return nil, fmt.Errorf("encountered network error while updating user with postup: %w", err)
	}

	var lists []*List
	if err := pu.decodeJSON(resp, &lists); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	return lists, err
}

func (pu *PostUp) GetListByTitle(ctx context.Context, title string) (*List, error) {
	lists, err := pu.GetLists(ctx)
	if err != nil {
		return nil, err
	}

	for _, l := range lists {
		if l.Title == title {
			return l, nil
		}
	}

	return nil, ErrListNotFound
}

func (pu *PostUp) GetListByExternalID(ctx context.Context, id string) (*List, error) {
	lists, err := pu.GetLists(ctx)
	if err != nil {
		return nil, err
	}

	for _, l := range lists {
		if l.ExternalID == id {
			return l, nil
		}
	}

	return nil, ErrListNotFound
}

func (pu *PostUp) GetListByListID(ctx context.Context, id int) (*List, error) {
	lists, err := pu.GetLists(ctx)
	if err != nil {
		return nil, err
	}

	for _, l := range lists {
		if l.ID == id {
			return l, nil
		}
	}

	return nil, ErrListNotFound
}
