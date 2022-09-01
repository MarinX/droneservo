package robot

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

type FakeDevice struct {
	errOn        error
	errStop      error
	countOn      int
	countPublish int
	countStop    int
}

func (f *FakeDevice) On(name string, fun func(s interface{})) (err error) {
	f.countOn++
	return f.errOn
}

func (f *FakeDevice) Publish(name string, data interface{}) {
	f.countPublish++
}

func (f *FakeDevice) Stop() error {
	f.countStop++
	return f.errStop
}

type FakeCtrl struct {
	count int
	err   error
}

func (f *FakeCtrl) Move(angle uint8) error {
	f.count++
	return f.err
}

func TestRobot(t *testing.T) {
	cases := []struct {
		it              string
		dev             *FakeDevice
		ctrl            *FakeCtrl
		expOnCalls      int
		expPublishCalls int
		expStopCalls    int
		expErrRegister  error
		expErrInit      error
		expErrStop      error
		handler         func(r *Robot)
	}{
		{
			it:           "run and set init angle",
			dev:          &FakeDevice{},
			ctrl:         &FakeCtrl{},
			expOnCalls:   1,
			expStopCalls: 1,
		},
		{
			it:              "publish and set angle",
			dev:             &FakeDevice{},
			ctrl:            &FakeCtrl{},
			expOnCalls:      1,
			expStopCalls:    1,
			expPublishCalls: 1,
			handler: func(r *Robot) {
				r.Publish("mock", nil)
			},
		},
		{
			it:              "publish multipe times",
			dev:             &FakeDevice{},
			ctrl:            &FakeCtrl{},
			expOnCalls:      1,
			expStopCalls:    1,
			expPublishCalls: 3,
			handler: func(r *Robot) {
				r.Publish("mock", nil)
				r.Publish("mock", nil)
				r.Publish("mock", nil)
			},
		},
		{
			it: "handle error on register",
			dev: &FakeDevice{
				errOn: errors.New("mock"),
			},
			ctrl:           &FakeCtrl{},
			expOnCalls:     1,
			expStopCalls:   1,
			expErrRegister: errors.New("mock"),
		},
		{
			it:  "handle error on init move",
			dev: &FakeDevice{},
			ctrl: &FakeCtrl{
				err: errors.New("mock"),
			},
			expOnCalls:   1,
			expStopCalls: 1,
			expErrInit:   errors.New("mock"),
		},
	}

	for _, c := range cases {
		t.Run(c.it, func(t *testing.T) {
			is := is.New(t)
			var err error

			robo := New(c.dev, c.ctrl, &ServoConfig{})

			err = robo.Register()
			is.Equal(c.expErrRegister, err)

			err = robo.Init()
			is.Equal(c.expErrInit, err)

			if c.handler != nil {
				c.handler(robo)
			}

			err = robo.Stop()
			is.Equal(c.expErrStop, err)

			is.Equal(c.expOnCalls, c.dev.countOn)
			is.Equal(c.expStopCalls, c.dev.countStop)
			is.Equal(c.expPublishCalls, c.dev.countPublish)

		})
	}

}

func TestRobotState(t *testing.T) {
	cases := []struct {
		it       string
		status   string
		expState uint8
	}{
		{
			it:       "set state to none",
			status:   "none",
			expState: 0,
		},
		{
			it:       "set state to running",
			status:   "running",
			expState: 1,
		},
		{
			it:       "set state to success",
			status:   "success",
			expState: 2,
		},
		{
			it:       "set state to failure",
			status:   "failure",
			expState: 3,
		},
		{
			it:       "no state set",
			status:   "bla",
			expState: 0,
		},
	}

	for _, c := range cases {
		t.Run(c.it, func(t *testing.T) {
			is := is.New(t)

			robo := New(&FakeDevice{}, &FakeCtrl{}, &ServoConfig{
				AngleNone:    0,
				AngleRunning: 1,
				AngleSuccess: 2,
				AngleFailure: 3,
			})

			robo.onMove(c.status)

			is.Equal(c.expState, robo.state.state)
		})
	}
}
