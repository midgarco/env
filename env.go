package env

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

var AppPrefix = ""

// Load will load an ini formatted file into
// the os ENVIRONMENT
func Load(prefix, path string, opts ...Option) error {
	AppPrefix = prefix
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
		key := strings.TrimSpace(AppPrefix + "_" + parts[0])
		if val := os.Getenv(key); val != "" {
			return nil
		}
		err = os.Setenv(key, strings.TrimSpace(parts[1]))
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

// Get returns the value from the environment config
func Get(key string) string {
	if !strings.HasPrefix(key, AppPrefix) {
		key = AppPrefix + "_" + key
	}
	return os.Getenv(key)
}

// GetWithDefault returns the value from the environment config
// or returns a default value if the setting is empty
func GetWithDefault(key, def string) string {
	if !strings.HasPrefix(key, AppPrefix) {
		key = AppPrefix + "_" + key
	}
	s := os.Getenv(key)
	if s == "" {
		s = def
	}
	return s
}

// Override environment config with additional options
func Override(opts ...Option) {
	for _, opt := range opts {
		if !strings.HasPrefix(opt.Key, AppPrefix) {
			opt.Key = AppPrefix + "_" + opt.Key
		}
		os.Setenv(strings.TrimSpace(opt.Key), strings.TrimSpace(opt.Value))
	}
}

// GetBool returns a boolean configuration setting
func GetBool(key string) bool {
	if !strings.HasPrefix(key, AppPrefix) {
		key = AppPrefix + "_" + key
	}
	v := GetWithDefault(key, "false")
	return (strings.ToLower(v) == "true" || v == "1")
}

// GetInt returns the integer value from the environment config
func GetInt(key string) int {
	if !strings.HasPrefix(key, AppPrefix) {
		key = AppPrefix + "_" + key
	}
	str := Get(key)
	v, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return v
}

// GetIntWithDefault returns the integer value from the environment config
// or returns a default value if the setting is empty
func GetIntWithDefault(key string, def int) int {
	if !strings.HasPrefix(key, AppPrefix) {
		key = AppPrefix + "_" + key
	}
	str := GetWithDefault(key, strconv.Itoa(def))
	v, err := strconv.Atoi(str)
	if err != nil {
		return def
	}
	return v
}
