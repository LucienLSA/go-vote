package jwt

import (
	"errors"
	"govote/app/config"
	"time"

	"github.com/golang-jwt/jwt"
)

// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	Id                 int64  `json:"id"`
	Name               string `json:"name"`
	jwt.StandardClaims        // 内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(id int64, name string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(config.Conf.AppConfig.JwtExpireTime * time.Hour)
	// 创建一个自己的声明数据
	claims := MyClaims{
		Id:   id,
		Name: name, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    config.Conf.AppConfig.JwtIssuer,  // 签发人
			Subject:   config.Conf.AppConfig.JwtSubject, // 签发对象
		},
	}
	// 使用指定的签名方法创建签名对象
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	token, err := tokenClaims.SignedString([]byte(config.Conf.AppConfig.JwtSecret))
	return token, err
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (myclaims *MyClaims, err error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	myclaims = new(MyClaims) // 必须手动初始化 分配内存空间
	tokenClaims, err := jwt.ParseWithClaims(tokenString, myclaims, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.Conf.AppConfig.JwtSecret), nil
	})
	if tokenClaims != nil {
		// 对token对象中的Claim进行类型断言
		if claims, ok := tokenClaims.Claims.(*MyClaims); ok && tokenClaims.Valid { // 校验toke
			return claims, nil
		}
	}
	return nil, errors.New("invalid token")
}
