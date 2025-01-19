package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/blessnetwork/b7s/models/bls"
)

type ConfigOption struct {
	Name     string         `yaml:"name,omitempty"`
	FullPath string         `yaml:"full_path,omitempty"`
	CLI      CLIFlag        `yaml:"cli,omitempty"`
	Env      string         `yaml:"env-var,omitempty"`
	Children []ConfigOption `yaml:"children,omitempty"`

	kind        reflect.Kind
	elementKind reflect.Kind // if kind is a slice, tell us what the elements of the slice are.
}

func (c ConfigOption) Type() string {

	// For slices say something like "list (string)", for primitive types print the type, and skip structs.
	if c.kind == reflect.Slice {
		return fmt.Sprintf("list (%s)", c.elementKind.String())
	}

	if c.kind != reflect.Struct {
		return c.kind.String()
	}

	return ""
}

func (c ConfigOption) Info() ConfigOptionInfo {

	info := ConfigOptionInfo{
		FullPath: c.FullPath,
		Name:     c.Name,
		CLI:      c.CLI,
		Env:      c.Env,
		Children: c.Children,
		Type:     c.Type(),
	}

	return info
}

type CLIFlag struct {
	Flag        string `yaml:"flag,omitempty"`
	Shorthand   string `yaml:"shorthand,omitempty"`
	Default     any    `yaml:"default,omitempty"`
	Description string `yaml:"description,omitempty"`
}

func getConfigOptions() []ConfigOption {
	cliDefaults := getDefaultFlagValues()
	return getStructInfo(reflect.TypeOf(Config{}), cliDefaults)
}

func getStructInfo(typ reflect.Type, cliDefaults map[string]any, parents ...string) []ConfigOption {

	out := make([]ConfigOption, 0)
	for _, field := range reflect.VisibleFields(typ) {

		var (
			kind     = field.Type.Kind()
			koanfTag = field.Tag.Get("koanf")
			parts    = fullPath(koanfTag, parents...)
			fullPath = strings.Join(parts, ".")
		)

		fi := ConfigOption{
			FullPath: fullPath,
			kind:     kind,
			Name:     koanfTag,
			// Env variable is set later, after we determine the type
		}

		ft := field.Tag.Get("flag")
		if ft != "" {
			flag, shorthand := getFlagFromTag(ft)

			cli := CLIFlag{
				Flag:        flag,
				Shorthand:   shorthand,
				Default:     cliDefaults[fullPath],
				Description: getFlagDescription(flag),
			}

			fi.CLI = cli
		}

		switch kind {

		case reflect.Struct:
			children := getStructInfo(field.Type, cliDefaults, parts...)
			fi.Children = children

		case reflect.Slice, reflect.Array:
			fi.elementKind = field.Type.Elem().Kind()
			fi.Env = envName(koanfTag, parents...)

		default:
			fi.Env = envName(koanfTag, parents...)
		}

		out = append(out, fi)
	}

	return out
}

func envName(name string, parents ...string) string {

	parts := make([]string, 0)
	for i := len(parents) - 1; i >= 0; i-- {
		title := strings.Title(parents[i])
		parts = append(parts, title)
	}

	nameFields := strings.Split(name, "-")
	var formattedName string
	for _, field := range nameFields {
		titled := strings.Title(field)
		formattedName += titled
	}

	var components []string
	components = append(components, strings.TrimSuffix(bls.EnvPrefix, EnvDelimiter)) // Trim trailing underscore so we don't repeat it.
	components = append(components, parts...)
	components = append(components, formattedName)

	return strings.Join(components, EnvDelimiter)
}
