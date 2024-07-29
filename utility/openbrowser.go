package utility

import (
	"os/exec"
	"runtime"
)

// OpenBrowser from gist https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8
func OpenBrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
