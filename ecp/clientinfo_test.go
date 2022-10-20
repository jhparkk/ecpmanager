package ecp

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sinsiway.com/golang/ecp_manager/prdebug"
)

func TestCreatClientInfo(t *testing.T) {
	prdebug.PrDebug = true
	db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&ClientInfo{})

	ci := ClientInfo{InfoId: 1, UserId: "test_uid", IpMacAddr: "38-7A-0E-2F-4A-D1"}

	tx := db.Create(&ci)
	if tx.Error != nil {
		t.Fatal(err)
	}

	var clientInfo ClientInfo
	if err = db.Table(TableNameClientInfo).First(&clientInfo).Error; err != nil {
		t.Fatal(err)
	}
	prdebug.Println(clientInfo)

}
