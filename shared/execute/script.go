package exec

import (
	"os/exec"
	"strings"
)

func Script(path string) (string, error) {
	out, err := exec.Command("bash", path).Output()
	return strings.TrimSpace(string(out)), err
}
