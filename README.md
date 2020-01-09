# Upload File Service With SSO 
    This project is a demo web service, which is supply upload file service with identity verification.
 
 ## Requirements
 Golang 1.11 and above
 
 ## Installation && start up
 ### Installation
 * `go get github.com/weizhe0422/WebServiceWithLoginAndUpload.git`
 
 ### Start up
  1. change disk to project folder
  2. type `go run main.go` in terminal to launch
  3. use `Account.txt` as database to store account and password, please save your information to log in.
   `ACCOUNT=YOUR_EMAIL|PASSWORD=YOUR_PASSWORD`
  4. use web browser and type `http://localhost:8080`
  5. Input E-Mail and Password to log in
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/master/img/Login.png" width="460" height="75" alt="Login">
  6. Multi-select target and press UPLOAD button to upload files
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/master/img/Multi-select-files.png" width="650" height="450" alt="Multi-select-files">
  7. You will get the last upload information, and also can do the next upload action
  <img src="https://github.com/weizhe0422/WebServiceWithLoginAndUpload/blob/master/img/UploadResult.png" width="650" height="450" alt="UploadResult">
  
  