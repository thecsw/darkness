package emilia

import (
	"strings"

	"github.com/thecsw/darkness/internals"
)

func SourceCodeTrimLeftWhitespace(page *internals.Page) {
	for i, content := range page.Contents {
		if !content.IsSourceCode() {
			continue
		}
		lines := strings.Split(content.SourceCode, "\n")
		if len(lines) < 1 {
			continue
		}
		offset := len(lines[0]) - len(strings.TrimLeft(lines[0], " "))
		for i, line := range lines {
			// if the initial offset is bigger, then abort the whole thing
			if len(line)-len(strings.TrimLeft(line, " ")) < offset {
				return
			}
			lines[i] = line[internals.Min(len(lines[i]), offset):]
		}
		(&page.Contents[i]).SourceCode = strings.Join(lines, "\n")
	}
}
