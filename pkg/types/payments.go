package types

type PaymentProviders interface {
	PaymentToken(name string) string
}
