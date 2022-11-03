package ecp

import (
	"errors"
	"net"
	"strconv"
	"sync"
	"time"

	"sinsiway.com/golang/ecpmanager/prdebug"
)

const (
	minPortNum uint16 = 23301
	maxPortNum        = 23500
)

type PortStat struct {
	Port        uint16
	ScanTime    time.Time
	IsUsingFlag bool
}

var LsnrPortIp string
var PortStatList []PortStat
var currIdxPSL int

func init() {
	currIdxPSL = 0
	for i := minPortNum; i < maxPortNum; i++ {
		p := PortStat{Port: i, IsUsingFlag: false}
		PortStatList = append(PortStatList, p)
	}
}

func IsUsingPort(ip string, port uint16) bool {
	if port == 0 {
		return true
	}
	to, err := net.DialTimeout("tcp", ip+":"+strconv.Itoa(int(port)), time.Second*1)
	if err != nil {
		//prdebug.Println(err)
		return false
	}
	to.Close()
	return true
}

func scanPort(lsnrIp string, start, end int) bool {
	isPortFound := false
	wg := sync.WaitGroup{}
	for l := start; l < end; l++ {
		if l >= int(maxPortNum-minPortNum)-1 {
			break
		}
		prdebug.Println(l, " : ", PortStatList[l])

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			PortStatList[i].IsUsingFlag = IsUsingPort(lsnrIp, PortStatList[i].Port)

			if PortStatList[i].IsUsingFlag == false {
				isPortFound = true
			}
			PortStatList[i].ScanTime = time.Now()
		}(l)

	}
	wg.Wait()

	prdebug.Println("isPortFound : ", isPortFound)
	return isPortFound
}

func GetAvailablePort(scanIp string) (port uint16, err error) {
	idx := currIdxPSL

	for idx < len(PortStatList)-1 {
		if PortStatList[idx].IsUsingFlag {
			idx++
			continue
		}
		if time.Since(PortStatList[idx].ScanTime) < 60*time.Second {
			currIdxPSL = idx
			PortStatList[currIdxPSL].IsUsingFlag = true
			port = PortStatList[currIdxPSL].Port
			break
		}

		if scanPort(scanIp, idx, idx+10) {
			for l := idx; l < idx+10; l++ {
				if PortStatList[l].IsUsingFlag == false {
					currIdxPSL = l
					break
				}
			}
			PortStatList[currIdxPSL].IsUsingFlag = true
			port = PortStatList[currIdxPSL].Port
			break
		}
		idx++
	}
	if port == 0 {
		err = errors.New("no available port")
		return 0, err
	}
	currIdxPSL++
	return port, nil
}
