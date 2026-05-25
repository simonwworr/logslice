package aggregate

import "errors"

// errInvalidDuration is returned when a non-positive bucket duration is given.
var errInvalidDuration = errors.New("aggregate: bucket duration must be positive")
