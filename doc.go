/*
Package cmdy is a self-contained library for implementing CLI programs.

cmdy combines the features I like from the `flag` stdlib package with the
features I like from https://github.com/google/subcommands.

cmdy focuses on minimalism and tries to imitate and leverage the stdlib as much as
possible. It does not attempt to replace `flag.Flag`, though it does extend it slightly.
It somehow still ended up being bigger than I'd hoped.

cmdy has no dependencies.

cmdy does not prioritise performance (beyond the fact that Go is already
pretty fast out of the box).

*/
package cmdy
