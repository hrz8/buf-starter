package greeter

type Repositor interface {
	GetGreeterTemplate(name string) string
}
