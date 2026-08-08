package lexical

import "strings"

type Comment struct {
	Content []string
}

func NewComment(content ...string) Comment {
	return Comment{Content: content}
}

func (c Comment) String() string {
	sb := strings.Builder{}
	for _, line := range c.Content {
		sb.WriteString("// " + line + "\n")
	}
	return sb.String()
}