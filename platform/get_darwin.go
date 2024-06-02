package platform

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <AppKit/NSScreen.h>

NSScreen* getPrimaryScreen() {
	// This is how some things do it, but as far as I know
	// the shorter way is safe and correct.
//	NSArray* screens = [NSScreen screens];
//	if ([screens count] < 1) {
//	    CGSize empty;
//	    return empty;
//	}

	NSScreen *screen = [NSScreen mainScreen];
	return screen;
}

CGSize getScreenDpi() {
	NSScreen *screen = getPrimaryScreen();
	NSDictionary *description = [screen deviceDescription];
	// NSSize displayPixelSize = [[description objectForKey:NSDeviceSize] sizeValue];
	NSSize res = [[description objectForKey:NSDeviceResolution] sizeValue]; // dpi
	// I guess NSSize and CGSize are interchangeable?
	return res;
}

CGFloat getScreenBackingScale() {
	NSScreen *screen = getPrimaryScreen();
	CGFloat scale = [screen backingScaleFactor];
	return scale;
}

// https://stackoverflow.com/questions/12589198/how-to-read-the-physical-screen-size-of-osx
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hackborn/onefunc/math/geo"
)

func get(keys ...string) (Info, error) {
	info := Info{}
	errs := []error{}
	for _, key := range keys {
		switch key {
		case AppDataPath:
			path, err := appDataPathFn()
			errs = append(errs, err)
			info.AppDataPath = path
		case Dpi:
			size := C.getScreenDpi()
			if size.width == 0 || size.height == 0 {
				errs = append(errs, fmt.Errorf("Key \"%v\" not valid", key))
			}
			info.Dpi = geo.Pt(float64(size.width), float64(size.height))
		case Scale:
			scale := float64(C.getScreenBackingScale())
			info.Scale = scale
		default:
			errs = append(errs, fmt.Errorf("Unknown key \"%v\"", key))
		}
	}
	return info, errors.Join(errs...)
}

func makeAppDataPath() (string, error) {
	if appName == "" {
		return "", fmt.Errorf("No app name for appdatapath, must first Set(SetAppName)")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(homeDir, "Library", "Application Support", appName)
	err = os.Mkdir(path, 0750)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	return path, nil
}

var appDataPathFn = sync.OnceValues(makeAppDataPath)
