package service

import (
	"encoding/base64"
	"errors"

	"nas-manager/internal/repository"
	"nas-manager/pkg/crypto"
)

const (
	settingCryptoSalt    = "crypto_salt"
	settingCryptoVerify = "crypto_verify"
)

// EncryptService - 加密服务
type EncryptService struct {
	repo   *repository.SettingRepository
	crypto *crypto.Crypto
}

// NewEncryptService - 创建加密服务
func NewEncryptService(repo *repository.SettingRepository) *EncryptService {
	return &EncryptService{
		repo:   repo,
		crypto: crypto.NewCrypto(),
	}
}

// SetupPasswordRequest - 设置密码请求
type SetupPasswordRequest struct {
	Password string `json:"password"`
}

// VerifyPasswordRequest - 验证密码请求
type VerifyPasswordRequest struct {
	Password string `json:"password"`
}

// ChangePasswordRequest - 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// SetupPassword - 设置加密密码
func (s *EncryptService) SetupPassword(req *SetupPasswordRequest) error {
	// 验证密码长度
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// 检查是否已设置密码
	hasPassword, err := s.HasPassword()
	if err != nil {
		return err
	}
	if hasPassword {
		return errors.New("password already set")
	}

	// 生成盐值
	salt, err := s.crypto.GenerateSalt()
	if err != nil {
		return err
	}

	// 生成验证值
	verifyValue, err := s.crypto.GenerateVerifyValue(req.Password, salt)
	if err != nil {
		return err
	}

	// 保存盐值和验证值
	if err := s.repo.SetSetting(settingCryptoSalt, base64.StdEncoding.EncodeToString(salt)); err != nil {
		return err
	}
	if err := s.repo.SetSetting(settingCryptoVerify, verifyValue); err != nil {
		return err
	}

	return nil
}

// VerifyPassword - 验证密码
func (s *EncryptService) VerifyPassword(req *VerifyPasswordRequest) (bool, error) {
	// 获取盐值和验证值
	saltStr, err := s.repo.GetSetting(settingCryptoSalt)
	if err != nil {
		return false, nil // 没有设置密码
	}

	verifyValue, err := s.repo.GetSetting(settingCryptoVerify)
	if err != nil {
		return false, nil
	}

	salt, err := base64.StdEncoding.DecodeString(saltStr)
	if err != nil {
		return false, err
	}

	return s.crypto.VerifyPassword(req.Password, salt, verifyValue), nil
}

// ChangePassword - 修改密码
func (s *EncryptService) ChangePassword(req *ChangePasswordRequest) error {
	// 验证新密码长度
	if len(req.NewPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	// 验证旧密码
	valid, err := s.VerifyPassword(&VerifyPasswordRequest{Password: req.OldPassword})
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("invalid old password")
	}

	// 生成新的盐值
	newSalt, err := s.crypto.GenerateSalt()
	if err != nil {
		return err
	}

	// 生成新的验证值
	newVerifyValue, err := s.crypto.GenerateVerifyValue(req.NewPassword, newSalt)
	if err != nil {
		return err
	}

	// 保存新的盐值和验证值
	if err := s.repo.SetSetting(settingCryptoSalt, base64.StdEncoding.EncodeToString(newSalt)); err != nil {
		return err
	}
	if err := s.repo.SetSetting(settingCryptoVerify, newVerifyValue); err != nil {
		return err
	}

	// TODO: 重新加密已有凭证（等 Story 1.2 的云存储功能实现后）
	// 如果有云存储凭证，需要用新密钥重新加密

	return nil
}

// HasPassword - 检查是否已设置密码
func (s *EncryptService) HasPassword() (bool, error) {
	_, err := s.repo.GetSetting(settingCryptoSalt)
	if err != nil {
		// 没有设置密码
		return false, nil
	}
	return true, nil
}
