package ecp

import (
	"testing"

	"sinsiway.com/golang/ecp_manager/prdebug"
)

func TestIsLsnrProcAlive(t *testing.T) {
	prdebug.PrDebug = true

	pid := 20428
	prdebug.Printf("pid[%d] isAlive[%v]\n", pid, IsLsnrProcAlive(pid))
}

func TestStartLsnrProc(t *testing.T) {
	prdebug.PrDebug = true
	homePath := "C:/Users/dbmas/workspace/golang/ecp_manager"
	pid, err := StartLsnrProc(homePath, "127.0.0.1", "25432", "192.168.10.214", "15432")
	if err != nil {
		prdebug.Println("StartLsnrProc failed : ", err)
	}
	prdebug.Println("pid : ", pid)

	prdebug.Printf("pid[%d] isAlive[%v]\n", pid, IsLsnrProcAlive(pid))
}

func TestStopLsnrProc(t *testing.T) {
	prdebug.PrDebug = true
	pid := 18428
	err := StopLsnrProc(pid)
	if err != nil {
		t.Fatal(err)
	}
}
