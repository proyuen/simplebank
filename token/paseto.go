package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker 是一个基于 PASETO v4 的令牌制造器
type PasetoMaker struct {
	paseto       *paseto.V4SymmetricKey // 核心变化：不再是 []byte，而是专门的 Key 对象
	symmetricKey []byte
}

// NewPasetoMaker 创建一个新的 PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	// 1. 将字符串转为 v4 专用的密钥对象
	// V4SymmetricKeyFromBytes 会自动处理 32 字节的校验
	key, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, err
	}

	maker := &PasetoMaker{
		paseto:       &key,
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// 1. 创建一个空的 Token 对象
	token := paseto.NewToken()

	// 2. 填入数据 (Claims)
	token.Set("username", username)               // 自定义字段
	token.SetIssuedAt(time.Now())                 // 签发时间
	token.SetNotBefore(time.Now())                // 生效时间
	token.SetExpiration(time.Now().Add(duration)) // 过期时间

	// 3. 加密并签名 (Encrypt)
	// 使用 v4.Local 模式 (对应原来的 v2.Local)
	// nil 表示没有 Footer (注脚)
	encryptedToken := token.V4Encrypt(*maker.paseto, nil)

	return encryptedToken, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	// 1. 创建解析器
	parser := paseto.NewParser()

	// 可以在这里添加额外的规则，例如：
	// parser.AddRule(paseto.NotExpired()) // 默认已经包含此规则

	// 2. 解析 Token
	// ParseV4Local 会自动解密并验证签名、过期时间
	parsedToken, err := parser.ParseV4Local(*maker.paseto, token, nil)
	if err != nil {
		// 将库的错误转换为我们自己的错误类型，方便上层处理
		if err.Error() == "token has expired" { // 注意：实际错误判断可能需要根据库的定义微调
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// 3. 从 Token 中提取数据还原成我们的 Payload 结构体
	// 注意：新库返回的是 parsedToken 对象，我们需要手动映射回原来的 Payload
	username, err := parsedToken.GetString("username")
	if err != nil {
		return nil, ErrInvalidToken
	}

	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}

	expiredAt, err := parsedToken.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	jti, ok := parsedToken.Claims()["jti"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	// 2. 将字符串解析为 UUID 对象
	tokenID, err := uuid.Parse(jti)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 构造返回 Payload
	payload := &Payload{
		ID:        tokenID, // 库自动生成的唯一 ID
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiredAt: expiredAt,
	}

	return payload, nil
}
