package domain

type Service interface {
	Name() string

	Init(ctx Context) error

}