package repository

import (
	"go-rest-api/model"

	"gorm.io/gorm"
)

// 下記functionのインスタンスを作成。それぞれエラーを返す
type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}

// IUserRepository(interface)を実装するための構造体
type userRepository struct {
	db *gorm.DB
}

// 上記構造体にIUserRepositoryを実装
// repository 初期化時DB接続を含める
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

// (ur *userRepository)にGetUserByEmailメソッドを実装
// user *model.user はDBから取得したユーザー情報が格納
// emailは一意に識別するためのemailアドレス
func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	// DBからemailのユーザーを検索、結果をuserに格納
	// エラーの場合はエラーerrを返却
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

// methodの実装
func (ur *userRepository) CreateUser(user *model.User) error {
	// GORMmethodのCreate、user構造体をDBに新規作成
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}