// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/webhook"
	"github.com/matryer/is"
)

type FakeEvent struct {
	count int
}

func (f *FakeEvent) Publish(name string, data interface{}) {
	f.count++
}

func TestPlugin(t *testing.T) {
	cases := []struct {
		it  string
		in  *webhook.Request
		out int
	}{
		{
			it: "should publish event when build is not nil",
			in: &webhook.Request{
				Build: &drone.Build{},
			},
			out: 1,
		},
		{
			it: "should not publish event",
			in: &webhook.Request{},
		},
	}
	for _, c := range cases {
		t.Run(c.it, func(t *testing.T) {
			is := is.New(t)

			ev := &FakeEvent{}
			p := New(ev)
			p.Deliver(context.TODO(), c.in)
			is.Equal(ev.count, c.out)
		})
	}
}
