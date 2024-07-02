//go:build android || ios
// +build android ios

package gui

import (
	"context"
	"errors"
	"net/url"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gocnpan/go2tv/devices"
	"github.com/gocnpan/go2tv/soapcalls"
	"github.com/gocnpan/go2tv/soapcalls/utils"
)

type deviceList struct {
	widget.List
}

func (c *deviceList) FocusGained() {}

func newDeviceList(dd *[]devType) *deviceList {
	list := &deviceList{}

	list.Length = func() int {
		return len(*dd)
	}

	list.CreateItem = func() fyne.CanvasObject {
		intListCont := container.NewHBox(widget.NewIcon(theme.NavigateNextIcon()), widget.NewLabel(""))
		return intListCont
	}

	list.UpdateItem = func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*fyne.Container).Objects[1].(*widget.Label).SetText((*dd)[i].name)
	}

	list.ExtendBaseWidget(list)
	return list
}

func mainWindow(s *NewScreen) fyne.CanvasObject {
	w := s.Current
	var data []devType
	list := newDeviceList(&data)

	go func() {
		var err error
		data, err = getDevices(1)
		if err != nil {
			data = nil
		}

		sort.Slice(data, func(i, j int) bool {
			return (data)[i].name < (data)[j].name
		})

		list.Refresh()
	}()

	mfiletext := widget.NewEntry()
	sfiletext := widget.NewEntry()

	mfile := widget.NewButton("Select Media File", func() {
		go mediaAction(s)
	})

	mfiletext.Disable()

	sfile := widget.NewButton("Select Subtitles File", func() {
		go subsAction(s)
	})

	sfiletext.Disable()

	playpause := widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {
		go playAction(s)
	})

	stop := widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), func() {
		go stopAction(s)
	})

	volumeup := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		go volumeAction(s, true)
	})

	muteunmute := widget.NewButtonWithIcon("", theme.VolumeUpIcon(), func() {
		go muteAction(s)
	})

	volumedown := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		go volumeAction(s, false)
	})

	clearmedia := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		go clearmediaAction(s)
	})

	clearsubs := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		go clearsubsAction(s)
	})

	externalmedia := widget.NewCheck("Media from URL", func(b bool) {})
	medialoop := widget.NewCheck("Loop Selected", func(b bool) {})

	mediafilelabel := canvas.NewText("File:", nil)
	subsfilelabel := canvas.NewText("Subtitles:", nil)
	devicelabel := canvas.NewText("Select Device:", nil)

	s.PlayPause = playpause
	s.Stop = stop
	s.MuteUnmute = muteunmute
	s.ExternalMediaURL = externalmedia
	s.MediaText = mfiletext
	s.SubsText = sfiletext
	s.DeviceList = list

	actionbuttons := container.New(&mainButtonsLayout{buttonHeight: 1.5, buttonPadding: theme.Padding()},
		playpause,
		volumedown,
		muteunmute,
		volumeup,
		stop)

	checklists := container.NewHBox(externalmedia, medialoop)
	mediasubsbuttons := container.New(layout.NewGridLayout(2), mfile, sfile)
	sfiletextArea := container.New(layout.NewBorderLayout(nil, nil, nil, clearsubs), clearsubs, sfiletext)
	mfiletextArea := container.New(layout.NewBorderLayout(nil, nil, nil, clearmedia), clearmedia, mfiletext)
	viewfilescont := container.New(layout.NewFormLayout(), mediafilelabel, mfiletextArea, subsfilelabel, sfiletextArea)
	buttons := container.NewVBox(mediasubsbuttons, viewfilescont, checklists, actionbuttons, container.NewPadded(devicelabel))
	content := container.New(layout.NewBorderLayout(buttons, nil, nil, nil), buttons, list)

	// Widgets actions
	list.OnSelected = func(id widget.ListItemID) {
		playpause.Enable()
		t, err := soapcalls.DMRextractor(context.Background(), data[id].addr)
		check(w, err)
		if err == nil {
			s.selectedDevice = data[id]
			s.controlURL = t.AvtransportControlURL
			s.eventlURL = t.AvtransportEventSubURL
			s.renderingControlURL = t.RenderingControlURL
			s.connectionManagerURL = t.ConnectionManagerURL
			if s.tvdata != nil {
				s.tvdata.RenderingControlURL = s.renderingControlURL
			}
		}
	}

	var mediafileOld fyne.URI
	var mediafileOldText string

	externalmedia.OnChanged = func(b bool) {
		if b {
			mfile.Disable()

			// rename the label
			mediafilelabel.Text = "URL:"
			mediafilelabel.Refresh()

			// keep old values
			mediafileOld = s.mediafile
			mediafileOldText = s.MediaText.Text

			// Clear the Media Text Area
			clearmediaAction(s)

			// Set some Media text defaults
			// to indicate that we're expecting a URL
			mfiletext.SetPlaceHolder("Enter URL here")
			mfiletext.Enable()
			return
		}

		medialoop.Enable()
		mfile.Enable()
		mediafilelabel.Text = "File:"
		mfiletext.SetPlaceHolder("")
		s.MediaText.Text = mediafileOldText
		s.mediafile = mediafileOld
		mediafilelabel.Refresh()
		mfiletext.Disable()
	}

	medialoop.OnChanged = func(b bool) {
		s.Medialoop = b
	}

	// Device list auto-refresh
	go refreshDevList(s, &data)

	// Check mute status for selected device
	go checkMutefunc(s)

	return content
}

