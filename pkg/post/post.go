package post

import (
	"database/sql"
	"fmt"
	"github.com/elolpuer/Blog/pkg/models"
	"time"
)

func Page() *models.IndexResp {
	Resp := new(models.IndexResp)
	Resp.Index = "Add Post"
	Resp.Title = "Add"
	return Resp
}

func Add(db *sql.DB, text string) error {
	year, month, day := time.Now().Date()
	hour, min, _ := time.Now().Clock()
	date:=fmt.Sprintf("%d:%d %d %s %d", hour, min, day, month, year)
	fmt.Println(date)
	fmt.Println(text)
	_, err := db.Exec("INSERT INTO posts (body, time_date) VALUES ($1, $2)", text, date)
	if err != nil {
		return err
	}
	return nil
}

func Posts(db *sql.DB) ([]*models.Post,error) {
	rows, err :=db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	posts := make([]*models.Post, 0)
	defer rows.Close()
	for rows.Next(){
		post := new(models.Post)
		err := rows.Scan(&post.ID, &post.Body, &post.Date)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func DeletePost(db *sql.DB, ID string) error {
	_, err := db.Exec("DELETE FROM  posts WHERE id=$1", ID)
	if err != nil {
		return err
	}
	return nil
}