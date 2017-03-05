package main

import (
	"fmt"
	r "github.com/dancannon/gorethink"
	"os"
)

func main() {
	session, err := r.Connect(r.ConnectOpts{
		Address:  os.Getenv("DBHOST"),
		Database: os.Getenv("DBNAME"),
	})
	if err != nil {
		panic("Can't init db.")
	}
	if os.Getenv("APP_ENV") == "production" {
		initProduction(session)
		fmt.Println("Production db inited.")
	} else {
		initDevelopment(session)
		fmt.Println("Development db inited.")
	}
}

func initProduction(session *r.Session) {
	r.DBCreate("hugito").Exec(session)
	r.TableCreate("user", r.TableCreateOpts{
		PrimaryKey: "login",
	}).Exec(session)
	r.TableCreate("repository", r.TableCreateOpts{
		PrimaryKey: "id",
	}).Exec(session)
	r.TableCreate("content", r.TableCreateOpts{
		PrimaryKey: "id",
	}).Exec(session)
}

func initDevelopment(session *r.Session) {
	r.DBCreate("development").Exec(session)
	r.TableCreate("user", r.TableCreateOpts{
		PrimaryKey: "login",
	}).Exec(session)
	r.TableCreate("repository", r.TableCreateOpts{
		PrimaryKey: "id",
	}).Exec(session)
	r.TableCreate("content", r.TableCreateOpts{
		PrimaryKey: "id",
	}).Exec(session)
}
