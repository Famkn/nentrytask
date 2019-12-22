package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_userHttpDeliver "github.com/famkampm/nentrytask/internal/user/delivery/http"
	"github.com/famkampm/nentrytask/internal/user/repository"
	"github.com/famkampm/nentrytask/internal/user/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func init() {
	godotenv.Load()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("init main success")
}

func main() {
	db := initDB()
	defer func() {
		closeDB(db)
		log.Println("CLOSING DB CON")
	}()
	// conn := initRedis()
	// defer func() {
	// 	log.Println("closing redis conection")
	// 	conn.Close()
	// }()
	redisPool := initRedisPool()
	userRepoMysql := repository.NewMysqlUserRepository(db)
	userRepoRedis := repository.NewRedisUserRepository(redisPool)
	userUsecase := usecase.NewUserUsecase(userRepoMysql, userRepoRedis)
	router := httprouter.New()

	_userHttpDeliver.NewUserHandler(router, userUsecase)

	// run server
	log.Fatal(http.ListenAndServe(":8080", router))

}

// func main() {
// 	drivername := os.Getenv("DB_DRIVER")
// 	pathname := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")

// 	db, err := sql.Open(drivername, pathname)
// 	if err != nil {
// 		log.Println("ERROR OPEN DB:", err.Error())
// 		panic(err.Error())
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		log.Println("DB PING ERROR:", err.Error())
// 		panic(err.Error())
// 	}
// 	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user ( id int not null AUTO_INCREMENT, username varchar(240) not null, password varchar(240) not null, nickname varchar(240), profile_image varchar(240), PRIMARY KEY (id) )")
// 	if err != nil {
// 		log.Println("CREATE TABLE USER err:", err.Error())
// 		panic(err.Error())
// 	}
// 	defer func() {
// 		err := db.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	conn, err := redis.Dial("tcp", "localhost:6379")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer func() {
// 		err = conn.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}()
// 	userRepoMysql := repository.NewMysqlUserRepository(db)
// 	userRepoRedis := repository.NewRedisUserRepository(conn)
// 	userUsecase := usecase.NewUserUsecase(userRepoMysql, userRepoRedis)
// 	router := httprouter.New()
// 	log.Println("MAIN INIT DONE")
// 	_userHttpDeliver.NewUserHandler(router, userUsecase)
// 	log.Fatal(http.ListenAndServe(":8080", router))

// }

func initDB() *sql.DB {
	drivername := os.Getenv("DB_DRIVER")
	pathname := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")

	db, err := sql.Open(drivername, pathname)
	if err != nil {
		log.Println("ERROR OPEN DB:", err.Error())
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Println("DB PING ERROR:", err.Error())
		panic(err.Error())
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user ( id int not null AUTO_INCREMENT, username varchar(240) not null, password varchar(240) not null, nickname varchar(240), profile_image varchar(240), PRIMARY KEY (id) )")
	if err != nil {
		log.Println("CREATE TABLE USER err:", err.Error())
		panic(err.Error())
	}
	return db
}

func closeDB(db *sql.DB) {
	log.Println("closing db connection")
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func initRedis() redis.Conn {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal(err)
	}
	// Importantly, use defer to ensure the connection is always
	// properly closed before exiting the main() function.
	return conn
}

func initRedisPool() *redis.Pool {

	pool := &redis.Pool{
		MaxIdle:     80,
		MaxActive:   12000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
	return pool
}
