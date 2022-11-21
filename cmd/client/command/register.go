package command

import (
	"context"
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc/status"
)

func (c *Command) Register(ctx context.Context) error {
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
	secret := string(b)
	if len(secret) < 4 {
		fmt.Println("Registration error: secret: length should be greater or equal 4")
	}

	token, err := c.client.Register(ctx, login, password)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			fmt.Printf("Registration error: %s\n", s.Message())
			return nil
		}
		return err
	}

	if err := saveCredentials(token, secret); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}

	return nil
}
