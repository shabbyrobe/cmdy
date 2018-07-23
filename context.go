package cmdy

import "context"

// Context implements context.Context.
//
// NOTE: this wraps context.Context for future-proofing reasons;
// the API is not yet well tested and this gives us the option to
// extend the interface without breaking BC.
type Context interface {
	context.Context
}

type commandContext struct {
	context.Context
}
