package utils

func DerefString(v *string, d string) string {
	if v == nil {
		return d
	}
	return *v
}
