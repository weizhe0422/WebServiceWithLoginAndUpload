package serviceFunc

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/weizhe0422/WebServiceWithLoginAndUpload/Utility"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Stat interface {
	Stat() (os.FileInfo, error)
}
type Size interface {
	Size() int64
}

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
	email := request.FormValue("email")
	password := request.FormValue("password")
	redirectTarget := "/"
	log.Println(email,password)
	if !Utility.IsEmpty(email) && !Utility.IsEmpty(password) {
		if Utility.IsValidUser(email, password) {
			log.Println("Login OK!")
			SetCookie(email, resp)
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
	compile := regexp.MustCompile(`(\w+([-+.]\w+)*)@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	submatch := compile.FindAllStringSubmatch(userName, -1)
	userName = submatch[0][1]
	if !Utility.IsEmpty(userName) {
		indexBody, _ := Utility.LoadFile("templates/index.html")
		fmt.Fprintf(resp, indexBody, userName, "")
	}else{
		http.Redirect(resp, request, "/", http.StatusFound)
	}
}

func Upload(resp http.ResponseWriter, request *http.Request){
	if request.Method == "POST" {
		file, head, err := request.FormFile("userfile")
		defer file.Close()
		if err != nil {
			log.Printf("failed to load file: %v",err)
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("success to load file")

		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("failed to read file: %v",err)
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("success to read file")

		err = ioutil.WriteFile(head.Filename, content, os.ModeAppend)
		if err != nil {
			log.Printf("failed to get upload file: %v", err)
		}

		userName := GetUserName(request)
		compile := regexp.MustCompile(`(\w+([-+.]\w+)*)@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
		submatch := compile.FindAllStringSubmatch(userName, -1)
		userName = submatch[0][1]
		indexBody, _ := Utility.LoadFile("templates/index.html")
		if sizeInterface, ok := file.(Size); ok{
			fmt.Fprintf(resp, indexBody, userName, fmt.Sprintf(`上次上傳資訊： 檔名:%s 檔案大小為：%d`, head.Filename, sizeInterface.Size()))
		}

	}
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
		}
		http.SetCookie(resp, cookie)
	}
}

func GetUserName(req *http.Request) (userName string){
	cookie, err := req.Cookie("cookie")
	if err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandle.Decode("cookie", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["email"]
		}
	}
	return userName
}