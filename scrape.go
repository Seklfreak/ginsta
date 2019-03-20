package ginsta

import (
	"errors"
	"strings"
)

func extractSharedData(pageContent string) (sharedData string, err error) {
	parts := strings.Split(pageContent, "window._sharedData = ")

	if len(parts) < 2 {
		return sharedData, errors.New("unable to parse shared data")
	}

	subParts := strings.Split(parts[1], "};")

	if len(subParts) < 1 {
		return sharedData, errors.New("unable to parse shared data")
	}

	return subParts[0] + "}", nil
}
