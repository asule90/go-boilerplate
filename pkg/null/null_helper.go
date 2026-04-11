package null

import "strconv"

// StrToNilBool converts a string to *bool, returns nil if empty.
func StrToNilBool(s string) *bool {
	if s == "" {
		return nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &b
}

// StrToNilInt converts a string to *int, returns nil if empty.
func StrToNilInt(s string) *int {
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}

// PtrBool returns a pointer to a bool value.
func PtrBool(b bool) *bool { return &b }

// PtrString returns a pointer to a string value.
func PtrString(s string) *string { return &s }

// PtrInt returns a pointer to an int value.
func PtrInt(i int) *int { return &i }

// PtrFloat64 returns a pointer to a float64 value.
func PtrFloat64(f float64) *float64 { return &f }

// NilIfEmpty returns nil if the string is empty, otherwise returns a pointer to the string.
func NilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
