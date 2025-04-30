package template

import (
	"os"
	"path/filepath"
)

type Loader struct {
	defaultTemplateContent []byte
}

func NewTemplateLoader(defaultTemplate []byte) *Loader {
	return &Loader{
		defaultTemplateContent: defaultTemplate,
	}
}

func (t *Loader) Load(path string) ([]byte, error) {
	if path == "" {
		return t.defaultTemplateContent, nil
	}

	return os.ReadFile(filepath.Clean(path))
}
