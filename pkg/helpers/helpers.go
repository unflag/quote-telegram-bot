package helpers

import "regexp"

var (
	forbiddenChars = regexp.MustCompile(`[^a-zA-Z0-9\-_\.=]`)
)

func Sanitize(text string) string {
	return forbiddenChars.ReplaceAllLiteralString(text, "")
}

func Retry(retry int, f func() error) error {
	var (
		errorCount int
		err        error
	)

	for errorCount < retry {
		if err = f(); err == nil {
			return nil
		}
		errorCount++
	}

	return err
}
