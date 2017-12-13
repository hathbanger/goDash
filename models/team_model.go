package models

import (
	"time"
	"fmt"

	"labix.org/v2/mgo/bson"
	"github.com/hathbanger/goDash/store"
)

type Team struct {
	Id 				bson.ObjectId		`json:"id",bson:"_id"`
	Timestamp 		time.Time	       	`json:"time",bson:"time"`
	TeamName		string          	`json:"teamName",bson:"teamName"`
	TeamType		string				`json:"teamType",bson:"teamType"`
	Organization 	*bson.ObjectId 		`json:"organization",bson:"organization"`
	Users 			[]*bson.ObjectId 	`json:"users",bson:"users,omitempty"`
	Campaigns		[]*bson.ObjectId 	`json:"campaigns",bson:"campaigns,omitempty"`
	Surveys 		[]*bson.ObjectId 	`json:"surveys",bson:"surveys,omitempty"`
}



func NewTeamModel(teamName string, teamType string, organizationId string) *Team {
	orgId := bson.ObjectIdHex(organizationId)

	t := new(Team)
	t.Id = bson.NewObjectId()
	t.Timestamp = time.Now()
	t.Organization = &orgId
	t.TeamName = teamName
	t.TeamType = teamType

	return t
}

func (t *Team) Save() error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	// org, err := FindOrganizationModel(t.Organization.Hex())

	collection, err := store.ConnectToCollection(
		session, "teams", []string{"id"})
	if err != nil {
		panic(err)
	}

	err = collection.Insert(&Team{
		Id: t.Id,
		Timestamp: t.Timestamp,
		Organization: t.Organization,
		TeamName: t.TeamName,
		TeamType: t.TeamType})
	if err != nil {
		return err
	}

	AddTeamToOrganization(t.Organization.Hex(), t.Id.Hex())


	return nil
}

func FindTeamModel(teamId string) (Team, error) {

	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(
		session, "teams", []string{"id"})
	if err != nil {
		panic(err)
	}

	team := Team{}
	err = collection.Find(bson.M{"id": bson.ObjectIdHex(teamId)}).One(&team)
	if err != nil {
		return team, err
	}

	return team, err
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

// func AddOrganizationToUser(userId bson.ObjectId, organizationId bson.ObjectId) error {
// 	session, err := store.ConnectToDb()
// 	defer session.Close()
// 	if err != nil {
// 		fmt.Println("err1", err)
// 	}

// 	collection, err := store.ConnectToCollection(
// 		session, "users", []string{"username"})
// 	if err != nil {
// 		fmt.Println("err2", err)
// 	}

// 	fmt.Println("FINDING ORG", organizationId)
// 	organization, err := FindOrganizationModel(organizationId.Hex())
// 	// err = collection.Update(
// 	// 	    bson.M{"id": bson.ObjectIdHex(userId)}, 
// 	// 	    bson.M{"$push": bson.M{"Organization": organization.Id}},
// 	// 	)

// 	query := bson.M{"id": userId}
// 	update := bson.M{"$push": bson.M{"organizations": &organization.Id}}

// 	// Update
// 	err = collection.Update(query, update)


// 	if err != nil {
// 		fmt.Println("err3", err)
// 	}

// 	return nil
// }

func AddTeamToOrganization(organizationId string, teamId string) error {
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

	organization, err := FindOrganizationModel(organizationId)
	team, err := FindTeamModel(teamId)
	
	query := bson.M{"id": organization.Id}
	update := bson.M{"$push": bson.M{"teams": &team.Id}}

	// Update
	err = collection.Update(query, update)


	if err != nil {
		fmt.Println("err3", err)
	}

	return nil
}
func AddUserToTeam(organizationId string, teamId string) error {
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

	organization, err := FindOrganizationModel(organizationId)
	team, err := FindTeamModel(teamId)

	query := bson.M{"id": organization.Id}
	update := bson.M{"$push": bson.M{"teams": &team.Id}}

	// Update
	err = collection.Update(query, update)


	if err != nil {
		fmt.Println("err3", err)
	}

	return nil
}
