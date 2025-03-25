package tools

func SliceInSlice(big, small []string) bool {
	for _, v := range small {
		if Contains(big, v) {
			return true
		}
	}
	return false
}

func Contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
