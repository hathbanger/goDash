package models

import (
	"time"
	"fmt"

	"labix.org/v2/mgo/bson"
	"github.com/hathbanger/goDash/store"
	// "github.com/butterfli-api/models"
	//"github.com/labstack/gommon/log"

)

type User struct {
	Id 				bson.ObjectId			`json:"id",bson:"_id,omitempty"`
	Timestamp 		time.Time	       		`json:"time",bson:"time,omitempty"`
	Username		string           		`json:"username",bson:"username,omitempty"`
	Password		string           		`json:"-",bson:"password,omitempty"`
	PhoneNumber		string           		`json:"phonenumber",bson:"phonenumber"`
	Organization 	*bson.ObjectId 			`json:"organization",bson:"organization,omitempty"`
	Surveys 		[]*bson.ObjectId 		`json:"surveys",bson:"surveys"`
}


type UserLogin struct {
	Username string `json:username`
	Password string `json:password`
}

type UserCreate struct {
	Username string `json:username,omitempty`
	Password string `json:password,omitempty`
	PhoneNumber string `json:phonenumber,omitempty`
	Organization string `json:organization,omitempty`
}

type BulkUserCreate struct {
	Organization string `json:organization`
	PhoneNumbers []string `json:"phoneNumbers"` 
}

func NewUserModel(username string, password string, phoneNumber string, organizationId string) *User {
	
	org := bson.ObjectIdHex(organizationId)

	u := new(User)
	u.Id = bson.NewObjectId()
	u.Username = username
	u.PhoneNumber = phoneNumber
	u.Password = password
	u.Organization = &org
	u.Timestamp = time.Now()

	return u
}
func NewAdminUserModel(username string, password string, phoneNumber string, organizationId string) *User {
	
	u := new(User)
	u.Id = bson.NewObjectId()
	u.Username = username
	u.PhoneNumber = phoneNumber
	u.Password = password
	// u.Organization = orgArr
	u.Timestamp = time.Now()

	return u
}

func (u *User) Save() error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("new user model!")
	collection, err := store.ConnectToCollection(
		session, "users", []string{"phonenumber"})
	if err != nil {
		panic(err)
	}

	err = collection.Insert(&User{
		Id: u.Id,
		Timestamp: u.Timestamp,
		Username: u.Username,
		PhoneNumber: u.PhoneNumber,
		Organization: u.Organization,
		Password: u.Password})
	fmt.Println("new user saved!", err)
	if err != nil {
		return err
	}

	if u.Organization != nil {
		orgId := u.Organization
		AddUserToOrganization(u.Id.Hex(), orgId.Hex())
	}

	fmt.Println("new user saved!")
	return nil
}

func FindUserModel(userId string) (User, error) {

	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "users", []string{"phonenumber"})
	if err != nil {
		panic(err)
	}

	user := User{}
	err = collection.Find(bson.M{"id": bson.ObjectIdHex(userId)}).One(&user)
	if err != nil {
		return user, err
	}

	return user, err
}

func FindUserModelByPhoneNumber(phoneNumber string) (User, error) {

	fmt.Println("FINDING USER BY PHONE ", phoneNumber)
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "users", []string{"phonenumber"})
	if err != nil {
		panic(err)
	}

	user := User{}
	err = collection.Find(bson.M{"phonenumber": phoneNumber}).One(&user)
	if err != nil {
		return user, err
	}

	fmt.Println("USER FOUND BY PHONENUMBER ", user)

	return user, err
}

func FindByUsernameModel(username string) (User, error) {

	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "users", []string{"phonenumber"})
	if err != nil {
		panic(err)
	}

	user := User{}
	err = collection.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return user, err
	}

	return user, err
}


func UpdateUserModel(userId string, username string, password string) (User, error) {

	user, err := FindUserModel(userId)
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection := session.DB("dash").C("users")
	colQuerier := bson.M{"id": user.Id}
	change := bson.M{"$set": bson.M{ "password": password }}
	err = collection.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}

	return user, err
}

func DeleteUserModel(userId string) error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "users", []string{"phonenumber"})
	if err != nil {
		panic(err)
	}

	err = collection.Remove(bson.M{"id": bson.ObjectIdHex(userId)})
	if err != nil {
		panic(err)
	}

	return nil
}
