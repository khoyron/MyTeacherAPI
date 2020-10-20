package presenter

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"myTeacherEndPoint/model"
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/spf13/cast"

	"bytes"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"os"
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pusher/pusher-http-go"
)

type InDB struct {
	DB *gorm.DB
}


func (idb *InDB) LoginHandler(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	password 	:= c.PostForm("password")
	email 		:= c.PostForm("email")

	println(password)

	var (
		person model.User
	)

	err := idb.DB.Where("Email = ?", email).First(&person)
	if err != nil {
		fmt.Println(err.Error)
	}

	if person.Email!=email {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":"failed",
			"message": email+" email not found",
		})
		c.Abort()
	}else {
		if person.Password != password {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":"failed",
				"message": "password wrong",
			})
			c.Abort()
		}else {

			sign := jwt.New(jwt.GetSigningMethod("HS256"))
			claims := make(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
			claims["id"] = cast.ToString(person.ID)
			sign.Claims = claims
			token, err := sign.SignedString([]byte("khoironKey"))

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":"failed",
					"message": err.Error(),
				})
				c.Abort()
			}else {
				c.JSON(http.StatusOK, gin.H{
					"status":"success",
					"token": token,
					"data":person,
				})
			}
		}
	}
}


func (idb *InDB) LoginSosmed(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	email 		:= c.PostForm("email")

	var (
		person model.User
	)


	err := idb.DB.Where("email = ?", email).First(&person)
	if err != nil {
		fmt.Println(err.Error)
	}

	if person.Email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":"failed",
			"message": "please register first",
		})
		c.Abort()
	}else {
		sign := jwt.New(jwt.GetSigningMethod("HS256"))
		claims := make(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		claims["id"] = cast.ToString(person.ID)
		sign.Claims = claims
		token, err := sign.SignedString([]byte("khoironKey"))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":"failed",
				"message": err.Error(),
			})
			c.Abort()
		}else {
			c.JSON(http.StatusOK, gin.H{
				"status":"success",
				"token": token,
				"data":person,
			})
		}
	}

}

func (idb *InDB) GetuserHandler(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		person model.User
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			errr := idb.DB.Where("id = ?", id).First(&person)
			if errr != nil {
				fmt.Println(errr.Error)
			}

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":person,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) GetUser(c *gin.Context) {
	var (
		person model.User
		result gin.H
	)
	id := c.Param("id")
	err := idb.DB.Where("id = ?", id).First(&person).Error
	if err != nil {
		result = gin.H{
			"result": err.Error(),
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": person,
			"count":  1,
		}
	}

	c.JSON(http.StatusOK, result)
}




func (idb *InDB) Getusers(c *gin.Context) {
	var (
		persons []model.User
		result  gin.H
	)

	idb.DB.Find(&persons)
	if len(persons) <= 0 {
		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": persons,
			"count":  len(persons),
		}
	}

	c.JSON(http.StatusOK, result)


}

func (idb *InDB) CreateUser(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	var (
		user model.User
	)

	first_name 	:= c.PostForm("first_name")
	last_name 	:= c.PostForm("last_name")
	address 	:= c.PostForm("address")
	zip_code 	:= c.PostForm("zip_code")
	user_name 	:= c.PostForm("user_name")
	pasword 	:= c.PostForm("pasword")
	email 		:= c.PostForm("email")
	mobile 		:= c.PostForm("mobile")
	gender 		:= c.PostForm("gender")
	description := c.PostForm("description")
	tipe 		:= c.PostForm("tipe")
	location 	:= c.PostForm("location")

	user.Firstname 	= first_name
	user.Lastname 	= last_name
	user.Location  	= location
	user.Zipcode   	= cast.ToInt(zip_code)
	user.Username  	= user_name
	user.Password  	= pasword
	user.Email      = email
	user.Hp      	= mobile
	print("--> "+mobile)
	user.Gender     = gender
	user.Description= description
	user.Tipe       = tipe
	user.Address	= address


	var (
		person model.User
	)

	err := idb.DB.Where("email LIKE ?", email).Find(&person)
	if err != nil {
		fmt.Println(err.Error)
	}

	if len(person.Email) <= 0 {
		idb.DB.Create(&user)
		c.JSON(http.StatusOK,gin.H{
			"status":"success",
			"data": user,
		})
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": "Email all ready exist",
		})
	}

	

}
/*
func StringToInt(value string) int {
	str, err := strconv.Atoi(value)
	if err!= nil {
		fmt.Print(err.Error())
	}
	return str
}
*/


func (idb *InDB) UpdateUser(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	var (
		person model.User
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			address 	:= c.PostForm("address")
			zip_code 	:= c.PostForm("zip_code")
			user_name 	:= c.PostForm("user_name")
			email 		:= c.PostForm("email")
			mobile 		:= c.PostForm("mobile")
			description := c.PostForm("description")

			var (
				newUser   model.User
				result    gin.H
			)

			err := idb.DB.First(&person, id).Error
			if err != nil {
				result = gin.H{
					"result": "data not found",
				}
			}

			newUser.Firstname 	= person.Firstname
			newUser.Lastname 	= person.Lastname
			newUser.Location  	= person.Location
			newUser.Tipe       	= person.Tipe
			newUser.Password  	= person.Password
			newUser.Gender     	= person.Gender
			newUser.Zipcode   	= cast.ToInt(zip_code)
			newUser.Username  	= user_name
			newUser.Email      	= email
			newUser.Hp   	    = mobile
			newUser.Description	= description
			newUser.Address     = address

			err = idb.DB.Model(&person).Updates(newUser).Error
			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"status":"success",
					"message": "successfully updated data",
					"data":newUser,
					"token":tokenString,

				}
			}

			c.JSON(http.StatusOK, result)

			/*c.JSON(http.StatusOK,gin.H{

				//"data":person,
			})*/

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}


}



