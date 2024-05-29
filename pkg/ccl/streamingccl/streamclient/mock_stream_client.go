// Copyright 2024 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt

package streamclient

import (
	"context"

	"github.com/cockroachdb/cockroach/pkg/ccl/streamingccl"
	"github.com/cockroachdb/cockroach/pkg/repstream/streampb"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/span"
	"github.com/cockroachdb/errors"
)

// MockStreamClient will return the slice of events associated to the stream
// partition being consumed. Stream partitions are identified by unique
// partition addresses.
type MockStreamClient struct {
	PartitionEvents map[string][]streamingccl.Event
	DoneCh          chan struct{}
	HeartbeatErr    error
	HeartbeatStatus streampb.StreamReplicationStatus
	OnHeartbeat     func() (streampb.StreamReplicationStatus, error)
}

var _ Client = &MockStreamClient{}

// Create implements the Client interface.
func (m *MockStreamClient) Create(
	_ context.Context, _ roachpb.TenantName, _ streampb.ReplicationProducerRequest,
) (streampb.ReplicationProducerSpec, error) {
	panic("unimplemented")
}

// Dial implements the Client interface.
func (m *MockStreamClient) Dial(_ context.Context) error {
	panic("unimplemented")
}

// Heartbeat implements the Client interface.
func (m *MockStreamClient) Heartbeat(
	_ context.Context, _ streampb.StreamID, _ hlc.Timestamp,
) (streampb.StreamReplicationStatus, error) {
	if m.OnHeartbeat != nil {
		return m.OnHeartbeat()
	}
	return m.HeartbeatStatus, m.HeartbeatErr
}

// Plan implements the Client interface.
func (m *MockStreamClient) Plan(_ context.Context, _ streampb.StreamID) (Topology, error) {
	panic("unimplemented mock method")
}

type mockSubscription struct {
	eventsCh chan streamingccl.Event
}

// Subscribe implements the Subscription interface.
func (m *mockSubscription) Subscribe(_ context.Context) error {
	return nil
}

// Events implements the Subscription interface.
func (m *mockSubscription) Events() <-chan streamingccl.Event {
	return m.eventsCh
}

// Err implements the Subscription interface.
func (m *mockSubscription) Err() error {
	return nil
}

// Subscribe implements the Client interface.
func (m *MockStreamClient) Subscribe(
	ctx context.Context,
	_ streampb.StreamID,
	_ int32,
	token SubscriptionToken,
	initialScanTime hlc.Timestamp,
	_ span.Frontier,
	_ ...SubscribeOption,
) (Subscription, error) {
	var events []streamingccl.Event
	var ok bool
	if events, ok = m.PartitionEvents[string(token)]; !ok {
		return nil, errors.Newf("no events found for partition %s", string(token))
	}
	log.Infof(ctx, "%q beginning subscription from time %v ", string(token), initialScanTime)

	log.Infof(ctx, "%q emitting %d events", string(token), len(events))
	eventCh := make(chan streamingccl.Event, len(events))
	for _, event := range events {
		log.Infof(ctx, "%q emitting event %v", string(token), event)
		eventCh <- event
	}
	log.Infof(ctx, "%q done emitting %d events", string(token), len(events))
	go func() {
		if m.DoneCh != nil {
			log.Infof(ctx, "%q waiting for doneCh", string(token))
			<-m.DoneCh
			log.Infof(ctx, "%q received event on doneCh", string(token))
		}
		close(eventCh)
	}()
	return &mockSubscription{eventsCh: eventCh}, nil
}

// Close implements the Client interface.
func (m *MockStreamClient) Close(_ context.Context) error {
	return nil
}

// Complete implements the streamclient.Client interface.
func (m *MockStreamClient) Complete(_ context.Context, _ streampb.StreamID, _ bool) error {
	return nil
}

// PriorReplicationDetails implements the streamclient.Client interface.
func (m *MockStreamClient) PriorReplicationDetails(
	_ context.Context, _ roachpb.TenantName,
) (string, string, hlc.Timestamp, error) {
	return "", "", hlc.Timestamp{}, nil
}

// ErrorStreamClient always returns an error when consuming a partition.
type ErrorStreamClient struct{ MockStreamClient }

var _ Client = &ErrorStreamClient{}

// ConsumePartition implements the streamclient.Client interface.
func (m *ErrorStreamClient) Subscribe(
	_ context.Context,
	_ streampb.StreamID,
	_ int32,
	_ SubscriptionToken,
	_ hlc.Timestamp,
	_ span.Frontier,
	_ ...SubscribeOption,
) (Subscription, error) {
	return nil, errors.New("this client always returns an error")
}

// Complete implements the streamclient.Client interface.
func (m *ErrorStreamClient) Complete(_ context.Context, _ streampb.StreamID, _ bool) error {
	return nil
}
