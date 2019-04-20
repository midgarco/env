package env

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Load will load an ini formatted file into
// the os ENVIRONMENT
func Load(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		ln, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// ignore comments
		if strings.HasPrefix(ln, "#") {
			continue
		}

		parts := strings.SplitN(ln, "=", 2)
		if len(parts) != 2 {
			continue
		}
		err = os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		if err != nil {
			return err
		}
	}
	return nil
}
