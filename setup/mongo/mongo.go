package mongo

import (
	"fmt"
	"gomix/utilities/paramstore"
	"gomix/utilities/system"
	"strings"

	"github.com/globalsign/mgo"
)

// ConnectToDB() - Create a new mongo connection without specifying credentials
func ConnectToDB() {
	path := fmt.Sprintf("/%s/db/", system.GetEnv("stage", "staging"))

	name, _ := paramstore.GetConfig(path + "name")
	replicaSet, _ := paramstore.GetJSONArrayConfig(path + "replicaSetArray")
	username, _ := paramstore.GetConfig(path + "username")
	password, _ := paramstore.GetConfig(path + "password")

	return ManuallyConnect(
		replicaSet,
		name,
		username,
		password,
	)
}

// ManuallyConnectToDB - Create a new mongo connection
func ManuallyConnectToDB(replicaSet []map[string]string, db string, username string, password string) (*mgo.Database, error) {
	var hostsArray []string
	for _, replica := range replicaSet {
		hostsArray = append(hostsArray, replica["host"]+":"+replica["port"])
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
