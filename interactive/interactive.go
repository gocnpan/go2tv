package interactive

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gocnpan/go2tv/httphandlers"
	"github.com/gocnpan/go2tv/soapcalls"
)

type Controller struct {
	TV         *soapcalls.TVPayload
	httpServer *httphandlers.HTTPserver
	lastAction string
	mu         sync.RWMutex
	flipflop   bool
}

func NewController(tv *soapcalls.TVPayload, server *httphandlers.HTTPserver) *Controller {
	c := Controller{
		TV:         tv,
		httpServer: server,
		lastAction: "",
		flipflop:   true,
	}

	return &c
}

func (c *Controller) Init() error {
	time.Sleep(10 * time.Millisecond)
	// Sending the Play1 action sooner may result
	// in a panic error since we need to properly
	// initialize the tcell window.
	if err := c.TV.SendtoTV("Play1"); err != nil {
		err = fmt.Errorf("controller send play1 to tv error: %w", err)
		return err
	}
	return nil
}

const (
	ActStop      string = "Stop"
	ActP         string = "P" // play | pause
	ActPlay      string = "Play"
	ActPause     string = "Pause"
	ActM         string = "M"
	ActVolumeUp  string = "VolumeUp"
	ActVolumeDwn string = "VolumeDown"
)

func (c *Controller) HandleEvent(ev string) error {
	switch ev {
	case ActStop: // 终止
		c.Stop()
	case ActP: // 暂停或开始
		c.PlayOrPause()
	case ActM: // 静音或恢复声音
		c.MuteOrUn()
	case ActVolumeUp, ActVolumeDwn: // 音量调节
		c.Volume(ev)
	default:
		return fmt.Errorf("unknown dlna action event: %s", ev)
	}

	return nil
}

// 终止当前设备投屏播放
func (c *Controller) Stop() {
	c.TV.SendtoTV(ActStop)
	c.httpServer.StopServer()

	c.TV = nil
	c.httpServer = nil
	c.EmitMsg("Stopped")
}

// 开始播放或暂停播放
func (c *Controller) PlayOrPause() {
	if c.flipflop {
		c.flipflop = false
		c.TV.SendtoTV(ActPause)
		return
	}

	c.flipflop = true
	c.TV.SendtoTV(ActPlay)
}

// 静音或恢复声音
func (c *Controller) MuteOrUn() error {
	currentMute, err := c.TV.GetMuteSoapCall()
	if err != nil {
		return fmt.Errorf("controller get mute soap call error: %w", err)
	}

	m := "0"
	if currentMute == "0" {
		m = "1"
	}
	if err = c.TV.SetMuteSoapCall(m); err != nil {
		return fmt.Errorf("controller set mute soap call error: %w", err)
	}

	return nil
}

// 音量调节 VolumeUp | VolumeDown
func (c *Controller) Volume(upDown string) error {
	currentVolume, err := c.TV.GetVolumeSoapCall()
	if err != nil {
		return fmt.Errorf("controller get volume soap call error: %w", err)
	}
	setVolume := currentVolume - 1
	if upDown == ActVolumeUp {
		setVolume = currentVolume + 1
	}

	strVolume := strconv.Itoa(setVolume)
	if err = c.TV.SetVolumeSoapCall(strVolume); err != nil {
		return fmt.Errorf("controller set volume soap call error: %w", err)
	}

	return nil
}

// EmitMsg displays the actions to the interactive terminal.
// Method to implement the screen interface
func (c *Controller) EmitMsg(inputtext string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastAction = inputtext
}

// Fini Method to implement the screen interface
func (c *Controller) Fini() {
}
