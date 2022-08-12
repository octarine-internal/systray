package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/skratchdot/open-golang/open"
)

func main() {
	onExit := func() {
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Lantern")
	mQuitOrig, err := systray.AddMenuItem("Quit", "Quit the whole app")
	if err != nil {
		panic(err)
	}
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		if err = systray.SetTitle("Awesome App"); err != nil {
			panic(err)
		}
		if err = systray.SetTooltip("Pretty awesome棒棒嗒"); err != nil {
			panic(err)
		}
		mChange, err := systray.AddMenuItem("Change Me", "Change Me")
		if err != nil {
			panic(err)
		}
		mChecked, err := systray.AddMenuItemCheckbox("Unchecked", "Check Me", true)
		if err != nil {
			panic(err)
		}
		mEnabled, err := systray.AddMenuItem("Enabled", "Enabled")
		if err != nil {
			panic(err)
		}
		// Sets the icon of a menu item. Only available on Mac.
		mEnabled.SetTemplateIcon(icon.Data, icon.Data)

		systray.AddMenuItem("Ignored", "Ignored")

		subMenuTop, err := systray.AddMenuItem("SubMenuTop", "SubMenu Test (top)")
		if err != nil {
			panic(err)
		}
		subMenuMiddle, err := subMenuTop.AddSubMenuItem("SubMenuMiddle", "SubMenu Test (middle)")
		if err != nil {
			panic(err)
		}
		subMenuBottom, err := subMenuMiddle.AddSubMenuItemCheckbox("SubMenuBottom - Toggle Panic!", "SubMenu Test (bottom) - Hide/Show Panic!", false)
		if err != nil {
			panic(err)
		}
		subMenuBottom2, err := subMenuMiddle.AddSubMenuItem("SubMenuBottom - Panic!", "SubMenu Test (bottom)")
		if err != nil {
			panic(err)
		}

		mUrl, err := systray.AddMenuItem("Open UI", "my home")
		if err != nil {
			panic(err)
		}
		mQuit, err := systray.AddMenuItem("退出", "Quit the whole app")
		if err != nil {
			panic(err)
		}

		// Sets the icon of a menu item. Only available on Mac.
		mQuit.SetIcon(icon.Data)

		systray.AddSeparator()
		mToggle, err := systray.AddMenuItem("Toggle", "Toggle the Quit button")
		if err != nil {
			panic(err)
		}
		shown := true
		toggle := func() {
			if shown {
				subMenuBottom.Check()
				subMenuBottom2.Hide()
				mQuitOrig.Hide()
				mEnabled.Hide()
				shown = false
			} else {
				subMenuBottom.Uncheck()
				subMenuBottom2.Show()
				mQuitOrig.Show()
				mEnabled.Show()
				shown = true
			}
		}

		for {
			select {
			case <-mChange.ClickedCh:
				mChange.SetTitle("I've Changed")
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					mChecked.Uncheck()
					mChecked.SetTitle("Unchecked")
				} else {
					mChecked.Check()
					mChecked.SetTitle("Checked")
				}
			case <-mEnabled.ClickedCh:
				mEnabled.SetTitle("Disabled")
				mEnabled.Disable()
			case <-mUrl.ClickedCh:
				open.Run("https://www.getlantern.org")
			case <-subMenuBottom2.ClickedCh:
				panic("panic button pressed")
			case <-subMenuBottom.ClickedCh:
				toggle()
			case <-mToggle.ClickedCh:
				toggle()
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}
