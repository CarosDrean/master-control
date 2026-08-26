package hotkey

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"master-control/model"
)

var (
	ErrNoPermission = errors.New("permiso de Accesibilidad requerido: Ajustes del Sistema > Privacidad y seguridad > Accesibilidad")
	ErrTimeout      = errors.New("se agotó el tiempo de espera de captura")
	ErrNotSupported = errors.New("la captura nativa solo está disponible en macOS")
)

const defaultTimeout = 15 * time.Second

func Capture(timeout time.Duration) (model.Hotkey, error) {
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	keycode, flags, err := captureKey(timeout)
	if err != nil {
		return model.Hotkey{}, err
	}

	mods := modifiersFromFlags(flags)

	return model.Hotkey{
		Keycode:   keycode,
		Modifiers: mods,
		Command:   buildCommand(keycode, mods),
	}, nil
}

func modifiersFromFlags(flags int64) []string {
	var mods []string
	if flags&(1<<17) != 0 {
		mods = append(mods, "shift")
	}
	if flags&(1<<18) != 0 {
		mods = append(mods, "control")
	}
	if flags&(1<<19) != 0 {
		mods = append(mods, "alt")
	}
	if flags&(1<<20) != 0 {
		mods = append(mods, "meta")
	}
	return mods
}

func buildCommand(keycode int64, mods []string) string {
	names := map[string]string{
		"shift":   "shift",
		"control": "control",
		"alt":     "option",
		"meta":    "command",
	}

	var using []string
	for _, m := range []string{"shift", "control", "alt", "meta"} {
		if contains(mods, m) {
			using = append(using, names[m]+" down")
		}
	}

	clause := ""
	if len(using) > 0 {
		clause = " using {" + strings.Join(using, ", ") + "}"
	}

	return fmt.Sprintf("tell application \"System Events\"\n\tkey code %d%s\nend tell", keycode, clause)
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}