func (idb *InDB) DeleteUser(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	var (
		person model.User
		result gin.H
	)

	id := c.Param("id")
	err := idb.DB.First(&person, id).Error
	if err != nil {
		result = gin.H{
			"result": "data not found",
		}
	}
	err = idb.DB.Delete(&person).Error
	if err != nil {
		result = gin.H{
			"result": "delete failed",
		}
	} else {
		result = gin.H{
			"result": "Data deleted successfully",
		}
	}

	c.JSON(http.StatusOK, result)

}

func (idb *InDB) CreateCategory(c *gin.Context) {
	var (
		subcategory model.Category
	)
	name 	:= c.PostForm("nama")

	subcategory.Nama_category 	= name

	idb.DB.Create(&subcategory)
	c.JSON(http.StatusOK,gin.H{
		"status":"success",
		"data": subcategory,
	})

}

func (idb *InDB) CreateServices(c *gin.Context)  {

	var (
		services model.Services
	)

	nama_services   := c.PostForm("nama_services")
	id_category		:= c.PostForm("id_category")
	description		:= c.PostForm("description")
	//id_user			:= c.PostForm("id_user")
	verification	:= c.PostForm("verification")
	salary 			:= c.PostForm("salary")
	level 			:= c.PostForm("level")
	experiance		:= c.PostForm("experiance")


	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			services.Nama 		  		= nama_services
			services.Id_category  		= cast.ToInt(id_category)
			services.Description  		= description
			services.Id_user	  		= cast.ToInt(id)
			services.Verification 		= verification
			services.Educational_Level  = level
			services.Salary				= salary
			services.Experiance			= experiance

			idb.DB.Create(&services)
			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data": services,
			})
		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}



}


func (idb *InDB) CreatedReview(c *gin.Context)  {

	var (
		review model.Review
	)

	id_services   		:= c.PostForm("id_services")
	star				:= c.PostForm("star")
	attitude			:= c.PostForm("attitude")
	comment 			:= c.PostForm("comment")
	communication		:= c.PostForm("communication")

	review.Id_services		= cast.ToInt(id_services)
	review.Star				= cast.ToInt(star)
	review.Attitude			= attitude
	review.Communication    = communication
	review.Comment			= comment

	if id_services!="" {
		idb.DB.Create(&review)
		c.JSON(http.StatusOK,gin.H{
			"status":"success",
			"data": review,
		})
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": "Please correctly id_services",
		})
	}

}

func (idb *InDB) GetListHome(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		categorys []model.Category
		services  []model.Services
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			err := idb.DB.Find(&categorys)
			if err != nil {
				fmt.Println(err.Error)
			}

			fmt.Println(id)

			result := []model.AllCategory{}

			for k, n := range categorys {
				err := idb.DB.Where("id_category  = ?", cast.ToString(cast.ToInt(categorys[k].ID))).Find(&services)
				if err != nil {
					fmt.Println(err.Error)
				}

				resultServices := []model.ServicesDetail{}
				for w, v := range services  {
					var (
						media	  model.Media
						user	  model.User
					)

					er := idb.DB.Where("id_services = ?", cast.ToString(services[w].ID)).First(&media)
					println(w)
					println("---> "+cast.ToString(services[w].ID))
					if er != nil {
						fmt.Println(er.Error)
					}

					e := idb.DB.Where(" id  = ?",cast.ToString(v.Id_user)).First(&user)
					if er != nil {
						fmt.Println(e.Error)
					}

					newServices := model.ServicesDetail{
						Services: v,
						Image:media,
						User:user,
					}
					resultServices =  append(resultServices,newServices)

				}

				new := model.AllCategory{
					Nama:  n.Nama_category,
					Id  :  cast.ToString(n.ID),
					ServicesDetail: resultServices,
				}

				result = append(result, new)
			}

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":result,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) GetMyService(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		service []model.Services
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			err := idb.DB.Where("id_user = ?",id).Find(&service)
			if err != nil {
				fmt.Println(err.Error)
			}

			fmt.Println(id)

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":service,
			})

		}
	}else {
		c.JSON(http.StatusBadRequest,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) GetListCategory(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		categorys []model.Category
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			err := idb.DB.Find(&categorys)
			if err != nil {
				fmt.Println(err.Error)
			}

			fmt.Println(id)

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":categorys,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}


func (idb *InDB) UploadImage(c *gin.Context)  {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		media model.Media
	)
	id_user   		:= c.PostForm("id_user")
	id_services		:= c.PostForm("id_services")
	image			:= c.PostForm("image")


	if image!="" {


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

/*
		c.JSON(http.StatusOK, gin.H{
			"status":"success",
			"link": "https://s3.us-east-2.amazonaws.com/myteacherdrive"+path,
		})*/


		media.Id_user 		= id_user
		media.Id_services   = id_services
		media.Image  		= "https://s3.us-east-2.amazonaws.com/myteacherdrive"+path

		idb.DB.Create(&media)

		c.JSON(http.StatusOK,gin.H{
			"status":"success",
			"data": media,
		})

		err = os.Remove(patFileLocal)

	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": "Please correctly media",
		})
	}

}