func refreshDevList(s *NewScreen, data *[]devType) {
	refreshDevices := time.NewTicker(5 * time.Second)

	w := s.Current

	_, err := getDevices(2)
	if err != nil && !errors.Is(err, devices.ErrNoDeviceAvailable) {
		check(w, err)
	}

	for range refreshDevices.C {
		datanew, _ := getDevices(2)

	outer:
		for _, old := range *data {
			oldAddress, _ := url.Parse(old.addr)
			for _, new := range datanew {
				newAddress, _ := url.Parse(new.addr)
				if newAddress.Host == oldAddress.Host {
					continue outer
				}
			}

			if utils.HostPortIsAlive(oldAddress.Host) {
				datanew = append(datanew, old)
			}

			sort.Slice(datanew, func(i, j int) bool {
				return (datanew)[i].name < (datanew)[j].name
			})
		}

		// check to see if the new refresh includes
		// one of the already selected devices
		var includes bool
		u, _ := url.Parse(s.controlURL)
		for _, d := range datanew {
			n, _ := url.Parse(d.addr)
			if n.Host == u.Host {
				includes = true
			}
		}

		*data = datanew

		if !includes && !utils.HostPortIsAlive(u.Host) {
			s.controlURL = ""
			s.DeviceList.UnselectAll()
		}

		var found bool
		for n, a := range *data {
			if s.selectedDevice.addr == a.addr {
				found = true
				s.DeviceList.Select(n)
			}
		}

		if !found {
			s.DeviceList.UnselectAll()
		}

		s.DeviceList.Refresh()
	}
}

func checkMutefunc(s *NewScreen) {
	checkMute := time.NewTicker(1 * time.Second)

	var checkMuteCounter int
	for range checkMute.C {

		// Stop trying after 5 failures
		// to get the mute status
		if checkMuteCounter == 5 {
			s.renderingControlURL = ""
			checkMuteCounter = 0
		}

		if s.renderingControlURL == "" {
			continue
		}

		if s.tvdata == nil {
			s.tvdata = &soapcalls.TVPayload{RenderingControlURL: s.renderingControlURL}
		}

		isMuted, err := s.tvdata.GetMuteSoapCall()
		if err != nil {
			checkMuteCounter++
			continue
		}

		checkMuteCounter = 0

		switch isMuted {
		case "1":
			setMuteUnmuteView("Unmute", s)
		case "0":
			setMuteUnmuteView("Mute", s)
		}
	}
}
