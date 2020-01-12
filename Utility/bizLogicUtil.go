package Utility

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gorilla/securecookie"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var cookieHandle = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getMongoDB() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	return mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
}

func IsValidUser(name string, password string) bool {
	/*dbUtil, err := getMongoDB()
	if err != nil {
		log.Printf("failed to connect to mongodb: %v", err)
		return false
	}
	collection := dbUtil.Database("test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil { log.Fatal("failed", err) }
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		log.Println("start to print data")
		var result bson.M
		err := cur.Decode(&result)
		if err != nil { log.Fatal(err) }
		// do something with result....
		log.Println(result)
	}
	log.Println("end to print data")*/

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