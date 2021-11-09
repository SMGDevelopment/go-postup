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
	if len(os.Args) < 3 {
		log.Fatalf("usage: %s PATH ADDRESS", os.Args[0])
	}

	var (
		ctx = context.Background()

		address = os.Args[2]

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

	recipient, err := pu.GetRecipientByAddress(ctx, address)
	if err != nil {
		log.Fatal(err)
	}

	writer.Encode(recipient)

	lists, err := pu.GetLists(ctx)
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write([]byte("\n- - - -\n\n"))

	writer.Encode(lists)
}
