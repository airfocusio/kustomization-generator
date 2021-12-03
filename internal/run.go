package internal

func Run(dir string, config Generator) error {
	return config.Generate(dir)
}
