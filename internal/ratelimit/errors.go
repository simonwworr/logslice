package ratelimit

import "errors"

// ErrInvalidRate is returned when a non-positive rate is provided to New.
var ErrInvalidRate = errors.New("ratelimit: rate must be positive")

// ErrInvalidBurst is returned when a burst value less than 1 is provided to New.
var ErrInvalidBurst = errors.New("ratelimit: burst must be at least 1")
