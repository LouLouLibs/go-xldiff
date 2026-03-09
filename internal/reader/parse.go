package reader

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseFileArg splits "file.xlsx:SheetName" into file path and sheet identifier.
// If no colon, sheet is empty (meaning first sheet).
func ParseFileArg(arg string) (file, sheet string) {
	idx := strings.LastIndex(arg, ":")
	if idx == -1 {
		return arg, ""
	}
	return arg[:idx], arg[idx+1:]
}

// ParseSkipFlag parses "--skip N" or "--skip N,M" into two skip values.
func ParseSkipFlag(flag string) (skip1, skip2 int, err error) {
	parts := strings.Split(flag, ",")
	switch len(parts) {
	case 1:
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid skip value %q: %w", parts[0], err)
		}
		return n, n, nil
	case 2:
		n1, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid skip value %q: %w", parts[0], err)
		}
		n2, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid skip value %q: %w", parts[1], err)
		}
		return n1, n2, nil
	default:
		return 0, 0, fmt.Errorf("invalid skip flag %q: expected N or N,M", flag)
	}
}
