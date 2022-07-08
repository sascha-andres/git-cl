package internal

// WithVersion allows to set a version  for the changelog
func WithVersion(version string) ChangeLogGeneratorOption {
	return func(generator *ChangeLogGenerator) {
		generator.version = version
	}
}
