//go:build darwin && !ios

package sys

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Foundation

#include <sys/types.h>
#include <sys/sysctl.h>
//#import <sys/utsname.h>
#include <Foundation/NSString.h>
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

void printInfo() {
}

// NSString* deviceName() {
const char* deviceName() {
//	printInfo();

	size_t size;
	sysctlbyname("hw.model", NULL, &size, NULL, 0);
	char *model = malloc(size);
	sysctlbyname("hw.model", model, &size, NULL, 0);
	//NSLog(@"%s", model); // You would probably copy it to a NSString for later use
	//free(model);
	return model;

//	struct utsname systemInfo;
//    uname(&systemInfo);
//    int NSUTF8StringEncoding = 4;
//    NSString* str = [NSString stringWithCString:systemInfo.machine
//        encoding:NSUTF8StringEncoding];
//    return [str UTF8String];
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
	"unsafe"

	"github.com/hackborn/onefunc/math/geo"
)

func get(keys ...string) (Info, error) {
	info := Info{}
	errs := []error{}
	for _, key := range keys {
		switch key {
		case AppPath:
			paths, err := appDataPathFn()
			errs = append(errs, err)
			info.AppPath = paths.AppPath
		case AppDocumentsPath:
			paths, err := appDataPathFn()
			errs = append(errs, err)
			info.AppDocumentsPath = paths.AppDocumentsPath
		case AppCachePath:
			paths, err := appDataPathFn()
			errs = append(errs, err)
			info.AppCachePath = paths.AppCachePath
		case Dpi:
			size := C.getScreenDpi()
			if size.width == 0 || size.height == 0 {
				errs = append(errs, fmt.Errorf("platform.Get(Dpi): Invalid response"))
			}
			info.Dpi = geo.Pt(float64(size.width), float64(size.height))
		case HardwareModel:
			cString := C.deviceName()
			defer C.free(unsafe.Pointer(cString))
			info.HardwareModel = C.GoString(cString)
		case Scale:
			scale := float64(C.getScreenBackingScale())
			info.Scale = scale
		default:
			errs = append(errs, fmt.Errorf("platform.Get: Unknown key \"%v\"", key))
		}
	}
	return info, errors.Join(errs...)
}

func makeAppDataPath() (Info, error) {
	info := Info{}
	if appName == "" {
		return info, fmt.Errorf("platform.Get: Missing app name, must first Set(SetAppName)")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return info, fmt.Errorf("platform.Get: %w", err)
	}

	// -- APP PATH
	appPath := filepath.Join(homeDir, "Library", "Application Support", appName)
	err = os.Mkdir(appPath, 0700)
	if err != nil && !os.IsExist(err) {
		return info, fmt.Errorf("platform.Get: %w", err)
	}
	info.AppPath = appPath

	// -- APP DOCUMENTS PATH
	// TODO: This is compatibility with how I currently have it, but pretty sure
	// this is not where macOS should be storing documents.
	info.AppDocumentsPath = appPath

	// -- APP CACHE PATH
	// TODO: This is just made-up, clients might want it somewhere else.
	appCachePath := filepath.Join(homeDir, "Library", "Application Support", appName, "cache")
	err = os.Mkdir(appCachePath, 0700)
	if err != nil && !os.IsExist(err) {
		return info, fmt.Errorf("platform.Get: %w", err)
	}
	info.AppCachePath = appCachePath

	return info, nil
}

var appDataPathFn = sync.OnceValues(makeAppDataPath)
