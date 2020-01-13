# Upload File Service With SSO 
    This project is a demo web service, which is supply upload file service with identity verification.
 
 ## Requirements
 Golang 1.11 and above
 
 ## Installation && start up
 ### Installation
 * `go get github.com/weizhe0422/WebServiceWithLoginAndUpload.git`
 
 ### Features
  * Already have existing web pages with web services function call 
  * Use cookie to support Single-Sign-on, you can upload file within 1 minutes without login again
  * Support multiple files upload function, and get the last result before next upload request
  * Support upload file to server side or AWS S3
  * Support upload big file to server side or AWS S3, but less than 100MB
  * Use Mongo to save user information
  
 ### Start up
  1. change disk to project folder
  2. `docker load < /dockerImg/mongodb.tar` to load mongo docker image
  3. type `go run main.go` in terminal to launch
  4. use web browser and type `http://localhost:8080`
  5. Input E-Mail and password to log in
  5.1 If the E-Mail and password is valid, then you can go through to upload web page
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/develop/img/Login.png" width="920" height="150" alt="Login">
  5.2 If invalid to find E-mail and password, then you need to registration page to registered.
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/develop/img/register.png" width="518" height="582" alt="Login">
  6. Multi-select target and press UPLOAD button to upload files
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/develop/img/Multi-select-files.png" width="406" height="220" alt="Multi-select-files">
  7. You will get the last upload information, and also can do the next upload action
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/develop/img/UploadResult_S3.png" width="627" height="222" alt="UploadResult_S3">
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/develop/img/UploadResult_Server.png" width="616" height="183" alt="UploadResult_Server">
 