package security

import (
	"fmt"
	"strings"

	"github.com/k3a/html2text"
)

func stripTags(content string) string {
	return html2text.HTML2Text(content)
}

func AllowedWord(text string, bannedWords map[string]bool) bool {
	strippedText := stripTags(text)
	for _, word := range strings.Split(strings.ToLower(strippedText), " ") {
		fmt.Println(word)
		if _, found := bannedWords[word]; found {
			return false
		}
	}

	return true
}
