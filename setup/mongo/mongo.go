package mongo

import (
	"fmt"

	"github.com/globalsign/mgo"
)

// ConnectToDB - Create a new mongo connection
func ConnectToDB(server string, port string, db string, username string, password string) (*mgo.Database, error) {
	session, err := mgo.Dial(server + ":" + port + "/" + db)
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
