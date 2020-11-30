package strip_test

import (
	"flag"
	"strings"
	"testing"

	"github.com/tajud99n/strip"
)

var (
	apiKey string
)

func init() {
	flag.StringVar(&apiKey, "key", "", "Your TEST secret key for the Stripe API. If present, integration tests will be run using this key.")
}

func TestClient_Customer(t *testing.T) {
	if apiKey == "" {
		t.Skip("No API key provided")
	}

	c := strip.Client{
		Key: apiKey,
	}
	tok := "tok_amex"
	email := "test@go.com"
	cus, err := c.Customer(tok, email)
	if err != nil {
		t.Errorf("Customer() err = %v; want %v", err, nil)
	}
	if cus == nil {
		t.Fatalf("Customer() = nil; want non-nil value")
	}
	if !strings.HasPrefix(cus.ID, "cus_") {
		t.Errorf("Customer() ID = %s; want prefix %q", cus.ID, "cus_")
	}
	if !strings.HasPrefix(cus.DefaultSource, "card_") {
		t.Errorf("Customer() DefaultSource = %s; want prefix %q", cus.DefaultSource, "card_")
	}
	if cus.Email != email {
		t.Errorf("Customer() Email = %s; want %s", cus.Email, email)
	}
}

func TestClient_Charge(t *testing.T) {
	if apiKey == "" {
		t.Skip("No API key provided")
	}

	c := strip.Client{
		Key: apiKey,
	}
	tok := "tok_amex"
	email := "test@go.com"
	cus, err := c.Customer(tok, email)
	if err != nil {
		t.Fatalf("Customer() err = %v; want %v", err, nil)
	}
	amount := 1234
	charge, err := c.Charge(cus.ID, amount)
	if err != nil {
		t.Errorf("Charge() err = %v; want %v", err, nil)
	}
	if charge == nil {
		t.Fatalf("Charge() = nil; want non-nil value")
	}
	if charge.Amount != amount {
		t.Fatalf("Charge() Amount = %d; want %d", charge.Amount, amount)
	}
}