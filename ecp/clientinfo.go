package ecp

const TableNameClientInfo = "client_infos"

type ClientInfo struct {
	//gorm.Model
	InfoId    int    `gorm:"column:info_id;primary key;autoincrement"`
	UserId    string `gorm:"column:user_id"`
	IpMacAddr string `gorm:"column:ip_mac_addr"`
}