func (idb *InDB) UploadImageProfile(c *gin.Context)  {

	var (
		media model.Media
	)

	image			:= c.PostForm("image")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})

	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id_user := claims["id"].(string)

			if image!="" {


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


				media.Id_user 		= id_user
				media.Image  		= "https://s3.us-east-2.amazonaws.com/myteacherdrive"+path

				idb.DB.Create(&media)

				c.JSON(http.StatusOK,gin.H{
					"status":"success",
					"data": media,
				})

				err = os.Remove(patFileLocal)

			}else {
				c.JSON(http.StatusOK,gin.H{
					"status":"failed",
					"message": "Please correctly media",
				})
			}



		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) ListImageProfile(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		images []model.Media
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			err := idb.DB.Where("id_user = ?", id).Find(&images)
			if err != nil {
				fmt.Println(err.Error)
			}

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data": images,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}


func (idb *InDB) DetailReview(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		reviews []model.Review
	)

	id_services := c.PostForm("id_services")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			println(id)
			err := idb.DB.Where("id_services = ?", id_services).Find(&reviews)
			if err != nil {
				fmt.Println(err.Error)
			}


			reviewModel := []model.ReviewModel{}
			for w, v := range reviews {
				println(w)
				var (
					Services model.Services
					Media model.Media
					User  model.User
					MediaUser model.Media
				)

				erq := idb.DB.Where("id = ?", v.Id_services).First(&Services)
				if erq != nil {
					println("data service not found")
				}

				er := idb.DB.Where("id_services = ?", v.Id_services).First(&Media)
				if er != nil {
					println("data service not found")
				}

				e := idb.DB.Where("id = ?", v.Id_user).First(&User)
				if e != nil {
					println("data service not found")
				}

				x := idb.DB.Where("id_user = ?", cast.ToString(User.ID)).First(&MediaUser)
				if x != nil {
					println("data service not found")
				}


				newReview := model.ReviewModel{

					Id_services 	: v.Id_services,
					Star            : v.Star,
					Id_user			: v.Id_user,
					Attitude		: v.Attitude,
					Comment			: v.Comment,
					Communication   : v.Communication,
					ImageServices   : Media.Image,
					ImageUser		: MediaUser.Image,
					NameUser		: User.Username,
				}
				reviewModel =  append(reviewModel,newReview)
			}




			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data": reviewModel,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) ListImageService(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		images []model.Media
	)

	id_services		:= c.PostForm("id_services")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			err := idb.DB.Where("id_services = ?", id_services).Find(&images)
			if err != nil {
				fmt.Println(err.Error)
			}

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data": images,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) CreateTransaction(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	var (
		transaction model.Transaction
	)

	id_services 	:= c.PostForm("id_services")
	status 			:= c.PostForm("status")
	id_teacher		:= c.PostForm("id_teacher")
	duration 		:= c.PostForm("duration")
	total_meet		:= c.PostForm("total_meet")
	total_prize		:= c.PostForm("total_prize")


	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			transaction.Id_services 	= id_services
			transaction.Status			= status
			transaction.Id_user			= id
			transaction.Id_teacher		= id_teacher
			transaction.Total_meet		= total_meet
			transaction.Duration		= duration
			transaction.Total_prize		= total_prize


			idb.DB.Create(&transaction)

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data": transaction,
			})
		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}


func (idb *InDB) CreateClass(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	var (
		clasess model.Classes
	)

	id_transaction 	:= c.PostForm("id_transaction")
	location		:= c.PostForm("location")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			println(id)

			clasess.Id_transaction	  = id_transaction
			clasess.Location		  = location

			idb.DB.Create(&clasess)

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data": clasess,
			})
		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) CreateSchedule(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	var (
		transaction model.Schedule
	)

	id_class 	    := c.PostForm("id_class")
	date_start		:= c.PostForm("date_start")
	date_end		:= c.PostForm("date_end")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			println(id)

			transaction.Id_class        = id_class
			transaction.Date		= date_start
			transaction.Time		= date_end
			transaction.Status			= "1"

			idb.DB.Create(&transaction)

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data": transaction,
			})
		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}


func (idb *InDB) CreateReceipt(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	var (
		receipt model.Receipt
	)

	id_transaction 		:= c.PostForm("id_transaction")
	BankName			:= c.PostForm("bank-name")
	AccountBankName		:= c.PostForm("account-bank")
	image		    	:= c.PostForm("image")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id	 := claims["id"].(string)
			println(id)

			if image!="" {

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

				receipt.Id_transaction	  = id_transaction
				receipt.Image		      = "https://s3.us-east-2.amazonaws.com/myteacherdrive"+path
				receipt.BankName		  = BankName
				receipt.AccountBankName	  = AccountBankName

				idb.DB.Create(&receipt)

				c.JSON(http.StatusOK,gin.H{
					"status":"success",
					"data": receipt,
				})

				err = os.Remove(patFileLocal)

			}else {
				c.JSON(http.StatusOK,gin.H{
					"status":"failed",
					"message": "Please correctly media",
				})
			}



		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}
}


