package pkg

import "time"

type CheckResult struct {
	// CheckResult is the result of checking
	StatusCode int
	Duration   time.Duration
	Err        error
}
