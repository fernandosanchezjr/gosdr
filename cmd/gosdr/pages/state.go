package pages

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/components"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/themes"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	log "github.com/sirupsen/logrus"
	"sync"
)

type State struct {
	mtx            sync.Mutex
	th             *themes.Theme
	sdrManager     *sdr.Manager
	deviceCards    map[devices.Id]layout.FlexChild
	connectButtons map[devices.Id]*widget.Clickable
	samplingMode   map[devices.Id]*widget.Enum
}

func NewState(th *themes.Theme, manager *sdr.Manager) *State {
	return &State{
		th:             th,
		sdrManager:     manager,
		deviceCards:    make(map[devices.Id]layout.FlexChild),
		connectButtons: make(map[devices.Id]*widget.Clickable),
		samplingMode:   make(map[devices.Id]*widget.Enum),
	}
}

func (s *State) RemoveDevice(id devices.Id) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	delete(s.deviceCards, id)
	delete(s.connectButtons, id)
}

func (s *State) AddDevice(id devices.Id) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, found := s.deviceCards[id]; !found {
		var device, deviceFound = s.sdrManager.GetInfo(id)
		if !deviceFound {
			log.WithFields(id.Fields()).Warn("AddDevice could not retrieve device")
			return
		}
		var connectButton = new(widget.Clickable)
		var samplingMode = new(widget.Enum)
		s.connectButtons[id] = connectButton
		s.samplingMode[id] = samplingMode
		s.deviceCards[id] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return components.DeviceCard(gtx, s.th, s.sdrManager, device, connectButton, samplingMode)
		})
	}
}

func (s *State) DeviceCards() (cards []layout.FlexChild) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for _, card := range s.deviceCards {
		cards = append(cards, card)
	}
	return
}