func (idb *InDB) GetLisClasessStudent(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		transaction []model.Transaction
		clasess     model.Classes
		schedule	model.Schedule
	)


	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			err := idb.DB.Where("id_user = ?", id).Find(&transaction)
			if err != nil {
				fmt.Println(err.Error)
			}

			result := []model.AllTransaction{}
			for k, v := range transaction {
				err := idb.DB.Where("id_transaction = ?", cast.ToString(v.ID)).First(&clasess)
				if err != nil {
					fmt.Println(err.Error)
				}
				println(k)
				er := idb.DB.Where("id_class = ?", cast.ToString(clasess.ID)).First(&schedule)
				if er != nil {
					fmt.Println(err.Error)
				}

				new := model.AllTransaction{
					DateStart : schedule.Date,
					Time	: schedule.Time,
					Transaction:  v,
					Classes : clasess,
				}

				result = append(result, new)
			}

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":result,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) GetlistMyClass(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		transaction 	[]model.Transaction

	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			err := idb.DB.Where("id_teacher = ?", id).Or("id_user = ?", id).Find(&transaction)
			if err != nil {
				fmt.Println(err.Error)
			}

			result := []model.AllTransaction{}
			for k, v := range transaction {
				var (
					clasess     	model.Classes
					schedule		model.Schedule
					services		model.Services
					mediaServices	model.Media
					mediaTeacher	model.Media
					mediaStudent	model.Media
					teacher			model.User
					student			model.User
				)
				err := idb.DB.Where("id_transaction = ?", cast.ToString(v.ID)).First(&clasess)
				if err != nil {
					fmt.Println(err.Error)
				}
				er := idb.DB.Where("id_class = ?", cast.ToString(clasess.ID)).First(&schedule)
				if er != nil {
					fmt.Println(er.Error)
				}

				e := idb.DB.Where("id = ?", cast.ToString(v.Id_services)).First(&services)
				if e != nil {
					fmt.Println(e.Error)
				}

				teach := idb.DB.Where("id = ?", cast.ToString(v.Id_teacher)).First(&teacher)
				if teach != nil {
					fmt.Println(teach.Error)
				}

				stud := idb.DB.Where("id = ?", cast.ToString(v.Id_user)).First(&student)
				if stud != nil {
					fmt.Println(stud.Error)
				}

				n := idb.DB.Where("id_services = ?", cast.ToString(v.Id_services)).First(&mediaServices)
				if n != nil {
					fmt.Println(e.Error)
				}

				m := idb.DB.Where("id_user = ?", cast.ToString(v.Id_teacher)).First(&mediaTeacher)
				if m != nil {
					fmt.Println(m.Error)
				}

				u := idb.DB.Where("id_user = ?", cast.ToString(v.Id_user)).First(&mediaStudent)
				if u != nil {
					fmt.Println(u.Error)
				}


				teachAndStudent := model.ClassDetail{
					Teacher:teacher,
					Student:student,
				}

				println(k)

				newModelMedia := model.MediaDetail{
					ImageServices:mediaServices.Image,
					ImageStudent:mediaStudent.Image,
					ImageTeacher:mediaTeacher.Image,
				}

				new := model.AllTransaction{
					DateStart :     schedule.Date,
					Time	  :		schedule.Time,
					Transaction:  v,
					Classes : clasess,
					Services: services,
					MediaDetail :newModelMedia,
					ClassDetail:teachAndStudent,
				}

				result = append(result, new)
			}

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":result,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}
}



func (idb *InDB) GetLisClasessTeacher(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		transaction []model.Transaction
		clasess     model.Classes
		schedule	model.Schedule
		services	model.Services
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			err := idb.DB.Where("id_teacher = ?", id).Find(&transaction)
			if err != nil {
				fmt.Println(err.Error)
			}

			result := []model.AllTransaction{}
			for k, v := range transaction {
				err := idb.DB.Where("id_transaction = ?", cast.ToString(v.ID)).First(&clasess)
				if err != nil {
					fmt.Println(err.Error)
				}
				er := idb.DB.Where("id_class = ?", cast.ToString(clasess.ID)).First(&schedule)
				if er != nil {
					fmt.Println(er.Error)
				}

				e := idb.DB.Where("id = ?", cast.ToString(v.Id_services)).First(&services)
				if e != nil {
					fmt.Println(e.Error)
				}


				println(k)

				new := model.AllTransaction{
					DateStart : schedule.Date,
					Time:schedule.Time,
					Transaction:  v,
					Classes : clasess,
					Services:services,
				}

				result = append(result, new)
			}

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":result,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}


func (idb *InDB) GetAllReceipt(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		receipt []model.Receipt
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			println(id)
			err := idb.DB.Find(&receipt)
			if err != nil {
				fmt.Println(err.Error)
			}

			result := []model.DataListReceipt{}
			for k, v := range receipt {
				var (
					transaction     	model.Transaction
				)
				err := idb.DB.Where("id = ?", cast.ToString(v.Id_transaction)).First(&transaction)
				if err != nil {
					fmt.Println(err.Error)
				}

				println(k)


				new := model.DataListReceipt{
					Receipt : v,
					Transaction:transaction,
				}

				result = append(result, new)
			}


			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":result,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) GetAllServices(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		services []model.Services
	)


	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			println(id)
			err := idb.DB.Find(&services)
			if err != nil {
				fmt.Println(err.Error)
			}


			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":services,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) GetServices(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	id_services := c.Param("id")

	var (
		result    			gin.H
		Services			model.Services
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			erq := idb.DB.Where("id = ?", id_services).Find(&Services)
			if erq != nil {
				println("data service not found")
			}

			var (
				media	  model.Media
				user	  model.User
			)

			er := idb.DB.Where("id_services = ?", cast.ToString(Services.ID)).First(&media)
			if er != nil {
				fmt.Println(er.Error)
			}

			e := idb.DB.Where(" id  = ?",cast.ToString(Services.Id_user)).First(&user)
			if er != nil {
				fmt.Println(e.Error)
			}

			newServices := model.ServicesSearch{

				Title 				: Services.Nama,
				Id_category 		: Services.Id_category,
				Description			: Services.Description,
				Id_user				: Services.Id_user,
				Verification		: Services.Verification,
				Salary 				: Services.Salary,
				Educational_Level 	: Services.Educational_Level,
				Experiance			: Services.Experiance,
				Image 				: media.Image,
				Nama				: user.Username,
			}

			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"status":"success",
					"data":newServices,
				}
			}

			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}

}


