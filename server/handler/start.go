package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"os"

	"github.com/Ubbo-Sathla/anylink/admin"
	"github.com/Ubbo-Sathla/anylink/base"
	"github.com/Ubbo-Sathla/anylink/dbdata"
	"github.com/Ubbo-Sathla/anylink/sessdata"
)

func Start() {
	dbdata.Start()
	sessdata.Start()

	switch base.Cfg.LinkMode {
	case base.LinkModeTUN:
		checkTun()
	case base.LinkModeTAP:
		checkTap()
	case base.LinkModeMacvtap:
		checkMacvtap()
	default:
		base.Fatal("LinkMode is err")
	}

	// 计算profile.xml的hash
	b, err := os.ReadFile(base.Cfg.Profile)
	if err != nil {
		panic(err)
	}
	ha := sha1.Sum(b)
	profileHash = hex.EncodeToString(ha[:])

	go admin.StartAdmin()
	go startTls()
	go startDtls()
}

func Stop() {
	_ = dbdata.Stop()
	destroyVtap()
}
