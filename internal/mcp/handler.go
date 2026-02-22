package mcp

import (
	"context"
	"encoding/json"
)

// ToolHandler processes a tool call and returns a result.
type ToolHandler func(ctx context.Context, args json.RawMessage) (ToolCallResult, error)

// ToolRegistry maps tool names to their definitions and handlers.
type ToolRegistry struct {
	tools    []Tool
	handlers map[string]ToolHandler
}

// NewToolRegistry creates an empty tool registry.
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		handlers: make(map[string]ToolHandler),
	}
}

// Register adds a tool with its handler.
func (r *ToolRegistry) Register(tool Tool, handler ToolHandler) {
	r.tools = append(r.tools, tool)
	r.handlers[tool.Name] = handler
}

// Tools returns all registered tools.
func (r *ToolRegistry) Tools() []Tool {
	return r.tools
}

// Handle calls the handler for the named tool.
func (r *ToolRegistry) Handle(ctx context.Context, name string, args json.RawMessage) (ToolCallResult, error) {
	handler, ok := r.handlers[name]
	if !ok {
		return ToolCallResult{
			Content: []ContentBlock{NewTextContent("unknown tool: " + name)},
			IsError: true,
		}, nil
	}
	return handler(ctx, args)
}