func (idb *InDB) UpdateServices(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H

		services			model.Services
		servicesUpdate 		model.Services
	)

	id_services := c.Param("id")

	err := idb.DB.Where("id = ?", id_services).First(&services)
	if err != nil {
		fmt.Println(err.Error)
	}


	servicesUpdate.Id_user						= services.Id_user
	servicesUpdate.Experiance					= services.Experiance
	servicesUpdate.Salary						= services.Salary
	servicesUpdate.Educational_Level			= services.Educational_Level
	servicesUpdate.Verification      			= "verified"
	servicesUpdate.Description					= services.Description
	servicesUpdate.Id_category					= services.Id_category
	servicesUpdate.Nama							= services.Nama

	n := idb.DB.Model(&services).Updates(servicesUpdate).Error
	if n != nil {
		result = gin.H{
			"status":"error",
			"message": "update failed",
		}
	} else {
		result = gin.H{
			"status":"success",
			"data": servicesUpdate,
		}
	}

	c.JSON(http.StatusOK, result)
}




func (idb *InDB) GetReceipt(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		receipt model.Receipt
	)


	id_receipt 		:= c.PostForm("id_receipt")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			println(id)
			err := idb.DB.Where("id = ?", id_receipt).First(&receipt)
			if err != nil {
				fmt.Println(err.Error)
			}

			result := []model.DataListReceipt{}
			var (
				transaction     	model.Transaction
			)
			en := idb.DB.Where("id = ?", cast.ToString(&receipt.Id_transaction)).First(&transaction)
			if en != nil {
				fmt.Println(en.Error)
			}


			new := model.DataListReceipt{
				Receipt : receipt,
				Transaction:transaction,
			}

			result = append(result, new)


			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":result,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}


func (idb *InDB) GetDataReceipt(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		receipt model.Receipt
	)

	id_receipt := c.Param("id")

	ek := idb.DB.Where("id = ?", id_receipt).First(&receipt)
	if ek != nil {
		fmt.Println(ek.Error)
	}

	var (
		transaction     	model.Transaction
	)

	result := []model.DataListReceipt{}
	en := idb.DB.Where("id = ?", cast.ToString(&receipt.Id_transaction)).First(&transaction)
	if en != nil {
		fmt.Println(en.Error)
	}


	new := model.DataListReceipt{
		Receipt : receipt,
		Transaction:transaction,
	}

	result = append(result, new)


	c.JSON(http.StatusOK,gin.H{
		"status":"success",
		"data":result,
	})

}



func (idb *InDB) GetAllSchedule(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		schedule []model.Schedule
	)

	id_class := c.PostForm("id_class")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			println(id)
			err := idb.DB.Where("id_class = ?", id_class).Find(&schedule)
			if err != nil {
				fmt.Println(err.Error)
			}

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":schedule,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}

//GetDataAllSchedule


func (idb *InDB) GetDataAllSchedule(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		schedules []model.Schedule
		allSchedule  []model.AllScehdule

	)

	err := idb.DB.Where("status = ?","2").Find(&schedules)
	if err != nil {
		fmt.Println(err.Error)
	}

	for k, v := range schedules {

		var (
			clases   model.Classes
			transaction model.Transaction
			user	model.User
			bank  model.BankDetail
		)

		err := idb.DB.Where("id = ?", cast.ToString(v.Id_class)).First(&clases)
		if err != nil {
			fmt.Println(err.Error)
		}
		println(k)
		er := idb.DB.Where("id = ?", cast.ToString(clases.Id_transaction)).First(&transaction)
		if er != nil {
			fmt.Println(err.Error)
		}

		e := idb.DB.Where("id = ?", cast.ToString(transaction.Id_teacher)).First(&user)
		if e != nil {
			fmt.Println(err.Error)
		}

		b := idb.DB.Where("id_user = ?", cast.ToString(user.ID)).First(&bank)
		if b != nil {
			fmt.Println(err.Error)
		}

		new := model.AllScehdule{
			Schedule : v,
			Clasess	: clases,
			Transaction:  transaction,
			BankDetail : bank,
		}

		allSchedule = append(allSchedule, new)
	}

	c.JSON(http.StatusOK,gin.H{
		"status":"success",
		"data":allSchedule,
	})

}


func (idb *InDB) GetCancelledSchedule(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		schedules []model.Schedule
		allSchedule  []model.ResponseCancelledModel

	)

	err := idb.DB.Where("status = ?","3").Find(&schedules)
	if err != nil {
		fmt.Println(err.Error)
	}
	for k, v := range schedules {

		var (
			clases   model.Classes
			transaction model.Transaction
			user	model.User
			bank    model.Receipt
		)

		err := idb.DB.Where("id = ?", cast.ToString(v.Id_class)).First(&clases)
		if err != nil {
			fmt.Println(err.Error)
		}
		println(k)
		er := idb.DB.Where("id = ?", cast.ToString(clases.Id_transaction)).First(&transaction)
		if er != nil {
			fmt.Println(err.Error)
		}

		e := idb.DB.Where("id = ?", cast.ToString(transaction.Id_user)).First(&user)
		if e != nil {
			fmt.Println(err.Error)
		}

		b := idb.DB.Where("id_transaction = ?", cast.ToString(transaction.ID)).First(&bank)
		if b != nil {
			fmt.Println(err.Error)
		}

		new := model.ResponseCancelledModel{
			Schedule : v,
			Clasess	: clases,
			Transaction:  transaction,
			BankDetail : bank,
		}

		allSchedule = append(allSchedule, new)
	}

	c.JSON(http.StatusOK,gin.H{
		"status":"success",
		"data":allSchedule,
	})

}



//UpdaDatateReceipt

func (idb *InDB) UpdateReceipt(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		receipt 			model.Receipt
		result    			gin.H

		transaction			model.Transaction
		transactionUpdate 	model.Transaction
	)

	id_receipt := c.PostForm("id_receipt")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			err := idb.DB.First(&receipt, id_receipt).Error
			if err != nil {
				result = gin.H{
					"result": "data receipt not found",
				}
			}

			er := idb.DB.First(&transaction, receipt.Id_transaction).Error
			if er != nil {
				result = gin.H{
					"result": "data transaction not found",
				}
			}

			transactionUpdate.Id_teacher		= transaction.Id_teacher
			transactionUpdate.Id_services		= transaction.Id_services
			transactionUpdate.Status			= "2"
			transactionUpdate.Id_user			= transaction.Id_user


			er = idb.DB.Model(&transaction).Updates(transactionUpdate).Error
			if er != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"status":"success",
					"data": transactionUpdate,
				}
			}

			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}

}


