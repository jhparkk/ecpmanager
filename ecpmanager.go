package main

import (
	"os"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sinsiway.com/golang/ecp_manager/ecp"
	"sinsiway.com/golang/ecp_manager/prdebug"
	"sinsiway.com/golang/ecp_manager/worker"
)

type EcpManager struct {
	worker       *worker.Worker
	homePath     string
	clientInfo   ecp.ClientInfo
	ecpLsnrInfos []ecp.EcpLsnrInfo
	db           *gorm.DB
}

func NewEcpManager(dbFilePath string) *EcpManager {
	em := &EcpManager{}
	em.worker = worker.NewWorker("ecp_manager", em)
	em.homePath = dbFilePath
	return em
}

func (em *EcpManager) In() (int, error) {
	db, err := gorm.Open(sqlite.Open(em.homePath+"/petra.db"), &gorm.Config{})
	if err != nil {
		return -1, err
	}
	em.db = db
	// get a client_infos row
	if err = em.db.Table(ecp.TableNameClientInfo).First(&em.clientInfo).Error; err != nil {
		return -1, err
	}
	// get ecp_lsnr_infos table rows
	if err = em.db.Table(ecp.TableNameEcpLsnrInfo).Find(&em.ecpLsnrInfos).Error; err != nil {
		return -1, err
	}

	var eli ecp.EcpLsnrInfo
	for _, eli = range em.ecpLsnrInfos {
		//
		// reserved - I:new / D:deleted / R:running
		// skip deleted
		//
		if eli.Reserved == "D" {
			continue
		}

		prdebug.Println(eli)

		// check remain process
		if eli.EcpLsnrPid != 0 && ecp.IsLsnrProcAlive(int(eli.EcpLsnrPid)) {
			continue
		}

		// get lsnr port
		eli.EcpIp = "127.0.0.1"
		eli.EcpLsnrPid = 0
		prdebug.Println("1. eli.ecpPort : ", eli.EcpPort)
		if eli.EcpPort == 0 {
			port, err := ecp.GetAvailablePort("127.0.0.1")
			if err != nil {
				worker.Logger.Printf("[%d] GetAvailablePort failed : %s\n", os.Getpid(), err)
				continue
			}
			if port == 0 {
				worker.Logger.Printf("[%d] port is not allocated for ecplocallsnr.\n", os.Getpid())
				continue
			}
			eli.EcpPort = port
		}
		prdebug.Println("2. eli.ecpPort : ", eli.EcpPort)
		err = em.startEcpLsnrProc(&eli)
		if err != nil {
			worker.Logger.Printf("[%d] startEcpLsnrProc failed : %s\n", os.Getpid(), err)
			continue
		}
	}

	worker.Logger.Printf("[%d] manager starts.\n", os.Getpid())
	return 0, nil
}

