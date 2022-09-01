package robot

import (
	"sync"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

// GoRobot is the interface for the gobot.Robot
type GoRobot interface {
	On(name string, f func(s interface{})) (err error)
	Publish(name string, data interface{})
	Stop() error
}

type ServoController interface {
	Move(angle uint8) (err error)
}

// Robot is a gobot wrapper
type Robot struct {
	robo        GoRobot
	servo       ServoController
	servoMapper map[string]uint8
	state       servoState
}

// ServoConfig is the configuration for the servo
type ServoConfig struct {
	AngleNone    uint8 // when the build is not running
	AngleRunning uint8 // when the build is running
	AngleSuccess uint8 // when the build is successful
	AngleFailure uint8 // when the build is failed
}

type servoState struct {
	sync.Mutex
	state uint8
}

// New returns a new robot
func New(robo GoRobot, controller ServoController, config *ServoConfig) *Robot {
	return &Robot{
		robo:  robo,
		servo: controller,
		servoMapper: map[string]uint8{
			"none":    config.AngleNone,
			"running": config.AngleRunning,
			"success": config.AngleSuccess,
			"failure": config.AngleFailure,
		},
		state: servoState{},
	}
}

// Create creates a new robot from gobot
func Create(gpioPin string, config *ServoConfig) (*Robot, error) {
	r := raspi.NewAdaptor()
	servo := gpio.NewServoDriver(r, gpioPin)
	robot := gobot.NewRobot("bot",
		[]gobot.Connection{r},
		[]gobot.Device{servo},
	)
	if err := robot.Start(false); err != nil {
		return nil, err
	}
	return New(robot, servo, config), nil
}

// Register registers the robot to listen for the move event
func (r *Robot) Register() error {
	return r.robo.On("move", r.onMove)
}

func (r *Robot) onMove(s interface{}) {
	if status, ok := s.(string); ok {
		rotation := r.servoMapper[status]
		if rotation == 0 {
			return
		}

		r.state.Lock()
		if r.state.state == rotation {
			r.state.Unlock()
			return
		}
		r.state.state = rotation
		r.state.Unlock()

		// we are finished, move to the none position
		if rotation == r.servoMapper["success"] || rotation == r.servoMapper["failure"] {
			go func() {
				time.Sleep(3 * time.Second)
				r.Init()
			}()
		}
		r.servo.Move(rotation)
	}
}

// Init moves to robot to the none position
func (r *Robot) Init() error {
	r.state.Lock()
	defer r.state.Unlock()
	r.state.state = r.servoMapper["none"]
	return r.servo.Move(r.state.state)
}

// Publish publishes a message to the robot
func (r *Robot) Publish(name string, data interface{}) {
	r.robo.Publish(name, data)
}

// Stop stops the robot (closes connection and device)
func (r *Robot) Stop() error {
	return r.robo.Stop()
}
