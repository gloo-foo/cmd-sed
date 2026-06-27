// Package alias provides an unprefixed name for the sed command.
//
//	import sed "github.com/gloo-foo/cmd-sed/alias"
//	sed.Sed("s/old/new/g")
//
// sed's flags (g, i, p, N) are part of the s/// script string, not separate
// option values, so the constructor is the only re-export.
package alias

import command "github.com/gloo-foo/cmd-sed"

// Sed re-exports the substitution command constructor.
var Sed = command.Sed
