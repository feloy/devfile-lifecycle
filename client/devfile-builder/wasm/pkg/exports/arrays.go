package exports

import "syscall/js"

func getStringArray(value js.Value) []string {
	l := value.Length()
	result := make([]string, 0, l)
	for i := 0; i < l; i++ {
		s := value.Index(i).String()
		if len(s) > 0 {
			result = append(result, s)
		}
	}
	return result
}
