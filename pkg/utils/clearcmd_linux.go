package utils

import (
	"os"
	"os/exec"
)

func ClearCmd() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
