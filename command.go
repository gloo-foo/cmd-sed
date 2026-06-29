package command

import (
	"regexp"
	"strings"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// Error is the sentinel error type for every failure this package can emit.
// Comparing with errors.Is against the constants below is the supported test.
type Error string

func (e Error) Error() string { return string(e) }

const (
	// ErrUnsupportedExpression is returned when the script is not an s/// command.
	ErrUnsupportedExpression Error = "sed: unsupported expression"
	// ErrIncompleteSubstitution is returned when the s/// form lacks its parts.
	ErrIncompleteSubstitution Error = "sed: incomplete substitution"
	// ErrInvalidPattern is returned when the regular expression does not compile.
	ErrInvalidPattern Error = "sed: invalid pattern"
	// ErrInvalidFlags is returned when the trailing flag string is not understood.
	ErrInvalidFlags Error = "sed: invalid flags"
)

// expression is the part of an s/// script after the delimiter has been
// stripped: the raw pattern, the raw replacement, and the raw flag string.
type expression struct {
	pattern     string
	replacement string
	flagString  string
}

// substitution is a compiled, ready-to-apply s/// command.
type substitution struct {
	re      *regexp.Regexp
	pattern string
	repl    []byte
	flags   substFlags
}

// substFlags is the parsed trailing-flag set of an s/// command.
//
//   - global:     g — replace every match on the line, not just the first.
//   - ignoreCase: i — match case-insensitively.
//   - printIfSub: p — emit the line a second time when a substitution occurred.
//   - nth:        N — replace only the Nth match (1-based); 0 means "first".
type substFlags struct {
	global     bool
	ignoreCase bool
	printIfSub bool
	nth        int
}

// Sed returns a command that applies a single GNU-style s/// substitution to
// each input line.
//
// Supported script form: s<delim>pattern<delim>replacement<delim>flags
// Any single byte may be the delimiter (e.g. s|old|new|). Supported flags:
//
//   - g — global: replace every match on the line.
//   - i — ignore case.
//   - p — print: emit the line again when a substitution occurred.
//   - N — a decimal number: replace only the Nth match.
//
// An invalid script makes the command fail on its first line so the error
// surfaces through the pipeline.
func Sed(expr string) gloo.Command[[]byte, []byte] {
	sub, err := compile(expr)
	if err != nil {
		return patterns.Expand(func(_ []byte) ([][]byte, error) { return nil, err })
	}
	return patterns.Expand(sub.apply)
}

// compile parses and compiles a script into an applicable substitution.
func compile(expr string) (substitution, error) {
	parsed, err := split(expr)
	if err != nil {
		return substitution{}, err
	}
	flags, err := parseFlags(parsed.flagString)
	if err != nil {
		return substitution{}, err
	}
	re, err := compilePattern(parsed.pattern, flags.ignoreCase)
	if err != nil {
		return substitution{}, err
	}
	return substitution{re: re, repl: []byte(parsed.replacement), flags: flags, pattern: parsed.pattern}, nil
}

// split breaks an s<delim>...<delim>...<delim>flags script into its raw parts.
func split(expr string) (expression, error) {
	if len(expr) < 2 || expr[0] != 's' {
		return expression{}, ErrUnsupportedExpression
	}
	parts := strings.SplitN(expr[2:], expr[1:2], 3)
	if len(parts) < 3 {
		return expression{}, ErrIncompleteSubstitution
	}
	return expression{pattern: parts[0], replacement: parts[1], flagString: parts[2]}, nil
}

// compilePattern compiles the regular expression, optionally case-insensitively.
func compilePattern(pattern string, ignoreCase bool) (*regexp.Regexp, error) {
	source := pattern
	if ignoreCase {
		source = "(?i)" + pattern
	}
	re, err := regexp.Compile(source)
	if err != nil {
		return nil, ErrInvalidPattern
	}
	return re, nil
}

// parseFlags interprets the trailing flag string of an s/// command.
func parseFlags(s string) (substFlags, error) {
	flags := substFlags{}
	for _, r := range s {
		if err := applyFlagRune(r, &flags); err != nil {
			return substFlags{}, err
		}
	}
	return flags, nil
}

// applyFlagRune folds one flag character into flags, accumulating Nth digits.
func applyFlagRune(r rune, flags *substFlags) error {
	switch r {
	case 'g':
		flags.global = true
	case 'i':
		flags.ignoreCase = true
	case 'p':
		flags.printIfSub = true
	default:
		return applyDigit(r, flags)
	}
	return nil
}

// applyDigit folds a decimal digit into the running Nth-match selector.
func applyDigit(r rune, flags *substFlags) error {
	if r < '0' || r > '9' {
		return ErrInvalidFlags
	}
	flags.nth = flags.nth*10 + int(r-'0')
	return nil
}

// apply runs the substitution on one line, returning one line normally and two
// (the line repeated) when the p flag fired on a line that actually changed.
func (s substitution) apply(line []byte) ([][]byte, error) {
	result, substituted := s.substitute(line)
	if s.flags.printIfSub && substituted {
		return [][]byte{result, result}, nil
	}
	return [][]byte{result}, nil
}

// substitute rewrites line per the flags, reporting whether anything changed.
func (s substitution) substitute(line []byte) ([]byte, bool) {
	if s.flags.global {
		return s.replaceFrom(line, 0)
	}
	return s.replaceFrom(line, s.targetIndex())
}

// targetIndex is the 0-based match index a non-global substitution rewrites.
func (s substitution) targetIndex() int {
	if s.flags.nth > 0 {
		return s.flags.nth - 1
	}
	return 0
}

// replaceFrom rewrites matches at or after the 0-based start index. A global
// run (start 0 via the global flag) rewrites every match from start onward;
// a targeted run rewrites exactly the one match at start. It reports whether
// any rewrite happened.
func (s substitution) replaceFrom(line []byte, start int) ([]byte, bool) {
	index := 0
	changed := false
	rewrite := func(match []byte) []byte {
		current := index
		index++
		if s.selected(current, start) {
			changed = true
			return s.re.ReplaceAll(match, s.repl)
		}
		return match
	}
	return s.re.ReplaceAllFunc(line, rewrite), changed
}

// selected reports whether the match at position current should be rewritten:
// every position from start onward when global, otherwise exactly start.
func (s substitution) selected(current, start int) bool {
	if s.flags.global {
		return current >= start
	}
	return current == start
}
