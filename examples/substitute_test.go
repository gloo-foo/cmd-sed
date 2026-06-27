package sed_test

import (
	"fmt"

	command "github.com/gloo-foo/cmd-sed"
	"github.com/gloo-foo/testable"
)

func ExampleSed_substitute() {
	// echo "hello world" | sed 's/world/universe/'
	output, _ := testable.Test(command.Sed("s/world/universe/"), "hello world\n")
	fmt.Print(output)
	// Output:
	// hello universe
}
