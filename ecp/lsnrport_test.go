package ecp

import (
	"testing"

	"sinsiway.com/golang/ecp_manager/prdebug"
)

func TestGetAvailablePort(t *testing.T) {
	prdebug.PrDebug = true

	// for i := 0; i < int(maxPortNum-minPortNum); i++ {
	// 	PortStatList[i].IsUsingFlag = true
	// }

	p, err := GetAvailablePort("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	prdebug.Println("1. port : ", p)

	p, err = GetAvailablePort("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	prdebug.Println("2. port : ", p)
}
