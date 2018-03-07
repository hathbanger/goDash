package store

import (
	"fmt"
	"labix.org/v2/mgo"
)

func ConnectToDb() (*mgo.Session, error) {

	url := "mongo:27017"
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	return session, err
}

func ConnectToCollection(
	session *mgo.Session,
	collection_str string,
	keyString []string) (*mgo.Collection, error) {

	collection := session.DB("dash").C(collection_str)
	index := mgo.Index{
		Key:        keyString,
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		fmt.Print("\nThere was an error in connecting to collection!\n")
		fmt.Print(err)
	}

	return collection, err
}
