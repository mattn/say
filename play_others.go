// +build !darwin, !windows

package main

import (
	"errors"
	"os"
	"os/exec"
)

func play(filename string) error {
	var cmd *exec.Cmd
	for _, player := range []string{"ffplay", "avplay"} {
		f, err := exec.LookPath(player)
		if err == nil {
			args := []string{"-autoexit", "-nodisp", filename}
			cmd = exec.Command(f, args...)
			cmd.Stdin = os.Stdin
			return cmd.Run()
		}
	}

	f, err := exec.LookPath("mplayer")
	if err == nil {
		cmd = exec.Command(f, filename)
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}
	return errors.New("player not found")
}
