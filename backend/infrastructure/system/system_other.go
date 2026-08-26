//go:build !darwin

package system

import (
	"context"

	"master-control/model"
)

func executeToken(_ context.Context, _ string) error {
	return ErrNotSupported
}

func getVolume(_ context.Context) (model.SystemVolume, error) {
	return model.SystemVolume{}, ErrNotSupported
}

func getBrightness() (model.SystemBrightness, error) {
	return model.SystemBrightness{}, ErrNotSupported
}

func setVolume(_ context.Context, _ int) error {
	return ErrNotSupported
}

func setBrightness(_ context.Context, _ float64) error {
	return ErrNotSupported
}