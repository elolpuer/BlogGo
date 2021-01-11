package post

import (
	"database/sql"
	"fmt"
	"github.com/elolpuer/Blog/pkg/models"
	"time"
)


func Add(db *sql.DB, userID int, text string) error {
	year, month, day := time.Now().Date()
	hour, min, _ := time.Now().Clock()
	date:=fmt.Sprintf("%d:%d %d %s %d", hour, min, day, month, year)
	_, err := db.Exec("INSERT INTO posts (user_id,body, time_date) VALUES ($1, $2,$3)", userID,text, date)
	if err != nil {
		return err
	}
	return nil
}

func Posts(db *sql.DB, userID int) ([]*models.Post,error) {
	rows, err :=db.Query("SELECT * FROM posts WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	posts := make([]*models.Post, 0)
	defer rows.Close()
	for rows.Next(){
		post := new(models.Post)
		err := rows.Scan(&post.ID, &post.UserID,&post.Body, &post.Date)
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

func DeletePost(db *sql.DB, ID string, userID int) error {
	_, err := db.Exec("DELETE FROM  posts WHERE id=$1 AND user_id=$2", ID, userID)
	if err != nil {
		return err
	}
	return nil
}