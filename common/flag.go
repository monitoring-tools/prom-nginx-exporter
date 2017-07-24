package common

import "strings"

// ArrFlags is custom flat struct for ability to pass array of values(strings)(example: --flag_name="value1" --flag_name="value2")
type ArrFlags []string

// String returns joined values through the separator ";"
func (f *ArrFlags) String() string {
	return strings.Join(*f, ";")
}

// Set sets new passed value
func (f *ArrFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

// AsArray returns values as array
func (f *ArrFlags) AsArray() []string {
	return *f
}
