package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
)

// Server is an MCP server that communicates over stdio (newline-delimited JSON-RPC).
type Server struct {
	info     ServerInfo
	registry *ToolRegistry
	logger   *slog.Logger
}

// NewServer creates a new MCP server.
func NewServer(name, version string, registry *ToolRegistry, logger *slog.Logger) *Server {
	return &Server{
		info:     ServerInfo{Name: name, Version: version},
		registry: registry,
		logger:   logger,
	}
}

// Serve reads JSON-RPC requests from reader and writes responses to writer.
func (s *Server) Serve(ctx context.Context, reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			resp := NewErrorResponse(nil, ErrCodeParse, "parse error: "+err.Error())
			s.writeResponse(writer, resp)
			continue
		}

		resp := s.handleRequest(ctx, req)
		s.writeResponse(writer, resp)
	}

	return scanner.Err()
}

func (s *Server) handleRequest(ctx context.Context, req Request) Response {
	s.logger.Info("handling request", "method", req.Method)

	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "initialized":
		return NewResponse(req.ID, map[string]string{})
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(ctx, req)
	default:
		return NewErrorResponse(req.ID, ErrCodeMethodNotFound, "method not found: "+req.Method)
	}
}

func (s *Server) handleInitialize(req Request) Response {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		ServerInfo:      s.info,
		Capabilities: ServerCaps{
			Tools: &ToolsCap{},
		},
	}
	return NewResponse(req.ID, result)
}

func (s *Server) handleToolsList(req Request) Response {
	result := ToolsListResult{
		Tools: s.registry.Tools(),
	}
	return NewResponse(req.ID, result)
}

func (s *Server) handleToolsCall(ctx context.Context, req Request) Response {
	var params ToolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewErrorResponse(req.ID, ErrCodeInvalidParams, "invalid params: "+err.Error())
	}

	result, err := s.registry.Handle(ctx, params.Name, params.Arguments)
	if err != nil {
		return NewErrorResponse(req.ID, ErrCodeInternal, fmt.Sprintf("tool error: %v", err))
	}

	return NewResponse(req.ID, result)
}

func (s *Server) writeResponse(w io.Writer, resp Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("failed to marshal response", "error", err)
		return
	}
	data = append(data, '\n')
	w.Write(data)
}
