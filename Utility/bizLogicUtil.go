package Utility

import (
	"context"
	"github.com/gorilla/securecookie"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strings"
	"time"
)

type User struct {
	EMAIL string "bson:`email`"
	PASSWORD string "bson:`password`"
}

var cookieHandle = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getMongoDB() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	return mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
}

func InsertUserInfo(name string, password string) bool{
	dbUtil, err := getMongoDB()
	if err != nil {
		log.Printf("failed to connect to mongodb: %v", err)
		return false
	}
	err = dbUtil.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("failed to ping mongo: %v", err)
		return false
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := dbUtil.Database("test").Collection("users")
	record, err := collection.InsertOne(ctx, bson.M{"email": strings.ToUpper(name), "password": password})
	if err != nil {
		log.Printf("failed to insert user info: %v", err)
		return false
	}
	log.Printf("ok to insert user info, record id: %v", record.InsertedID)
	return true
}

func ChkIsUserExist(name string) bool{
	dbUtil, err := getMongoDB()
	if err != nil {
		log.Printf("failed to connect to mongodb: %v", err)
		return false
	}
	err = dbUtil.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("failed to ping mongo: %v", err)
		return false
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := dbUtil.Database("test").Collection("users")
	filter := bson.M{"email":strings.ToUpper(name)}
	err = collection.FindOne(ctx, filter).Decode(&User{})
	if err != nil {
		log.Printf("failed to find user record: %v", err)
		return false
	}
	log.Println("ok to find the user info")
	return true
}

func IsValidUser(name string, password string) bool {
	dbUtil, err := getMongoDB()
	if err != nil {
		log.Printf("failed to connect to mongodb: %v", err)
		return false
	}
	err = dbUtil.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("failed to ping mongo: %v", err)
		return false
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := dbUtil.Database("test").Collection("users")
	filter := bson.M{"email":strings.ToUpper(name),"password":password}
	err = collection.FindOne(ctx, filter).Decode(&User{})
	if err != nil {
		log.Printf("failed to find user record: %v", err)
		return false
	}
	log.Println("ok to find the user info")
	return true
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