package sed_test

import (
	"fmt"
	"os"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-sed"
)

// This example demonstrates reading from a file instead of inline input.
func ExampleSed_fromFile_substitute() {
	// sed 's/world/universe/' testdata/text.txt
	data, err := os.ReadFile("testdata/text.txt")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "read testdata: %v\n", err)
		return
	}
	output, _ := testable.Test(command.Sed("s/world/universe/"), string(data))
	fmt.Print(output)
	// Output:
	// hello universe
}
