package system

import (
	"context"
	"errors"

	"master-control/model"
)

var (
	ErrNotSupported = errors.New("el control nativo de sistema solo está disponible en macOS")
	ErrUnknownToken = errors.New("comando de sistema desconocido")
	// ErrNoAccessibility: inyectar eventos (scroll) exige el permiso de Accesibilidad.
	ErrNoAccessibility = errors.New("permiso de Accesibilidad requerido: Ajustes del Sistema > Privacidad y seguridad > Accesibilidad")
)

// MaxScrollDelta acota cada eje de un evento de scroll (píxeles), para que un
// cliente roto no pueda mandar la página al otro extremo de un golpe.
const MaxScrollDelta = 2000

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

// Scroll desplaza la ventana bajo el puntero dx/dy píxeles (positivo = hacia
// arriba/izquierda, la convención de la rueda de CoreGraphics).
func (s *System) Scroll(_ context.Context, dx, dy int) error {
	dx, dy = clampDelta(dx), clampDelta(dy)
	if dx == 0 && dy == 0 {
		return nil
	}
	return postScroll(dx, dy)
}

func clampDelta(v int) int {
	return min(max(v, -MaxScrollDelta), MaxScrollDelta)
}
