package command_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gloo-foo/testable"
	"github.com/gloo-foo/testable/assertion"

	command "github.com/gloo-foo/cmd-sed"
)

func TestSed_BasicSubstitution(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/a/x/"), "abc\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"xbc"})
}

func TestSed_FirstOccurrenceOnly(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/a/x/"), "aaa\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"xaa"})
}

func TestSed_GlobalSubstitution(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/a/x/g"), "aaa\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"xxx"})
}

func TestSed_NthOccurrence(t *testing.T) {
	// The 2 flag rewrites only the second match, leaving the first and third.
	lines, err := testable.TestLines(command.Sed("s/a/x/2"), "aaa\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"axa"})
}

func TestSed_NthBeyondMatchCountChangesNothing(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/a/x/5"), "aaa\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"aaa"})
}

func TestSed_IgnoreCase(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/abc/x/i"), "ABC\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"x"})
}

func TestSed_IgnoreCaseGlobal(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/a/x/gi"), "AaA\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"xxx"})
}

func TestSed_PrintFlagDuplicatesChangedLine(t *testing.T) {
	// p emits the line a second time, but only because a substitution occurred.
	lines, err := testable.TestLines(command.Sed("s/a/x/p"), "abc\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"xbc", "xbc"})
}

func TestSed_PrintFlagSkipsUnchangedLine(t *testing.T) {
	// No match means no substitution, so p must not duplicate the line.
	lines, err := testable.TestLines(command.Sed("s/z/x/p"), "abc\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"abc"})
}

func TestSed_NoMatch(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/z/x/"), "abc\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"abc"})
}

func TestSed_MultipleLines(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/hello/world/"), "hello foo\nhello bar\nno match\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"world foo", "world bar", "no match"})
}

func TestSed_RegexPatternGlobal(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/[0-9]+/NUM/g"), "abc 123 def 456\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"abc NUM def NUM"})
}

func TestSed_AlternateDelimiter(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s|/usr/local|/opt|"), "/usr/local/bin\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"/opt/bin"})
}

func TestSed_EmptyReplacement(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/foo//"), "foobar\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"bar"})
}

func TestSed_CaptureGroup(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/(\\w+) (\\w+)/${2} ${1}/"), "hello world\n")
	assertion.NoError(t, err)
	assertion.Lines(t, lines, []string{"world hello"})
}

func TestSed_EmptyInput(t *testing.T) {
	lines, err := testable.TestLines(command.Sed("s/a/x/"), "")
	assertion.NoError(t, err)
	assertion.Empty(t, lines)
}

func TestSed_UnsupportedExpressionErrors(t *testing.T) {
	_, err := testable.TestLines(command.Sed("x"), "abc\n")
	assertion.Error(t, err)
	if !errors.Is(err, command.ErrUnsupportedExpression) {
		t.Fatalf("got %v, want ErrUnsupportedExpression", err)
	}
}

func TestSed_EmptyExpressionErrors(t *testing.T) {
	_, err := testable.TestLines(command.Sed(""), "abc\n")
	assertion.Error(t, err)
	if !errors.Is(err, command.ErrUnsupportedExpression) {
		t.Fatalf("got %v, want ErrUnsupportedExpression", err)
	}
}

func TestSed_IncompleteSubstitutionErrors(t *testing.T) {
	// Only one delimiter after the pattern: no replacement/flags section.
	_, err := testable.TestLines(command.Sed("s/a"), "abc\n")
	assertion.Error(t, err)
	if !errors.Is(err, command.ErrIncompleteSubstitution) {
		t.Fatalf("got %v, want ErrIncompleteSubstitution", err)
	}
}

func TestSed_InvalidPatternErrors(t *testing.T) {
	_, err := testable.TestLines(command.Sed("s/[invalid/x/"), "abc\n")
	assertion.Error(t, err)
	if !errors.Is(err, command.ErrInvalidPattern) {
		t.Fatalf("got %v, want ErrInvalidPattern", err)
	}
}

func TestSed_InvalidFlagsError(t *testing.T) {
	_, err := testable.TestLines(command.Sed("s/a/x/z"), "abc\n")
	assertion.Error(t, err)
	if !errors.Is(err, command.ErrInvalidFlags) {
		t.Fatalf("got %v, want ErrInvalidFlags", err)
	}
}

func TestSed_ErrorMessage(t *testing.T) {
	if command.ErrUnsupportedExpression.Error() != "sed: unsupported expression" {
		t.Fatalf("unexpected error string: %q", command.ErrUnsupportedExpression.Error())
	}
}

func ExampleSed() {
	lines, _ := testable.TestLines(command.Sed("s/hello/world/"), "hello foo\nhello bar\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// world foo
	// world bar
}

func ExampleSed_global() {
	lines, _ := testable.TestLines(command.Sed("s/a/x/g"), "banana\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// bxnxnx
}
