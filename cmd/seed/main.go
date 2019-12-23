package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	// "fmt"
	"log"
	"os"

	// "strings"

	"github.com/famkampm/nentrytask/internal/models"
	"github.com/famkampm/nentrytask/pkg/helper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gopkg.in/guregu/null.v3"
	// "gopkg.in/guregu/null.v3"
)

func main() {
	// SEEDER DB
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err.Error())
	}

	db, err := OpenDB()
	if err != nil {
		log.Println("gagal open db. err:", err.Error())
		panic(err.Error())
	}
	defer CloseDB(db)
	err = CreateDB(db)
	if err != nil {
		log.Println("gagal create db. err:", err.Error())
		panic(err.Error())
	}
	err = ChooseDB(db)
	if err != nil {
		log.Println("gagal choose db. err:", err.Error())
		panic(err.Error())
	}

	err = CreateUserTable(db)
	if err != nil {
		log.Println("gagal create user db. err:", err.Error())
		panic(err.Error())
	}

	log.Println("DB aman")
	hashedPassword, err := helper.Hash("pass")
	if err != nil {
		log.Println("hashing pass fail")
		panic(err.Error())
	}
	pass := string(hashedPassword)
	totalRow := 5000000
	pembagi := 10000
	iterateOver := totalRow / pembagi
	rows := MakeRows(totalRow, &pass)
	for i := 1; i <= iterateOver; i++ {
		startIndex := i
		endIndex := i + pembagi
		err = BulkInsert(rows[startIndex:endIndex], db)
		if err != nil {
			log.Println("ERORR GAN: ", err.Error())
		}
	}

	log.Println("len rows:", len(rows))
}

func OpenDB() (*sql.DB, error) {
	drivername := os.Getenv("DB_DRIVER")
	// pathname := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")
	pathname := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/"

	log.Println("drivername", drivername)
	log.Println("pathname", pathname)
	db, err := sql.Open(drivername, pathname)
	if err != nil {
		log.Println("opendb open. err:", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Println("opendb ping. err", err.Error())
		return nil, err
	}
	return db, nil
}

func CloseDB(db *sql.DB) {
	db.Close()
}

func CreateDB(db *sql.DB) error {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS sys;")
	if err != nil {
		log.Println("createdb. err:", err.Error())
		return err
	}
	return nil
}

func ChooseDB(db *sql.DB) error {
	_, err := db.Exec("USE sys")
	if err != nil {
		log.Println("choosedb. err:", err.Error())
		return err
	}
	return nil
}

func CreateUserTable(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS user (id int not null auto_increment, username varchar(40) CHARACTER SET utf8mb4 not null, password varchar(240) not null, nickname varchar(240), profile_image varchar(240), PRIMARY KEY (id), index(username) )")
	if err != nil {
		log.Println("createuser table. prepare error:", err.Error())
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("createuser table. exec error:", err.Error())
		return err
	}
	return nil
}

func BulkInsert(unsavedRows []*models.User, db *sql.DB) error {
	valueStrings := make([]string, 0, len(unsavedRows))
	valueArgs := make([]interface{}, 0, len(unsavedRows)*4)
	for _, post := range unsavedRows {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, post.Username, post.Password, post.Nickname, post.ProfileImage)
	}
	stmt := fmt.Sprintf("INSERT INTO user (username, password, nickname, profile_image) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	if err != nil {
		log.Println("err bulkinsert:", err.Error())
	}
	return err
}

func MakeRows(n int, pass *string) []*models.User {
	rows := make([]*models.User, 0, n)
	for i := 1; i <= n; i++ {
		temp := strconv.FormatInt(int64(i), 10)
		row := &models.User{
			Username:     temp,
			Password:     *pass,
			Nickname:     null.StringFrom(temp),
			ProfileImage: null.StringFrom(temp),
		}
		rows = append(rows, row)
	}
	return rows
}
