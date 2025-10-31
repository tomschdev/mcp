// Copyright 2025 Redpanda Data, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/redpanda-data/protoc-gen-go-mcp/pkg/runtime"
	pb "github.com/tomschdev/mcp/gen/mcp/tom/v1"
	mcpPb "github.com/tomschdev/mcp/gen/mcp/tom/v1/v1mcp"
)

// Ensure our interface and the official gRPC interface are grpcClient
var (
	grpcClient pb.TomServiceClient
	_          = pb.TomServiceClient(grpcClient)
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Example auto-generated gRPC-MCP with runtime LLM provider selection",
		"1.0.0",
	)

	srv := TomServer{}

	// Get LLM provider from environment variable, default to standard
	providerStr := os.Getenv("LLM_PROVIDER")
	var provider runtime.LLMProvider
	switch providerStr {
	case "openai":
		provider = runtime.LLMProviderOpenAI
		fmt.Printf("Using OpenAI-compatible MCP handlers\n")
	case "standard":
		fallthrough
	default:
		provider = runtime.LLMProviderStandard
		fmt.Printf("Using standard MCP handlers\n")
	}

	// Register handlers for the selected provider
	mcpPb.RegisterTomServiceHandlerWithProvider(s, &srv, provider)

	// Alternative: Register specific handlers directly
	// mcpPb.RegisterTestServiceHandler(s, &srv)        // Standard
	// mcpPb.RegisterTestServiceHandlerOpenAI(s, &srv)  // OpenAI

	// Alternative: Register both for different tool names
	// mcpPb.RegisterTestServiceHandler(s, &srv)
	// mcpPb.RegisterTestServiceHandlerOpenAI(s, &srv)

	mcpPb.ForwardToTomServiceClient(s, grpcClient)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

type TomServer struct{}

func (t *TomServer) CreateItem(ctx context.Context, in *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	return &pb.CreateItemResponse{
		Id: "item-123",
	}, nil
}

func (t *TomServer) GetItem(ctx context.Context, in *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	return &pb.GetItemResponse{
		Item: &pb.Item{
			Id:   in.GetId(),
			Name: "Retrieved item",
		},
	}, nil
}

func (t *TomServer) ProcessWellKnownTypes(ctx context.Context, in *pb.ProcessWellKnownTypesRequest) (*pb.ProcessWellKnownTypesResponse, error) {
	return &pb.ProcessWellKnownTypesResponse{
		Message: "Processed well-known types",
	}, nil
}
