package system

import (
	"context"
	"errors"

	"master-control/model"
)

var (
	ErrNotSupported = errors.New("el control nativo de sistema solo está disponible en macOS")
	ErrUnknownToken = errors.New("comando de sistema desconocido")
)

var catalog = []model.SystemAction{
	{ID: "volume_up", Label: "Subir volumen", Icon: "🔊", Command: "volume_up", Category: "volumen"},
	{ID: "volume_down", Label: "Bajar volumen", Icon: "🔉", Command: "volume_down", Category: "volumen"},
	{ID: "volume_mute", Label: "Silenciar", Icon: "🔇", Command: "volume_mute", Category: "volumen"},
	{ID: "volume_unmute", Label: "Quitar silencio", Icon: "🔈", Command: "volume_unmute", Category: "volumen"},
	{ID: "brightness_up", Label: "Subir brillo", Icon: "☀️", Command: "brightness_up", Category: "brillo"},
	{ID: "brightness_down", Label: "Bajar brillo", Icon: "🌙", Command: "brightness_down", Category: "brillo"},
	{ID: "media_play_pause", Label: "Play / Pausa", Icon: "⏯️", Command: "media_play_pause", Category: "medios"},
	{ID: "media_next", Label: "Siguiente pista", Icon: "⏭️", Command: "media_next", Category: "medios"},
	{ID: "media_prev", Label: "Pista anterior", Icon: "⏮️", Command: "media_prev", Category: "medios"},
}

type System struct{}

func New() *System {
	return &System{}
}

func (s *System) GetCatalog(_ context.Context) ([]model.SystemAction, error) {
	return catalog, nil
}

func (s *System) Execute(ctx context.Context, token string) error {
	return executeToken(ctx, token)
}

func (s *System) GetVolume(ctx context.Context) (model.SystemVolume, error) {
	return getVolume(ctx)
}

func (s *System) GetBrightness(ctx context.Context) (model.SystemBrightness, error) {
	return getBrightness()
}

func (s *System) SetVolume(ctx context.Context, volume int) error {
	if volume < 0 {
		volume = 0
	}
	if volume > 100 {
		volume = 100
	}
	return setVolume(ctx, volume)
}

func (s *System) SetBrightness(ctx context.Context, brightness float64) error {
	if brightness < 0 {
		brightness = 0
	}
	if brightness > 1 {
		brightness = 1
	}
	return setBrightness(ctx, brightness)
}