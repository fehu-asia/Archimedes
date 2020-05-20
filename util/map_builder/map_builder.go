package map_builder

func BuilderMap(kvs ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if len(kvs) >= 2 {
		for i := 0; i < len(kvs); i += 2 {
			m[kvs[i].(string)] = kvs[i+1]
		}
	}
	return m
}
