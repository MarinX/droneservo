// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"

	"github.com/drone/drone-go/plugin/webhook"
)

type Eventer interface {
	Publish(name string, data interface{})
}

// New returns a new webhook extension.
func New(ev Eventer) webhook.Plugin {
	return &plugin{
		ev: ev,
	}
}

type plugin struct {
	ev Eventer
}

func (p *plugin) Deliver(ctx context.Context, req *webhook.Request) error {
	if req.Build != nil {
		p.ev.Publish("move", req.Build.Status)
	}
	return nil
}
