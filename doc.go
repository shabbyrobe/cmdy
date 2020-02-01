/*
Package cmdy is a self-contained library for implementing CLI programs.

cmdy combines the features I like from the flag stdlib package with the
features I like from https://github.com/google/subcommands.

cmdy probably doesn't really need to exist, but I like it and use it for
my own projects. There are a lot of CLI libraries for Go but this one is mine.

cmdy focuses on minimalism and tries to imitate and leverage the stdlib as
much as possible. It does not attempt to replace flag.Flag, though it does
extend it slightly.

cmdy has no dependencies beyond the stdlib.
*/
package cmdy
