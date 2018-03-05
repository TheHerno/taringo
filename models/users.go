package models

import (
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)

const MysqlString = "taringo:taringo*23@tcp(127.0.0.1:3306)/taringo"

type User struct {
	UserId		int `json:"id_user"`
	Username	string `json:"username"`
	Password	string `json:"password"`
	Email		string `json:"email"`
	DateReg		string `json:"date_reg"`
	Rango		int `json:"rango"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegistrarUsuario(username,password,email string) error {
	db , err := sql.Open("mysql", MysqlString)
	if err != nil {
		log.Println(err)
		return errors.New("Error conectando a la base de datos.")
	}
	defer db.Close()
	password , _ = HashPassword(password)
	fecha := time.Now().Local().Format("2006-01-02")
	user := User{
		Username:	username,
		Password:	password,
		Email:		email,
		Rango:		1,
		DateReg:	fecha,
	}

	rows, err := db.Query("SELECT * FROM users WHERE username='"+user.Username+"'")
	if err != nil {
		log.Println(err)
		return errors.New("Error consultando la base de datos.")
	}
	
	if rows.Next() {
		return errors.New("Ese usuario ya existe.")
	}

	stmt, err := db.Prepare("INSERT INTO users SET username=? , password=?, email=? , date_reg=? , rango=?")
	if err != nil {
		log.Println(err)
		return errors.New("Error preparando la solicitud a la base de datos.")
	}
	_ , err = stmt.Exec(user.Username,user.Password,user.Email,user.DateReg,user.Rango)
	if err != nil {
		log.Println(err)
		return errors.New("Error insertando en la base de datos.")
	}
	return nil
}

func Login(username,password string) (int , error) {
	db , err := sql.Open("mysql", MysqlString)
	if err != nil {
		log.Println(err)
		return 0,  errors.New("Error insertando en la base de datos.")
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE username='"+username+"'")
	if err != nil {
		log.Println(err)
		return 0, errors.New("Error consultando la base de datos.")
	}
	
	if !rows.Next() {
		return 0, errors.New("Usuario no encontrado.")
	}
	var user User
	err = rows.Scan(&user.UserId,&user.Username,&user.Password,&user.Email,&user.DateReg,&user.Rango)
	if err != nil {
		log.Println(err)
		return 0, errors.New("Error inesperado.")
	}
	if !CheckPasswordHash(password, user.Password){
		return 0, errors.New("Contraseña Incorrecta.")
	}
	return user.UserId, nil
}