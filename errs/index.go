package errs

import (
	"fmt"
	"github.com/pkg/errors"
)

type userError struct {
	error
}

func Raise(err error) {
	panic(userError{err})
}

func RaiseF(format string, args ...interface{}) {
	Raise(fmt.Errorf(format, args...))
}

func CatchInstance(instance error, cb func()) {
	if err := recover(); err != nil {
		if user_err, type_matches := err.(userError); type_matches {
			if instance == user_err.error {
				cb()
			} else {
				panic(user_err)
			}
		} else {
			panic(err)
		}
	}
}

func Catch[T error](cb func(T)) {
	if err := recover(); err != nil {
		if user_err, type_matches := err.(userError); type_matches {
			if err, type_matches := user_err.error.(T); type_matches {
				cb(err)
			} else {
				panic(user_err)
			}
		} else {
			panic(err)
		}
	}
}

func Propagate(err error) {
	if err != nil {
		Raise(err)
	}
}

const assumption_crash_rationale = "Deliberatly crashing to avoid progressing based on a " +
	"faulty/dangerous/undefined state, potentially causing even more trouble."

func AssumptionHasBeenBroken() {
	panic("Something that was assumed to never happen happened. " + assumption_crash_rationale)
}

func Assume(b bool) {
	if !b {
		AssumptionHasBeenBroken()
	}
}

func AssumeNoError(err error) {
	if err != nil {
		panic(errors.Wrap(err, "The error was assumed never to happen, yet it has. "+assumption_crash_rationale))
	}
}

type ErrorString string

func (self ErrorString) Error() string {
	return string(self)
}
