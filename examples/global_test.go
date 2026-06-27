package sed_test

import (
	"fmt"

	command "github.com/gloo-foo/cmd-sed"
	"github.com/gloo-foo/testable"
)

func ExampleSed_global() {
	// echo "foo foo foo" | sed 's/foo/bar/g'
	output, _ := testable.Test(command.Sed("s/foo/bar/g"), "foo foo foo\n")
	fmt.Print(output)
	// Output:
	// bar bar bar
}
