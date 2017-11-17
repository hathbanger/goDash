package models

import (
	"time"
	"fmt"

	"labix.org/v2/mgo/bson"
	"github.com/hathbanger/goDash/store"
)

type Survey struct {
	Id 				bson.ObjectId			`json:"id",bson:"_id"`
	Timestamp 		time.Time	       		`json:"time",bson:"time"`
	Organization 	string          	`json:"organizationId",bson:"organizationId"`
	Content 		[][]map[string]string 	`json:"content",bson:"content"`
}


func NewSurveyModel(organizationId string, userId string, content [][]map[string]string) *Survey {
	s := new(Survey)
	s.Id = bson.NewObjectId()
	s.Timestamp = time.Now()
	s.Organization = organizationId
	s.Content = content

	return s
}

func (s *Survey) Save() error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "surveys", []string{"time"})
	if err != nil {
		panic(err)
	}

	err = collection.Insert(&Survey{
		Id: s.Id,
		Timestamp: s.Timestamp,
		Organization: s.Organization,
		Content: s.Content})
	if err != nil {
		return err
	}

	// organization, err := FindOrganizationModel(s.Organization.Hex())
	// AddSurveyToOrganization(s.Id.Hex(), organization.Id.Hex())


	return nil
}

func FindSurveyModel(surveyId string) (Survey, error) {

	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "surveys", []string{"surveyName"})
	if err != nil {
		panic(err)
	}

	survey := Survey{}
	err = collection.Find(bson.M{"id": bson.ObjectIdHex(surveyId)}).One(&survey)
	if err != nil {
		return survey, err
	}
	return survey, err
}


// func UpdateOrganizationModel(organizationId string, organizationName string) (Organization, error) {

// 	organization, err := FindOrganizationModel(organizationId)
// 	session, err := store.ConnectToDb()
// 	defer session.Close()
// 	if err != nil {
// 		panic(err)
// 	}

// 	collection := session.DB("butterfli").C("organizations")
// 	colQuerier := bson.M{"id": organization.Id}
// 	change := bson.M{"$set": bson.M{ "organizationName": organizationName }}
// 	err = collection.Update(colQuerier, change)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return organization, err
// }

// func DeleteOrganizationModel(organizationId string) error {
// 	session, err := store.ConnectToDb()
// 	defer session.Close()
// 	if err != nil {
// 		panic(err)
// 	}

// 	collection, err := store.ConnectToCollection(
// 		session, "organizations", []string{"organizationName"})
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = collection.Remove(bson.M{"id": bson.ObjectIdHex(organizationId)})
// 	if err != nil {
// 		panic(err)
// 	}

// 	return nil
// }

func AddSurveyToOrganization(surveyId string, organizationId string) error {
	fmt.Println("AddSurveyToOrganization FIRED")
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

	survey, err := FindSurveyModel(surveyId)
	
	organization, err := FindOrganizationModel(organizationId)
	query := bson.M{"id": organization.Id}
	update := bson.M{"$push": bson.M{"surveys": &survey.Id}}

	fmt.Println("organization", organization)
	// Update
	err = collection.Update(query, update)


	if err != nil {
		fmt.Println("err3", err)
	}

	return nil
}
