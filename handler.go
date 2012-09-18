package main

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func newPassword() string {
	password := bson.NewObjectId().Hex()
	salt := []byte("mongoapi")
	return fmt.Sprintf("%x", pbkdf2.Key([]byte(password), salt, 4096, len(salt)*8, sha512.New))
}

func Add(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func Bind(w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Query().Get(":name")
	database := Session.DB(name)
	err := database.AddUser(name, "", false)
	if err != nil {
		return err
	}
	data := map[string]string{
		"MONGO_URI":           "127.0.0.1:27017",
		"MONGO_USER":          name,
		"MONGO_PASSWORD":      newPassword(),
		"MONGO_DATABASE_NAME": name,
	}
	b, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	fmt.Fprint(w, string(b))
	w.WriteHeader(http.StatusCreated)
	return nil
}

func Unbind(w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Query().Get(":name")
	database := Session.DB(name)
	err := database.RemoveUser(name)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func Remove(w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Query().Get(":name")
	err := Session.DB(name).DropDatabase()
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func Status(w http.ResponseWriter, r *http.Request) error {
	_, err := mgo.Dial("localhost:27017")
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

type Handler func(http.ResponseWriter, *http.Request) error

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
