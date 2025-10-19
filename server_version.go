package nuxeo

import (
	"fmt"
	"strings"
)

// ServerVersion represents the Nuxeo server version.
type ServerVersion struct {
	Major int
	Minor int
	Patch int
}

// ParseServerVersion parses a version string (e.g., "10.2.3") into a ServerVersion.
func ParseServerVersion(version string) *ServerVersion {
	parts := strings.Split(version, ".")
	v := &ServerVersion{}
	if len(parts) > 0 {
		fmt.Sscanf(parts[0], "%d", &v.Major)
	}
	if len(parts) > 1 {
		fmt.Sscanf(parts[1], "%d", &v.Minor)
	}
	if len(parts) > 2 {
		fmt.Sscanf(parts[2], "%d", &v.Patch)
	}
	return v
}
