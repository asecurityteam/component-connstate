package connstate

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	stat "github.com/asecurityteam/component-stat"
	"github.com/asecurityteam/settings"
)

const (
	statCounterClientNew      = "http.server.connstate.new"
	statGaugeClientNew        = "http.server.connstate.new.gauge"
	statCounterClientActive   = "http.server.connstate.active"
	statGaugeClientActive     = "http.server.connstate.active.gauge"
	statCounterClientIdle     = "http.server.connstate.idle"
	statGaugeClientIdle       = "http.server.connstate.idle.gauge"
	statCounterClientClosed   = "http.server.connstate.closed"
	statCounterClientHijacked = "http.server.connstate.hijacked"
	interval                  = 5 * time.Second
)

// ConnState plugs into the http.Server.ConnState attribute to track the number of
// client connections to the server.
type ConnState struct {
	Stat                      stat.Stat
	Tracking                  *sync.Map
	NewClientCounterName      string
	NewClientGaugeName        string
	ActiveClientCounterName   string
	ActiveClientGaugeName     string
	IdleClientCounterName     string
	IdleClientGaugeName       string
	ClosedClientCounterName   string
	HijackedClientCounterName string
	Interval                  time.Duration
	statMut                   *sync.Mutex
	stopCh                    chan interface{}
}

// Report loops on a time interval and pushes a set of gauge metrics.
func (c *ConnState) Report() {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.report()
		case <-c.stopCh:
			return
		}
	}
}

// Close the reporting loop.
func (c *ConnState) Close() error {
	close(c.stopCh)
	return nil
}

func (c *ConnState) report() {
	var n float64
	var a float64
	var i float64
	c.Tracking.Range(func(key interface{}, value interface{}) bool {
		switch value.(http.ConnState) {
		case http.StateNew:
			n = n + 1
		case http.StateActive:
			a = a + 1
		case http.StateIdle:
			i = i + 1
		}
		return true
	})
	c.statMut.Lock()
	defer c.statMut.Unlock()
	c.Stat.Gauge(c.NewClientGaugeName, n)
	c.Stat.Gauge(c.ActiveClientGaugeName, a)
	c.Stat.Gauge(c.IdleClientGaugeName, i)
}

// HandleEvent tracks state changes of a connection.
func (c *ConnState) HandleEvent(conn net.Conn, state http.ConnState) {
	c.statMut.Lock()
	defer c.statMut.Unlock()
	switch state {
	case http.StateNew:
		c.Stat.Count(c.NewClientCounterName, 1)
		c.Tracking.Store(conn, state)
	case http.StateActive:
		c.Stat.Count(c.ActiveClientCounterName, 1)
		c.Tracking.Store(conn, state)
	case http.StateIdle:
		c.Stat.Count(c.IdleClientCounterName, 1)
		c.Tracking.Store(conn, state)
	case http.StateHijacked:
		c.Stat.Count(c.HijackedClientCounterName, 1)
		c.Tracking.Delete(conn)
	case http.StateClosed:
		c.Stat.Count(c.ClosedClientCounterName, 1)
		c.Tracking.Delete(conn)
	}
}

// Config is a container for internal metrics settings.
type Config struct {
	NewCounter      string        `description:"Name of the counter metric tracking new clients."`
	NewGauge        string        `description:"Name of the gauge metric tracking new clients."`
	ActiveCounter   string        `description:"Name of the counter metric tracking active clients."`
	ActiveGauge     string        `description:"Name of the gauge metric tracking active clients."`
	IdleCounter     string        `description:"Name of the counter metric tracking idle clients."`
	IdleGauge       string        `description:"Name of the gauge metric tracking idle clients."`
	ClosedCounter   string        `description:"Name of the counter metric tracking closed clients."`
	HijackedCounter string        `description:"Name of the counter metric tracking hijacked clients."`
	ReportInterval  time.Duration `description:"Interval on which gauges are reported."`
}

// Name of the configuration root.
func (*Config) Name() string {
	return "connstate"
}

// Description returns the help information for the configuration root.
func (*Config) Description() string {
	return "Connection state metric names."
}

// Component implements the settings.Component interface for connection
// state monitoring.
type Component struct {
	Stat stat.Stat
}

// NewComponent populates the default values.
func NewComponent() *Component {
	return &Component{}
}

// WithStat returns a copy of the Component bound to a given Stat instance.
func (*Component) WithStat(s stat.Stat) *Component {
	return &Component{Stat: s}
}

// Settings returns a configuration with all defaults set.
func (*Component) Settings() *Config {
	return &Config{
		NewCounter:      statCounterClientNew,
		NewGauge:        statGaugeClientNew,
		ActiveCounter:   statCounterClientActive,
		ActiveGauge:     statGaugeClientActive,
		IdleCounter:     statCounterClientIdle,
		IdleGauge:       statGaugeClientIdle,
		ClosedCounter:   statCounterClientClosed,
		HijackedCounter: statCounterClientHijacked,
		ReportInterval:  interval,
	}
}

// New produces a ServerFn bound to the given configuration.
func (c *Component) New(_ context.Context, conf *Config) (*ConnState, error) {
	return &ConnState{
		Stat:                      c.Stat,
		Tracking:                  &sync.Map{},
		NewClientCounterName:      conf.NewCounter,
		NewClientGaugeName:        conf.NewGauge,
		ActiveClientCounterName:   conf.ActiveCounter,
		ActiveClientGaugeName:     conf.ActiveGauge,
		IdleClientCounterName:     conf.IdleCounter,
		IdleClientGaugeName:       conf.IdleGauge,
		ClosedClientCounterName:   conf.ClosedCounter,
		HijackedClientCounterName: conf.HijackedCounter,
		Interval:                  conf.ReportInterval,
		statMut:                   &sync.Mutex{},
		stopCh:                    make(chan interface{}),
	}, nil
}

// Load is a convenience method for binding the source to the component.
func Load(ctx context.Context, source settings.Source, c *Component) (*ConnState, error) {
	dst := new(ConnState)
	err := settings.NewComponent(ctx, source, c, dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
