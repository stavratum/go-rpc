package main

import (
	"encoding/binary"
	"strings"
	"syscall"
	"time"

	"github.com/winlabs/gowin32"
	"golang.org/x/sys/windows"

	"github.com/hugolgst/rich-go/client"
)

type User32 struct {
	*windows.LazyDLL

	getForegroundWindow *windows.LazyProc
}

func (module *User32) Load() {
	module.getForegroundWindow = module.NewProc("GetForegroundWindow")
}

func (module *User32) GetForegroundWindow() syscall.Handle {
	r, _, _ := module.getForegroundWindow.Call()
	return syscall.Handle(r)
}

func NewUser32() (user32 *User32) {
	user32 = &User32{
		LazyDLL: windows.NewLazyDLL("user32.dll"),
	}
	user32.Load()

	return
}

func main() {
	user32 := NewUser32()

	err := client.Login("1298882588911599617")
	if err != nil {
		panic(err)
	}
	defer client.Logout()

	for {
		pID, _ := gowin32.GetWindowProcessID(user32.GetForegroundWindow())

		// name
		fullpath, _ := gowin32.GetProcessFullPathName(pID, gowin32.ProcessNameNative)
		split := strings.Split(fullpath, "\\")

		name, _ := strings.CutSuffix(
			split[len(split)-1],
			".exe",
		)

		// rpc
		counters, _ := gowin32.GetProcessTimeCounters(pID)
		start := FiletimeToUnix(counters.Creation)

		err = client.SetActivity(client.Activity{
			State:      name,
			Details:    name,
			LargeImage: name,
			LargeText:  name,
			// SmallImage: name,
			// SmallText:  name,

			Timestamps: &client.Timestamps{
				Start: &start,
			},
		})

		if err != nil {
			panic(err)
		}

		time.Sleep(time.Minute)
	}
}

// that was annoying
func FiletimeToUnix(ft uint64) time.Time {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, ft)

	return time.Unix(0,
		(&syscall.Filetime{
			LowDateTime:  binary.LittleEndian.Uint32(buffer[:4]),
			HighDateTime: binary.LittleEndian.Uint32(buffer[4:]),
		}).Nanoseconds(),
	)
}
