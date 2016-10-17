package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type UserError struct {
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}

type UserServer struct {
	db       *gorm.DB
	userRepo *UserRepo
}

func NewUserServer(db *gorm.DB) *UserServer {
	return &UserServer{
		db:       db,
		userRepo: NewUserRepo(db),
	}
}

func (u *UserServer) HandlePost(w http.ResponseWriter, r *http.Request) {

	ninja := User{}
	json.NewDecoder(r.Body).Decode(&ninja)
	ninja.UserType = ROLE_NINJA
	user := context.Get(r, "user")
	ninja.Token = user.(*jwt.Token).Claims.(jwt.MapClaims)["session_token"].(string)
	u.db.Create(&ninja)
	json.NewEncoder(w).Encode(ninja)
}

func (u *UserServer) HandlePut(w http.ResponseWriter, r *http.Request) {
	var err error

	ninjaId, err := parseIdFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserError{
			Code:   http.StatusBadRequest,
			Errors: []string{"Ninja ID must be a number"},
		})
		return
	}

	user := context.Get(r, "user")
	ninja, err := u.userRepo.GetNinja(user.(*jwt.Token).Claims.(jwt.MapClaims)["session_token"].(string), ninjaId)
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(UserError{
			Code:   404,
			Errors: []string{fmt.Sprintf("JSON formatting error: %s", err.Error())},
		})
		return
	}

	updates := &User{}

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(UserError{
			Code:   400,
			Errors: []string{err.Error()},
		})
		return
	}

	updates.Token = user.(*jwt.Token).Claims.(jwt.MapClaims)["session_token"].(string)
	u.db.Model(&ninja).Updates(&updates)
	json.NewEncoder(w).Encode(ninja)
}

func (u *UserServer) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	ninjas, err := u.userRepo.GetNinjas(user.(*jwt.Token).Claims.(jwt.MapClaims)["session_token"].(string))
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(UserError{
			Code:   404,
			Errors: []string{err.Error()},
		})
	}
	json.NewEncoder(w).Encode(ninjas)
}

func (u *UserServer) HandleGet(w http.ResponseWriter, r *http.Request) {
	ninjaId, err := parseIdFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserError{
			Code:   http.StatusBadRequest,
			Errors: []string{"Ninja ID must be a number"},
		})
		return
	}

	user := context.Get(r, "user")
	ninja, err := u.userRepo.GetNinja(user.(*jwt.Token).Claims.(jwt.MapClaims)["session_token"].(string), ninjaId)
	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(UserError{
			Code:   404,
			Errors: []string{err.Error()},
		})
		return
	}

	json.NewEncoder(w).Encode(ninja)
}

func (u *UserServer) HandleDelete(w http.ResponseWriter, r *http.Request) {
	ninjaId, err := parseIdFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserError{
			Code:   http.StatusBadRequest,
			Errors: []string{"Ninja ID must be a number"},
		})
		return
	}

	user := context.Get(r, "user")

	if u.db.Where("token = ? AND id = ?", user.(*jwt.Token).Claims.(jwt.MapClaims)["session_token"].(string), ninjaId).Delete(&User{}).RecordNotFound() {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}

func parseIdFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, errors.New("Ninja ID must be a number")
	}

	return id, nil
}
