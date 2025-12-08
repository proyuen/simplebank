package token

import (
	"fmt"
	"time"

	// V2 版本常用的库和依赖
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto" // V2 版本的核心库
)

// Maker 接口定义省略

type PasetoMaker struct {
	paseto       *paseto.V2 // 关键点：V2 实例
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	// V2 Local 密钥也必须是 32 字节
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(), // 关键：初始化 V2 实例
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// CreateToken: V2 加密实现
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	// 1. 调用 NewPayload 创建结构体数据
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	// 2. 使用 V2 实例对 Payload 结构体进行加密
	// V2 的 Encrypt 方法是通用的，会自动处理 JSON 序列化
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

// VerifyToken: V2 解密实现
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	// 1. 使用 V2 实例进行解密
	// 如果解密失败（签名不匹配、密钥错误），会返回底层错误
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)

	if err != nil {
		// 返回通用错误，这里就是你之前看到的 "token is invalid" 的来源
		return nil, ErrInvalidToken
	}

	// 2. 手动检查 Payload 结构体本身的有效性（是否过期）
	err = payload.Valid()
	if err != nil {
		return nil, err // 返回 ErrExpiredToken
	}

	return payload, nil
}
