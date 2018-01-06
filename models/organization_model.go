package models

import (
	"fmt"
	"time"

	"github.com/hathbanger/goDash/store"
	"labix.org/v2/mgo/bson"
)

type Organization struct {
	Id               bson.ObjectId    `json:"id",bson:"_id"`
	Timestamp        time.Time        `json:"time",bson:"time"`
	OrganizationName string           `json:"organizationName",bson:"organizationName"`
	Teams            []*bson.ObjectId `json:"teams",bson:"teams,omitempty"`
	Users            []*bson.ObjectId `json:"users",bson:"users,omitempty"`
	Campaigns        []*bson.ObjectId `json:"campaigns",bson:"campaigns"`
	Surveys          []*bson.ObjectId `json:"surveys",bson:"surveys"`
}

func NewOrganizationModel(organizationName string, userId string) *Organization {
	objId := bson.ObjectIdHex(userId)

	o := new(Organization)
	o.Id = bson.NewObjectId()
	o.Timestamp = time.Now()
	o.OrganizationName = organizationName
	o.Users = []*bson.ObjectId{&objId}

	return o
}

func (o *Organization) Save() error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "organizations", []string{"organizationName"})
	if err != nil {
		panic(err)
	}

	err = collection.Insert(&Organization{
		Id:               o.Id,
		Timestamp:        o.Timestamp,
		OrganizationName: o.OrganizationName,
		Users:            o.Users})
	if err != nil {
		return err
	}

	userArr := o.Users[0]
	user, err := FindUserModel(userArr.Hex())

	AddOrganizationToUser(user.Id, o.Id)

	return nil
}

func FindOrganizationModel(organizationId string) (Organization, error) {

	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "organizations", []string{"organizationName"})
	if err != nil {
		panic(err)
	}

	organization := Organization{}
	err = collection.Find(bson.M{"id": bson.ObjectIdHex(organizationId)}).One(&organization)
	if err != nil {
		return organization, err
	}
	return organization, err
}

func UpdateOrganizationModel(organizationId string, organizationName string) (Organization, error) {

	organization, err := FindOrganizationModel(organizationId)
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection := session.DB("butterfli").C("organizations")
	colQuerier := bson.M{"id": organization.Id}
	change := bson.M{"$set": bson.M{"organizationName": organizationName}}
	err = collection.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}

	return organization, err
}

func DeleteOrganizationModel(organizationId string) error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "organizations", []string{"organizationName"})
	if err != nil {
		panic(err)
	}

	err = collection.Remove(bson.M{"id": bson.ObjectIdHex(organizationId)})
	if err != nil {
		panic(err)
	}

	return nil
}

func AddOrganizationToUser(userId bson.ObjectId, organizationId bson.ObjectId) error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		fmt.Println("err1", err)
	}

	collection, err := store.ConnectToCollection(
		session, "users", []string{"username"})
	if err != nil {
		fmt.Println("err2", err)
	}

	fmt.Println("FINDING ORG", organizationId)
	organization, err := FindOrganizationModel(organizationId.Hex())
	// err = collection.Update(
	// 	    bson.M{"id": bson.ObjectIdHex(userId)},
	// 	    bson.M{"$push": bson.M{"Organization": organization.Id}},
	// 	)

	query := bson.M{"id": userId}
	update := bson.M{"$set": bson.M{"organization": &organization.Id}}

	// Update
	err = collection.Update(query, update)

	if err != nil {
		fmt.Println("err3", err)
	}

	return nil
}

func AddUserToOrganization(userId string, organizationId string) error {
	fmt.Println("AddUserToOrganization FIRED")
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		fmt.Println("err1", err)
	}

	collection, err := store.ConnectToCollection(
		session, "organizations", []string{"organizationName"})
	if err != nil {
		fmt.Println("err2", err)
	}

	user, err := FindUserModel(userId)
	organization, err := FindOrganizationModel(organizationId)

	fmt.Println("FINDING ORG", organization)
	query := bson.M{"id": organization.Id}
	update := bson.M{"$push": bson.M{"users": &user.Id}}
	organizationAFter, err := FindOrganizationModel(organizationId)

	fmt.Println("organizationAFter", organizationAFter)
	// Update
	err = collection.Update(query, update)

	if err != nil {
		fmt.Println("err3", err)
	}

	return nil
}
