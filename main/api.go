package main

import (
	"github.com/gin-gonic/gin"
	"myTeacherEndPoint/database"
	"myTeacherEndPoint/presenter"
	"net/http"
	_ "github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/awsutil"
	_ "github.com/aws/aws-sdk-go/aws/credentials"
	_ "github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/aws/aws-sdk-go/aws/session"
	"os"
	"encoding/base64"
	"time"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/spf13/cast"
	_ "github.com/pusher/pusher-http-go"
	"github.com/pusher/pusher-http-go"
)

const GIN_MODE = "GIN_MODE"
const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
)
const (
	debugCode   = iota
	releaseCode = iota
)
var gin_mode int = debugCode
func SetMode(value string) {
	switch value {
	case DebugMode:
		gin_mode = debugCode
	case ReleaseMode:
		gin_mode = releaseCode
	default:
		panic("gin mode unknown, the allowed modes are: " + DebugMode + " and " + ReleaseMode)
	}
}


func main() {
	db := database.DBInit("user")
	category := database.DBInit("category")
	services := database.DBInit("services")
	media 	 := database.DBInit("media")
	review	 := database.DBInit("review")
	transaction := database.DBInit("transaction")
	schedule := database.DBInit("schedule")
	classes  := database.DBInit("clasess")
	receipt  := database.DBInit("receipt")
	city     := database.DBInit("city")
	profincy := database.DBInit("profincy")
	bank 	 := database.DBInit("bank")

	inDB := &presenter.InDB{DB: db}
	categoryDB := &presenter.InDB{DB:category}
	servicesDB := &presenter.InDB{DB:services}
	mediaDB    := &presenter.InDB{DB:media}
	reviewDB   := &presenter.InDB{DB:review}
	transactionDB := &presenter.InDB{DB:transaction}
	scheduleDB  := &presenter.InDB{DB:schedule}
	clasessDB   := &presenter.InDB{DB:classes}
	receiptDB   := &presenter.InDB{DB:receipt}
	profincyDB	:= &presenter.InDB{DB:profincy}
	cityDB		:= &presenter.InDB{DB:city}
	bankDB		:= &presenter.InDB{DB:bank}

/*	gin.SetMode(gin.ReleaseMode)
	SetMode(ReleaseMode)*/

	EndPoint := gin.Default()
	EndPoint.Use(gin.Recovery())

	//OAuth
	EndPoint.POST		( "/api/login", inDB.LoginHandler)
	EndPoint.POST		( "/api/register", inDB.CreateUser)
	EndPoint.POST		( "/api/getuser",inDB.GetuserHandler)
	EndPoint.POST		( "/api/login-google",inDB.LoginSosmed)
	EndPoint.POST		( "/api/login-facebook",inDB.LoginSosmed)
	EndPoint.POST		( "/api/cek",test)
	EndPoint.POST		( "/api/upload",testUploadAws)
	EndPoint.POST		( "/api/testpusher",pusherTest)


	//services
	EndPoint.POST		( "/api/category",categoryDB.CreateCategory)
	EndPoint.POST           ("/api/create-services",servicesDB.CreateServices)
	EndPoint.POST		("/api/add-review",reviewDB.CreatedReview)
	EndPoint.POST		("/api/list-home",categoryDB.GetListHome)
	EndPoint.POST		("/api/list-category",categoryDB.GetListCategory)
	EndPoint.POST		("/api/myservice",servicesDB.GetMyService)
	EndPoint.POST		("/api/upload-image-server",mediaDB.UploadImage)
	EndPoint.POST		("/api/detail-review",reviewDB.DetailReview)
	EndPoint.POST		("/api/list-image-services",mediaDB.ListImageService)
	EndPoint.POST		("/api/search_service",servicesDB.SearchService)
	EndPoint.POST		("/api/service_by_category",servicesDB.ServiceByCategory)
	EndPoint.POST		("/api/service_by_user",servicesDB.ServiceByUser)

	//order
	EndPoint.POST		("/api/create-transaction",transactionDB.CreateTransaction)
	EndPoint.POST		("/api/create-schedule",scheduleDB.CreateSchedule)
	EndPoint.POST		("/api/create-class",clasessDB.CreateClass)
	EndPoint.POST		("/api/create-receipt",receiptDB.CreateReceipt)

	//detailOrder
	EndPoint.POST		("/api/get-clases-teacher",transactionDB.GetLisClasessTeacher)
	EndPoint.POST		("/api/get-clases-student",transactionDB.GetLisClasessStudent)
	EndPoint.POST		("/api/get-my-class-activ",transactionDB.GetlistMyClass)
	EndPoint.POST		("/api/get-all-schedule",scheduleDB.GetAllSchedule)
	EndPoint.POST		("/api/update-schedule",scheduleDB.UpdateSchedule)
	EndPoint.POST		("/api/update-cancelled-schedule",scheduleDB.UpdateCancelledSchedule)
	EndPoint.POST		("/api/update-nocancelled-schedule",scheduleDB.UpdateNoCancelledSchedule)
	EndPoint.POST		("/api/update-allready-paid",scheduleDB.UpdateAllreadyPaidSchedule)


	EndPoint.POST		("/api/upload-image-profile",mediaDB.UploadImageProfile)
	EndPoint.POST		("/api/list-image-profile",mediaDB.ListImageProfile)
	EndPoint.GET 		( "/api/user",  inDB.Getusers)
	EndPoint.GET 		( "/api/user/:id",  inDB.GetUser)
	EndPoint.POST 		( "/api/update",  inDB.UpdateUser)
	EndPoint.DELETE		( "/api/person/:id", inDB.DeleteUser)
	EndPoint.POST		("/api/get-all-receipt",receiptDB.GetAllReceipt)
	EndPoint.POST		("/api/getreceipt",receiptDB.GetReceipt)
	EndPoint.POST		("/api/update-receipt",receiptDB.UpdateReceipt)
	EndPoint.GET		("/api/getreceipt/:id",receiptDB.GetDataReceipt)
	EndPoint.GET		("/api/update-receipt/:id",receiptDB.UpdaDataReceipt)
	EndPoint.GET		("/api/get-paid-schedule",scheduleDB.GetDataAllSchedule)
	EndPoint.GET		("/api/get-cancelled",scheduleDB.GetCancelledSchedule)
	EndPoint.POST		("/api/update-transaction",scheduleDB.UpdateTransaction)
	EndPoint.POST		("/api/create-cancelled",scheduleDB.CreateCancelled)



	EndPoint.POST		("/api/get-profince",profincyDB.GetProfincy)
	EndPoint.POST		("/api/get-city",cityDB.GetCity)

	EndPoint.POST		("/api/getallservices",receiptDB.GetAllServices)
	EndPoint.GET		("/api/getservices/:id",receiptDB.GetServices)
	EndPoint.GET		("/api/update-services-status/:id",receiptDB.UpdateServices)

	EndPoint.POST		("/api/create-bank",bankDB.CreateBankDetail)

	EndPoint.POST		("/api/check-date",transactionDB.CheckDate)


	//EndPoint.Use(cors.Default())
	EndPoint.Use(CORSMiddleware())
	EndPoint.Run		(			":8080")

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}



