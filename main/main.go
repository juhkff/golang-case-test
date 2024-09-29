package main

import (
	ginmodeltest "golang-case-test/gin-model-test"
	"time"
)

func main() {
	go func() {
		time.Sleep(3 * time.Second)
		ginmodeltest.Request()
	}()
	ginmodeltest.WebThread()
}
