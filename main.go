package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/usecase"
	"go-rest-api/validator"
)

func main() {
	// 返却されたインスタンスを通じてDBの操作を行う
	db := db.NewDB()
	// userValidatorにuserValidatorのインスタンスを代入
	userValidator := validator.NewUserValidator()
	// taskValidatorにtaskValidatorのインスタンスを代入
	taskValidator := validator.NewTaskValidator()
	// それぞれRepositoryにrepositoryインスタンスを
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	// それぞれのinstanceを作成(repositoryはDBアクセスを含める)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)
	e := router.NewRouter(userController, taskController)
	e.Logger.Fatal(e.Start(":8080"))
}