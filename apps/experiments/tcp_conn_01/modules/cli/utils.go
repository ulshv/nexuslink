package cli

import (
	"fmt"
	"strings"
)

func parseHostPort(hostPort string) (string, string, error) {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid <host>:<port> param")
	}
	return parts[0], parts[1], nil
}
