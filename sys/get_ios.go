//go:build ios

package sys

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <UiKit/UIScreen.h>

//float scale = 1;
//  if ([[UIScreen mainScreen] respondsToSelector:@selector(scale)]) {
//    scale = [[UIScreen mainScreen] scale];
//  }
//  float dpi;
//  if (UI_USER_INTERFACE_IDIOM() == UIUserInterfaceIdiomPad) {
//    dpi = 132 * scale;
//  } else if (UI_USER_INTERFACE_IDIOM() == UIUserInterfaceIdiomPhone) {
//    dpi = 163 * scale;
//  } else {
//    dpi = 160 * scale;
//  }

UIScreen* getPrimaryScreen() {
	UIScreen *screen = [UIScreen mainScreen];
	return screen;
}

CGSize getScreenDpi() {
	float scale = 1;
	UIScreen *screen = getPrimaryScreen();
	if ([screen respondsToSelector:@selector(scale)]) {
		scale = [screen scale];
	}
	float dpi;
	if ([UIDevice currentDevice].userInterfaceIdiom == UIUserInterfaceIdiomPad) {
		dpi = 132 * scale;
	} else if ([UIDevice currentDevice].userInterfaceIdiom == UIUserInterfaceIdiomPhone) {
		dpi = 163 * scale;
	} else {
		dpi = 160 * scale;
	}
	return CGSizeMake(dpi, dpi);
}

CGFloat getScreenBackingScale() {
	UIScreen *screen = getPrimaryScreen();
	CGFloat scale = [screen scale];
	return scale;
}

*/
import "C"

import (
	"errors"
	"fmt"
	"io/fs"
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
		case Scale:
			scale := float64(C.getScreenBackingScale())
			info.Scale = scale
		default:
			errs = append(errs, fmt.Errorf("platform.Get: Unknown key \"%v\"", key))
		}
	}
	return info, errors.Join(errs...)
}

// https://stackoverflow.com/questions/20401519/right-place-to-store-the-application-data-in-ios
// https://developer.apple.com/library/archive/documentation/FileManagement/Conceptual/FileSystemProgrammingGuide/FileSystemOverview/FileSystemOverview.html
// Backed-up data goes in
// ./Documents/
// Caches go in
// ./Library/Caches/
func makeAppDataPath() (Info, error) {
	// iOS is sandboxed.
	info := Info{}
	// Not a requirement for iOS. Possibly should
	// enforce it for consistency?
	//	if appName == "" {
	//		return info, fmt.Errorf("platform.Get: Missing app name, must first Set(SetAppName)")
	//	}
	// homeDir should be "." on iOS. Note that the full path
	// to this location is available at os.Getenv("HOME").
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return info, fmt.Errorf("platform.Get: %w", err)
	}

	// -- APP PATH
	info.AppPath = homeDir

	// -- APP DOCUMENTS PATH
	appDocsPath := filepath.Join(homeDir, "Documents")
	if _, err := os.Stat(appDocsPath); os.IsNotExist(err) {
		return info, fmt.Errorf("platform.Get: AppDocumentsPath err %w", err)
	}
	info.AppDocumentsPath = appDocsPath

	// -- APP CACHE PATH
	appCachePath := filepath.Join(homeDir, "Library", "Caches")
	if _, err := os.Stat(appDocsPath); os.IsNotExist(err) {
		return info, fmt.Errorf("platform.Get: AppCachePath err %w", err)
	}
	info.AppCachePath = appCachePath

	return info, nil
}

func walk(dir string) {
	walkFn := func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)
		if path == dir {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		return fs.SkipDir
	}
	filepath.WalkDir(dir, walkFn)
}

var appDataPathFn = sync.OnceValues(makeAppDataPath)
