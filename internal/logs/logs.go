package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"sponsorbucks-client/internal/localconfig"
)

func Path() (string, error) {
	dir, err := localconfig.ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "logs", "sponsorbucks.log"), nil
}

func Append(event string, fields map[string]string) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString(time.Now().UTC().Format(time.RFC3339))
	b.WriteString(" event=")
	b.WriteString(safeField(event))
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := fields[k]
		if v == "" {
			continue
		}
		b.WriteByte(' ')
		b.WriteString(safeField(k))
		b.WriteByte('=')
		b.WriteString(safeField(v))
	}
	b.WriteByte('\n')
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(b.String())
	return err
}

func Read() (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

func safeField(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\t", " ")
	return fmt.Sprintf("%q", value)
}
