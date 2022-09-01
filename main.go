// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/drone/drone-go/plugin/webhook"
	"github.com/marinx/droneservo/plugin"
	"github.com/marinx/droneservo/robot"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type spec struct {
	Bind   string `envconfig:"DRONE_BIND"`
	Debug  bool   `envconfig:"DRONE_DEBUG"`
	Secret string `envconfig:"DRONE_SECRET"`
	// Pin for rpi
	GPIOPin string `envconfig:"DRONE_GPIO_PIN"`
	// servo configuration
	ServoAngleNone uint8 `envconfig:"DRONE_SERVO_ANGLE_NONE"`
	ServoAngleRun  uint8 `envconfig:"DRONE_SERVO_ANGLE_RUN"`
	ServoAngleFail uint8 `envconfig:"DRONE_SERVO_ANGLE_FAIL"`
	ServoAnglePass uint8 `envconfig:"DRONE_SERVO_ANGLE_PASS"`
}

func main() {
	flag.Parse()

	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	robo, err := robot.Create(spec.GPIOPin, &robot.ServoConfig{
		AngleNone:    spec.ServoAngleNone,
		AngleRunning: spec.ServoAngleRun,
		AngleSuccess: spec.ServoAnglePass,
		AngleFailure: spec.ServoAngleFail,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	if err := robo.Register(); err != nil {
		logrus.Fatal(err)
	}
	// sets the servo to the default position
	if err := robo.Init(); err != nil {
		logrus.Fatal(err)
	}

	handler := webhook.Handler(
		plugin.New(
			robo,
		),
		spec.Secret,
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)
	http.Handle("/", handler)
	go func() {
		logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logrus.Info("shutting down server")
	if err := robo.Stop(); err != nil {
		logrus.Error(err)
	}
}
