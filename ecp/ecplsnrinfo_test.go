package ecp

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sinsiway.com/golang/ecp_manager/prdebug"
)

func TestCreateELI(t *testing.T) {
	prdebug.PrDebug = true
	db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	if err != nil {
		prdebug.Println(err)
		return
	}
	db.AutoMigrate(&EcpLsnrInfo{})

	elis := []EcpLsnrInfo{
		{ListenAddrId: 10, DbmsIp: "dbmsip1"},
		{ListenAddrId: 11},
		{ListenAddrId: 100},
		{ListenAddrId: 400},
	}

	tx := db.Create(&elis)
	if tx.Error != nil {
		prdebug.Println(err)
		return
	}

	eli := EcpLsnrInfo{InfoId: 2}
	// if err = db.Table("ecp_lsnr_infos").Delete(&eli, eli.InfoId).Error; err != nil {
	// 	prdebug.Println(err)
	// }
	db.Delete(&EcpLsnrInfo{}, 2)
	if err = db.Table("ecp_lsnr_infos").Where(11).Delete(&EcpLsnrInfo{}).Error; err != nil {
		prdebug.Println(err)
	}
	prdebug.Println("delete : ", eli)

	var selectElis []EcpLsnrInfo
	db.Table("ecp_lsnr_infos").Find(&selectElis)
	for _, eli := range selectElis {
		eli.Reserved = "updated"
		if err = db.Table("ecp_lsnr_infos").Where(eli.InfoId).Save(&eli).Error; err != nil {
			prdebug.Println("save err : ", err)
		}

		tx := db.Table("ecp_lsnr_infos").Where(
			map[string]interface{}{
				"info_id": eli.InfoId,
			})

		if tx.Error != nil {
			prdebug.Println("find err : ", err)
		}
		tx.Save(&eli)

		prdebug.Println(eli)

	}

}