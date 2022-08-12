//go:build !windows
// +build !windows

package systray

// #include "systray.h"
import "C"

import (
	"unsafe"
)

func registerSystray() error {
	C.registerSystray()
	return nil
}

func nativeLoop() {
	C.nativeLoop()
}

func quit() {
	C.quit()
}

// SetIcon sets the systray icon.
// iconBytes should be the content of .ico for windows and .ico/.jpg/.png
// for other platforms.
func SetIcon(iconBytes []byte) error {
	cstr := (*C.char)(unsafe.Pointer(&iconBytes[0]))
	C.setIcon(cstr, (C.int)(len(iconBytes)), false)
	return nil
}

// SetTitle sets the systray title, only available on Mac and Linux.
func SetTitle(title string) error {
	C.setTitle(C.CString(title))
	return nil
}

// SetTooltip sets the systray tooltip to display on mouse hover of the tray icon,
// only available on Mac and Windows.
func SetTooltip(tooltip string) error {
	C.setTooltip(C.CString(tooltip))
	return nil
}

func addOrUpdateMenuItem(item *MenuItem) error {
	var disabled C.short
	if item.disabled {
		disabled = 1
	}
	var checked C.short
	if item.checked {
		checked = 1
	}
	var isCheckable C.short
	if item.isCheckable {
		isCheckable = 1
	}
	var parentID uint32 = 0
	if item.parent != nil {
		parentID = item.parent.id
	}
	C.add_or_update_menu_item(
		C.int(item.id),
		C.int(parentID),
		C.CString(item.title),
		C.CString(item.tooltip),
		disabled,
		checked,
		isCheckable,
	)
	return nil
}

func addSeparator(id uint32) error {
	C.add_separator(C.int(id))
	return nil
}

func hideMenuItem(item *MenuItem) error {
	C.hide_menu_item(
		C.int(item.id),
	)
	return nil
}

func showMenuItem(item *MenuItem) error {
	C.show_menu_item(
		C.int(item.id),
	)
	return nil
}

//export systray_ready
func systray_ready() {
	systrayReady()
}

//export systray_on_exit
func systray_on_exit() {
	systrayExit()
}

//export systray_menu_item_selected
func systray_menu_item_selected(cID C.int) {
	systrayMenuItemSelected(uint32(cID))
}