func (idb *InDB) UpdaDataReceipt(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		receipt 			model.Receipt
		result    			gin.H

		transaction			model.Transaction
		transactionUpdate 	model.Transaction
	)

	id_receipt := c.Param("id")

	err := idb.DB.First(&receipt, id_receipt).Error
	if err != nil {
		result = gin.H{
			"result": "data receipt not found",
		}
	}

	er := idb.DB.First(&transaction, receipt.Id_transaction).Error
	if er != nil {
		result = gin.H{
			"result": "data transaction not found",
		}
	}

	transactionUpdate.Id_teacher		= transaction.Id_teacher
	transactionUpdate.Id_services		= transaction.Id_services
	transactionUpdate.Status			= "2"
	transactionUpdate.Id_user			= transaction.Id_user


	er = idb.DB.Model(&transaction).Updates(transactionUpdate).Error
	if er != nil {
		result = gin.H{
			"status":"error",
			"message": "update failed",
		}
	} else {
		result = gin.H{
			"status":"success",
			"data": transactionUpdate,
		}
	}

	c.JSON(http.StatusOK, result)
}




func (idb *InDB) UpdateSchedule(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H

		schedule			model.Schedule
		scheduleUpdate	 	model.Schedule
		classes				model.Classes
		transaction			model.Transaction
	)

	id_schedule := c.PostForm("id_schedule")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			err := idb.DB.First(&schedule, id_schedule).Error
			if err != nil {
				result = gin.H{
					"result": "data receipt not found",
				}
			}

			scheduleUpdate.Date				= schedule.Date
			scheduleUpdate.Time				= schedule.Time
			scheduleUpdate.Status			= "2"
			scheduleUpdate.Id_class			= schedule.Id_class
			scheduleUpdate.Description		= schedule.Description

			er := idb.DB.Where("id = ?", schedule.Id_class).Find(&classes)
			if er != nil {
				println("data class not found")
			}
			erq := idb.DB.Where("id = ?", classes.Id_transaction).Find(&transaction)
			if erq != nil {
				println("data service not found")
			}

			client := pusher.Client{
				AppId: "",
				Key: "",
				Secret: "",
				Cluster: "ap1",
				Secure: true,
			}

			data := map[string]string{"message": "alhamdulillah"}
			client.Trigger("my-channel", "push."+transaction.Id_teacher, data)
			client.Trigger("my-channel", "push."+transaction.Id_user, data)



			err = idb.DB.Model(&schedule).Updates(scheduleUpdate).Error
			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"id_teacher":transaction.Id_teacher,
					"status":"success",
					"data": schedule,
				}
			}

			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) CreateCancelled(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H

		schedule			model.Schedule
		scheduleUpdate	 	model.Schedule
		classes				model.Classes
		transaction			model.Transaction
	)

	id_schedule := c.PostForm("id_schedule")
	description := c.PostForm("description")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			err := idb.DB.First(&schedule, id_schedule).Error
			if err != nil {
				result = gin.H{
					"result": "data receipt not found",
				}
			}else {
				scheduleUpdate.Date				= schedule.Date
				scheduleUpdate.Time				= schedule.Time
				scheduleUpdate.Status			= "3"
				scheduleUpdate.Id_class			= schedule.Id_class
				scheduleUpdate.Description		= description

				er := idb.DB.Where("id = ?", schedule.Id_class).Find(&classes)
				if er != nil {
					println("data class not found")
				}
				erq := idb.DB.Where("id = ?", classes.Id_transaction).Find(&transaction)
				if erq != nil {
					println("data service not found")
				}

				client := pusher.Client{
					AppId: "",
					Key: "",
					Secret: "",
					Cluster: "ap1",
					Secure: true,
				}

				data := map[string]string{"message": "cancell schedule"}
				client.Trigger("my-channel", "push."+transaction.Id_teacher, data)
				client.Trigger("my-channel", "push."+transaction.Id_user, data)



				err = idb.DB.Model(&schedule).Updates(scheduleUpdate).Error
				if err != nil {
					result = gin.H{
						"status":"error",
						"message": "update failed",
					}
				} else {
					result = gin.H{
						"id_teacher":transaction.Id_teacher,
						"status":"success",
						"data": schedule,
					}
				}

			}


			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}

}


