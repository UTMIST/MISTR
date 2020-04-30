package update

import (
	"bytes"
	"log"
	"os/exec"
)

const upToDate = "Already up to date."

// IsUpdated return whether there's a new commit to the repository.
func IsUpdated() bool {
	cmd := exec.Command("git", "pull")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	str := out.String()
	return len(upToDate) <= len(str) && str[:len(upToDate)] == upToDate
}
