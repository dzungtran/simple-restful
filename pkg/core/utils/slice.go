package utils

func IsStringSliceContains(strSlice []string, str string) bool {
	for _, s := range strSlice {
		if s == str {
			return true
		}
	}
	return false
}

func IsInt64SliceContains(intSlice []int64, num int64) bool {
	for _, s := range intSlice {
		if s == num {
			return true
		}
	}
	return false
}
