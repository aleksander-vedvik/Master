// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.7.0-devel
// 	protoc            v3.12.4
// source: zorums.proto

package dev

import (
	gorums "github.com/relab/gorums"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = gorums.EnforceVersion(7 - gorums.MinVersion)
	// Verify that the gorums runtime is sufficiently up-to-date.
	_ = gorums.EnforceVersion(gorums.MaxVersion - 7)
)

func (b *Broadcast) QuorumCallWithBroadcast(req *Request, opts ...gorums.BroadcastOption) {
	if b.metadata.BroadcastID == 0 {
		panic("broadcastID cannot be empty. Use srv.BroadcastQuorumCallWithBroadcast instead")
	}
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	options.ServerAddresses = append(options.ServerAddresses, b.srvAddrs...)
	b.orchestrator.BroadcastHandler("dev.ZorumsService.QuorumCallWithBroadcast", req, b.metadata.BroadcastID, b.enqueueBroadcast, options)
}

func (b *Broadcast) MulticastWithBroadcast(req *Request, opts ...gorums.BroadcastOption) {
	if b.metadata.BroadcastID == 0 {
		panic("broadcastID cannot be empty. Use srv.BroadcastMulticastWithBroadcast instead")
	}
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	options.ServerAddresses = append(options.ServerAddresses, b.srvAddrs...)
	b.orchestrator.BroadcastHandler("dev.ZorumsService.MulticastWithBroadcast", req, b.metadata.BroadcastID, b.enqueueBroadcast, options)
}

func (b *Broadcast) BroadcastInternal(req *Request, opts ...gorums.BroadcastOption) {
	if b.metadata.BroadcastID == 0 {
		panic("broadcastID cannot be empty. Use srv.BroadcastBroadcastInternal instead")
	}
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	options.ServerAddresses = append(options.ServerAddresses, b.srvAddrs...)
	b.orchestrator.BroadcastHandler("dev.ZorumsService.BroadcastInternal", req, b.metadata.BroadcastID, b.enqueueBroadcast, options)
}

func (b *Broadcast) BroadcastWithClientHandlerAndBroadcastOption(req *Request, opts ...gorums.BroadcastOption) {
	if b.metadata.BroadcastID == 0 {
		panic("broadcastID cannot be empty. Use srv.BroadcastBroadcastWithClientHandlerAndBroadcastOption instead")
	}
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	options.ServerAddresses = append(options.ServerAddresses, b.srvAddrs...)
	b.orchestrator.BroadcastHandler("dev.ZorumsService.BroadcastWithClientHandlerAndBroadcastOption", req, b.metadata.BroadcastID, b.enqueueBroadcast, options)
}
