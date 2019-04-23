package MyGRPCLib

func ConvFloat64ToFloat64Pointer(f float64) *float64 {
	if f == 0.0 {
		return nil
	}
	return &f
}

func ConvStringPointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ConvStringToStringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ConvFloat64PointerToFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func ConvInt64PointerToInt64(f *int64) int64 {
	if f == nil {
		return 0
	}
	return *f
}

func ConvInt64ToInt64Pointer(f int64) *int64 {
	if f == 0 {
		return nil
	}
	return &f
}
