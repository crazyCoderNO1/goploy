package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/zhenorzz/goploy/core"
	"github.com/zhenorzz/goploy/model"
)

// User 用户字段
type User struct{}

// Login user login api
func (user *User) Login(w http.ResponseWriter, r *http.Request) {
	type ReqData struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	type RepData struct {
		Token string `json:"token"`
	}
	var reqData ReqData
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &reqData)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	model := model.User{Account: reqData.Account, Password: reqData.Password}
	err = model.Vaildate()
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	token, err := user.createToken(model.ID)
	if err != nil {
		response := core.Response{Code: 1, Message: err.Error()}
		response.Json(w)
		return
	}
	data := RepData{Token: token}
	response := core.Response{Data: data}
	response.Json(w)
}

// Info get user info api
func (user *User) Info(w http.ResponseWriter, r *http.Request) {
	fmt.Println(core.GolbalUserID)
	type RepData struct {
		UserInfo struct {
			ID uint32 `json:"id"`
		} `json:"userInfo"`
	}
	data := RepData{}
	data.UserInfo.ID = core.GolbalUserID
	response := core.Response{Data: data}
	response.Json(w)
}

func (user *User) createToken(id uint32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"nbf": time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGN_KEY")))

	//Sign and get the complete encoded token as string
	return tokenString, err
}
