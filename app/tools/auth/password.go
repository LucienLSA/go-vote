package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// EncryptV2 对密码进行哈希加密
func EncryptV2(pwd string) string {
	newPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("密码加密失败：", err)
		return ""
	}
	newPwdStr := string(newPwd)
	fmt.Printf("加密后的密码：%s\n", newPwdStr)
	return newPwdStr
}
