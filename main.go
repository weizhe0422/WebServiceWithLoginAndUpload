package main

import (
	"github.com/gorilla/mux"
	"github.com/weizhe0422/WebServiceWithLoginAndUpload/serviceFunc"
	"log"
	"net/http"
)

var router = mux.NewRouter()

func main() {

	router.HandleFunc("/", serviceFunc.LoginPage)
	router.HandleFunc("/welcome", serviceFunc.Welcome)
	router.HandleFunc("/login", serviceFunc.LoginWithSvrVerify).Methods("POST")
	router.HandleFunc("/loginwithoauth", serviceFunc.LoginWithOAUTH).Methods("POST")
	router.HandleFunc("/registerPage", serviceFunc.RegisterPage)
	router.HandleFunc("/register", serviceFunc.Register)
	router.HandleFunc("/upload", serviceFunc.Upload)
	router.HandleFunc("/AWSupload", serviceFunc.AWSUpload)
	router.HandleFunc("/oauth/redirect", serviceFunc.OAuthLogin)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
