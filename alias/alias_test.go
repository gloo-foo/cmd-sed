package alias_test

import (
	"slices"
	"testing"

	"github.com/gloo-foo/testable"

	sed "github.com/gloo-foo/cmd-sed/alias"
)

// The alias package re-exports the constructor under an unprefixed name. A
// mis-wired re-export (Sed bound to the wrong function) compiles cleanly, so
// only behavior can prove the wiring. Each test drives one supported s///
// flag through the re-exported constructor and asserts the GNU sed output it
// must produce.

func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAlias_FirstOccurrence(t *testing.T) {
	lines, err := testable.TestLines(sed.Sed("s/a/x/"), "aaa\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"xaa"})
}

func TestAlias_GlobalFlag(t *testing.T) {
	lines, err := testable.TestLines(sed.Sed("s/a/x/g"), "aaa\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"xxx"})
}

func TestAlias_IgnoreCaseFlag(t *testing.T) {
	lines, err := testable.TestLines(sed.Sed("s/abc/x/i"), "ABC\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"x"})
}

func TestAlias_NthFlag(t *testing.T) {
	lines, err := testable.TestLines(sed.Sed("s/a/x/2"), "aaa\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"axa"})
}

func TestAlias_PrintFlag(t *testing.T) {
	lines, err := testable.TestLines(sed.Sed("s/a/x/p"), "abc\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"xbc", "xbc"})
}
