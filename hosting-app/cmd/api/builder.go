package main

import (
	"fmt"
	"os"
	"os/exec"
)

func BuildRepo(repoPath string) error {
	err := os.Chdir(repoPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("npm", "run", "build")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	fmt.Println("Build is successful")
	fmt.Println("Output", string(output), "\n\n\n\n\n\n\n\n\n\n\n\n\nOutput ends here")

	return nil
}
