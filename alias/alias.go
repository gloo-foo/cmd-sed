// Package alias provides an unprefixed name for the sed command.
//
//	import sed "github.com/gloo-foo/cmd-sed/alias"
//	sed.Sed("s/old/new/g")
//
// sed's flags (g, i, p, N) are part of the s/// script string, not separate
// option values, so the constructor is the only re-export.
package alias

import (
	gloo "github.com/gloo-foo/framework"

	command "github.com/gloo-foo/cmd-sed"
)

// Sed re-exports the substitution command constructor.
func Sed(script command.SedScript) gloo.Command[[]byte, []byte] { return command.Sed(script) }

// Script names the s/// substitution script argument.
type Script = command.SedScript
