package env

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Load will load an ini formatted file into
// the os ENVIRONMENT
func Load(path string, opts ...Option) error {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	err := load(path + ".env")
	if err != nil {
		return err
	}
	// Set override options
	Override(opts...)
	return nil
}

func load(name string) error {
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

// Option allows you to provide addtional settings to override
type Option struct {
	Key   string
	Value string
}

// Get returns the value from the envronment config
func Get(key string) string {
	return os.Getenv(key)
}

// GetWithDefault returns the value from the environment config
// or returns a default value if the setting is empty
func GetWithDefault(key, def string) string {
	s := os.Getenv(key)
	if s == "" {
		s = def
	}
	return s
}

// Override environment config with additional options
func Override(opts ...Option) {
	for _, opt := range opts {
		os.Setenv(opt.Key, opt.Value)
	}
}

// GetBool returns a boolean configuration setting
func GetBool(key string) bool {
	v := GetWithDefault(key, "false")
	return (strings.ToLower(v) == "true" || v == "1")
}
