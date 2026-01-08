// Copyright (C) 2024-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package views provides view models for the Lux TUI.
package views

import (
	"time"
)

// NodeStatus represents a node's status
type NodeStatus struct {
	Name      string
	NodeID    string
	Status    string
	Version   string
	Uptime    string
	Connected int
	Staking   bool
}

// ChainStatus represents a chain's status
type ChainStatus struct {
	ChainID    string
	SubnetID   string
	Name       string
	Type       string
	Height     uint64
	Status     string
	Validators int
	TPS        float64
}

// ValidatorStatus represents a validator's status
type ValidatorStatus struct {
	NodeID    string
	Stake     string
	StartTime time.Time
	EndTime   time.Time
	Uptime    string
	Connected bool
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Source    string
	Message   string
}
