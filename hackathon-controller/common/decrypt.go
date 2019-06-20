package common

import (
	"encoding/hex"
	"github.com/YouDail/golang_micro/hackathon-controller/crypto"
	log "github.com/golang/glog"
	"github.com/spf13/viper"
)

//解密mysql.PasswdSecret
func Decrypt(str string) (bool, string) {
	src := []byte(str)
	dst := make([]byte, hex.DecodedLen(len(src)))
	if _, err := hex.Decode(dst, src); err != nil {
		log.Errorln("Decode Failed: %v \n", err)
		return false, err.Error()
	}
	pwd, err := crypto.ECBDecrypt([]byte(viper.GetString("key")), dst)
	if err != nil {
		log.Errorln("Decrypted Failed: %v \n", err)
		return false, err.Error()
	}
	log.Infoln("解密Secret的结果是： ", string(pwd))
	return true, string(pwd)
}
