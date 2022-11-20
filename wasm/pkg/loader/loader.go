package loader

type Scene struct {
	ClearColor []float32        `yaml:"clearColor"`
	Layers     map[string]Layer `yaml:"layers"`
}

type Layer struct {
}
