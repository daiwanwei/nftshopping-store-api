package security

type Authentication interface {
	GetAuthorities() []string
	GetName() string
}
