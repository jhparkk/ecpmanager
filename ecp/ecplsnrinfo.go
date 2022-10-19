package ecp

const TableNameEcpLsnrInfo = "ecp_lsnr_infos"

type EcpLsnrInfo struct {
	//gorm.Model
	InfoId       int    `gorm:"column:info_id;primary key;autoincrement"`
	ListenAddrId int64  `gorm:"column:listen_addr_id"`
	LastUpdate   int64  `gorm:"column:last_update"`
	DbmsIp       string `gorm:"column:dbms_ip"`
	DbmsPort     uint16 `gorm:"column:dbms_port"`
	RelayIp      string `gorm:"column:relay_ip"`
	RelayPort    uint16 `gorm:"column:relay_port"`
	EcpIp        string `gorm:"column:ecp_ip"`
	EcpPort      uint16 `gorm:"column:ecp_port"`
	EcpLsnrPid   int    `gorm:"column:ecp_lsnr_pid"`
	Priority     uint8  `gorm:"column:priority"`
	SvrTimeout   int32  `gorm:"column:svr_timeout"`
	FailOpen     uint8  `gorm:"column:fail_open"`
	Reserved     string `gorm:"column:reserved"`
}
