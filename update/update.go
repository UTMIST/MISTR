package update

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

const upToDate = "Already up to date."

const fetch = "git fetch --all"
const reset = "git reset --hard origin/master"

// IsUpdated return whether there's a new commit to the repository.
func IsUpdated() bool {
	var out bytes.Buffer

	for _, cmd_str := range []string{fetch, reset} {
		cmd_strs := strings.Split(cmd_str, " ")
		cmd := exec.Command(cmd_strs[0], cmd_strs[1:]...)
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

	str := out.String()
	return len(upToDate) <= len(str) && str[:len(upToDate)] == upToDate
}
