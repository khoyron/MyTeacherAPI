package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	//ID 			int `gorm:"size:10"`
	Firstname 	string
	Lastname  	string
	Location  	string
	Zipcode   	int
	Username  	string
	Password  	string
	Email      	string
	Hp 			string
	Gender      string
	Description  string
	Tipe        string
	Address		string

}

type Category struct {
	gorm.Model
	//id_category 	string
	Nama_category   string

}

type AllCategory struct {
	Nama 		 	   string
	Id   			   string
	ServicesDetail     []ServicesDetail
}

type ClassDetail struct {
	Teacher User
	Student User
}

type MediaDetail struct {
	ImageServices string
	ImageTeacher string
	ImageStudent string
}

type ServicesDetail struct {

	Services			Services
	Image 				Media
	User				User

}

type ServicesSearch struct {

	Title 				string
	Id_category 		int
	Description			string
	Id_user				int
	Verification		string
	Salary 				string
	Educational_Level 	string
	Experiance			string
	Image 				string
	Nama				string
}


type Media struct {
	gorm.Model

	Id_services 	string
	Id_user			string
	Image		 	string
}

type ClassActivity struct {
	gorm.Model

	Id_transaksi 	string
	Status 			string
	Class		 	string
	Id_schedule		string
}

type Pricing struct {
	gorm.Model

	Id_services 	string
	Nama			string
	Pricing		 	string
	Price		 	string
	Time		 	string
	Detail_Pricing  string
}


type Review struct {
	gorm.Model

	Id_services 	int
	Star            int
	Id_user			int
	Attitude		string
	Comment			string
	Communication   string
}

type ReviewModel struct {

	Id_services 	int
	Star            int
	Id_user			int
	Attitude		string
	Comment			string
	Communication   string
	ImageServices   string
	ImageUser		string
	NameUser		string
}

type Services struct {
	gorm.Model

	Nama 				string
	Id_category 		int
	Description			string
	Id_user				int
	Verification		string
	Salary 				string
	Educational_Level 	string
	Experiance			string

}

type BankDetail struct {
	gorm.Model

	BankName 			string
	AccountName  		string
	Norek 				string
	Id_user				int
}

type Schedule struct {
	gorm.Model

	Id_class			string
	Date				string
	Time				string
	Status				string
	Description			string
}

type AllTransaction struct {

	DateStart			string
	Time				string
	Transaction			Transaction
	Classes				Classes
	Services			Services
	MediaDetail			MediaDetail
	ClassDetail			ClassDetail
}

type AllScehdule struct {
	Schedule Schedule
	Clasess Classes
	Transaction Transaction
	BankDetail BankDetail
}

type ResponseCancelledModel struct {
	Schedule Schedule
	Clasess Classes
	Transaction Transaction
	BankDetail Receipt
}


type Transaction struct {
	gorm.Model

	Id_services 	string
	Status			string
	Id_user			string
	Id_teacher		string
	Duration 		string
	Total_meet		string
	Total_prize		string
}

type Classes struct {
	gorm.Model

	Id_transaction  string
	Location		string
}

type Receipt struct {
	gorm.Model

	Id_transaction			string
	Image    				string
	BankName				string
	AccountBankName			string
}

type DataListReceipt struct {
	Receipt Receipt
	Transaction Transaction
}

type City struct {
	gorm.Model

	Provinsi_id  int
	Provinsi 	 string
	Nama		 string
	Type		 string
	Kodepos		 int

}

type Profincy struct {
	gorm.Model

	Name			string

}
