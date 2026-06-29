package sed_test

import (
	"fmt"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-sed"
)

func ExampleSed_substitute() {
	// echo "hello world" | sed 's/world/universe/'
	output, _ := testable.Test(command.Sed("s/world/universe/"), "hello world\n")
	fmt.Print(output)
	// Output:
	// hello universe
}
