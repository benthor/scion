// Copyright 2018 ETH Zurich, Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package env contains common command line and initialization code for SCION services.
// If something is specific to one app, it should go into that app's code and not here.
//
// During initialization, SIGHUPs are masked. To call a function on each
// SIGHUP, pass the function when calling Init.
package env

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/config"
	"github.com/scionproto/scion/go/lib/fatal"
	"github.com/scionproto/scion/go/lib/infra/modules/itopo"
	"github.com/scionproto/scion/go/lib/log"
	"github.com/scionproto/scion/go/lib/overlay"
	"github.com/scionproto/scion/go/lib/sciond"
	_ "github.com/scionproto/scion/go/lib/scrypto" // Make sure math/rand is seeded
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/topology"
	"github.com/scionproto/scion/go/lib/util"
)

const (
	DefaultLoggingLevel = "info"
	// Default max size of log files in MiB
	DefaultLoggingFileSize = 50
	// Default max age of log file in days
	DefaultLoggingFileMaxAge = 7
	// Default file name for topology file (only the last element of the path)
	DefaultTopologyPath = "topology.json"

	// SciondInitConnectPeriod is the default total amount of time spent
	// attempting to connect to sciond on start.
	SciondInitConnectPeriod = 20 * time.Second

	// ShutdownGraceInterval is the time applications wait after issuing a
	// clean shutdown signal, before forcerfully tearing down the application.
	ShutdownGraceInterval = 5 * time.Second
)

var sighupC chan os.Signal

func init() {
	os.Setenv("TZ", "UTC")
	sighupC = make(chan os.Signal, 1)
	signal.Notify(sighupC, syscall.SIGHUP)
}

var _ config.Config = (*General)(nil)

type General struct {
	// ID is the SCION element ID. This is used to choose the relevant
	// portion of the topology file for some services.
	ID string
	// ConfigDir for loading extra files (currently, only topology.json)
	ConfigDir string
	// Topology is the file path for the local topology JSON file.
	Topology string
	// ReconnectToDispatcher can be set to true to enable the snetproxy reconnecter.
	ReconnectToDispatcher bool
}

// InitDefaults sets the default value for Topology if not already set.
func (cfg *General) InitDefaults() {
	if cfg.Topology == "" {
		cfg.Topology = filepath.Join(cfg.ConfigDir, DefaultTopologyPath)
	}
}

func (cfg *General) Validate() error {
	if cfg.ID == "" {
		return common.NewBasicError("No element ID specified", nil)
	}
	return cfg.checkDir()
}

// checkDir checks that the config dir is a directory.
func (cfg *General) checkDir() error {
	if cfg.ConfigDir != "" {
		info, err := os.Stat(cfg.ConfigDir)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return common.NewBasicError(
				fmt.Sprintf("%v is not a directory", cfg.ConfigDir), nil)
		}
	}
	return nil
}

func (cfg *General) Sample(dst io.Writer, path config.Path, ctx config.CtxMap) {
	config.WriteString(dst, fmt.Sprintf(generalSample, ctx[config.ID]))
}

func (cfg *General) ConfigName() string {
	return "general"
}

var _ config.Config = (*SciondClient)(nil)

// SciondClient contains information for running snet with sciond.
type SciondClient struct {
	// Path is the sciond path. It defaults to sciond.DefaultSCIONDPath.
	Path string
	// InitialConnectPeriod is the maximum amount of time spent attempting to
	// connect to sciond on start.
	InitialConnectPeriod util.DurWrap
}

func (cfg *SciondClient) InitDefaults() {
	if cfg.Path == "" {
		cfg.Path = sciond.DefaultSCIONDPath
	}
	if cfg.InitialConnectPeriod.Duration == 0 {
		cfg.InitialConnectPeriod.Duration = SciondInitConnectPeriod
	}
}