func (em *EcpManager) Run() (int, error) {

	//
	// check ecp_lsnr_infos table rows
	// reserved - I:new / D:deleted / R:running
	//
	if err := em.db.Table(ecp.TableNameEcpLsnrInfo).Find(&em.ecpLsnrInfos).Error; err != nil {
		worker.Logger.Printf("[%d] failed to select table ecp_lsnr_infos : %s\n", os.Getpid(), err)
		em.worker.Sleep(5)
		return 0, nil
	}
	var eli ecp.EcpLsnrInfo
	for _, eli = range em.ecpLsnrInfos {
		switch eli.Reserved {
		case "R":
			{
				prdebug.Println("switch R : ", eli)
				//
				// check ecp lsnr processes
				//
				if ecp.IsLsnrProcAlive(int(eli.EcpLsnrPid)) == false {
					if ecp.IsUsingPort("127.0.0.1", eli.EcpPort) {
						port, err := ecp.GetAvailablePort("127.0.0.1")
						if err != nil {
							worker.Logger.Printf("[%d] GetAvailablePort failed : %s\n", os.Getpid(), err)
							continue
						}
						if port == 0 {
							worker.Logger.Printf("[%d] port is not allocated for ecplocallsnr.\n", os.Getpid())
							continue
						}
						eli.EcpPort = port
					}
					err := em.startEcpLsnrProc(&eli)
					if err != nil {
						worker.Logger.Printf("[%d] startEcpLsnrProc failed : %s\n", os.Getpid(), err)
						continue
					}
				}
				break
			}
		case "I":
			{
				//
				// start a new ecp lsnr process
				//
				prdebug.Println("switch I")
				port, err := ecp.GetAvailablePort("127.0.0.1")
				if err != nil {
					worker.Logger.Printf("[%d] GetAvailablePort failed : %s\n", os.Getpid(), err)
					continue
				}
				if port == 0 {
					worker.Logger.Printf("[%d] port is not allocated for ecplocallsnr.\n", os.Getpid())
					continue
				}
				eli.EcpPort = port
				eli.EcpIp = "127.0.0.1"
				err = em.startEcpLsnrProc(&eli)
				if err != nil {
					worker.Logger.Printf("[%d] startEcpLsnrProc failed : %s\n", os.Getpid(), err)
					continue
				}
				break
			}
		case "U":
			{
				//
				// restart the ecp lsnr process
				//
				prdebug.Println("switch U")
				err := ecp.StopLsnrProc(eli.EcpLsnrPid)
				if err != nil {
					worker.Logger.Printf("[%d] StopLsnrProc failed : %s\n", os.Getpid(), err)
				}
				eli.EcpLsnrPid = 0
				err = em.startEcpLsnrProc(&eli)
				if err != nil {
					worker.Logger.Printf("[%d] startEcpLsnrProc failed : %s\n", os.Getpid(), err)
					continue
				}

				break
			}
		case "D":
			{
				//
				// stop the ecp lsnr process
				//
				prdebug.Println("switch D")
				go func(pid int) {
					err := ecp.StopLsnrProc(eli.EcpLsnrPid)
					if err != nil {
						worker.Logger.Printf("[%d] StopLsnrProc[%d] failed : %s\n", os.Getpid(), eli.EcpLsnrPid, err)
					}
				}(eli.EcpLsnrPid)
				eli.EcpLsnrPid = 0
				//
				// delete ecp_lsnr_info table row
				//
				err := em.db.Table(ecp.TableNameEcpLsnrInfo).Where(eli.InfoId).Delete(&ecp.EcpLsnrInfo{}).Error
				if err != nil {
					worker.Logger.Printf("[%d] failed to update a row: %s\n", os.Getpid(), err)
					continue
				}
				break
			}
		default:
		}

	}

	em.worker.Sleep(5 * time.Second)
	return 0, nil
}

func (em *EcpManager) Out() (int, error) {
	//
	// stop processes
	//
	em.db.Table(ecp.TableNameEcpLsnrInfo).Find(&em.ecpLsnrInfos)
	var eli ecp.EcpLsnrInfo
	for _, eli = range em.ecpLsnrInfos {
		err := ecp.StopLsnrProc(eli.EcpLsnrPid)
		if err != nil {
			worker.Logger.Printf("[%d] StopLsnrProc failed : %s\n", os.Getpid(), err)
			continue
		}
		eli.EcpLsnrPid = 0
		err = em.db.Table(ecp.TableNameEcpLsnrInfo).Where(eli.InfoId).Save(&eli).Error
		if err != nil {
			worker.Logger.Printf("[%d] failed to update a row: %s\n", os.Getpid(), err)
			continue
		}
	}

	//
	// close db
	//
	sqlDB, err := em.db.DB()
	if err != nil {
		return -1, err
	}
	sqlDB.Close()

	worker.Logger.Printf("[%d] manager stops.\n", os.Getpid())
	return 0, nil
}

func (em *EcpManager) startEcpLsnrProc(eli *ecp.EcpLsnrInfo) error {
	pid, err := ecp.StartLsnrProc(
		em.homePath,
		eli.EcpIp,
		strconv.Itoa(int(eli.EcpPort)),
		eli.RelayIp,
		strconv.Itoa(int(eli.RelayPort)),
		em.clientInfo.UserId,
		em.clientInfo.IpMacAddr)
	if err != nil {
		return err
	}
	eli.EcpLsnrPid = pid

	//
	// update reserved 'R'
	//
	eli.Reserved = "R"
	err = em.db.Table(ecp.TableNameEcpLsnrInfo).Where(eli.InfoId).Save(&eli).Error
	if err != nil {
		return err
	}
	return nil
}
