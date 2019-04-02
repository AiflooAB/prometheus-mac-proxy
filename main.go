package main

import (
	"github.com/AiflooAB/prometheus-mac-proxy/pkg/mac"
	"github.com/AiflooAB/prometheus-mac-proxy/pkg/proxy"
)

func main() {
	proxy.Start(mac.NewFileTransformer)
}
