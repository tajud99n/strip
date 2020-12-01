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

const (
	tokenAmex        = "tok_amex"
	tokenInvalid     = "tok_invalid"
	tokenExpiredCard = "tok_chargeDeclinedExpiredCard"
	email            = "test@go.com"
)

func init() {
	flag.StringVar(&apiKey, "key", "", "Your TEST secret key for the Stripe API. If present, integration tests will be run using this key.")
}

func TestClient_Customer(t *testing.T) {
	if apiKey == "" {
		t.Skip("No API key provided")
	}
	type checkFn func(*testing.T, *strip.Customer, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoErr := func() checkFn {
		return func(t *testing.T, charge *strip.Customer, err error) {
			if err != nil {
				t.Fatalf(" err = %v; want nil", err)
			}
		}
	}

	hasErrType := func(s string) checkFn {
		return func(t *testing.T, charge *strip.Customer, err error) {
			se, ok := err.(strip.Error)
			if !ok {
				t.Fatalf("err isn't a strip.Error")
			}
			if se.Type != s {
				t.Fatalf("err.Type = %s; want %s", se.Type, s)
			}
		}
	}

	hasIDPrefix := func() checkFn {
		return func(t *testing.T, c *strip.Customer, e error) {
			if !strings.HasPrefix(c.ID, "cus_") {
				t.Errorf("ID = %s; want prefix %q", c.ID, "cus_")
			}
		}
	}

	hasCardDefaultSource := func() checkFn {
		return func(t *testing.T, c *strip.Customer, e error) {

			if !strings.HasPrefix(c.DefaultSource, "card_") {
				t.Errorf("Customer() DefaultSource = %s; want prefix %q", c.DefaultSource, "card_")
			}
		}
	}

	hasEmail := func (e string) checkFn {
		return func(t *testing.T, c *strip.Customer, err error) {
			if c.Email != e {
				t.Errorf("Email = %s; want %s", c.Email, e)
			}
		}
	}
	c := strip.Client{
		Key: apiKey,
	}
	tests := map[string]struct {
		token string
		email string
		checks []checkFn
	}{
		"valid customer with amex": {
			token: tokenAmex,
			email: email,
			checks: check(hasNoErr(), hasIDPrefix(), hasCardDefaultSource(), hasEmail(email)),
		},
		"invalid token": {
			token: tokenInvalid,
			email: email,
			checks: check(hasErrType(strip.ErrTypeInvalidRequest)),
		},
		"expired card": {
			token: tokenExpiredCard,
			email: email,
			checks: check(hasErrType(strip.ErrTypeCardError)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cus, err := c.Customer(tc.token, tc.email)
			for _, check := range tc.checks {
				check(t, cus, err)
			}
		})

	}
}

func TestClient_Charge(t *testing.T) {
	if apiKey == "" {
		t.Skip("No API key provided")
	}

	type checkFn func(*testing.T, *strip.Charge, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoErr := func() checkFn {
		return func(t *testing.T, charge *strip.Charge, err error) {
			if err != nil {
				t.Fatalf(" err = %v; want nil", err)
			}
		}
	}
	hasAmount := func(amount int) checkFn {
		return func(t *testing.T, charge *strip.Charge, err error) {
			if charge.Amount != amount {
				t.Errorf("Amount = %d; want %d", charge.Amount, amount)
			}
		}
	}
	hasErrType := func(s string) checkFn {
		return func(t *testing.T, charge *strip.Charge, err error) {
			se, ok := err.(strip.Error)
			if !ok {
				t.Fatalf("err isn't a strip.Error")
			}
			if se.Type != s {
				t.Fatalf("err.Type = %s; want %s", se.Type, s)
			}
		}
	}

	c := strip.Client{
		Key: apiKey,
	}

	cus, err := c.Customer(tokenAmex, email)
	if err != nil {
		t.Fatalf("Customer() err = %v; want %v", err, nil)
	}

	tests := map[string]struct {
		customerID string
		amount     int
		checks     []checkFn
	}{
		"valid charge": {
			customerID: cus.ID,
			amount:     1234,
			checks:     check(hasNoErr(), hasAmount(1234)),
		},
		"invalid customer id": {
			customerID: "cus_missing",
			amount:     1234,
			checks:     check(hasErrType(strip.ErrTypeInvalidRequest)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			charge, err := c.Charge(tc.customerID, tc.amount)
			for _, check := range tc.checks {
				check(t, charge, err)
			}
		})
	}
}
