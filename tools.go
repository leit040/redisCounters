package redisCounters

func inArray(key string, array []string) bool {
	for _, value := range array {
		if key == value {
			return true
		}
	}
	return false
}