func (idb *InDB) UpdateCancelledSchedule(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H

		schedule			model.Schedule
		scheduleUpdate	 	model.Schedule
		classes				model.Classes
		transaction			model.Transaction
	)

	id_schedule := c.PostForm("id_schedule")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			err := idb.DB.First(&schedule, id_schedule).Error
			if err != nil {
				result = gin.H{
					"result": "data receipt not found",
				}
			}

			scheduleUpdate.Date			    = schedule.Date
			scheduleUpdate.Time		        = schedule.Time
			scheduleUpdate.Status			= "4"
			scheduleUpdate.Id_class			= schedule.Id_class
			scheduleUpdate.Description		= schedule.Description

			er := idb.DB.Where("id = ?", schedule.Id_class).Find(&classes)
			if er != nil {
				println("data class not found")
			}
			erq := idb.DB.Where("id = ?", classes.Id_transaction).Find(&transaction)
			if erq != nil {
				println("data service not found")
			}

			client := pusher.Client{
				AppId: "641175",
				Key: "",
				Secret: "",
				Cluster: "ap1",
				Secure: true,
			}

			data := map[string]string{"message": "Schedule cancelled because "+schedule.Description}
			client.Trigger("my-channel", "push."+transaction.Id_teacher, data)
			client.Trigger("my-channel", "push."+transaction.Id_user, data)



			err = idb.DB.Model(&schedule).Updates(scheduleUpdate).Error
			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"id_teacher":transaction.Id_teacher,
					"status":"success",
					"data": schedule,
				}
			}

			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) UpdateNoCancelledSchedule(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H

		schedule			model.Schedule
		scheduleUpdate	 	model.Schedule
		classes				model.Classes
		transaction			model.Transaction
	)

	id_schedule := c.PostForm("id_schedule")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			err := idb.DB.First(&schedule, id_schedule).Error
			if err != nil {
				result = gin.H{
					"result": "data receipt not found",
				}
			}

			scheduleUpdate.Date			    = schedule.Date
			scheduleUpdate.Time		        = schedule.Time
			scheduleUpdate.Status			= "5"
			scheduleUpdate.Id_class			= schedule.Id_class
			scheduleUpdate.Description		= schedule.Description

			er := idb.DB.Where("id = ?", schedule.Id_class).Find(&classes)
			if er != nil {
				println("data class not found")
			}
			erq := idb.DB.Where("id = ?", classes.Id_transaction).Find(&transaction)
			if erq != nil {
				println("data service not found")
			}

			client := pusher.Client{
				AppId: "",
				Key: "",
				Secret: "",
				Cluster: "ap1",
				Secure: true,
			}

			data := map[string]string{"message": "cancell schedule di tolak"}
			client.Trigger("my-channel", "push."+transaction.Id_teacher, data)
			client.Trigger("my-channel", "push."+transaction.Id_user, data)



			err = idb.DB.Model(&schedule).Updates(scheduleUpdate).Error
			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"id_teacher":transaction.Id_teacher,
					"status":"success",
					"data": schedule,
				}
			}

			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}

}

func (idb *InDB) UpdateAllreadyPaidSchedule(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H

		schedule			model.Schedule
		scheduleUpdate	 	model.Schedule
		classes				model.Classes
		transaction			model.Transaction
	)

	id_schedule := c.PostForm("id_schedule")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			err := idb.DB.First(&schedule, id_schedule).Error
			if err != nil {
				result = gin.H{
					"result": "data receipt not found",
				}
			}

			scheduleUpdate.Date			    = schedule.Date
			scheduleUpdate.Time		        = schedule.Time
			scheduleUpdate.Status			= "10"
			scheduleUpdate.Id_class			= schedule.Id_class
			scheduleUpdate.Description		= schedule.Description

			er := idb.DB.Where("id = ?", schedule.Id_class).Find(&classes)
			if er != nil {
				println("data class not found")
			}
			erq := idb.DB.Where("id = ?", classes.Id_transaction).Find(&transaction)
			if erq != nil {
				println("data service not found")
			}

			client := pusher.Client{
				AppId: "",
				Key: "",
				Secret: "",
				Cluster: "ap1",
				Secure: true,
			}

			data := map[string]string{"message": "cancell schedule di tolak"}
			client.Trigger("my-channel", "push."+transaction.Id_teacher, data)
			client.Trigger("my-channel", "push."+transaction.Id_user, data)



			err = idb.DB.Model(&schedule).Updates(scheduleUpdate).Error
			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"id_teacher":transaction.Id_teacher,
					"status":"success",
					"data": schedule,
				}
			}

			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}

}





func (idb *InDB) UpdateTransaction(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H
		transaction			model.Transaction
		transactionUpdate   model.Transaction
	)

	id_transaction := c.PostForm("id_transaction")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			erq := idb.DB.Where("id = ?", id_transaction).Find(&transaction)
			if erq != nil {
				println("data service not found")
			}

			//update transaction
			transactionUpdate.Total_prize = transaction.Total_prize
			transactionUpdate.Id_user	  = transaction.Id_user
			transactionUpdate.Id_teacher  = transaction.Id_teacher
			transactionUpdate.Id_services = transaction.Id_services
			transactionUpdate.Duration    = transaction.Duration
			transactionUpdate.Total_meet  = transaction.Total_meet
			transactionUpdate.Status      = "2"



			err = idb.DB.Model(&transaction).Updates(transactionUpdate).Error

			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"id_teacher":transactionUpdate,
					"status":"success",
				}
			}

			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}

}



func (idb *InDB) GetProfincy(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		profincy []model.Profincy
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			println(id)
			err := idb.DB.Find(&profincy)
			if err != nil {
				fmt.Println(err.Error)
			}


			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":profincy,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}
}

func (idb *InDB) GetCity(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		city []model.City
	)

	id_profincy := c.PostForm("id_profince")



	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)

			println(id)
			err := idb.DB.Where("provinsi_id = ?", id_profincy).Find(&city)
			if err != nil {
				fmt.Println(err.Error)
			}

			c.JSON(http.StatusOK,gin.H{
				"status":"success",
				"data":city,
			})

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed",
			"message": err.Error(),
		})
	}
}



func (idb *InDB) CreateBankDetail(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()
	var (
		bankDetail model.BankDetail
	)

	id_user 	    := c.PostForm("id_user")
	bankName 		:= c.PostForm("bankName")
	accountName		:= c.PostForm("accountName")
	norek			:= c.PostForm("norek")


	bankDetail.Id_user 	    = cast.ToInt(id_user)
	bankDetail.BankName	    = bankName
	bankDetail.AccountName 	= accountName
	bankDetail.Norek		= norek


	idb.DB.Create(&bankDetail)

	c.JSON(http.StatusOK,gin.H{
		"status":"success",
		"data": bankDetail,
	})


}


