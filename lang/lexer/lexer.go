package lexer

import (
	"fmt"
)

func getDialect(source string) (Dialect, error) {
	for _, d := range dialects {
		matches := d.DetectionPattern().FindStringSubmatch(source)

		if len(matches) > 0 {
			return d, nil
		}
	}

	return nil, fmt.Errorf("you must define your dialect")
}

// GetDialectByName returns the dialect with the given name
func GetDialectByName(name string) (Dialect, bool) {
	d, ok := dialects[name]
	return d, ok
}
