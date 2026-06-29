package sed_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-sed"
)

func ExampleSed_global() {
	// echo "foo foo foo" | sed 's/foo/bar/g'
	output, _ := testable.Test(command.Sed("s/foo/bar/g"), "foo foo foo\n")
	fmt.Print(output)
	// Output:
	// bar bar bar
}
