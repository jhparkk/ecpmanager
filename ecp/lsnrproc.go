package ecp

import (
	"errors"
	"os"
	"os/exec"
	"strconv"

	ps "github.com/mitchellh/go-ps"

	"sinsiway.com/golang/ecp_manager/prdebug"
)

const BinNameEcpLocalLsnr = "ecplocallsnr.exe"

func IsLsnrProcAlive(pid int) bool {
	p, err := ps.FindProcess(pid)
	if err != nil {
		prdebug.Println("FindProcess failed : ", err)
		return false
	}

	if p == nil {
		// no process
		return false
	}

	prdebug.Println("found process - ", p.Pid(), ":", p.Executable())
	if p.Executable() != BinNameEcpLocalLsnr {
		return false
	}

	return true
}

func StartLsnrProc(homePath, lsnrIp, lsnrPort, relayIp, relayPort string) (pid int, err error) {
	cmd := exec.Command(homePath+"/"+BinNameEcpLocalLsnr, lsnrIp, lsnrPort, relayIp, relayPort)
	err = cmd.Start()
	if err != nil {
		return 0, err
	}
	pid = cmd.Process.Pid

	return pid, nil
}

func StopLsnrProc(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		prdebug.Printf("find process[%d] failed : %s", pid, err)
		return err
	}
	if p == nil {
		return errors.New("process[" + strconv.Itoa(pid) + "] not found")
	}

	//err = p.Signal(syscall.SIGTERM)  // not supported on windows
	// if err != nil {
	// 	prdebug.Println("Signal failed : ", err)
	// 	return err
	// }

	err = p.Kill()
	if err != nil {
		prdebug.Println("Kill failed : ", err)
		return err
	}
	return nil
}
