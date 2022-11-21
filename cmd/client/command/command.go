package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"

	"keeper/internal/services"
)

// ClientService interface describe requirements for client service
type ClientService interface {
	Register(ctx context.Context, login, password string) (string, error)
	Login(ctx context.Context, login, password string) (string, error)
	List(ctx context.Context, token string) ([]string, error)
	Get(ctx context.Context, token string, secret string, name string) (services.Item, error)
	Add(ctx context.Context, token string, secret string, item services.Item) error
	Delete(ctx context.Context, token string, name string) error
}

// Command implement logic for client commands.
type Command struct {
	log    *zap.SugaredLogger
	client ClientService
}

// NewCommand construct Command.
func NewCommand(log *zap.SugaredLogger, client ClientService) *Command {
	return &Command{
		log:    log,
		client: client,
	}
}

type credentials struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}

func saveCredentials(token string, secret string) error {
	cred := credentials{
		Token:  token,
		Secret: secret,
	}
	f, err := os.Create(".credentials")
	if err != nil {
		return fmt.Errorf("creating credentials file: %w", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()
	data, err := json.Marshal(cred)
	if err != nil {
		return fmt.Errorf("marshaling credentials: %w", err)
	}
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}
	return nil
}

func readCredentials() (string, string, error) {
	f, err := os.Open(".credentials")
	if err != nil {
		return "", "", fmt.Errorf("opening credentials file: %w", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return "", "", fmt.Errorf("reading credentials file: %w", err)
	}
	var cred credentials
	if err := json.Unmarshal(data, &cred); err != nil {
		return "", "", fmt.Errorf("unmarshaling credentials: %w", err)
	}
	return cred.Token, cred.Secret, nil
}

type typePassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type typeCard struct {
	Number      string `json:"number"`
	Holder      string `json:"holder"`
	ExpiryMonth string `json:"expiry_month"`
	ExpiryYear  string `json:"expiry_year"`
	Cvv         string `json:"cvv"`
}
