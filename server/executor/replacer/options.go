package replacer

type InjectorOption func(*Injector)

func WithVariableProvider(provider VariableProvider) InjectorOption {
	return func(i *Injector) {
		i.provider = provider
	}
}
