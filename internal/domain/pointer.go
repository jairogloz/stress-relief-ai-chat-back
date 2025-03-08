package domain

func IntPtr(i int) *int {
	return &i
}

func StrPtr(s string) *string {
	return &s
}
