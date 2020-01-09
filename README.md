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
  
 ### Start up
  1. change disk to project folder
  2. type `go run main.go` in terminal to launch
  3. use `Account.txt` as database to store account and password, please save your information to log in.
   `ACCOUNT=YOUR_EMAIL|PASSWORD=YOUR_PASSWORD`
  4. use web browser and type `http://localhost:8080`
  5. Input E-Mail and Password to log in
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/master/img/Login.png" width="920" height="150" alt="Login">
  6. Multi-select target and press UPLOAD button to upload files
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/master/img/Multi-select-files.png" width="406" height="220" alt="Multi-select-files">
  7. You will get the last upload information, and also can do the next upload action
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/master/img/UploadResult_S3.png" width="627" height="222" alt="UploadResult_S3">
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/master/img/UploadResult_Server.png" width="616" height="183" alt="UploadResult_Server">
 