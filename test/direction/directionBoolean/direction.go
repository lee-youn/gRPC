package directionboolean

// Convert a slice to a map for fast lookup
func SliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool)
	for _, v := range slice {
		m[v] = true
	}
	return m
}

// 첫글자가 대문자여야만 접근 가능하다..
func DirectionBoolean(key string, value string) bool {
	// Define arrays for each direction
	directions := map[string][]string{
		"Rs": {"Rs", "Ls", "Ur", "Dr"},
		"Rl": {"Rl", "Ll", "Dr", "Us"},
		"Rr": {"Rr", "Us", "Ds", "Lr", "Ls"},
		"Ls": {"Ls", "Rs", "Ur", "Dr"},
		"Ll": {"Ll", "Rl", "Dr", "Us"},
		"Lr": {"Lr", "Ds", "Us", "Rr", "Rs"},
		"Ds": {"Ds", "Us", "Rr", "Lr"},
		"Dl": {"Dl", "Ul", "Rr", "Ls"},
		"Dr": {"Dr", "Rs", "Ls", "Ur", "Us"},
		"Us": {"Us", "Ds", "Rr", "Lr"},
		"Ul": {"Ul", "Dl", "Ls", "Rr"},
		"Ur": {"Ur", "Ls", "Ds", "Rs", "Rr"},
	}

	// Convert direction arrays to maps for fast lookup
	directionMaps := make(map[string]map[string]bool)
	for k, v := range directions {
		directionMaps[k] = SliceToMap(v)
	}

	if directionMaps["Rs"]["Ls"] {
		return true
	} else {
		return false
	}
}
