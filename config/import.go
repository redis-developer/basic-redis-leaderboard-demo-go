package config

const (
	defaultImportPath = "./seed.json"

	envImportPath = "IMPORT_PATH"
)

type Import struct {
	path string
}

func (i Import) Path() string {
	return i.path
}
