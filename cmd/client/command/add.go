package command

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"keeper/internal/services"
)

var (
	allowedTypes = map[string]string{
		"password": "password",
		"text":     "text",
		"card":     "card",
		"binary":   "binary",
	}
)

// Add client command for creating and storing new item.
func (c *Command) Add(ctx context.Context, name string, itemType string) error {
	token, secret, err := readCredentials()
	if err != nil {
		return err
	}

	selectedType, ok := allowedTypes[itemType]
	if !ok {
		fmt.Println("Not allowed item type selected")
		return nil
	}
	fmt.Printf("Creating new item \"%s\" of type \"%s\":\n", name, selectedType)

	var data []byte
	switch selectedType {
	case "password":
		in := bufio.NewReader(os.Stdin)
		fmt.Print("Login: ")
		login, err := in.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Print("Password: ")
		password, err := in.ReadString('\n')
		if err != nil {
			return err
		}
		pt := typePassword{
			Login:    strings.Trim(login, "\n"),
			Password: strings.Trim(password, "\n"),
		}
		data, err = json.Marshal(&pt)
		if err != nil {
			return err
		}
	case "text":
		fmt.Println("Enter text for item: ")
		snr := bufio.NewScanner(os.Stdin)
		var text string
		for snr.Scan() {
			line := snr.Text()
			if len(line) == 0 {
				break
			}
			text = text + line + "\n"
		}
		data = []byte(strings.Trim(text, "\n"))
	case "card":
		in := bufio.NewReader(os.Stdin)
		fmt.Print("Number: ")
		number, err := in.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Print("Holder: ")
		holder, err := in.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Print("Expiry Month: ")
		expiryMonth, err := in.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Print("Expiry Year: ")
		expiryYear, err := in.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Print("CVV: ")
		cvv, err := in.ReadString('\n')
		if err != nil {
			return err
		}
		ct := typeCard{
			Number:      strings.Trim(number, "\n"),
			Holder:      strings.Trim(holder, "\n"),
			ExpiryMonth: strings.Trim(expiryMonth, "\n"),
			ExpiryYear:  strings.Trim(expiryYear, "\n"),
			Cvv:         strings.Trim(cvv, "\n"),
		}
		data, err = json.Marshal(&ct)
		if err != nil {
			return err
		}
	case "binary":
		fmt.Print("File path: ")
		var filePath string
		_, err := fmt.Scanln(&filePath)
		if err != nil {
			return err
		}
		f, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("File open error: %s\n", err.Error())
			return err
		}
		data, err = io.ReadAll(f)
		if err != nil {
			fmt.Printf("File read error: %s\n", err.Error())
		}
	}

	item := services.Item{
		Name: name,
		Type: selectedType,
		Data: data,
	}

	fmt.Println("Adding metadata (key: value pairs) for item (blank key for end):")
	for {
		var key string
		fmt.Print("Key: ")
		if _, err := fmt.Scanln(&key); err != nil || len(key) == 0 {
			break
		}
		fmt.Print("Value: ")
		in := bufio.NewReader(os.Stdin)
		value, err := in.ReadString('\n')
		if err != nil {
			return err
		}
		m := services.Metadata{
			Key:   key,
			Value: strings.Trim(value, "\n"),
		}
		item.Metadata = append(item.Metadata, m)
	}

	err = c.client.Add(ctx, token, secret, item)
	if err != nil {
		return err
	}

	return nil
}
