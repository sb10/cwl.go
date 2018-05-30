package cwl

import "sort"

// StringArrayable converts "xxx" to ["xxx"] if it's not slice.
func StringArrayable(i interface{}) []string {
	dest := []string{}
	switch x := i.(type) {
	case []interface{}:
		for _, s := range x {
			dest = append(dest, s.(string))
		}
	case string:
		dest = append(dest, x)
	}
	return dest
}

// sortKeys sorts the keys of a given map
func sortKeys(x map[string]interface{}) []string {
	keys := make([]string, len(x))
	i := 0
	for key := range x {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}
