package generator

import (
	"fmt"
	"strings"

	j "github.com/dave/jennifer/jen"
)

func renderStatement(s *j.Statement) string {
	return s.GoString()
}

// Fragment is a simple type for grouping comments with arbitrary jen statements
// Jen only has nice formatting for comments in the context of a File, which is a
// type we are trying to avoid in order to use Jen to only render more complex fragments
// on the assumption that go templates will be more readable for files.
type Fragment struct {
	name      string
	comments  []string
	statement *j.Statement
}

// Render concats the fragment comments with each other,
// and the rendered Statement, in a predictable way
func (f *Fragment) Render() string {
	return fmt.Sprintf("%s%s", f.FormattedComments(), renderStatement(f.statement))
}

// FormattedComments uses a simple rendering algo of
// prepending every line in the set of all comments with '//'
// and joining them together with '\n'
func (f *Fragment) FormattedComments() string {
	if f.comments == nil || len(f.comments) == 0 {
		return ""
	}
	fmtd := ""
	for _, c := range f.comments {
		for _, line := range strings.Split(c, "\n") {
			if line == "" {
				fmtd = fmtd + "\n"
				continue
			}
			fmtd = fmt.Sprintf("%s// %s\n", fmtd, line)
		}
	}
	return fmtd
}
