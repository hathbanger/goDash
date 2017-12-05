package models

import (
	"time"
	"fmt"

	"labix.org/v2/mgo/bson"
	"github.com/hathbanger/goDash/store"
)

type Campaign struct {
	Id 				bson.ObjectId		`json:"id",bson:"_id"`
	Timestamp 		time.Time	       	`json:"time",bson:"time"`
	Organization 	*bson.ObjectId      `json:"organizationId",bson:"organizationId"`
	Users			[]*bson.ObjectId 	`json:"users",bson:"users"`
	Questions 		[]string 			`json:"questions",bson:"questions"`
	Surveys 		[]*bson.ObjectId 	`json:"surveys",bson:"surveys"`
	Started			bool				`json:"started",bson:"started"`
}


func NewCampaignModel(organizationId string, questions []string) *Campaign {
	organization, err := FindOrganizationModel(organizationId)
	users := organization.Users
	if err != nil {
		fmt.Println("couldn't find organization", err)
	}

	s := new(Campaign)
	s.Id = bson.NewObjectId()
	s.Timestamp = time.Now()
	s.Organization = &organization.Id
	s.Users = users
	s.Questions = questions
	s.Started = false

	return s
}

func (s *Campaign) Save() error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "campaigns", []string{"time"})
	if err != nil {
		panic(err)
	}

	err = collection.Insert(&Campaign{
		Id: s.Id,
		Timestamp: s.Timestamp,
		Organization: s.Organization,
		Users: s.Users,
		Questions: s.Questions,
		Started: s.Started})
	if err != nil {
		return err
	}

	organization, err := FindOrganizationModel(s.Organization.Hex())
	AddCampaignToOrganization(s.Id.Hex(), organization.Id.Hex())


	return nil
}

func FindCampaignModel(campaignId string) (Campaign, error) {

	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "campaigns", []string{"time"})
	if err != nil {
		panic(err)
	}

	Campaign := Campaign{}
	err = collection.Find(bson.M{"id": bson.ObjectIdHex(campaignId)}).One(&Campaign)
	if err != nil {
		return Campaign, err
	}
	return Campaign, err
}


func StartCampaignModel(campaignId string) (Campaign, error) {
	
	var answers []string
	var err error

	campaign, err := FindCampaignModel(campaignId)
	
	if campaign.Started == false {
		session, err := store.ConnectToDb()
		defer session.Close()
		if err != nil {
			panic(err)
		}

		collection, err := store.ConnectToCollection(
			session, "campaigns", []string{"time"})
		if err != nil {
			panic(err)
		}


		for _, v := range campaign.Users {
			survey := NewSurveyModel(campaign.Organization.Hex(), campaignId, v.Hex(), answers)
			err = survey.Save()
		}

		if err != nil {
			fmt.Println(err)
		}

		colQuerier := bson.M{"id": campaign.Id}
		change := bson.M{"$set": bson.M{ "started": true }}
		err = collection.Update(colQuerier, change)
		if err != nil {
			panic(err)
		}

		campaign, err = FindCampaignModel(campaignId)

		return campaign, err
	}
	return campaign, err	
}

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

// func AddCampaignToOrganization(CampaignId string, organizationId string) error {
// 	fmt.Println("AddCampaignToOrganization FIRED")
// 	session, err := store.ConnectToDb()
// 	defer session.Close()
// 	if err != nil {
// 		fmt.Println("err1", err)
// 	}

// 	collection, err := store.ConnectToCollection(
// 		session, "organizations", []string{"organizationName"})
// 	if err != nil {
// 		fmt.Println("err2", err)
// 	}

// 	Campaign, err := FindCampaignModel(CampaignId)
	
// 	organization, err := FindOrganizationModel(organizationId)
// 	query := bson.M{"id": organization.Id}
// 	update := bson.M{"$push": bson.M{"Campaigns": &Campaign.Id}}

// 	fmt.Println("organization", organization)
// 	// Update
// 	err = collection.Update(query, update)


// 	if err != nil {
// 		fmt.Println("err3", err)
// 	}

// 	return nil
// }


func AddCampaignToOrganization(campaignId string, organizationId string) error {
	fmt.Println("AddCampaignToOrganization FIRED")
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

	campaign, err := FindCampaignModel(campaignId)
	organization, err := FindOrganizationModel(organizationId)

	query := bson.M{"id": organization.Id}
	update := bson.M{"$push": bson.M{"campaigns": &campaign.Id}}
	organizationAFter, err := FindOrganizationModel(organizationId)

	fmt.Println("organizationAFter", organizationAFter)
	// Update
	err = collection.Update(query, update)


	if err != nil {
		fmt.Println("err3", err)
	}

	return nil
}
