package repository

import (
	"fmt"
	"go-rest-api/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// taskに関するDB操作のためのmethod定義
type ITaskRepository interface {
	GetAllTasks(tasks *[]model.Task, userId uint) error
	GetTaskById(task *model.Task, userId uint, taskId uint) error
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task, userId uint, taskId uint) error
	DeleteTask(userId uint, taskId uint) error
}

// interfaceの実装、DBとの連携
type taskRepository struct {
	db *gorm.DB
}

// 新しいtaskRepositoryのinstanceを作成、DB接続を受け取り、
// taskRepository instanceを初期化
func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db}
}

// (tr *taskRepository)にmethodの実装 特定のユーザーのtaskをすべて取得
// tasksはタスクリストを格納するためのスライスポインタ
func (tr *taskRepository) GetAllTasks(tasks *[]model.Task, userId uint) error {
	// tr.dbでDB接続 JoinsでUserテーブルを結合、関連するデータも取得可能
	// Find(tasks)で合致するtaskをtasksに格納
	if err := tr.db.Joins("User").Where("user_id=?", userId).Order("created_at").Find(tasks).Error; err != nil {
		return err
	}
	return nil
}

// taskidに紐づくtaskを検索、taskに格納
func (tr *taskRepository) GetTaskById(task *model.Task, userId uint, taskId uint) error {
	// Gorm method First(task, taskId) taskIdに対応する最初のtaskがtaskに格納される
	if err := tr.db.Joins("User").Where("user_id=?", userId).First(task, taskId).Error; err != nil {
		return err
	}
	return nil
}

// CreateaaTask methodを実装
func (tr *taskRepository) CreateTask(task *model.Task) error {
	if err := tr.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

// taskを更新するmethodを実装
func (tr *taskRepository) UpdateTask(task *model.Task, userId uint, taskId uint) error {
	// Clauses(clause.Returning{})でRETURNINGを有効に、更新後のデータを取得できる
	result := tr.db.Model(task).Clauses(clause.Returning{}).Where("id=? AND user_id=?", taskId, userId).Update("title", task.Title)
	if result.Error != nil {
		return result.Error
	}
	// RowsAffectedクエリの実行により変更された行の数を返す
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}

func (tr *taskRepository) DeleteTask(userId uint, taskId uint) error {
	// Delete(&model.Task{})を実行しresultに格納
	result := tr.db.Where("id=? AND user_id=?", taskId, userId).Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}
	// RowsAffectedクエリの実行により変更された行の数を返す
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}