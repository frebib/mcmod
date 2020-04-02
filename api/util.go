package api

func StringInSlice(haystack []string, needle string) bool {
	for _, str := range haystack {
		if str == needle {
			return true
		}
	}
	return false
}
