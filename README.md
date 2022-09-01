A webhook extension for fun project using gobot, rpi and servo motor . _Please note this project requires Drone server version 1.4 or higher._

## Demo

// file

## Installation

Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```console
$ go build 
$ ./piservo
```

Set env variables (it's best first to check the angles using cli tool)

```sh
// set GPIO pin where servo motor is located on rpi
export DRONE_GPIO_PIN=17
// set the angle when the build is not running (default state)
export DRONE_SERVO_ANGLE_NONE=22
// set the angle when the build is running
export DRONE_SERVO_ANGLE_RUN=33
// set the angle when build is failed
export DRONE_SERVO_ANGLE_FAIL=45
// set the angle when build is success
export DRONE_SERVO_ANGLE_PASS=60

```

Update your Drone server configuration to include the plugin address and the shared secret.

```text
DRONE_WEBHOOK_ENDPOINT=http://1.2.3.4:3000
DRONE_WEBHOOK_SECRET=bea26a2221fd8090ea38720fc445eca6
```

## License

This software is licensed under the [Blue Oak Model License 1.0.0](https://spdx.org/licenses/BlueOak-1.0.0.html).
