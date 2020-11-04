package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/task_checker_api/src/model"
)

//GetTasksHandler タスク取得
var GetTasksHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userJwt := r.Context().Value("user")
	mailAddress := userJwt.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	tasks := GetTasks(mailAddress)
	json.NewEncoder(w).Encode(tasks)
})

//GetTasks 一覧取得
func GetTasks(mailAddress string) []model.Task {
	db := model.DBConnect()
	result, err := db.Query(`
		SELECT
			task.id, task.title, task.done, task.del_flag, task.created_at, task.updated_at
		FROM 
			task 
			INNER JOIN users 
			ON task.user_id = users.id 
		WHERE users.mail_address = $1 
			AND task.del_flag = false
		ORDER BY sort;`, mailAddress)
	if err != nil {
		panic(err.Error())
	}

	tasks := []model.Task{}
	for result.Next() {
		task := model.Task{}
		var id uint
		var title string
		var done, delFlag bool
		var createdAt, updatedAt time.Time

		err = result.Scan(&id, &title, &done, &delFlag, &createdAt, &updatedAt)
		if err != nil {
			panic(err.Error())
		}

		task.ID = id
		task.Title = title
		task.Done = done
		task.DelFlag = delFlag
		task.CreatedAt = createdAt
		task.UpdatedAt = updatedAt
		tasks = append(tasks, task)
	}
	db.Close()
	return tasks
}

//FindTaskByID タスク検索
func FindTaskByID(id uint) model.Task {
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

//PutDone チェックボックストグル
var PutDone = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userJwt := r.Context().Value("user")
	mailAddress := userJwt.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	if mailAddress == "" {
		panic("no authorized")
	}
	var task model.Task
	json.NewDecoder(r.Body).Decode(&task)
	db := model.DBConnect()
	_, err := db.Exec("UPDATE task SET done = $1 WHERE id = $2", task.Done, task.ID)
	if err != nil {
		panic(err.Error())
	}
	db.Close()
	json.NewEncoder(w).Encode("success")
})

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

	task := FindTaskByID(uint(id))

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
