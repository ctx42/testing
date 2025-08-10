// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"errors"
	"os/exec"

	"github.com/ctx42/testing/pkg/notice"
)

// ExitCode checks "err" is pointer to [exec.ExitError] with exit code equal to
// "want". Returns nil if it's, otherwise it returns an error with a message
// indicating the expected and actual values.
func ExitCode(want int, err error, opts ...any) error {
	if want == 0 && err == nil {
		return nil
	}
	if err == nil {
		return notice.New("expected *exec.ExitError got nil")
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		have := ee.ExitCode()
		if want != have {
			ops := DefaultOptions(opts...)
			msg := notice.New("expected exit code").
				Want("%d", want).
				Have("%d", have)
			return AddRows(ops, msg)
		}
		return nil
	}

	ops := DefaultOptions(opts...)
	msg := notice.New("expected err to have \"%T\" in its chain", ee)
	return AddRows(ops, msg)
}
