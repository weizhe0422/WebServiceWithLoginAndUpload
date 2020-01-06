package serviceFunc

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/weizhe0422/WebServiceWithLoginAndUpload/Utility"
	"log"
	"net/http"
)

var cookieHandle = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func LoginPage(resp http.ResponseWriter, request *http.Request){
	file, err := Utility.LoadFile("templates/login.html")
	if err != nil{
		log.Printf("failed to load login templates: %v", err)
	}
	fmt.Fprintf(resp, file)
}

func Login(resp http.ResponseWriter, request *http.Request){
	name := request.FormValue("name")
	password := request.FormValue("password")
	redirectTarget := "/"
	log.Println(name,password)
	if !Utility.IsEmpty(name) && !Utility.IsEmpty(password) {
		if Utility.IsValidUser(name, password) {
			log.Println("Login OK!")
			SetCookie(name, resp)
			redirectTarget = "/welcome"
		}else{
			log.Println("Login fail!")
			//TODO Wait to implement register service
			redirectTarget = "/register"
		}
	}else{
		log.Println("Empty!")
	}
	http.Redirect(resp, request, redirectTarget, http.StatusFound)
}

func Welcome(resp http.ResponseWriter, request *http.Request){
	userName := GetUserName(request)
	if !Utility.IsEmpty(userName) {
		indexBody, _ := Utility.LoadFile("templates/index.html")
		fmt.Fprintf(resp, indexBody, userName)
	}else{
		http.Redirect(resp, request, "/", http.StatusFound)
	}
}

func SetCookie(userName string, resp http.ResponseWriter){
	mapName := map[string] string{
		"name": userName,
	}

	encode, err := cookieHandle.Encode("cookie", mapName)
	if err == nil {
		cookie := &http.Cookie{
			Name:       "cookie",
			Value:      encode,
			Path:       "/",
		}
		http.SetCookie(resp, cookie)
	}
}

func GetUserName(req *http.Request) (userName string){
	cookie, err := req.Cookie("cookie")
	if err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandle.Decode("cookie", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}