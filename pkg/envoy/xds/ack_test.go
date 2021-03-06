// Copyright 2018 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !privileged_tests

package xds

import (
	"context"
	"time"

	"github.com/cilium/cilium/pkg/completion"
	envoy_api_v2_core "github.com/cilium/cilium/pkg/envoy/envoy/api/v2/core"

	. "gopkg.in/check.v1"
)

type AckSuite struct{}

var _ = Suite(&AckSuite{})

const (
	node0 = "10.0.0.0"
	node1 = "10.0.0.1"
	node2 = "10.0.0.2"
)

var nodes = map[string]*envoy_api_v2_core.Node{
	node0: {Id: "sidecar~10.0.0.0~node0~bar"},
	node1: {Id: "sidecar~10.0.0.1~node1~bar"},
	node2: {Id: "sidecar~10.0.0.2~node2~bar"},
}

// IsCompletedChecker checks that a Completion is completed without errors.
type IsCompletedChecker struct {
	*CheckerInfo
}

func (c *IsCompletedChecker) Check(params []interface{}, names []string) (result bool, error string) {
	comp, ok := params[0].(*completion.Completion)
	if !ok {
		return false, "completion must be a *completion.Completion"
	}
	if comp == nil {
		return false, "completion is nil"
	}

	select {
	case <-comp.Completed():
		return comp.Err() == nil, ""
	default:
		return false, ""
	}
}

// IsCompleted checks that a Completion is completed.
var IsCompleted Checker = &IsCompletedChecker{
	&CheckerInfo{Name: "IsCompleted", Params: []string{
		"completion"}},
}

func (s *AckSuite) TestUpsertSingleNode(c *C) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	typeURL := "type.googleapis.com/envoy.api.v2.DummyConfiguration"
	wg := completion.NewWaitGroup(ctx)

	// Empty cache is the version 1
	cache := NewCache()
	acker := NewAckingResourceMutatorWrapper(cache, IstioNodeToIP)

	// Create version 2 with resource 0.
	comp := wg.AddCompletion()
	defer comp.Complete(nil)
	acker.Upsert(typeURL, resources[0].Name, resources[0], []string{node0}, comp)
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for the right resource, from another node.
	acker.HandleResourceVersionAck(2, 2, nodes[node1], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for another resource, from the right node.
	acker.HandleResourceVersionAck(2, 2, nodes[node0], []string{resources[1].Name}, typeURL, "")
	c.Assert(comp, Not(IsCompleted))

	// Ack an older version, for the right resource, from the right node.
	acker.HandleResourceVersionAck(1, 1, nodes[node0], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for the right resource, from the right node.
	acker.HandleResourceVersionAck(2, 2, nodes[node0], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, IsCompleted)
}

func (s *AckSuite) TestUpsertMultipleNodes(c *C) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	typeURL := "type.googleapis.com/envoy.api.v2.DummyConfiguration"
	wg := completion.NewWaitGroup(ctx)

	// Empty cache is the version 1
	cache := NewCache()
	acker := NewAckingResourceMutatorWrapper(cache, IstioNodeToIP)

	// Create version 2 with resource 0.
	comp := wg.AddCompletion()
	defer comp.Complete(nil)
	acker.Upsert(typeURL, resources[0].Name, resources[0], []string{node0, node1}, comp)
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for the right resource, from another node.
	acker.HandleResourceVersionAck(2, 2, nodes[node2], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for the right resource, from one of the nodes (node0).
	// One of the nodes (node1) still needs to ACK.
	acker.HandleResourceVersionAck(2, 2, nodes[node0], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for the right resource, from the last remaining node (node1).
	acker.HandleResourceVersionAck(2, 2, nodes[node1], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, IsCompleted)
}

func (s *AckSuite) TestUpsertMoreRecentVersion(c *C) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	typeURL := "type.googleapis.com/envoy.api.v2.DummyConfiguration"
	wg := completion.NewWaitGroup(ctx)

	// Empty cache is the version 1
	cache := NewCache()
	acker := NewAckingResourceMutatorWrapper(cache, IstioNodeToIP)

	// Create version 2 with resource 0.
	comp := wg.AddCompletion()
	defer comp.Complete(nil)
	acker.Upsert(typeURL, resources[0].Name, resources[0], []string{node0}, comp)
	c.Assert(comp, Not(IsCompleted))

	// Ack an older version, for the right resource, from the right node.
	acker.HandleResourceVersionAck(1, 1, nodes[node0], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, Not(IsCompleted))

	// Ack a more recent version, for the right resource, from the right node.
	acker.HandleResourceVersionAck(123, 123, nodes[node0], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, IsCompleted)
}

func (s *AckSuite) TestUpsertMoreRecentVersionNack(c *C) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	typeURL := "type.googleapis.com/envoy.api.v2.DummyConfiguration"
	wg := completion.NewWaitGroup(ctx)

	// Empty cache is the version 1
	cache := NewCache()
	acker := NewAckingResourceMutatorWrapper(cache, IstioNodeToIP)

	// Create version 2 with resource 0.
	comp := wg.AddCompletion()
	defer comp.Complete(nil)
	acker.Upsert(typeURL, resources[0].Name, resources[0], []string{node0}, comp)
	c.Assert(comp, Not(IsCompleted))

	// Ack an older version, for the right resource, from the right node.
	acker.HandleResourceVersionAck(1, 1, nodes[node0], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, Not(IsCompleted))

	// NAck a more recent version, for the right resource, from the right node.
	acker.HandleResourceVersionAck(1, 2, nodes[node0], []string{resources[0].Name}, typeURL, "Detail")
	// IsCompleted is true only for completions without error
	c.Assert(comp, Not(IsCompleted))
	c.Assert(comp.Err(), Not(Equals), nil)
	c.Assert(comp.Err(), DeepEquals, &ProxyError{Err: ErrNackReceived, Detail: "Detail"})
}

