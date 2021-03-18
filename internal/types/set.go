package types

// StringSet is a set of strings.
type StringSet map[string]struct{}

// Insert inserts values into the set.
func (o StringSet) Insert(value ...string) StringSet {
	for _, v := range value {
		o[v] = struct{}{}
	}
	return o
}

// Remove removes values from the set.
func (o StringSet) Remove(value ...string) StringSet {
	for _, v := range value {
		delete(o, v)
	}
	return o
}

// Contains takes a value and returns true if the StringSet contains the value.
func (o StringSet) Contains(value string) bool {
	_, ok := o[value]
	return ok
}

// Slice converts the StringSet to a string slice.
func (o StringSet) Slice() []string {
	ss := make([]string, len(o))
	i := 0
	for s := range o {
		ss[i] = s
		i++
	}
	return ss
}

// StringSetFromSlice creates a StringSet from a string slice.
func StringSetFromSlice(ss []string) StringSet {
	o := StringSet{}
	for _, s := range ss {
		o[s] = struct{}{}
	}
	return o
}
