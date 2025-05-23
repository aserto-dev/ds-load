package kongyaml

import (
	"io"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

type YAMLResolver struct {
	yamlKey string
}

func NewYAMLResolver(yamlKey string) *YAMLResolver {
	return &YAMLResolver{
		yamlKey: yamlKey,
	}
}

// Loader is a Kong configuration loader for YAML.
//
//nolint:ireturn // loader returns a kong resolver interface
func (y *YAMLResolver) Loader(r io.Reader) (kong.Resolver, error) {
	decoder := yaml.NewDecoder(r)
	config := map[string]any{}

	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	if y.yamlKey != "" {
		var ok bool
		config, ok = config[y.yamlKey].(map[string]any)

		if !ok {
			return kong.ResolverFunc(func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
				return nil, nil
			}), nil
		}
	}

	return kong.ResolverFunc(func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
		// Build a string path up to this flag.
		path := []string{}
		path = append(path, flag.Name)
		path = strings.Split(strings.Join(path, "-"), "-")
		s := find(config, path)

		if s == nil {
			fullPath := []string{}
			for n := parent.Node(); n != nil && n.Type != kong.ApplicationNode; n = n.Parent {
				fullPath = append([]string{n.Name}, fullPath...)
			}

			fullPath = append(fullPath, path...)
			s = find(config, fullPath)
		}

		return s, nil
	}), nil
}

func find(config map[string]any, path []string) any {
	for i := range path {
		prefix := strings.Join(path[:i+1], "-")
		if child, ok := config[prefix].(map[string]any); ok {
			return find(child, path[i+1:])
		}
	}

	return config[strings.Join(path, "-")]
}
