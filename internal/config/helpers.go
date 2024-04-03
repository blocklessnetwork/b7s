package config

// By default we have flags in a tree structure, but most of the time we just want the leaves.
func flattenConfigOptions(cfgOptions []ConfigOption) []ConfigOption {
	out := make([]ConfigOption, 0, len(cfgOptions))
	for _, cfg := range cfgOptions {
		expandConfigOption(cfg, &out)
	}
	return out
}

func expandConfigOption(cfg ConfigOption, out *[]ConfigOption) {
	if len(cfg.Children) == 0 {
		*out = append(*out, cfg)
		return
	}
	for _, fc := range cfg.Children {
		expandConfigOption(fc, out)
	}
}

func flattenMap(prefix string, in map[string]any, flat map[string]any) {
	for k, v := range in {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch cv := v.(type) {
		default:
			flat[key] = v

		case map[string]any:
			flattenMap(key, cv, flat)
		}
	}
}

func fullPath(name string, parents ...string) []string {
	var full []string
	full = append(full, parents...)
	full = append(full, name)
	return full
}
