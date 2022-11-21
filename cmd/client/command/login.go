package command

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc/status"
)

// Login client command for user login.
func (c *Command) Login(ctx context.Context) error {
	var login, password string
	fmt.Print("Login: ")
	if _, err := fmt.Scanln(&login); err != nil {
		return err
	}
	fmt.Print("Password: ")
	b, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}
	fmt.Println()
	password = string(b)
	fmt.Print("Secret for encryption: ")
	b, err = terminal.ReadPassword(0)
	if err != nil {
		return err
	}
	fmt.Println()
	hash := md5.Sum(b)
	secret := hex.EncodeToString(hash[:])

	token, err := c.client.Login(ctx, login, password)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			fmt.Printf("Login error: %s\n", s.Message())
			return nil
		}
		return err
	}

	if err := saveCredentials(token, secret); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}

	return nil
}
