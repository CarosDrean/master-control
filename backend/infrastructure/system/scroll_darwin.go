//go:build darwin

package system

/*
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>

// Sin prompt: el diálogo de Accesibilidad ya lo dispara la captura de hotkeys;
// el scroll se llama decenas de veces por segundo y no debe abrir ventanas.
static int scroll_trusted(void) {
    return AXIsProcessTrusted() ? 1 : 0;
}

// Publica un evento de rueda en píxeles (scroll suave), como un trackpad.
// Va a la ventana que está bajo el puntero del ratón.
static int scroll_post(int32_t dy, int32_t dx) {
    CGEventRef ev = CGEventCreateScrollWheelEvent2(NULL, kCGScrollEventUnitPixel, 2, dy, dx, 0);
    if (ev == NULL) {
        return -1;
    }
    CGEventPost(kCGHIDEventTap, ev);
    CFRelease(ev);
    return 0;
}
*/
import "C"

import "fmt"

func postScroll(dx, dy int) error {
	if C.scroll_trusted() == 0 {
		return ErrNoAccessibility
	}
	if r := C.scroll_post(C.int32_t(dy), C.int32_t(dx)); r != 0 {
		return fmt.Errorf("CGEventCreateScrollWheelEvent2 failed: %d", r)
	}
	return nil
}
