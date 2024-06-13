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

func registerClientServerHandlers(srv *clientServerImpl) {

	srv.RegisterHandler("dev.ZorumsService.BroadcastWithClientHandler1", gorums.ClientHandler(srv.clientBroadcastWithClientHandler1))
	srv.RegisterHandler("dev.ZorumsService.BroadcastWithClientHandler2", gorums.ClientHandler(srv.clientBroadcastWithClientHandler2))
	srv.RegisterHandler("dev.ZorumsService.BroadcastWithClientHandlerAndBroadcastOption", gorums.ClientHandler(srv.clientBroadcastWithClientHandlerAndBroadcastOption))
}
