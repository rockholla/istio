//  Copyright 2019 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package galley

import (
	"testing"

	"istio.io/istio/pkg/test/framework2/core"
)

// Instance of Galley
type Instance interface {
	core.Resource

	// Address of the Galley MCP Server.
	Address() string

	// ApplyConfig applies the given config yaml file via Galley.
	ApplyConfig(ns core.Namespace, yamlText ...string) error

	// ApplyConfigOrFail applies the given config yaml file via Galley.
	ApplyConfigOrFail(t *testing.T, ns core.Namespace, yamlText ...string)

	// ApplyConfigDir recursively applies all the config files in the specified directory
	ApplyConfigDir(ns core.Namespace, configDir string) error

	// ClearConfig clears all applied config so far.
	ClearConfig() error

	// WaitForSnapshot waits until the given snapshot is observed for the given type URL.
	WaitForSnapshot(collection string, snapshot ...map[string]interface{}) error
}

// Configuration for Galley
type Config struct {
	// MeshConfig to use for this instance.
	MeshConfig string
}

// New returns a new instance of echo.
func New(ctx core.Context, cfg Config) (i Instance, err error) {
	err = core.UnsupportedEnvironment(ctx.Environment())
	ctx.Environment().Case(core.Native, func() {
		i, err = newNative(ctx, cfg)
	})
	ctx.Environment().Case(core.Kube, func() {
		i, err = newKube(ctx, cfg)
	})
	return
}

// NewOrFail returns a new Galley instance, or fails test.
func NewOrFail(t *testing.T, c core.Context, cfg Config) Instance {
	t.Helper()

	i, err := New(c, cfg)
	if err != nil {
		t.Fatalf("galley.NewOrFail: %v", err)
	}
	return i
}
