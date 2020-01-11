package serviceFunc

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/weizhe0422/WebServiceWithLoginAndUpload/Utility"
	"io"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
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

const (
	fileChunk = 10 * 1024 * 1024
	maxRetries = 3
)

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
		fmt.Fprintf(resp, indexBody, userName, "" , "")
	}else{
		http.Redirect(resp, request, "/", http.StatusFound)
	}
}

func Upload(resp http.ResponseWriter, request *http.Request){
	reader, err := request.MultipartReader()
	if err != nil {
		fmt.Fprintln(resp, err)
		return
	}
	values := make(map[string][]string, 0)
	maxValueBytes := int64(10000000)
	respString:=""
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		name := part.FormName()
		if name == "" {
			continue
		}
		fileName := part.FileName()

		var b bytes.Buffer

		if fileName == "" {
			n, err := io.CopyN(&b, part, maxValueBytes)
			if err != nil && err != io.EOF {
				fmt.Fprintln(resp, err)
				return
			}

			maxValueBytes -= n
			if maxValueBytes <= 0 {
				msg := "multipart message too large"
				fmt.Fprint(resp, msg)
				return
			}
			values[name] = append(values[name], b.String())
		}
		dst, err := os.Create("./" + fileName)
		defer dst.Close()

		fileSize := 0
		for {
			buffer := make([]byte, fileChunk)
			cBytes, err := part.Read(buffer)
			fileSize += cBytes
			if err == io.EOF {
				break
			}
			dst.Write(buffer[0:cBytes])
		}
		respString = respString + fmt.Sprintf(`檔名:%s 檔案大小為：%d bytes`, fileName, fileSize)+"<br>"
	}

	userName, err := Utility.GetUserName(request)
	if err != nil {
		http.Error(resp, "expired!", http.StatusForbidden)
	}
	compile := regexp.MustCompile(`(\w+([-+.]\w+)*)@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	userName = compile.FindAllStringSubmatch(userName, -1)[0][1]
	indexBody, _ := Utility.LoadFile("templates/index.html")
	fmt.Fprintf(resp, indexBody, userName, respString, "Complete upload!")

}
func AWSUpload(resp http.ResponseWriter, request *http.Request){
	userName, err := Utility.GetUserName(request)
	if err != nil {
		http.Error(resp, "expired!", http.StatusForbidden)
	}
	compile := regexp.MustCompile(`(\w+([-+.]\w+)*)@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	userName = compile.FindAllStringSubmatch(userName, -1)[0][1]
	indexBody, _ := Utility.LoadFile("templates/index.html")

	request.ParseMultipartForm(10000000)
	formData := request.MultipartForm
	files := formData.File["multiplefiles"]

	totalResult := ""
	for i, _ := range files{
		file, err := files[i].Open()
		if err != nil {
			log.Printf("failed to load file: %v",err)
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		resultString := ""
		if files[i].Size < fileChunk {
			resultString = uploadToS3(files[i], resp)
			totalResult = totalResult + resultString
		}else{
			_, resultString = uploadBigFileToS3(files[i], resp, indexBody, userName)
			totalResult = totalResult + resultString
		}
	}
	fmt.Fprintf(resp, indexBody, userName, "Complete upload", totalResult)
}

func uploadToS3(targetFile *multipart.FileHeader, resp http.ResponseWriter) (respString string) {
	file, err :=  targetFile.Open()
	if err != nil {
		log.Printf("failed to load file: %v",err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	config := aws.Config{
		Region: aws.String("ap-northeast-1"),
		Credentials: credentials.AnonymousCredentials,
	}
	sess, _ := session.NewSession(&config)
	svc := s3manager.NewUploader(sess)
	fmt.Println("uploading to S3...")
	upload, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String("weizheuploadservice"),
		Key:    aws.String(targetFile.Filename),
		Body:   file,
	})
	if err != nil {
		log.Printf("failed to upload to S3: %v\n", err)
	}
	log.Println("location:", upload.Location)
	log.Println("UploadID:", upload.UploadID)
	log.Println("VersionID:", upload.VersionID)

	if sizeInterface, ok := file.(Size); ok{
		respString = respString + fmt.Sprintf(`檔名:%s 檔案大小為：%d `, targetFile.Filename, sizeInterface.Size())+"<br>"
		respString = respString + fmt.Sprintf(`AWS位置: %s `, upload.Location)+"<br>"
	}
	return respString
}

func uploadBigFileToS3(targetFile *multipart.FileHeader, resp http.ResponseWriter, indexBody, userName string) (compResp *s3.CompleteMultipartUploadOutput,respString string){
	file, err := targetFile.Open()
	if err != nil {
		log.Printf("failed to load file: %v",err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	//get how many parts of file need to upload
	fileParts := uint64(math.Ceil(float64(targetFile.Size) / float64(fileChunk)))
	log.Printf("The file size is %d, and it will seperate as %d parts to upload.", targetFile.Size, fileParts)

	awsCfg := aws.NewConfig().WithRegion("ap-northeast-1").WithCredentials(credentials.NewStaticCredentials("AKIAIQPQ4HMWNQKZY7EQ","LgiGu7ECbEI05o2Fdt7vfo2s+Z3TmVJX3cU1zbbL",""))
	svc := s3.New(session.New(), awsCfg)

	buf, _ := ioutil.ReadAll(file)
	uploadInput := &s3.CreateMultipartUploadInput{
		Bucket:aws.String("weizheuploadservice"),
		Key:    aws.String(targetFile.Filename),
		ContentType: aws.String(http.DetectContentType(buf)),
	}
	multipartUpload, err := svc.CreateMultipartUpload(uploadInput)
	if err != nil {
		log.Printf("failed to create uploader: %v",err)
	}
	log.Println("ok to create uploader, and start to upload to S3...")

	var curr, partLength int64
	var totalSize = targetFile.Size
	var completedParts []*s3.CompletedPart
	buffer := make([]byte, targetFile.Size)
	partNumber := 1
	for curr=0 ; totalSize != 0; curr+=partLength {
		if totalSize < fileChunk {
			partLength = totalSize
		}else{
			partLength = fileChunk
		}

		completedPart, err := uploadPart(resp, indexBody, userName, svc, multipartUpload, buffer[curr:curr+partLength], partNumber)
		if err!=nil{
			log.Println(err.Error())
			err := abortMultipartUpload(svc, multipartUpload)
			if err != nil {
				log.Println("failed to abort upload: %v", err)
			}
			return
		}
		totalSize -= partLength
		partNumber++
		completedParts = append(completedParts, completedPart)
	}
	compResp, err = completeMultipartUpload(svc, multipartUpload, completedParts)

	respString = respString + fmt.Sprintf(`檔名:%s 檔案大小為：%d `, targetFile.Filename, targetFile.Size)+"<br>"
	respString = respString + fmt.Sprintf(`AWS位置: %s `, *compResp.Location)+"<br>"

	return compResp, respString
}

func X_AWSUpload(resp http.ResponseWriter, request *http.Request){
	flusher, ok := resp.(http.Flusher)
	if !ok {
		http.Error(resp, "Streaming not supported",http.StatusInternalServerError)
		return
	}

	userName, err := Utility.GetUserName(request)
	if err != nil {
		http.Error(resp, "expired!", http.StatusForbidden)
	}
	compile := regexp.MustCompile(`(\w+([-+.]\w+)*)@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	userName = compile.FindAllStringSubmatch(userName, -1)[0][1]
	indexBody, _ := Utility.LoadFile("templates/index.html")

	if request.Method == "POST" {
		flusher.Flush()

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



		}

		fmt.Fprintf(resp, indexBody, userName, "Complete upload", respString)
	}
}

func uploadPart(httpResp http.ResponseWriter, indexBody string, userName string, svc *s3.S3, resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNumber int) (*s3.CompletedPart, error) {
	flusher, _ := httpResp.(http.Flusher)

	tryNum := 1
	partInput := &s3.UploadPartInput{
		Body:          bytes.NewReader(fileBytes),
		Bucket:        resp.Bucket,
		Key:           resp.Key,
		PartNumber:    aws.Int64(int64(partNumber)),
		UploadId:      resp.UploadId,
		ContentLength: aws.Int64(int64(len(fileBytes))),
	}

	for tryNum <= maxRetries {
		log.Println(tryNum," retry...")
		uploadResult, err := svc.UploadPart(partInput)
		if err != nil {
			if tryNum == maxRetries {
				if aerr, ok := err.(awserr.Error); ok {
					return nil, aerr
				}
				return nil, err
			}
			fmt.Printf("Retrying to upload part #%v\n", partNumber)
			tryNum++
		} else {
			fmt.Printf("Uploaded part #%v\n", partNumber)
			//fmt.Fprintf(httpResp, indexBody, userName, fmt.Sprintf("Uploaded part #%v\n", partNumber),"")
			flusher.Flush()

			return &s3.CompletedPart{
				ETag:       uploadResult.ETag,
				PartNumber: aws.Int64(int64(partNumber)),
			}, nil
		}
	}
	return nil, nil
}

func abortMultipartUpload(svc *s3.S3, resp *s3.CreateMultipartUploadOutput) error {
	fmt.Println("Aborting multipart upload for UploadId#" + *resp.UploadId)
	abortInput := &s3.AbortMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
	}
	_, err := svc.AbortMultipartUpload(abortInput)
	return err
}

func completeMultipartUpload(svc *s3.S3, resp *s3.CreateMultipartUploadOutput, completedParts []*s3.CompletedPart) (*s3.CompleteMultipartUploadOutput, error) {
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	return svc.CompleteMultipartUpload(completeInput)
}