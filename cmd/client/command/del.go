package command

import (
	"context"
	"fmt"
)

// Del client command for removing item.
func (c *Command) Del(ctx context.Context, name string) error {
	token, _, err := readCredentials()
	err = c.client.Delete(ctx, token, name)
	if err != nil {
		return err
	}
	fmt.Printf("Item %s sucessfully deleted\n", name)
	return nil
}
