// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The govulncheckapi package is copied from x/vuln/internal/govulncheck
// and matches the output structure of govulncheck when ran in -json mode.
package govulncheckapi

import (
	"time"

	"golang.org/x/pkgsite-metrics/internal/osv"
)

// Message is an entry in the output stream. It will always have exactly one
// field filled in.
type Message struct {
	Config   *Config    `json:"config,omitempty"`
	Progress *Progress  `json:"progress,omitempty"`
	OSV      *osv.Entry `json:"osv,omitempty"`
	Finding  *Finding   `json:"finding,omitempty"`
}

// ProtocolVersion is the current protocol version this file implements
const ProtocolVersion = "v0.1.0"

type Config struct {
	// ProtocolVersion specifies the version of the JSON protocol.
	ProtocolVersion string `json:"protocol_version,omitempty"`

	// ScannerName is the name of the tool, for example, govulncheck.
	//
	// We expect this JSON format to be used by other tools that wrap
	// govulncheck, which will have a different name.
	ScannerName string `json:"scanner_name,omitempty"`

	// ScannerVersion is the version of the tool.
	ScannerVersion string `json:"scanner_version,omitempty"`

	// DB is the database used by the tool, for example,
	// vuln.go.dev.
	DB string `json:"db,omitempty"`

	// LastModified is the last modified time of the data source.
	DBLastModified *time.Time `json:"db_last_modified,omitempty"`

	// GoVersion is the version of Go used for analyzing standard library
	// vulnerabilities.
	GoVersion string `json:"go_version,omitempty"`

	// Consider only vulnerabilities that apply to this OS.
	GOOS string `json:"goos,omitempty"`

	// Consider only vulnerabilities that apply to this architecture.
	GOARCH string `json:"goarch,omitempty"`

	// ImportsOnly instructs vulncheck to analyze import chains only.
	// Otherwise, call chains are analyzed too.
	ImportsOnly bool `json:"imports_only,omitempty"`
}

type Progress struct {
	// A time stamp for the message.
	Timestamp *time.Time `json:"time,omitempty"`

	// Message is the progress message.
	Message string `json:"message,omitempty"`
}

// Vuln represents a single OSV entry.
type Finding struct {
	// OSV is the id of the detected vulnerability.
	OSV string `json:"osv,omitempty"`

	// FixedVersion is the module version where the vulnerability was
	// fixed. This is empty if a fix is not available.
	//
	// If there are multiple fixed versions in the OSV report, this will
	// be the fixed version in the latest range event for the OSV report.
	//
	// For example, if the range events are
	// {introduced: 0, fixed: 1.0.0} and {introduced: 1.1.0}, the fixed version
	// will be empty.
	//
	// For the stdlib, we will show the fixed version closest to the
	// Go version that is used. For example, if a fix is available in 1.17.5 and
	// 1.18.5, and the GOVERSION is 1.17.3, 1.17.5 will be returned as the
	// fixed version.
	FixedVersion string `json:"fixed_version,omitempty"`

	// Trace contains an entry for each frame in the trace.
	//
	// Frames are sorted starting from the imported vulnerable symbol
	// until the entry point. The first frame in Frames should match
	// Symbol.
	//
	// In binary mode, trace will contain a single-frame with no position
	// information.
	//
	// When a package is imported but no vulnerable symbol is called, the trace
	// will contain a single-frame with no symbol or position information.
	Trace []*Frame `json:"trace,omitempty"`
}

// Frame represents an entry in a finding trace.
type Frame struct {
	// Module is the module path of the module containing this symbol.
	//
	// Importable packages in the standard library will have the path "stdlib".
	Module string `json:"module"`

	// Version is the module version from the build graph.
	Version string `json:"version,omitempty"`

	// Package is the import path.
	Package string `json:"package,omitempty"`

	// Function is the function name.
	Function string `json:"function,omitempty"`

	// Receiver is the receiver type if the called symbol is a method.
	//
	// The client can create the final symbol name by
	// prepending Receiver to FuncName.
	Receiver string `json:"receiver,omitempty"`

	// Position describes an arbitrary source position
	// including the file, line, and column location.
	// A Position is valid if the line number is > 0.
	Position *Position `json:"position,omitempty"`
}

// Position is a copy of token.Position used to marshal/unmarshal
// JSON correctly.
type Position struct {
	Filename string `json:"filename,omitempty"` // filename, if any
	Offset   int    `json:"offset"`             // offset, starting at 0
	Line     int    `json:"line"`               // line number, starting at 1
	Column   int    `json:"column"`             // column number, starting at 1 (byte count)
}