func (s *AckSuite) TestDeleteSingleNode(c *C) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	typeURL := "type.googleapis.com/envoy.api.v2.DummyConfiguration"
	wg := completion.NewWaitGroup(ctx)

	// Empty cache is the version 1
	cache := NewCache()
	acker := NewAckingResourceMutatorWrapper(cache, IstioNodeToIP)

	// Create version 2 with resource 0.
	comp := wg.AddCompletion()
	defer comp.Complete(nil)
	acker.Upsert(typeURL, resources[0].Name, resources[0], []string{node0}, comp)
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for the right resource, from the right node.
	acker.HandleResourceVersionAck(2, 2, nodes[node0], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, IsCompleted)

	// Create version 3 with no resources.
	comp = wg.AddCompletion()
	defer comp.Complete(nil)
	acker.Delete(typeURL, resources[0].Name, []string{node0}, comp)
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for another resource, from another node.
	acker.HandleResourceVersionAck(3, 3, nodes[node1], []string{resources[2].Name}, typeURL, "")
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for another resource, from the right node.
	acker.HandleResourceVersionAck(3, 3, nodes[node0], []string{resources[2].Name}, typeURL, "")
	// The resource name is ignored. For delete, we only consider the version.
	c.Assert(comp, IsCompleted)
}

func (s *AckSuite) TestDeleteMultipleNodes(c *C) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	typeURL := "type.googleapis.com/envoy.api.v2.DummyConfiguration"
	wg := completion.NewWaitGroup(ctx)

	// Empty cache is the version 1
	cache := NewCache()
	acker := NewAckingResourceMutatorWrapper(cache, IstioNodeToIP)

	// Create version 2 with resource 0.
	comp := wg.AddCompletion()
	defer comp.Complete(nil)
	acker.Upsert(typeURL, resources[0].Name, resources[0], []string{node0}, comp)
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for the right resource, from the right node.
	acker.HandleResourceVersionAck(2, 2, nodes[node0], []string{resources[0].Name}, typeURL, "")
	c.Assert(comp, IsCompleted)

	// Create version 3 with no resources.
	comp = wg.AddCompletion()
	defer comp.Complete(nil)
	acker.Delete(typeURL, resources[0].Name, []string{node0, node1}, comp)
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for another resource, from one of the nodes.
	acker.HandleResourceVersionAck(3, 3, nodes[node1], []string{resources[2].Name}, typeURL, "")
	c.Assert(comp, Not(IsCompleted))

	// Ack the right version, for another resource, from the remaining node.
	acker.HandleResourceVersionAck(3, 3, nodes[node0], []string{resources[2].Name}, typeURL, "")
	// The resource name is ignored. For delete, we only consider the version.
	c.Assert(comp, IsCompleted)
}

