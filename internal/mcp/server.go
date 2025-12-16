// Package mcp 提供 Model Context Protocol (MCP) 服务器实现
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"cursor2api/internal/tools"
)

// Server MCP 服务器
type Server struct {
	executor *tools.Executor
	input    io.Reader
	output   io.Writer
}

// NewServer 创建 MCP 服务器
func NewServer() *Server {
	return &Server{
		executor: tools.NewExecutor(),
		input:    os.Stdin,
		output:   os.Stdout,
	}
}

// JSONRPCRequest JSON-RPC 请求
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse JSON-RPC 响应
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError JSON-RPC 错误
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ServerInfo 服务器信息
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult 初始化结果
type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

// ServerCapabilities 服务器能力
type ServerCapabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

// ToolsCapability 工具能力
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// Tool MCP 工具定义
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema 输入模式
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

// Property 属性定义
type Property struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// ToolsListResult 工具列表结果
type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

// CallToolParams 调用工具参数
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResult 调用工具结果
type CallToolResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ContentItem 内容项
type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// Run 运行 MCP 服务器 (stdio 模式)
func (s *Server) Run() error {
	log.Println("[MCP] 服务器启动 (stdio 模式)")
	scanner := bufio.NewScanner(s.input)
	// 增加缓冲区大小以处理大请求
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, -32700, "Parse error", err.Error())
			continue
		}

		s.handleRequest(req)
	}

	return scanner.Err()
}

// handleRequest 处理请求
func (s *Server) handleRequest(req JSONRPCRequest) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "initialized":
		// 客户端确认初始化完成，无需响应
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolsCall(req)
	case "ping":
		s.sendResult(req.ID, map[string]string{})
	default:
		s.sendError(req.ID, -32601, "Method not found", req.Method)
	}
}

// handleInitialize 处理初始化请求
func (s *Server) handleInitialize(req JSONRPCRequest) {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    "cursor2api-mcp",
			Version: "1.0.0",
		},
	}
	s.sendResult(req.ID, result)
}

// handleToolsList 处理工具列表请求
func (s *Server) handleToolsList(req JSONRPCRequest) {
	tools := []Tool{
		{
			Name:        "bash",
			Description: "执行 bash 命令",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"command": {Type: "string", Description: "要执行的命令"},
					"cwd":     {Type: "string", Description: "工作目录（可选）"},
				},
				Required: []string{"command"},
			},
		},
		{
			Name:        "read_file",
			Description: "读取文件内容",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"path": {Type: "string", Description: "文件路径"},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "write_file",
			Description: "写入文件内容",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"path":    {Type: "string", Description: "文件路径"},
					"content": {Type: "string", Description: "文件内容"},
				},
				Required: []string{"path", "content"},
			},
		},
		{
			Name:        "list_dir",
			Description: "列出目录内容",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"path": {Type: "string", Description: "目录路径"},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "edit",
			Description: "编辑文件（查找替换）",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"path":        {Type: "string", Description: "文件路径"},
					"old_string":  {Type: "string", Description: "要替换的内容"},
					"new_string":  {Type: "string", Description: "替换后的内容"},
					"replace_all": {Type: "boolean", Description: "是否替换所有匹配"},
				},
				Required: []string{"path", "old_string", "new_string"},
			},
		},
	}

	s.sendResult(req.ID, ToolsListResult{Tools: tools})
}

// handleToolsCall 处理工具调用请求
func (s *Server) handleToolsCall(req JSONRPCRequest) {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, "Invalid params", err.Error())
		return
	}

	output, err := s.executor.Execute(params.Name, params.Arguments)

	if err != nil {
		s.sendResult(req.ID, CallToolResult{
			Content: []ContentItem{{Type: "text", Text: fmt.Sprintf("错误: %v\n%s", err, output)}},
			IsError: true,
		})
		return
	}

	s.sendResult(req.ID, CallToolResult{
		Content: []ContentItem{{Type: "text", Text: output}},
	})
}

// sendResult 发送成功响应
func (s *Server) sendResult(id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.sendResponse(resp)
}

// sendError 发送错误响应
func (s *Server) sendError(id interface{}, code int, message string, data interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	s.sendResponse(resp)
}

// sendResponse 发送响应
func (s *Server) sendResponse(resp JSONRPCResponse) {
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[MCP] 序列化响应失败: %v", err)
		return
	}
	fmt.Fprintln(s.output, string(data))
}
