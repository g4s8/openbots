package spec

import (
	"errors"
)

// Invoice is a payment invoice.
type Invoice struct {
	// Provider is a payment provider name.
	Provider string `yaml:"provider"`
	// Title of the product
	Title string `yaml:"title"`
	// Description of the product
	Description string `yaml:"description"`
	// Unique bot deep-linking parameter that can be used to generate this invoice
	Payload string `yaml:"payload"`
	// Three-letter ISO 4217 currency code
	Currency string `yaml:"currency"`
	// Price breakdown, a list of components (e.g. product price, tax, discount, delivery cost, delivery tax, bonus, etc.)
	Prices []Price `yaml:"prices"`
}

// Validate invoice
func (i *Invoice) validate() (errs []error) {
	errs = make([]error, 0)
	if i.Title == "" {
		errs = append(errs, errors.New("empty invoice title"))
	}
	if i.Description == "" {
		errs = append(errs, errors.New("empty invoice description"))
	}
	if i.Payload == "" {
		errs = append(errs, errors.New("empty invoice payload"))
	}
	if i.Currency == "" {
		errs = append(errs, errors.New("empty invoice currency"))
	}
	if len(i.Prices) == 0 {
		errs = append(errs, errors.New("empty invoice prices"))
	}
	for _, p := range i.Prices {
		errs = append(errs, p.validate()...)
	}
	return
}

// Price of the product
type Price struct {
	// Label of the price
	Label string `yaml:"label"`
	// Price in the smallest units of the currency (integer, not float/double).
	Amount int `yaml:"amount"`
}

// Validate price
func (p *Price) validate() (errs []error) {
	errs = make([]error, 0)
	if p.Label == "" {
		errs = append(errs, errors.New("empty price label"))
	}
	if p.Amount <= 0 {
		errs = append(errs, errors.New("invalid price amount"))
	}
	return
}

type PreCheckoutTrigger struct {
	InvoicePayload string `yaml:"invoicePayload"`
}

type PreCheckoutAnswer struct {
	Ok           bool   `yaml:"ok"`
	ErrorMessage string `yaml:"errorMessage"`
}

type PostCheckoutTrigger struct {
	InvoicePayload string `yaml:"invoicePayload"`
}
