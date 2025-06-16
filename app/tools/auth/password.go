package auth

import (
	"govote/app/tools/log"

	"golang.org/x/crypto/bcrypt"
)

// EncryptV2 对密码进行哈希加密
func EncryptV2(pwd string) string {
	newPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.L.Errorf("加密密码失败, err:%s\n", err)
		return ""
	}
	newPwdStr := string(newPwd)
	// fmt.Printf("加密后的密码：%s\n", newPwdStr)
	return newPwdStr
}
