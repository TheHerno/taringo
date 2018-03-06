package controllers

import (
	"strconv"
	"errors"
	"log"
	"database/sql"
	"../models"
	_ "github.com/go-sql-driver/mysql"
)

const MysqlString = "taringo:taringo*23@tcp(127.0.0.1:3306)/taringo"

func getPostsList(page int) ([]models.PostI,error) {
	var posts []models.PostI

	db , err := sql.Open("mysql", MysqlString)
	if err != nil {
		log.Println(err)
		return nil , errors.New("Error conectando a la base de datos.")
	}
	defer db.Close()
	var offset string
	if page == 0 {
		offset = "0"
	} else {
		offset = strconv.Itoa(page*20-1)
	}
	rows, err := db.Query("SELECT id_post,title,id_categoria FROM posts ORDER BY id_post DESC LIMIT "+offset+",20")
	if err != nil {
		log.Println(err)
		return nil , errors.New("Error obteniendo posts.")
	}

	for rows.Next() {
		var postI models.PostI
		err = rows.Scan(&postI.IdPost,&postI.Titulo,&postI.IdCategoria)
		if err != nil {
			log.Println(err)
			return nil , errors.New("Error inesperado.")
		}
		posts = append(posts, postI)
	}

	return posts , nil
} 