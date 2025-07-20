package goround

func WithStructInject(enable bool) option {
	return func(config serviceConfig) serviceConfig {
		config.enableStructInject = enable
		return config
	}
}

func WithStructTag(tag string) option {
	return func(config serviceConfig) serviceConfig {
		config.structTag = tag
		return config
	}
}

func WithFunctionInject(enable bool) option {
	return func(config serviceConfig) serviceConfig {
		config.enableFunctionInject = enable
		return config
	}
}
func WithFunctionInjectName(name string) option {
	return func(config serviceConfig) serviceConfig {
		config.functionInjectName = name
		return config
	}
}

func WithStructDefaultInject(enable bool) option {
	return func(config serviceConfig) serviceConfig {
		config.structDefaultInject = enable
		return config
	}
}
