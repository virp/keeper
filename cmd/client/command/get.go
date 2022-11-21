package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/grpc/status"
)

// Get client command for receiving item details.
func (c *Command) Get(ctx context.Context, name string) error {
	token, secret, err := readCredentials()
	if err != nil {
		return err
	}
	item, err := c.client.Get(ctx, token, secret, name)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			fmt.Printf("Get item error: %s", s.Message())
			return nil
		}
		return err
	}

	fmt.Println("Item details")
	fmt.Printf("Name: %s\n", item.Name)
	fmt.Printf("Type: %s\n", item.Type)
	if len(item.Metadata) == 0 {
		fmt.Println("Metadata: empty")
	} else {
		fmt.Println("Metadata:")
		for _, m := range item.Metadata {
			fmt.Printf("\t%s: %s\n", m.Key, m.Value)
		}
	}
	switch item.Type {
	case "text":
		fmt.Printf("Text: %s\n", string(item.Data))
	case "password":
		var pass typePassword
		if err := json.Unmarshal(item.Data, &pass); err != nil {
			return err
		}
		fmt.Printf("Login: %s\n", pass.Login)
		fmt.Printf("Password: %s\n", pass.Password)
	case "card":
		var card typeCard
		if err := json.Unmarshal(item.Data, &card); err != nil {
			return err
		}
		fmt.Printf("Number: %s\n", card.Number)
		fmt.Printf("Holder: %s\n", card.Holder)
		fmt.Printf("Expired: %s/%s\n", card.ExpiryMonth, card.ExpiryYear)
		fmt.Printf("CVV: %s\n", card.Cvv)
	case "binary":
		f, err := os.Create(item.Name)
		if err != nil {
			return err
		}
		if _, err := f.Write(item.Data); err != nil {
			return err
		}
		//goland:noinspection GoUnhandledErrorResult
		f.Close()
		fmt.Printf("File %s saved on disk\n", item.Name)
	}

	return nil
}
