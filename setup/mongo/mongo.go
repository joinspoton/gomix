package mongo

import (
	"strings"
	"fmt"

	"github.com/globalsign/mgo"
)

// ConnectToDB - Create a new mongo connection
func ConnectToDB(replicaSet []map[string]string, db string, username string, password string) (*mgo.Database, error) {
	hostsArray := []string
	for _, replica := range replicaSet {
		hostsArray = append(hostsArray, replica["host"] + ":" + replica["port"])
	}
	hosts := strings.Join(hostsArray, ",")

	session, err := mgo.Dial(fmt.Sprintf("mongodb://%s/%s", hosts, db))
	var credentials mgo.Credential
	session.SetSafe(&mgo.Safe{})
	credentials.Username = username
	credentials.Password = password
	err = session.Login(&credentials)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Session created")
	}
	c := session.DB("calendar")
	return c, nil
}

// ConnectToCollection - Create a new collection connection
func ConnectToCollection(db *mgo.Database, collection string) (*mgo.Collection, error) {
	c := db.C(collection)
	return c, nil
}
