package filestorage

import (
	"strconv"
	"strings"
)

func generateFileName(name string, ownerID int, isPrivate bool) string {
	nameLen := len(name) + 5

	if isPrivate {
		nameLen += 8
	}

	builder := strings.Builder{}
	builder.Grow(nameLen)

	if isPrivate {
		builder.WriteString("private/")
	}

	builder.WriteString(strconv.Itoa(ownerID))
	builder.WriteString("/")
	builder.WriteString(name)

	return builder.String()
}
