package internal

// WithVersion allows to set a version  for the changelog
func WithVersion(version string) ChangeLogGeneratorOption {
	return func(generator *ChangeLogGenerator) {
		generator.version = version
	}
}

// WithPrintConfiguration will print out currently used configuration
func WithPrintConfiguration() ChangeLogGeneratorOption {
	return func(generator *ChangeLogGenerator) {
		generator.printConfiguration = true
	}
}

// WithConfiguration can be used to set the configuration from a file
func WithConfiguration(configuration *ChangeLogGenerator) ChangeLogGeneratorOption {
	return func(generator *ChangeLogGenerator) {
		generator.configuration = configuration
	}
}
