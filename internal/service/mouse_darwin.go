//go:build darwin

package service

/*
#cgo LDFLAGS: -framework CoreGraphics

#include <CoreGraphics/CoreGraphics.h>
#include <stdlib.h>

CGPoint getMouseLocation() {
    CGEventSourceRef source = CGEventSourceCreate(kCGEventSourceStateCombinedSessionState);
    CGEventRef event = CGEventCreate(source);
    CGPoint loc = CGEventGetLocation(event);
    CFRelease(event);
    CFRelease(source);
    return loc;
}

void moveMouseTo(int x, int y) {
    CGEventSourceRef source = CGEventSourceCreate(kCGEventSourceStateCombinedSessionState);
    CGPoint loc = CGPointMake(x, y);
    CGEventRef event = CGEventCreateMouseEvent(source, kCGEventMouseMoved, loc, kCGMouseButtonLeft);
    CGEventPost(kCGHIDEventTap, event);
    CFRelease(event);
    CFRelease(source);
}
*/
import "C"

func init() {
	moveMouseDarwinImpl = func(offsetX, offsetY int) error {
		loc := C.getMouseLocation()
		newX := int(loc.x) + offsetX
		newY := int(loc.y) + offsetY
		C.moveMouseTo(C.int(newX), C.int(newY))
		return nil
	}
}
