package Utility

import (
	"bufio"
	"fmt"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var cookieHandle = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func IsValidUser(name string, password string) bool {
	_name:=strings.ToUpper(name)

	file, err := os.Open("Account.txt")
	if err != nil {
		log.Fatalf("failed to access acount info: %v", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		log.Println(scanner.Text())
		if scanner.Text() == fmt.Sprintf("ACCOUNT=%s|PASSWORD=%s",_name, password){
			return true
		}
	}
	return false
}


func SetCookie(eMail string, resp http.ResponseWriter){
	mapName := map[string] string{
		"email": eMail,
	}

	encode, err := cookieHandle.Encode("cookie", mapName)
	if err == nil {
		cookie := &http.Cookie{
			Name:       "cookie",
			Value:      encode,
			Path:       "/",
			MaxAge:		60,
			Expires:	time.Now().Add(60*time.Second),
		}
		http.SetCookie(resp, cookie)
	}
}

func GetUserName(req *http.Request) (userName string, err error){
	cookie, err := req.Cookie("cookie")
	if err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandle.Decode("cookie", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["email"]
		}
	}else{
		log.Printf("failed to get user name: %v",err)
		return "", err
	}
	return userName, nil
}