func (idb *InDB) SearchService(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H
		Services			[]model.Services
	)

	search := c.PostForm("search")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			erq := idb.DB.Where("nama LIKE ?", "%"+search+"%").Find(&Services)
			if erq != nil {
				println("data service not found")
			}

			resultServices := []model.ServicesSearch{}
			for w, v := range Services  {
				var (
					media	  model.Media
					user	  model.User
				)

				er := idb.DB.Where("id_services = ?", cast.ToString(Services[w].ID)).First(&media)
				println(w)
				if er != nil {
					fmt.Println(er.Error)
				}

				e := idb.DB.Where(" id  = ?",cast.ToString(v.Id_user)).First(&user)
				if er != nil {
					fmt.Println(e.Error)
				}

				newServices := model.ServicesSearch{

					Title 				: v.Nama,
					Id_category 		: v.Id_category,
					Description			: v.Description,
					Id_user				: v.Id_user,
					Verification		: v.Verification,
					Salary 				: v.Salary,
					Educational_Level 	: v.Educational_Level,
					Experiance			: v.Experiance,
					Image 				: media.Image,
					Nama				: user.Username,
				}
				resultServices =  append(resultServices,newServices)

			}

			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"status":"success",
					"data":resultServices,
				}
			}

			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}
}

//service_by_category
func (idb *InDB) ServiceByCategory(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H
		Services			[]model.Services
	)

	id_category := c.PostForm("id_category")

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			erq := idb.DB.Where("id_category = ?", id_category).Find(&Services)
			if erq != nil {
				println("data service not found")
			}

			resultServices := []model.ServicesSearch{}
			for w, v := range Services  {
				var (
					media	  model.Media
					user	  model.User
				)

				er := idb.DB.Where("id_services = ?", cast.ToString(Services[w].ID)).First(&media)
				println(w)
				println("---> "+cast.ToString(Services[w].ID))
				if er != nil {
					fmt.Println(er.Error)
				}

				e := idb.DB.Where(" id  = ?",cast.ToString(v.Id_user)).First(&user)
				if er != nil {
					fmt.Println(e.Error)
				}

				newServices := model.ServicesSearch{

					Title 				: v.Nama,
					Id_category 		: v.Id_category,
					Description			: v.Description,
					Id_user				: v.Id_user,
					Verification		: v.Verification,
					Salary 				: v.Salary,
					Educational_Level 	: v.Educational_Level,
					Experiance			: v.Experiance,
					Image 				: media.Image,
					Nama				: user.Username,
				}
				resultServices =  append(resultServices,newServices)

			}

			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"status":"success",
					"data":resultServices,
				}
			}

			c.JSON(http.StatusOK, result)
		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}
}

func (idb *InDB) ServiceByUser(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H
		Services			[]model.Services
	)

	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			erq := idb.DB.Where("id_user = ?", id).Find(&Services)
			if erq != nil {
				println("data service not found")
			}

			resultServices := []model.ServicesSearch{}
			for w, v := range Services  {
				var (
					media	  model.Media
					user	  model.User
				)

				er := idb.DB.Where("id_services = ?", cast.ToString(Services[w].ID)).First(&media)
				println(w)
				if er != nil {
					fmt.Println(er.Error)
				}

				e := idb.DB.Where(" id  = ?",cast.ToString(v.Id_user)).First(&user)
				if er != nil {
					fmt.Println(e.Error)
				}

				newServices := model.ServicesSearch{

					Title 				: v.Nama,
					Id_category 		: v.Id_category,
					Description			: v.Description,
					Id_user				: v.Id_user,
					Verification		: v.Verification,
					Salary 				: v.Salary,
					Educational_Level 	: v.Educational_Level,
					Experiance			: v.Experiance,
					Image 				: media.Image,
					Nama				: user.Username,
				}
				resultServices =  append(resultServices,newServices)

			}

			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				result = gin.H{
					"status":"success",
					"data":resultServices,
				}
			}

			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}
}


func (idb *InDB) CheckDate(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Next()

	var (
		result    			gin.H
		transaction			[]model.Transaction
		resultCheck	        []string
	)

	id_teacher := c.PostForm("id_teacher")
	date	   := c.PostForm("date")
	time 	   := c.PostForm("time")


	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("khoironKey"), nil
	})


	if token.Valid && err == nil  {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := claims["id"].(string)
			println(id)

			erq := idb.DB.Where("id_teacher = ?", id_teacher).Find(&transaction)
			if erq != nil {
				println("data service not found")
			}

			for k, v := range transaction {
				println(k)
				var (
					clasess       model.Classes
					schedule      []model.Schedule
				)
				err := idb.DB.Where("id_transaction = ?", cast.ToString(v.ID)).First(&clasess)
				if err != nil {
					fmt.Println(err.Error)
				}
				er := idb.DB.Where("id_class = ?", cast.ToString(clasess.ID)).Find(&schedule)
				if er != nil {
					fmt.Println(er.Error)
				}

				for k, n := range schedule {
					println(k)
					println(cast.ToString(n.ID))

					if n.Date == date && n.Time == time {
						resultCheck = append(resultCheck, "1")
					}
				}

			}


			if err != nil {
				result = gin.H{
					"status":"error",
					"message": "update failed",
				}
			} else {
				if len(resultCheck)>0 {
					result = gin.H{
						"message":"schedule cannot be used",
						"status":"failed",
					}
				}else {
					result = gin.H{
						"message":"schedule can be used",
						"status":"success",
					}
				}

			}

			c.JSON(http.StatusOK, result)

		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":"failed token",
			"message": err.Error(),
		})
	}

}



func HashID() string {
	now := time.Now().UnixNano()
	/*buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, now)*/

	return cast.ToString(now)
}