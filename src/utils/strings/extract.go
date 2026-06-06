package utils_strings

import "regexp"

func Extract(src string, regExp string) *string {
	re := regexp.MustCompile(regExp)
	matches := re.FindStringSubmatch(src)
	if len(matches) > 1 {
		return &matches[1]
	}
	return nil
}
