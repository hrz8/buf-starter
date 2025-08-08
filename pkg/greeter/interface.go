package greeter

type Repositor interface {
	GetGreeterTemplate(name string) string
	GetAllowedNames(page, limit int32) []string
	GetAllowedNamesWithTotal(page, limit int32) ([]string, int32)
	GetTotalAllowedNames() int32
}
