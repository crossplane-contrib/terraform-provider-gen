package generator

import "fmt"

type ResourceNamer interface {
	TypeName() string
	TypeListName() string
	SpecTypeName() string
	StatusTypeName() string
	AtProviderTypeName() string
	ForProviderTypeName() string
}

type defaultNamer struct {
	resourceName string
}

func (n defaultNamer) TypeName() string {
	return n.resourceName
}

func (n defaultNamer) TypeListName() string {
	return fmt.Sprintf("%sList", n.TypeName())
}

func (n defaultNamer) SpecTypeName() string {
	return fmt.Sprintf("%sSpec", n.TypeName())
}

func (n defaultNamer) StatusTypeName() string {
	return fmt.Sprintf("%sStatus", n.TypeName())
}

func (n defaultNamer) AtProviderTypeName() string {
	return fmt.Sprintf("%sObservation", n.TypeName())
}

func (n defaultNamer) ForProviderTypeName() string {
	return fmt.Sprintf("%sParameters", n.TypeName())
}

func NewDefaultNamer(resourceName string) ResourceNamer {
	return defaultNamer{resourceName: resourceName}
}
