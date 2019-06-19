package connstate

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/asecurityteam/settings"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestConnState_HandleEventAdd(t *testing.T) {
	type fields struct {
		Tracking *sync.Map
	}
	type args struct {
		state    http.ConnState
		statName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "new",
			fields: fields{
				Tracking: &sync.Map{},
			},
			args: args{
				state:    http.StateNew,
				statName: statCounterClientNew,
			},
		},
		{
			name: "active",
			fields: fields{
				Tracking: &sync.Map{},
			},
			args: args{
				state:    http.StateActive,
				statName: statCounterClientActive,
			},
		},
		{
			name: "idle",
			fields: fields{
				Tracking: &sync.Map{},
			},
			args: args{
				state:    http.StateIdle,
				statName: statCounterClientIdle,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctrl = gomock.NewController(t)
			defer ctrl.Finish()

			var stat = NewMockStat(ctrl)
			var conn = NewMockConn(ctrl)
			c := &ConnState{
				Stat:                      stat,
				Tracking:                  tt.fields.Tracking,
				NewClientCounterName:      statCounterClientNew,
				NewClientGaugeName:        statGaugeClientNew,
				ActiveClientCounterName:   statCounterClientActive,
				ActiveClientGaugeName:     statGaugeClientActive,
				IdleClientCounterName:     statCounterClientIdle,
				IdleClientGaugeName:       statGaugeClientIdle,
				ClosedClientCounterName:   statCounterClientClosed,
				HijackedClientCounterName: statCounterClientHijacked,
				stopCh:                    make(chan interface{}),
				statMut:                   &sync.Mutex{},
			}
			stat.EXPECT().Count(tt.args.statName, 1.0)
			c.HandleEvent(conn, tt.args.state)
			var _, ok = tt.fields.Tracking.Load(conn)
			if !ok {
				t.Error("did not store connection in map with state")
			}
		})
	}
}

func TestConnState_HandleEventRemove(t *testing.T) {
	type fields struct {
		Tracking *sync.Map
	}
	type args struct {
		state    http.ConnState
		statName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "closed",
			fields: fields{
				Tracking: &sync.Map{},
			},
			args: args{
				state:    http.StateClosed,
				statName: statCounterClientClosed,
			},
		},
		{
			name: "hijacked",
			fields: fields{
				Tracking: &sync.Map{},
			},
			args: args{
				state:    http.StateHijacked,
				statName: statCounterClientHijacked,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctrl = gomock.NewController(t)
			defer ctrl.Finish()

			var stat = NewMockStat(ctrl)
			var conn = NewMockConn(ctrl)
			c := &ConnState{
				Stat:                      stat,
				Tracking:                  tt.fields.Tracking,
				NewClientCounterName:      statCounterClientNew,
				NewClientGaugeName:        statGaugeClientNew,
				ActiveClientCounterName:   statCounterClientActive,
				ActiveClientGaugeName:     statGaugeClientActive,
				IdleClientCounterName:     statCounterClientIdle,
				IdleClientGaugeName:       statGaugeClientIdle,
				ClosedClientCounterName:   statCounterClientClosed,
				HijackedClientCounterName: statCounterClientHijacked,
				stopCh:                    make(chan interface{}),
				statMut:                   &sync.Mutex{},
			}
			c.Tracking.Store(conn, tt.args.state)
			stat.EXPECT().Count(tt.args.statName, 1.0)
			c.HandleEvent(conn, tt.args.state)
			var _, ok = tt.fields.Tracking.Load(conn)
			if ok {
				t.Error("did not remove connection from map when closed")
			}
		})
	}
}

func TestConnStateReport(t *testing.T) {
	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var stat = NewMockStat(ctrl)
	var smap = &sync.Map{}
	c := &ConnState{
		Stat:                      stat,
		Tracking:                  smap,
		NewClientCounterName:      statCounterClientNew,
		NewClientGaugeName:        statGaugeClientNew,
		ActiveClientCounterName:   statCounterClientActive,
		ActiveClientGaugeName:     statGaugeClientActive,
		IdleClientCounterName:     statCounterClientIdle,
		IdleClientGaugeName:       statGaugeClientIdle,
		ClosedClientCounterName:   statCounterClientClosed,
		HijackedClientCounterName: statCounterClientHijacked,
		stopCh:                    make(chan interface{}),
		statMut:                   &sync.Mutex{},
	}
	smap.Store(NewMockConn(ctrl), http.StateNew)
	smap.Store(NewMockConn(ctrl), http.StateActive)
	smap.Store(NewMockConn(ctrl), http.StateActive)
	smap.Store(NewMockConn(ctrl), http.StateIdle)
	smap.Store(NewMockConn(ctrl), http.StateIdle)
	smap.Store(NewMockConn(ctrl), http.StateIdle)

	stat.EXPECT().Gauge(statGaugeClientNew, 1.0)
	stat.EXPECT().Gauge(statGaugeClientActive, 2.0)
	stat.EXPECT().Gauge(statGaugeClientIdle, 3.0)
	c.report()
}

func TestConnStateReportLoop(t *testing.T) {
	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var stat = NewMockStat(ctrl)
	var smap = &sync.Map{}
	c := &ConnState{
		Stat:                      stat,
		Interval:                  time.Millisecond,
		Tracking:                  smap,
		NewClientCounterName:      statCounterClientNew,
		NewClientGaugeName:        statGaugeClientNew,
		ActiveClientCounterName:   statCounterClientActive,
		ActiveClientGaugeName:     statGaugeClientActive,
		IdleClientCounterName:     statCounterClientIdle,
		IdleClientGaugeName:       statGaugeClientIdle,
		ClosedClientCounterName:   statCounterClientClosed,
		HijackedClientCounterName: statCounterClientHijacked,
		stopCh:                    make(chan interface{}),
		statMut:                   &sync.Mutex{},
	}
	smap.Store(NewMockConn(ctrl), http.StateNew)
	smap.Store(NewMockConn(ctrl), http.StateActive)
	smap.Store(NewMockConn(ctrl), http.StateActive)
	smap.Store(NewMockConn(ctrl), http.StateIdle)
	smap.Store(NewMockConn(ctrl), http.StateIdle)
	smap.Store(NewMockConn(ctrl), http.StateIdle)

	stat.EXPECT().Gauge(statGaugeClientNew, 1.0).MinTimes(1)
	stat.EXPECT().Gauge(statGaugeClientActive, 2.0).MinTimes(1)
	stat.EXPECT().Gauge(statGaugeClientIdle, 3.0).MinTimes(1)
	go c.Report()
	time.Sleep(2 * time.Millisecond)
	c.Close()
}

func TestComponent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	src := settings.NewMapSource(map[string]interface{}{})
	s := NewMockStat(ctrl)
	cmp := NewComponent().WithStat(s)
	cs, err := Load(context.Background(), src, cmp)
	require.Nil(t, err)
	require.NotNil(t, cs)
}
