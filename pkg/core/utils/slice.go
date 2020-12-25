package utils

func IsStringSliceContains(strSlice []string, str string) bool {
	for _, s := range strSlice {
		if s == str{
			return true
		}
	}
	return false
}