package payments

type MapProvider map[string]string

var EmptyProvider = MapProvider{}

func NewMapProvider(providers map[string]string) (out MapProvider) {
	out = make(map[string]string, len(providers))
	for k, v := range providers {
		out[k] = v
	}
	return
}

func (p MapProvider) PaymentToken(name string) string {
	return p[name]
}
