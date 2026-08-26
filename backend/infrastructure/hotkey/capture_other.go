//go:build !darwin

package hotkey

import "time"

func captureKey(timeout time.Duration) (int64, int64, error) {
	return 0, 0, ErrNotSupported
}