package port

import (
	"context"
	"fmt"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetConnection(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func CreateUser(mongoconn *mongo.Database, collection string, userdata User) interface{} {
	// Hash the password before storing it
	hashedPassword, err := HashPassword(userdata.Password)
	if err != nil {
		return err
	}
	privateKey, publicKey := watoken.GenerateKey()
	userid := userdata.Username
	tokenstring, err := watoken.Encode(userid, privateKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokenstring)
	// decode token to get userid
	useridstring := watoken.DecodeGetId(publicKey, tokenstring)
	if useridstring == "" {
		fmt.Println("expire token")
	}
	fmt.Println(useridstring)
	userdata.Private = privateKey
	userdata.Public = publicKey
	userdata.Password = hashedPassword

	// Insert the user data into the database
	return atdb.InsertOneDoc(mongoconn, collection, userdata)
}

func GCFGetHandle(mongoconn *mongo.Database, collection string) []User {
	user := atdb.GetAllDoc[[]User](mongoconn, collection)
	return user
}

func DeleteUser(mongoconn *mongo.Database, collection string, userdata User) interface{} {
	filter := bson.M{"username": userdata.Username}
	return atdb.DeleteOneDoc(mongoconn, collection, filter)
}

func ReplaceOneDoc(mongoconn *mongo.Database, collection string, filter bson.M, userdata User) interface{} {
	return atdb.ReplaceOneDoc(mongoconn, collection, filter, userdata)
}

func FindUser(mongoconn *mongo.Database, collection string, userdata User) User {
	filter := bson.M{"username": userdata.Username}
	return atdb.GetOneDoc[User](mongoconn, collection, filter)
}

func FindUserUser(mongoconn *mongo.Database, collection string, userdata User) User {
	filter := bson.M{
		"username": userdata.Username,
	}
	return atdb.GetOneDoc[User](mongoconn, collection, filter)
}

func FindUserUserr(mongoconn *mongo.Database, collection string, userdata User) (User, error) {
	filter := bson.M{
		"username": userdata.Username,
	}

	var user User
	err := mongoconn.Collection(collection).FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func IsPasswordValid(mongoconn *mongo.Database, collection string, userdata User) bool {
	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mongoconn, collection, filter)
	return CheckPasswordHash(userdata.Password, res.Password)
}

func IsPasswordValidd(mconn *mongo.Database, collection string, userdata User) (User, bool) {
	filter := bson.M{"username": userdata.Username}
	var foundUser User
	err := mconn.Collection(collection).FindOne(context.Background(), filter).Decode(&foundUser)
	if err != nil {
		return User{}, false
	}
	// Verify password here
	if CheckPasswordHash(userdata.Password, foundUser.Password) {
		return foundUser, true
	}
	return User{}, false
}

func FindUserByUsername(mongoconn *mongo.Database, collection string, username string) (User, error) {
	var user User
	filter := bson.M{"username": username}
	err := mongoconn.Collection(collection).FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