func test(ctx *gin.Context)  {
	nama 		:= ctx.PostForm("nama")
	ctx.JSON(http.StatusOK, gin.H{
		"status":"success",
		"token": "oke",
		"test":nama,
	})
}

func pusherTest(c *gin.Context)  {
	client := pusher.Client{
		AppId: "641175",
		Key: "",
		Secret: "",
		Cluster: "ap1",
		Secure: true,
	}

	tester := c.PostForm("tester")

	data := map[string]string{"message": tester}
	client.Trigger("my-channel", "my-event", data)
}

func testUploadAws(c *gin.Context) {
	image 		:= c.PostForm("image")

	dec, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		panic(err)
	}

	id := HashID()
	patFileLocal := "image_"+id+".png";
	f, err := os.Create(patFileLocal)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	//upload to AWS
	aws_access_key_id := ""
	aws_secret_access_key := ""
	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)

	_, err = creds.Get()
	if err != nil {
		// handle error
	}

	cfg := aws.NewConfig().WithRegion("us-east-2").WithCredentials(creds)

	svc := s3.New(session.New(), cfg)

	file, err := os.Open(patFileLocal)

	if err != nil {
		// handle error
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	size := fileInfo.Size()

	buffer := make([]byte, size) // read file content to buffer

	file.Read(buffer)

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := "/media/"+ file.Name()

	acl := "public-read"

	params := &s3.PutObjectInput{
		Bucket: aws.String("myteacherdrive"),
		Key: aws.String(path),
		Body: fileBytes,
		ContentLength: aws.Int64(size),
		ContentType: aws.String(fileType),
		ACL: aws.String(acl),
	}

	resp, err := svc.PutObject(params)
	if err != nil {
		// handle error
	}

	fmt.Printf("response %s", awsutil.StringValue(resp))


	c.JSON(http.StatusOK, gin.H{
		"status":"success",
		"link": "https://s3.us-east-2.amazonaws.com/myteacherdrive"+path,
	})

	err = os.Remove(patFileLocal)
}

func HashID() string {
	now := time.Now().UnixNano()
	/*buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, now)*/

	return cast.ToString(now)
}
