package session

import "gitlab.com/iotTracker/brain/tracker/device/zx303"

type Session struct {
	LoggedIn    bool
	ZX303Device zx303.ZX303
}
