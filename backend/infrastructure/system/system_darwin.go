//go:build darwin

package system

/*
#cgo LDFLAGS: -framework CoreGraphics
#include <dlfcn.h>
#include <CoreGraphics/CoreGraphics.h>

typedef int (*ds_get_fn)(int, float *);
typedef int (*ds_set_fn)(int, float);

static int sys_brightness_get(CGDirectDisplayID display, float *out) {
    void *h = dlopen("/System/Library/PrivateFrameworks/DisplayServices.framework/DisplayServices", RTLD_NOW | RTLD_LOCAL);
    if (h == NULL) {
        return -1;
    }
    ds_get_fn get = (ds_get_fn)dlsym(h, "DisplayServicesGetBrightness");
    if (get == NULL) {
        return -1;
    }
    return get((int)display, out);
}

static int sys_brightness_set(CGDirectDisplayID display, float value) {
    void *h = dlopen("/System/Library/PrivateFrameworks/DisplayServices.framework/DisplayServices", RTLD_NOW | RTLD_LOCAL);
    if (h == NULL) {
        return -1;
    }
    ds_set_fn set = (ds_set_fn)dlsym(h, "DisplayServicesSetBrightness");
    if (set == NULL) {
        return -1;
    }
    return set((int)display, value);
}

typedef void (*mr_fn)(int, void *);

static int sys_media(int command) {
    void *h = dlopen("/System/Library/PrivateFrameworks/MediaRemote.framework/MediaRemote", RTLD_NOW | RTLD_LOCAL);
    if (h == NULL) {
        return -1;
    }
    mr_fn send = (mr_fn)dlsym(h, "MRMediaRemoteSendCommand");
    if (send == NULL) {
        return -1;
    }
    send(command, NULL);
    return 0;
}
*/
import "C"

import (
	"context"
	"fmt"
	"math"
	"os/exec"
	"strings"

	"master-control/model"
)

func executeToken(ctx context.Context, token string) error {
	switch token {
	case "volume_up":
		return runOSAScript(ctx, "set volume output volume ((output volume of (get volume settings)) + 10)")
	case "volume_down":
		return runOSAScript(ctx, "set volume output volume ((output volume of (get volume settings)) - 10)")
	case "volume_mute":
		return runOSAScript(ctx, "set volume output muted true")
	case "volume_unmute":
		return runOSAScript(ctx, "set volume output muted false")
	case "brightness_up":
		return adjustBrightness(0.1)
	case "brightness_down":
		return adjustBrightness(-0.1)
	case "media_play_pause":
		return mediaCommand(2)
	case "media_next":
		return mediaCommand(4)
	case "media_prev":
		return mediaCommand(5)
	default:
		return fmt.Errorf("%w: %s", ErrUnknownToken, token)
	}
}

func runOSAScript(ctx context.Context, script string) error {
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("osascript: %s: %w", string(out), err)
	}
	return nil
}

func adjustBrightness(delta float32) error {
	display := C.CGMainDisplayID()

	var current C.float
	if r := C.sys_brightness_get(display, &current); r != 0 {
		return fmt.Errorf("DisplayServicesGetBrightness failed: %d", r)
	}

	next := float32(math.Round(float64(float32(current)+delta)*100)) / 100
	if next < 0 {
		next = 0
	}
	if next > 1 {
		next = 1
	}

	if r := C.sys_brightness_set(display, C.float(next)); r != 0 {
		return fmt.Errorf("DisplayServicesSetBrightness failed: %d", r)
	}

	return nil
}

func mediaCommand(command C.int) error {
	if r := C.sys_media(command); r != 0 {
		return fmt.Errorf("MRMediaRemoteSendCommand failed: %d", r)
	}
	return nil
}

func getVolume(ctx context.Context) (model.SystemVolume, error) {
	cmd := exec.CommandContext(ctx, "osascript", "-e", "get volume settings")
	out, err := cmd.Output()
	if err != nil {
		return model.SystemVolume{}, err
	}

	vol := 0
	muted := false
	for _, part := range strings.Split(string(out), ",") {
		kv := strings.SplitN(strings.TrimSpace(part), ":", 2)
		if len(kv) != 2 {
			continue
		}
		switch strings.TrimSpace(kv[0]) {
		case "output volume":
			_, _ = fmt.Sscanf(strings.TrimSpace(kv[1]), "%d", &vol)
		case "output muted":
			muted = strings.TrimSpace(kv[1]) == "true"
		}
	}

	return model.SystemVolume{Volume: vol, Muted: muted}, nil
}

func getBrightness() (model.SystemBrightness, error) {
	display := C.CGMainDisplayID()

	var b C.float
	if r := C.sys_brightness_get(display, &b); r != 0 {
		return model.SystemBrightness{}, fmt.Errorf("DisplayServicesGetBrightness failed: %d", r)
	}

	return model.SystemBrightness{Brightness: float64(b)}, nil
}

func setVolume(ctx context.Context, volume int) error {
	return runOSAScript(ctx, fmt.Sprintf("set volume output volume %d", volume))
}

func setBrightness(_ context.Context, brightness float64) error {
	display := C.CGMainDisplayID()

	if r := C.sys_brightness_set(display, C.float(brightness)); r != 0 {
		return fmt.Errorf("DisplayServicesSetBrightness failed: %d", r)
	}

	return nil
}