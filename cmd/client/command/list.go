package command

import (
	"context"
	"fmt"

	"google.golang.org/grpc/status"
)

// List client command for receiving items list.
func (c *Command) List(ctx context.Context) error {
	token, _, err := readCredentials()
	if err != nil {
		return err
	}
	items, err := c.client.List(ctx, token)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			fmt.Printf("Get items list error: %s", s.Message())
			return nil
		}
		return err
	}

	if len(items) == 0 {
		fmt.Println("Your items list is empty now")
		return nil
	}

	fmt.Printf("Your items [%d]:\n", len(items))
	for _, item := range items {
		fmt.Printf("\t%s\n", item)
	}
	return nil
}
