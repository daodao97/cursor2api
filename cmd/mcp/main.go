// MCP 服务器入口
// 用于 stdio 模式运行 MCP 服务器
package main

import (
	"log"
	"os"

	"cursor2api/internal/mcp"
)

func main() {
	// 禁用日志输出到 stdout（MCP 使用 stdout 通信）
	log.SetOutput(os.Stderr)

	server := mcp.NewServer()
	if err := server.Run(); err != nil {
		log.Fatalf("[MCP] 服务器错误: %v", err)
	}
}
