package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

<<<<<<< HEAD
	"../model"
=======
	"github.com/task_checker_api/src/model"
>>>>>>> 225fbed8e2604733903d14e234c2e48bea9df3b4

	"github.com/gin-gonic/gin"
)

//TasksGET 一覧取得
func TasksGET(c *gin.Context) {
<<<<<<< HEAD
	print("aaa")
=======
>>>>>>> 225fbed8e2604733903d14e234c2e48bea9df3b4
	db := model.DBConnect()
	result, err := db.Query("SELECT * FROM task ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	// json返却用
	tasks := []model.Task{}
	for result.Next() {
		task := model.Task{}
		var id uint
		var createdAt, updatedAt time.Time
		var title string

		err = result.Scan(&id, &createdAt, &updatedAt, &title)
		if err != nil {
			panic(err.Error())
		}

		task.ID = id
		task.CreatedAt = createdAt
		task.UpdatedAt = updatedAt
		task.Title = title
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

//FindByID タスク検索
func FindByID(id uint) model.Task {
	db := model.DBConnect()
	result, err := db.Query("SELECT * FROM task WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	task := model.Task{}
	for result.Next() {
		var createdAt, updatedAt time.Time
		var title string

		err = result.Scan(&id, &createdAt, &updatedAt, &title)
		if err != nil {
			panic(err.Error())
		}

		task.ID = id
		task.CreatedAt = createdAt
		task.UpdatedAt = updatedAt
		task.Title = title
	}
	return task
}

//TaskPOST タスク登録
func TaskPOST(c *gin.Context) {
	db := model.DBConnect()

	title := c.PostForm("title")
	now := time.Now()

	_, err := db.Exec("INSERT INTO task (title, created_at, updated_at) VALUES(?, ?, ?)", title, now, now)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("post sent. title: %s", title)
}

//TaskPATCH タスク更新
func TaskPATCH(c *gin.Context) {
	db := model.DBConnect()

	id, _ := strconv.Atoi(c.Param("id"))
	title := c.PostForm("title")
	now := time.Now()

	_, err := db.Exec("UPDATE task SET title = ?, updated_at=? WHERE id = ?", title, now, id)
	if err != nil {
		panic(err.Error())
	}

	task := FindByID(uint(id))

	fmt.Println(task)
	c.JSON(http.StatusOK, gin.H{"task": task})
}

//TaskDELETE タスク削除
func TaskDELETE(c *gin.Context) {
	db := model.DBConnect()

	id, _ := strconv.Atoi(c.Param("id"))

	_, err := db.Query("DELETE FROM task WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, "deleted")
}
