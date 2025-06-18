package test

import (
	"fmt"
	"govote/app/tools/jwt"
	"testing"
)

func TestGetToken(t *testing.T) {
	str, _ := jwt.GenToken(1, "lucien")
	fmt.Printf("str:%s", str)
}

func TestParseToken(t *testing.T) {
	str := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwibmFtZSI6Imx1Y2llbiIsImV4cCI6MTc1MDIzNTE1MCwiaXNzIjoibHVjaWVuIiwic3ViIjoiZ28tdm90ZSJ9.B25gJ3QbYW7BznDIktxednMipaoa8BWHcCSqP733lM8"
	token, _ := jwt.ParseToken(str)
	fmt.Printf("token:%+v", token)
}
