package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/creack/pty"
	"github.com/rcarmo/go-te/pkg/te"
)

func main() {
	python := os.Getenv("ESCTEST2_PYTHON")
	if python == "" {
		python = "python3"
	}

	esctestPath := filepath.Join("vendor", "esctest2", "esctest", "esctest.py")
	args := append([]string{esctestPath}, os.Args[1:]...)
	cmd := exec.Command(python, args...)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to start esctest2:", err)
		os.Exit(1)
	}
	defer func() {
		_ = ptmx.Close()
	}()

	screen := te.NewScreen(80, 24)
	screen.WriteProcessInput = func(data string) {
		_, _ = ptmx.Write([]byte(data))
	}
	stream := te.NewStream(screen, false)

	buf := make([]byte, 4096)
	for {
		n, readErr := ptmx.Read(buf)
		if n > 0 {
			if feedErr := stream.FeedBytes(buf[:n]); feedErr != nil {
				_, _ = fmt.Fprintln(os.Stderr, "feed error:", feedErr)
			}
		}
		if readErr != nil {
			if readErr != io.EOF {
				_, _ = fmt.Fprintln(os.Stderr, "read error:", readErr)
			}
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		_, _ = fmt.Fprintln(os.Stderr, "esctest2 exit error:", err)
		os.Exit(1)
	}
}
