package cwl

// Namespaces ...
type Namespaces []Namespace

// New constructs "Namespaces" struct.
func (n Namespaces) New(i interface{}) Namespaces {
	dest := []Namespace{}
	switch x := i.(type) {
	case []interface{}:
		for _, v := range x {
			dest = append(dest, Namespace{}.New(v))
		}
	case map[string]interface{}:
		for _, key := range sortKeys(x) {
			tmp := map[string]interface{}{}
			tmp[key] = x[key]
			dest = append(dest, Namespace{}.New(tmp))
		}
	default:
		dest = append(dest, Namespace{}.New(x))
	}
	return dest
}

// Namespace ...
type Namespace map[string]interface{}

// New constructs a Namespace from any interface.
func (n Namespace) New(i interface{}) Namespace {
	dest := Namespace{}
	switch x := i.(type) {
	case map[string]interface{}:
		for key, v := range x {
			dest[key] = v
		}
	}
	return dest
}
