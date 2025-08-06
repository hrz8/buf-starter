package greeter

import "fmt"

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) GetGreeterTemplate(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}
