package types

import "regexp"

// StringSlice is a string slice.
type StringSlice []string

// Contains determines if the string slice has a given value.
func (o StringSlice) Contains(value string) bool {
	for _, s := range o {
		if s == value {
			return true
		}
	}
	return false
}

// Filter filters down to items that return true from the condition function.
func (o StringSlice) Filter(cond func(value string) bool) StringSlice {
	ss := []string{}
	for _, s := range o {
		if cond(s) {
			ss = append(ss, s)
		}
	}
	return StringSlice(ss)
}

// Filter filters down to items that match the regular expression.
func (o StringSlice) FilterRegexp(re *regexp.Regexp) StringSlice {
	return o.Filter(func(value string) bool {
		return re.MatchString(value)
	})
}
