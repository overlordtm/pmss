package utils

func Int64Ptr(i int64) *int64 {
	return &i
}

func Uint32Ptr(i uint32) *uint32 {
	return &i
}

func StringPtr(s string) *string {
	return &s
}

func ErrToStrPtr(err error) *string {
	if err == nil {
		return nil
	}
	s := err.Error()
	return &s
}
