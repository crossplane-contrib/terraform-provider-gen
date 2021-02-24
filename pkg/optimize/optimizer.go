package optimize

import "github.com/crossplane-contrib/terraform-provider-gen/pkg/generator"

type Optimizer func(*generator.ManagedResource) (*generator.ManagedResource, error)

func NewOptimizerChain(optimizers ...Optimizer) Optimizer {
	return func(mr *generator.ManagedResource) (*generator.ManagedResource, error) {
		var err error
		for _, o := range optimizers {
			mr, err = o(mr)
			if err != nil {
				return nil, err
			}
		}
		return mr, err
	}
}
