package strings

import (
	"fmt"
	"strings"
)

func Metric(chunks ...any) string {
	chunkStrings := []string{}
	for _, chunk := range chunks {
		chunkStrings = append(chunkStrings, toString(chunk))
	}
	return strings.Join(chunkStrings, `.`)
}

func toString(a any) string {
	stringer, ok := a.(fmt.Stringer)
	if ok {
		return stringer.String()
	}
	return fmt.Sprintf(`%v`, a)
}
