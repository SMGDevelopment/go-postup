package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/cheddartv/go-postup"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("usage: %s PATH ACTION ADDRESS [ARGS...]", os.Args[0])
	}

	var (
		ctx = context.Background()

		action  = os.Args[2]
		address = os.Args[3]
		args    = os.Args[4:]

		username string
		password string
	)

	{
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &struct {
			U *string `json:"username"`
			P *string `json:"password"`
		}{
			U: &username,
			P: &password,
		})

		if err != nil {
			log.Fatal(err)
		}
	}

	var (
		pu     = postup.NewPostUp(username, password, nil)
		writer = json.NewEncoder(os.Stdout)
	)
	writer.SetIndent("", "\t")

	var err error
	switch action {
	case "create":
		err = createRecipient(ctx, address, pu, writer)
	case "subscribe":
		err = subscribeRecipient(ctx, address, args[0], pu, writer)
	case "delete":
		err = deleteRecipient(ctx, address, pu, writer)
	case "get-list":
		err = getList(ctx, address, pu, writer)
	case "get-lists":
		err = getLists(ctx, pu, writer)
	case "get":
		fallthrough
	default:
		err = getRecipient(ctx, address, pu, writer)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func getRecipient(ctx context.Context, addr string, p *postup.PostUp, w *json.Encoder) error {
	recipient, err := p.GetRecipientByAddress(ctx, addr)
	if err != nil {
		return err
	}

	return w.Encode(recipient)
}

func createRecipient(ctx context.Context, addr string, p *postup.PostUp, w *json.Encoder) error {
	res, err := p.RecipientCreate(ctx, &postup.CreateRecipientRequest{
		ExternalID: addr,
		Address:    addr,
		Channel:    "E",
	})

	if err != nil {
		return err
	}

	return w.Encode(res)
}

func subscribeRecipient(ctx context.Context, addr, mailingList string, p *postup.PostUp, w *json.Encoder) error {
	list, err := p.GetListByTitle(ctx, mailingList)
	if err != nil {
		return err
	}

	res, err := p.SubscribeRecipientToList(ctx, &postup.SubscribeRecipientRequest{
		Recipient: &postup.Recipient{
			Address: addr,
		},
		List: list,
	})

	if err != nil {
		return err
	}

	return w.Encode(res)
}

func deleteRecipient(ctx context.Context, addr string, p *postup.PostUp, w *json.Encoder) error {
	res, err := p.DeleteRecipientByAddress(ctx, addr)
	if err != nil {
		return err
	}

	return w.Encode(res)
}

func getLists(ctx context.Context, p *postup.PostUp, w *json.Encoder) error {
	res, err := p.GetLists(ctx)
	if err != nil {
		return err
	}

	return w.Encode(res)
}

func getList(ctx context.Context, addr string, p *postup.PostUp, w *json.Encoder) error {
	res, err := p.GetListByTitle(ctx, addr)
	if err != nil {
		return err
	}

	return w.Encode(res)
}
