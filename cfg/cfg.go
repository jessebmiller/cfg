package cfg

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func init() {
	http.HandleFunc(
		"/cfg-req",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, os.Getenv("CFG_REQFILE"))
	})
}

// panicOn panics if err is not nil
func panicOn(err error) {
	if err != nil { panic(err) }
}

// find returns the configuration value set for a key or not found error
func find(key string) (string, error) {
	val := os.Getenv(key)
	var err error
	if val == "" {
		err = fmt.Errorf("Missing config. Key %s not found.", key)
	} else {
		err = nil
	}
	return val, err
}

// Find remembers the key then calls Find
func Find(key string) (string, error) {
	defer rememberPair(key, "")
	return find(key)
}


// Get returns a set value for a key or a default value
func Get(key string, defaultVal string) (string) {
	defer rememberPair(key, defaultVal)
	val, err := find(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func parseFile() map[string][]string {
	file, err := os.Open(os.Getenv("CFG_REQFILE"))
	panicOn(err)
	defer file.Close()
	m := make(map[string][]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		kv := strings.Split(scanner.Text(), "=")
		m[kv[0]] = strings.Split(kv[1], ",")
	}
	return m
}

func writeFile(pairs map[string][]string) {
	file, err := os.Create(os.Getenv("CFG_REQFILE"))
	panicOn(err)
	defer file.Close()
	for k, vs := range pairs {
		file.WriteString(k)
		file.WriteString("=")
		file.WriteString(strings.Join(vs, ","))
		file.WriteString("\n")
	}
}

func rememberPair(key string, val string) error {
	kvs := parseFile()
	kvs[key] = append(kvs[key], val)
	writeFile(kvs)
	return nil
}
