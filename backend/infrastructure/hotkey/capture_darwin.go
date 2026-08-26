//go:build darwin

package hotkey

/*
#cgo LDFLAGS: -framework ApplicationServices -framework CoreFoundation
#include <Carbon/Carbon.h>
#include <ApplicationServices/ApplicationServices.h>
#include <CoreFoundation/CoreFoundation.h>

static int64_t g_keycode = -1;
static int64_t g_flags = 0;

static int hasAccessibilityPermission(void) {
    CFStringRef key = CFSTR("AXTrustedCheckOptionPrompt");
    CFDictionaryRef options = CFDictionaryCreate(
        kCFAllocatorDefault,
        (const void **)&key,
        (const void **)&kCFBooleanTrue,
        1,
        &kCFTypeDictionaryKeyCallBacks,
        &kCFTypeDictionaryValueCallBacks);
    Boolean trusted = AXIsProcessTrustedWithOptions(options);
    CFRelease(options);
    return trusted ? 1 : 0;
}

static CGEventRef tapCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
    (void)proxy;
    (void)refcon;
    if (type != kCGEventKeyDown) {
        return event;
    }

    CGKeyCode keycode = (CGKeyCode)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
    switch (keycode) {
        case 55: case 56: case 57: case 58: case 59: case 60: case 61: case 62: case 63:
            return event;
    }

    g_keycode = (int64_t)keycode;
    g_flags = (int64_t)CGEventGetFlags(event);
    CFRunLoopStop(CFRunLoopGetCurrent());
    return event;
}

// returns 0 on success, -1 if accessibility permission is missing, -2 on timeout
static int captureNextKey(int64_t *outKeycode, int64_t *outFlags, double timeout) {
    if (!hasAccessibilityPermission()) {
        return -1;
    }

    CFMachPortRef tap = CGEventTapCreate(
        kCGSessionEventTap,
        kCGHeadInsertEventTap,
        kCGEventTapOptionListenOnly,
        CGEventMaskBit(kCGEventKeyDown),
        tapCallback,
        NULL);
    if (tap == NULL) {
        return -1;
    }

    CFRunLoopSourceRef source = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, tap, 0);
    CFRunLoopAddSource(CFRunLoopGetCurrent(), source, kCFRunLoopCommonModes);
    CGEventTapEnable(tap, true);

    int result = -2;
    double elapsed = 0.0;
    while (elapsed < timeout) {
        CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0.1, false);
        if (g_keycode >= 0) {
            result = 0;
            break;
        }
        elapsed += 0.1;
    }

    CGEventTapEnable(tap, false);
    CFRunLoopRemoveSource(CFRunLoopGetCurrent(), source, kCFRunLoopCommonModes);
    CFRelease(source);
    CFRelease(tap);

    if (result == 0) {
        *outKeycode = g_keycode;
        *outFlags = g_flags;
        g_keycode = -1;
        g_flags = 0;
    }

    return result;
}
*/
import "C"

import (
	"fmt"
	"time"
)

func captureKey(timeout time.Duration) (int64, int64, error) {
	var keycode C.int64_t
	var flags C.int64_t

	code := C.captureNextKey(&keycode, &flags, C.double(timeout.Seconds()))
	switch code {
	case 0:
		return int64(keycode), int64(flags), nil
	case -1:
		return 0, 0, ErrNoPermission
	case -2:
		return 0, 0, ErrTimeout
	default:
		return 0, 0, fmt.Errorf("captureNextKey returned %d", code)
	}
}