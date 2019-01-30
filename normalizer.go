package zenbun

type Normalizer interface {
	Normalize(*string) string
}

type NormalizeFunc func(*string) string

func (f NormalizeFunc) Normalize(document *string) string {
	return f(document)
}
