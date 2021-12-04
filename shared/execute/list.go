package exec

import (
	"os/exec"
	"strings"
)

func List(list []string) (string, error) {
	out, err := exec.Command("bash", list...).Output()
	return strings.TrimSpace(string(out)), err
}