func (cfg *SciondClient) Validate() error {
	if cfg.InitialConnectPeriod.Duration == 0 {
		return common.NewBasicError("InitialConnectPeriod must not be zero", nil)
	}
	return nil
}

func (cfg *SciondClient) Sample(dst io.Writer, path config.Path, _ config.CtxMap) {
	config.WriteString(dst, sciondClientSample)
}

func (cfg *SciondClient) ConfigName() string {
	return "sd_client"
}

// SetupEnv initializes a basic environment for applications. If reloadF is not
// nil, the application will call reloadF whenever it receives a SIGHUP signal.
func SetupEnv(reloadF func()) {
	setupSignals(reloadF)
}

// setupSignals sets up a goroutine that on received SIGTERM/SIGINT signals
// informs the application that it should shut down. If reloadF is not nil,
// setupSignals also calls reloadF on SIGHUP.
func setupSignals(reloadF func()) {
	fatal.Check()
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)
	go func() {
		defer log.LogPanicAndExit()
		s := <-sig
		log.Info("Received signal, exiting...", "signal", s)
		fatal.Shutdown(ShutdownGraceInterval)
	}()
	if reloadF != nil {
		go func() {
			defer log.LogPanicAndExit()
			for range sighupC {
				log.Info("Received config reload signal")
				reloadF()
			}
		}()
	}
}

func ReloadTopology(topologyPath string) {
	topo, err := topology.LoadFromFile(topologyPath)
	if err != nil {
		log.Error("Unable to reload topology", "err", err)
		return
	}
	if _, _, err := itopo.SetStatic(topo, true); err != nil {
		log.Error("Unable to set topology", "err", err)
		return
	}
	log.Info("Reloaded topology")
}

func GetPublicSnetAddress(ia addr.IA, topoAddr *topology.TopoAddr) *snet.Addr {
	// snet only supports udp4 for now
	if topoAddr.Overlay != overlay.UDPIPv4 {
		panic("unsupported overlay")
	}
	pub := topoAddr.PublicAddr(topoAddr.Overlay)
	if pub == nil {
		return nil
	}
	return &snet.Addr{IA: ia, Host: pub}
}

func GetBindSnetAddress(ia addr.IA, topoAddr *topology.TopoAddr) *snet.Addr {
	// snet only supports udp4 for now
	if topoAddr.Overlay != overlay.UDPIPv4 {
		panic("unsupported overlay")
	}
	bind := topoAddr.BindAddr(topoAddr.Overlay)
	if bind == nil {
		return nil
	}
	return &snet.Addr{IA: ia, Host: bind}
}

var _ config.Config = (*Metrics)(nil)

type Metrics struct {
	config.NoDefaulter
	config.NoValidator
	// Prometheus contains the address to export prometheus metrics on. If
	// not set, metrics are not exported.
	Prometheus string
}

func (cfg *Metrics) Sample(dst io.Writer, path config.Path, _ config.CtxMap) {
	config.WriteString(dst, metricsSample)
}

func (cfg *Metrics) ConfigName() string {
	return "metrics"
}

func (cfg *Metrics) StartPrometheus() {
	fatal.Check()
	if cfg.Prometheus != "" {
		http.Handle("/metrics", promhttp.Handler())
		log.Info("Exporting prometheus metrics", "addr", cfg.Prometheus)
		go func() {
			defer log.LogPanicAndExit()
			if err := http.ListenAndServe(cfg.Prometheus, nil); err != nil {
				fatal.Fatal(common.NewBasicError("HTTP ListenAndServe error", err))
			}
		}()
	}
}

// QUIC contains configuration for control-plane speakers.
type QUIC struct {
	ResolutionFraction float64
	Address            string
	CertFile           string
	KeyFile            string
}

func (cfg *QUIC) Sample(dst io.Writer, path config.Path, _ config.CtxMap) {
	config.WriteString(dst, quicSample)
}

func (cfg *QUIC) ConfigName() string {
	return "quic"
}
