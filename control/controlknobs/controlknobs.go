// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

// Package controlknobs contains client options configurable from control which can be turned on
// or off. The ability to turn options on and off is for incrementally adding features in.
package controlknobs

import (
	"sync/atomic"

	"tailscale.com/syncs"
	"tailscale.com/tailcfg"
	"tailscale.com/types/opt"
)

// Knobs is the set of knobs that the control plane's coordination server can
// adjust at runtime.
type Knobs struct {
	// DisableUPnP indicates whether to attempt UPnP mapping.
	DisableUPnP atomic.Bool

	// DisableDRPO is whether control says to disable the
	// DERP route optimization (Issue 150).
	DisableDRPO atomic.Bool

	// KeepFullWGConfig is whether we should disable the lazy wireguard
	// programming and instead give WireGuard the full netmap always, even for
	// idle peers.
	KeepFullWGConfig atomic.Bool

	// RandomizeClientPort is whether control says we should randomize
	// the client port.
	RandomizeClientPort atomic.Bool

	// OneCGNAT is whether the the node should make one big CGNAT route
	// in the OS rather than one /32 per peer.
	OneCGNAT syncs.AtomicValue[opt.Bool]

	// ForceBackgroundSTUN forces netcheck STUN queries to keep
	// running in magicsock, even when idle.
	ForceBackgroundSTUN atomic.Bool

	// DisableDeltaUpdates is whether the node should not process
	// incremental (delta) netmap updates and should treat all netmap
	// changes as "full" ones as tailscaled did in 1.48.x and earlier.
	DisableDeltaUpdates atomic.Bool
}

// UpdateFromNodeAttributes updates k (if non-nil) based on the provided self
// node attributes (Node.Capabilities).
func (k *Knobs) UpdateFromNodeAttributes(selfNodeAttrs []string) {
	if k == nil {
		return
	}
	var (
		keepFullWG          bool
		disableDRPO         bool
		disableUPnP         bool
		randomizeClientPort bool
		disableDeltaUpdates bool
		oneCGNAT            opt.Bool
		forceBackgroundSTUN bool
	)
	for _, attr := range selfNodeAttrs {
		switch attr {
		case tailcfg.NodeAttrDebugDisableWGTrim:
			keepFullWG = true
		case tailcfg.NodeAttrDebugDisableDRPO:
			disableDRPO = true
		case tailcfg.NodeAttrDisableUPnP:
			disableUPnP = true
		case tailcfg.NodeAttrRandomizeClientPort:
			randomizeClientPort = true
		case tailcfg.NodeAttrOneCGNATEnable:
			oneCGNAT.Set(true)
		case tailcfg.NodeAttrOneCGNATDisable:
			oneCGNAT.Set(false)
		case tailcfg.NodeAttrDebugForceBackgroundSTUN:
			forceBackgroundSTUN = true
		case tailcfg.NodeAttrDisableDeltaUpdates:
			disableDeltaUpdates = true
		}
	}
	k.KeepFullWGConfig.Store(keepFullWG)
	k.DisableDRPO.Store(disableDRPO)
	k.DisableUPnP.Store(disableUPnP)
	k.RandomizeClientPort.Store(randomizeClientPort)
	k.OneCGNAT.Store(oneCGNAT)
	k.ForceBackgroundSTUN.Store(forceBackgroundSTUN)
	k.DisableDeltaUpdates.Store(disableDeltaUpdates)
}

// AsDebugJSON returns k as something that can be marshalled with json.Marshal
// for debug.
func (k *Knobs) AsDebugJSON() map[string]any {
	if k == nil {
		return nil
	}
	return map[string]any{
		"DisableUPnP":         k.DisableUPnP.Load(),
		"DisableDRPO":         k.DisableDRPO.Load(),
		"KeepFullWGConfig":    k.KeepFullWGConfig.Load(),
		"RandomizeClientPort": k.RandomizeClientPort.Load(),
		"OneCGNAT":            k.OneCGNAT.Load(),
		"ForceBackgroundSTUN": k.ForceBackgroundSTUN.Load(),
		"DisableDeltaUpdates": k.DisableDeltaUpdates.Load(),
	}
}
