package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/rcarmo/go-te/pkg/te"
)

func main() {
	screen := te.NewScreen(80, 24)
	screen.WriteProcessInput = func(data string) {
		_, _ = os.Stdout.Write([]byte(data))
	}
	stream := te.NewStream(screen, false)

	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if feedErr := stream.FeedBytes(buf[:n]); feedErr != nil {
				_, _ = fmt.Fprintln(os.Stderr, "feed error:", feedErr)
			}
		}
		if err != nil {
			if err == io.EOF {
				return
			}
			_, _ = fmt.Fprintln(os.Stderr, "read error:", err)
			return
		}
	}
}
