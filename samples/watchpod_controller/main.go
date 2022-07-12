package main

import (
	"fmt"

	"github.com/tingshua-yts/BetterOperator/watchpod_controller/pkg"
)

func main() {
	fmt.Println("Hello World")
	pkg.Start("/Users/tingshuai.yts/kube/ai_studio_kube.config")
}
