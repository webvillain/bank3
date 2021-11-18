package memdb

import (
	"bank/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
)

var database *sql.DB

func ConnectDb() {
	os.Remove("./bank.db")
	db, err := sql.Open("sqlite3", "./bank.db")
	if err != nil {
		log.Fatal(err)
	}
	database = db

	fmt.Println("Connection To Database Is Success !")
	db.Ping()
	ST, err := db.Prepare("CREATE TABLE IF NOT EXISTS USERS(ID INTEGER PRIMARY KEY, Name TEXT NOT NULL , Email TEXT NOT NULL );")
	if err != nil {
		log.Fatal(err)
	}
	ST.Exec()
	fmt.Println("Table Name USERS Is Created Successfully !")

}

type Database struct {
	*sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		DB: db,
	}
}

var ErrNotImplemented = errors.New("not implemnted")

func (d *Database) User(ctx context.Context, id string) (*db.User, error) {
	var singleuser db.User

	rows, err := database.Query("SELECT * FROM USERS WHERE ID = ?")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&singleuser.ID, &singleuser.Name, &singleuser.Email)
	}
	return &singleuser, nil

}

/*
{"Name":"x","Email":"Y"}
*/
func (d *Database) CreateUser(ctx context.Context, u *db.User) error {
	// var id uuid.UUID
	// var name string
	// var email string
	// fmt.Println("ENTER USER NAME : ")
	// fmt.Scanf("%g", &name)
	// fmt.Println("ENTER USER EMAIL: ")
	// fmt.Scanf("%g", &email)

	stmt, err := database.Prepare("INSERT INTO USERS (ID , Name , Email) VALUES (?,?,?) ;")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(u.ID, u.Name, u.Email)
	if err != nil {
		fmt.Println("error while creating new user ")
	}
	fmt.Println("user Created Successfully !")
	fmt.Println(res.RowsAffected())
	return nil
}

/*
{"Name":"x1","Email":"Y1"}
*/
func (d *Database) UpdateUser(ctx context.Context, u *db.User) error {
	// _, ok := d.users[u.ID]
	// if !ok {
	// 	return db.ErrNotFound
	// }
	// d.users[u.ID] = u
	// return nil

	myquery4 := `
	UPDATE USERS
	SET Name = ?, AGE = ?, Email = ?
	WHERE ID = ?;
	
	`
	var user *db.User
	fmt.Println("Enter User ID For Which you Want To Update Data ")
	fmt.Scanf("%d", &user.ID)
	fmt.Println("Enter User Name :")
	fmt.Scanf("%s", &user.Name)
	fmt.Println("Enter User Email :")
	fmt.Scanf("%s", &user.Email)
	stmt, err := database.Prepare(myquery4)
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(user.Name, user.Email, &user.ID)
	if err != nil {
		log.Fatal(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("No.Of Rows Affected !", n)
	return nil
}

func (d *Database) DeleteUser(ctx context.Context, id string) error {
	// _, ok := d.users[id]
	// if !ok {
	// 	return db.ErrNotFound
	// }
	// delete(d.users, id)
	// return nil
	// var I int
	// fmt.Println("Please Enter The Id Of Which User You Want To Delete From Database .")
	// fmt.Scanf("%d", &I)
	myquery3 := `
	DELETE FROM USERS WHERE ID = (?);
	`

	stmt, err := database.Prepare(myquery3)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(id)
	stmt.Close()
	fmt.Println("User Successfully Deleted ")
	return nil
}

// UPDATING EXISTING USER IN DATABASE

func (d *Database) ListUsers(ctx context.Context) ([]*db.User, error) {
	var user *db.User
	rows, er := database.Query("SELECT * FROM USERS ;")
	if er != nil {
		log.Fatal(er)
	}
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Email)

	}
	var users []*db.User
	users = append(users, user)
	return users, nil

}