func (s *AckSuite) TestRevertInsert(c *C) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	typeURL := "type.googleapis.com/envoy.api.v2.DummyConfiguration"
	wg := completion.NewWaitGroup(ctx)

	cache := NewCache()
	acker := NewAckingResourceMutatorWrapper(cache, IstioNodeToIP)

	// Create version 1 with resource 0.
	// Insert.
	revert := acker.Upsert(typeURL, resources[0].Name, resources[0], []string{node0}, nil)

	// Insert another resource.
	_ = acker.Upsert(typeURL, resources[2].Name, resources[2], []string{node0}, nil)

	res, err := cache.Lookup(typeURL, resources[0].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[0])

	res, err = cache.Lookup(typeURL, resources[2].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[2])

	comp := wg.AddCompletion()
	defer comp.Complete(nil)
	revert(comp)

	res, err = cache.Lookup(typeURL, resources[0].Name)
	c.Assert(err, IsNil)
	c.Assert(res, IsNil)

	res, err = cache.Lookup(typeURL, resources[2].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[2])
}

func (s *AckSuite) TestRevertUpdate(c *C) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	typeURL := "type.googleapis.com/envoy.api.v2.DummyConfiguration"
	wg := completion.NewWaitGroup(ctx)

	cache := NewCache()
	acker := NewAckingResourceMutatorWrapper(cache, IstioNodeToIP)

	// Create version 1 with resource 0.
	// Insert.
	acker.Upsert(typeURL, resources[0].Name, resources[0], []string{node0}, nil)

	// Insert another resource.
	_ = acker.Upsert(typeURL, resources[2].Name, resources[2], []string{node0}, nil)

	res, err := cache.Lookup(typeURL, resources[0].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[0])

	res, err = cache.Lookup(typeURL, resources[2].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[2])

	// Update.
	revert := acker.Upsert(typeURL, resources[0].Name, resources[1], []string{node0}, nil)

	res, err = cache.Lookup(typeURL, resources[0].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[1])

	comp := wg.AddCompletion()
	defer comp.Complete(nil)
	revert(comp)

	res, err = cache.Lookup(typeURL, resources[0].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[0])

	res, err = cache.Lookup(typeURL, resources[2].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[2])
}

func (s *AckSuite) TestRevertDelete(c *C) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	typeURL := "type.googleapis.com/envoy.api.v2.DummyConfiguration"
	wg := completion.NewWaitGroup(ctx)

	cache := NewCache()
	acker := NewAckingResourceMutatorWrapper(cache, IstioNodeToIP)

	// Create version 1 with resource 0.
	// Insert.
	acker.Upsert(typeURL, resources[0].Name, resources[0], []string{node0}, nil)

	// Insert another resource.
	_ = acker.Upsert(typeURL, resources[2].Name, resources[2], []string{node0}, nil)

	res, err := cache.Lookup(typeURL, resources[0].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[0])

	res, err = cache.Lookup(typeURL, resources[2].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[2])

	// Delete.
	revert := acker.Delete(typeURL, resources[0].Name, []string{node0}, nil)

	res, err = cache.Lookup(typeURL, resources[0].Name)
	c.Assert(err, IsNil)
	c.Assert(res, IsNil)

	res, err = cache.Lookup(typeURL, resources[2].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[2])

	comp := wg.AddCompletion()
	defer comp.Complete(nil)
	revert(comp)

	res, err = cache.Lookup(typeURL, resources[0].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[0])

	res, err = cache.Lookup(typeURL, resources[2].Name)
	c.Assert(err, IsNil)
	c.Assert(res, Equals, resources[2])
}
