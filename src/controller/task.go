package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"main/src/model"

	"github.com/dgrijalva/jwt-go"
)

var layout = "2006-01-02 15:04:05"

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
			task.id, task.user_id, task.title, task.done, task.del_flag, task.sort, 
			task.created_at, task.updated_at
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
		var userId int
		var title string
		var done, delFlag bool
		var sort int
		var createdAt, updatedAt time.Time

		err = result.Scan(&id, &userId, &title, &done, &delFlag, &sort, &createdAt, &updatedAt)
		if err != nil {
			panic(err.Error())
		}

		task.ID = id
		task.UserID = userId
		task.Title = title
		task.Done = done
		task.DelFlag = delFlag
		task.Sort = sort
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

//PostTask タスク新規追加
var PostTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userJwt := r.Context().Value("user")
	mailAddress := userJwt.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	if mailAddress == "" {
		panic("no authorized")
	}
	var task model.Task
	json.NewDecoder(r.Body).Decode(&task)
	db := model.DBConnect()
	_, err := db.Exec("INSERT INTO task (user_id, title, sort, created_at, updated_at) SELECT id, $1, $2, now(), now() FROM users WHERE mail_address=$3;", task.Title, task.
		Sort, mailAddress)
	if err != nil {
		panic(err.Error())
	}
	db.Close()
	tasks := GetTasks(mailAddress)
	if len(tasks) > 0 {
		updateSort(tasks)
	}
	json.NewEncoder(w).Encode("success")
})

//PutTask タスクタイトル更新
var PutTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userJwt := r.Context().Value("user")
	mailAddress := userJwt.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	if mailAddress == "" {
		panic("no authorized")
	}
	var task model.Task
	json.NewDecoder(r.Body).Decode(&task)
	db := model.DBConnect()
	_, err := db.Exec("UPDATE task SET title = $1 WHERE id = $2", task.Title, task.ID)
	if err != nil {
		panic(err.Error())
	}
	db.Close()
	json.NewEncoder(w).Encode("success")
})

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

//DeleteTask タスク削除
var DeleteTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userJwt := r.Context().Value("user")
	mailAddress := userJwt.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	if mailAddress == "" {
		panic("no authorized")
	}
	var task model.Task
	json.NewDecoder(r.Body).Decode(&task)
	db := model.DBConnect()
	_, err := db.Exec("DELETE FROM task WHERE id = $1", task.ID)
	if err != nil {
		panic(err.Error())
	}
	db.Close()
	tasks := GetTasks(mailAddress)
	if len(tasks) > 0 {
		updateSort(tasks)
	}
	json.NewEncoder(w).Encode("success")
})

//UpSort 並べ替え
var UpSort = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	userJwt := r.Context().Value("user")
	mailAddress := userJwt.(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	if mailAddress == "" {
		panic("no authorized")
	}
	var tasks []model.Task
	json.NewDecoder(r.Body).Decode(&tasks)
	updateSort(tasks)
	json.NewEncoder(w).Encode("success")
})

//updateSort タスクのソート更新
func updateSort(tasks []model.Task) {
	sql := "INSERT INTO task VALUES "
	for index, task := range tasks {
		value := ""
		if index != 0 {
			value = " ,"
		}
		done := "false"
		if task.Done {
			done = "true"
		}
		delFlag := "false"
		if task.DelFlag {
			delFlag = "true"
		}
		value += fmt.Sprintf("(%d,%d,'%s','%s',%s,%d,'%s','%s')", task.ID, task.UserID, task.Title, done, delFlag, index+1, task.CreatedAt.Format(layout), task.UpdatedAt.Format(layout))
		sql = sql + value
	}
	sql = sql + " ON CONFLICT (id) DO UPDATE SET sort = EXCLUDED.sort"
	fmt.Println(sql)
	db := model.DBConnect()
	_, err := db.Exec(sql)
	if err != nil {
		panic(err.Error())
	}
	db.Close()
}
