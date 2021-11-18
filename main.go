package main

import (
	"bank/db"
	"bank/db/memdb"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func main() {
	memdb.ConnectDb()
	http.HandleFunc("/users", userHandler)
	http.ListenAndServe(":8080", nil)
}

/*
{
	"Name":"vijay",
	"Email":"vijay@test.com"
}
*/
// for what this DB used for
var DB memdb.Database

func userHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// TODO: crete user
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		u := &db.User{}
		err = json.Unmarshal(data, u)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}
		u, err = DB.CreateUser(r.Context(), u)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		data, err = json.Marshal(u)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(data)
	}
	if r.Method == http.MethodGet {
		users, err := DB.ListUsers(r.Context())
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		data, err := json.Marshal(users)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(data)
	}
	if r.Method == http.MethodDelete {
		id := r.URL.Query().Get("ID")
		err := DB.DeleteUser(r.Context(), id)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write([]byte("user is deleted"))
	}
	if r.Method == http.MethodPatch {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		u := &db.User{}
		err = json.Unmarshal(data, u)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}
		us, err := DB.User(r.Context(), u.ID)
		if errors.Is(err, db.ErrNotFound) {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(err.Error()))
			return
		}
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		changed := false
		if u.Name != us.Name {
			changed = true
		}
		if u.Email != us.Email {
			changed = true
		}
		if !changed {
			return
		}
		err = DB.UpdateUser(r.Context(), u)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		data, err = json.Marshal(u)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(data)
	}
}
