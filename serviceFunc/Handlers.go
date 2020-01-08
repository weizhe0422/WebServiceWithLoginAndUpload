package serviceFunc

import (
	"fmt"
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
			Utility.SetCookie(email, resp)
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
	userName, err := Utility.GetUserName(request)
	if err != nil {
		http.Error(resp, "expired!", http.StatusForbidden)
		fmt.Fprintln(resp, err)
	}
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
		request.ParseMultipartForm(10000000)
		formData := request.MultipartForm
		files := formData.File["multiplefiles"]

		respString := ""
		for i, _ := range files{
			file, err := files[i].Open()
			if err != nil {
				log.Printf("failed to load file: %v",err)
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			content, err := ioutil.ReadAll(file)
			if err != nil {
				log.Printf("failed to read file: %v",err)
				http.Error(resp, err.Error(), http.StatusInternalServerError)
				return
			}
			err = ioutil.WriteFile(files[i].Filename, content, os.ModeAppend)
			if err != nil {
				log.Printf("failed to get upload file: %v", err)
			}

			if sizeInterface, ok := file.(Size); ok{
				respString = respString + fmt.Sprintf(`檔名:%s 檔案大小為：%d`, files[i].Filename, sizeInterface.Size())+"<br>"
			}
		}

		userName, err := Utility.GetUserName(request)
		if err != nil {
			http.Error(resp, "expired!", http.StatusForbidden)
		}
		compile := regexp.MustCompile(`(\w+([-+.]\w+)*)@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
		userName = compile.FindAllStringSubmatch(userName, -1)[0][1]
		indexBody, _ := Utility.LoadFile("templates/index.html")
		fmt.Fprintf(resp, indexBody, userName, respString)
	}
}

