package indent

// Indent inserts prefix before all non-empty line
func Indent(s string, prefix string) string {
	return string(BytesIndent([]byte(s), []byte(prefix)))
}

// BytesIndent inserts prefix before all lines
func BytesIndent(s []byte, prefix []byte) []byte {
	var res []byte
	first := true
	for _, c := range s {
		if first && c != '\n' {
			res = append(res, prefix...)
		}
		res = append(res, c)
		first = c == '\n'
	}
	return res
}
