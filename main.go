package main

import (
	"log"
	"os"

	"sinsiway.com/golang/ecpmanager/prdebug"
	"sinsiway.com/golang/ecpmanager/worker"
)

func main() {
	//prdebug.PrDebug = true
	prdebug.Println("start program")

	if len(os.Args) != 2 {
		log.Fatal("ecpmanager <ecp_home_path>")
	}
	homePath := os.Args[1]

	logFile, err := os.OpenFile(homePath+"/ecp_manager.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		prdebug.Println("OpenFile() failed : ", err)
	}
	worker.Logger = log.New(logFile, "" /*no prefix*/, log.LstdFlags|log.Lshortfile)
	manager := NewEcpManager(homePath)

	if _, err = manager.worker.Start(); err != nil {
		log.Fatal("Start() failed : ", err)
	}
	manager.worker.Wait()
}
