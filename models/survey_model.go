package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hathbanger/goDash/store"
	"labix.org/v2/mgo/bson"
)

type Survey struct {
	Id           bson.ObjectId       `json:"id",bson:"_id"`
	Timestamp    time.Time           `json:"time",bson:"time"`
	Organization *bson.ObjectId      `json:"organizationId",bson:"organizationId"`
	Campaign     *bson.ObjectId      `json:"campaign",bson:"campaign"`
	User         *bson.ObjectId      `json:"user",bson:"user"`
	Answers      []string            `json:"answers",bson:"answers"`
	QA           []map[string]string `json:"qanda",bson:"qanda"`
	Finished     bool                `json:"finished",bson:"finished"`
}

type QandA struct {
	Question string `json:"question",bson:"question"`
	Answer   string `json:"answer",bson:"answer"`
}

func NewSurveyModel(organizationId string, campaignId string, userId string, answers []string) *Survey {
	objId := bson.ObjectIdHex(organizationId)
	camObjId := bson.ObjectIdHex(campaignId)
	usrObjId := bson.ObjectIdHex(userId)
	s := new(Survey)
	s.Id = bson.NewObjectId()
	s.Timestamp = time.Now()
	s.Organization = &objId
	s.Campaign = &camObjId
	s.User = &usrObjId
	s.Answers = answers
	s.Finished = false

	return s
}

type BulkResponse struct {
	Responses []*ReceiveRes `json:responses`
}

type ReceiveRes struct {
	Body string `json:body`
	From string `json:from`
}

func (s *Survey) Save() error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("starting to save")
	collection, err := store.ConnectToCollection(
		session, "surveys", []string{"time"})
	if err != nil {
		panic(err)
	}

	err = collection.Insert(&Survey{
		Id:           s.Id,
		Timestamp:    s.Timestamp,
		Organization: s.Organization,
		Campaign:     s.Campaign,
		User:         s.User,
		Answers:      s.Answers,
		Finished:     s.Finished})
	fmt.Println("survey saved", err)
	if err != nil {
		return err
	}

	// organization, err := FindOrganizationModel(s.Organization.Hex())
	AddSurveyToCampaign(s.Campaign.Hex(), s.Id.Hex())
	AddSurveyToOrganization(s.Organization.Hex(), s.Id.Hex())
	AddSurveyToUser(s.User.Hex(), s.Id.Hex())

	return nil
}

func FindSurveyModel(surveyId string) (Survey, error) {

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

	survey := Survey{}
	err = collection.Find(bson.M{"id": bson.ObjectIdHex(surveyId)}).One(&survey)
	if err != nil {
		return survey, err
	}
	return survey, err
}

func FindSurveyModelByPhoneNumber(phoneNumber string) (Survey, error) {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	survey := Survey{}
	user, _ := FindUserModelByPhoneNumber(phoneNumber)

	for _, v := range user.Surveys {
		surv, _ := FindSurveyModel(v.Hex())
		if surv.Finished != true {
			survey = surv
		}
	}
	fmt.Println("\nSURVEY found via phonenumber: ", survey)
	return survey, err
}

func UpdateSurveyModel(surveyId string, answers []string) (Survey, error) {
	survey, err := FindSurveyModel(surveyId)
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

	colQuerier := bson.M{"id": survey.Id}
	change := bson.M{"$set": bson.M{"answers": answers}}
	err = collection.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}
	campaign, err := FindCampaignModel(survey.Campaign.Hex())

	if len(campaign.Questions) == len(answers) {
		change := bson.M{"$set": bson.M{"finished": true}}
		err = collection.Update(colQuerier, change)
		if err != nil {
			panic(err)
		}
	} else {
		change := bson.M{"$set": bson.M{"finished": false}}
		err = collection.Update(colQuerier, change)
		if err != nil {
			panic(err)
		}
	}

	return survey, err
}

// Send message to phone number via Twilio
func SendQuestionToPhone(numberTo string, message string) (map[string]interface{}, error) {

	accountSid := os.Getenv("ACCOUNT_SID")
	authToken := os.Getenv("AUTH_TOKEN")
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}

	msgData.Set("To", numberTo)
	msgData.Set("From", "+17204087088")

	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Create HTTP request client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, _ := client.Do(req)
	AddQuestionToSurveyModel(numberTo, message)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data)
		}
		return data, err
	} else {
		fmt.Println(resp)
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data)
		}
		return data, err
	}
}

func AddQuestionToSurveyModel(phoneNumber string, question string) (Survey, error) {
	survey, err := FindSurveyModelByPhoneNumber(phoneNumber)

	questionMap := map[string]string{
		question: "",
	}

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

	colQuerier := bson.M{"id": survey.Id}
	change := bson.M{"$push": bson.M{"qanda": questionMap}}
	err = collection.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}
	survey, err = FindSurveyModelByPhoneNumber(phoneNumber)

	return survey, err
}

func AddResponseToSurveyModel(phoneNumber string, response string) (Survey, error) {

	fmt.Println("adding survey to response")
	survey, err := FindSurveyModelByPhoneNumber(phoneNumber)

	if err != nil {
		fmt.Println(err)
	}

	if survey.Finished == true {
		fmt.Println("THE SURVEYs DONE!")
	}
	if survey.Finished == false {
		fmt.Println("adding response to survey")
		answerArray := survey.Answers
		answerArray = append(answerArray, response)
		surveyId := survey.Id.Hex()
		UpdateSurveyModel(surveyId, answerArray)
		fmt.Println("Just added response to survey!")

		// TODO: send response to dialogflow

	}
	survey, err = FindSurveyModelByPhoneNumber(phoneNumber)
	return survey, err
}

func AddSurveyToOrganization(organizationId string, surveyId string) error {
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

func AddSurveyToCampaign(campaignId string, surveyId string) error {
	fmt.Println("AddSurveyToCampaign FIRED")
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		fmt.Println("err1", err)
	}

	collection, err := store.ConnectToCollection(
		session, "campaigns", []string{"time"})
	if err != nil {
		fmt.Println("err2", err)
	}

	survey, err := FindSurveyModel(surveyId)

	campaign, err := FindCampaignModel(campaignId)
	query := bson.M{"id": campaign.Id}
	update := bson.M{"$push": bson.M{"surveys": &survey.Id}}

	fmt.Println("CAMPid", campaign.Id)
	// Update
	err = collection.Update(query, update)

	if err != nil {
		fmt.Println("err3", err)
	}

	return nil
}

func AddSurveyToUser(userId string, surveyId string) error {
	fmt.Println("AddSurveyToUser FIRED")
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		fmt.Println("err1", err)
	}

	collection, err := store.ConnectToCollection(
		session, "users", []string{"phonenumber"})
	if err != nil {
		fmt.Println("err2", err)
	}

	survey, err := FindSurveyModel(surveyId)

	user, err := FindUserModel(userId)
	query := bson.M{"id": user.Id}
	update := bson.M{"$push": bson.M{"surveys": &survey.Id}}

	// Update
	err = collection.Update(query, update)

	if err != nil {
		fmt.Println("err3", err)
	}

	return nil
}
