package port

import (
	"context"
	"fmt"
	"os"

	"github.com/aiteung/atdb"
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

func ReplaceOneDoc(mongoconn *mongo.Database, collection string, filter bson.M, userdata User) interface{} {
	return atdb.ReplaceOneDoc(mongoconn, collection, filter, userdata)
}

func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

func GetAllUser(MongoConn *mongo.Database, colname string) []User {
	data := atdb.GetAllDoc[[]User](MongoConn, colname)
	return data
}

func GetOneUser(MongoConn *mongo.Database, colname string, userdata User) User {
	filter := bson.M{"nipp": userdata.Nipp}
	data := atdb.GetOneDoc[User](MongoConn, colname, filter)
	return data
}

func GCFGetHandle(mongoconn *mongo.Database, collection string) []User {
	user := atdb.GetAllDoc[[]User](mongoconn, collection)
	return user
}

func UpdatePassword(mongoconn *mongo.Database, user User) (Updatedid interface{}) {
	filter := bson.M{"nipp": user.Nipp}
	pass, _ := HashPassword(user.Password)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: pass}}}}
	res, err := mongoconn.Collection("user").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "gagal update data"
	}
	return res
}
func InsertUserdata(MongoConn *mongo.Database, nipp, nama, jabatan, divisi, bidang, password, role string) (InsertedID interface{}) {
	req := new(User)
	req.Nipp = nipp
	req.Nama = nama
	req.Jabatan = jabatan
	req.Divisi = divisi
	req.Bidang = bidang
	req.Password = password
	req.Role = role
	return InsertOneDoc(MongoConn, "user", req)
}

func PasswordValidator(MongoConn *mongo.Database, colname string, userdata User) bool {
	filter := bson.M{"nipp": userdata.Nipp}
	data := atdb.GetOneDoc[User](MongoConn, colname, filter)
	hashChecker := CompareHashPass(userdata.Password, data.Password)
	return hashChecker
}

func CompareNipp(MongoConn *mongo.Database, Colname, nipp string) bool {
	filter := bson.M{"nipp": nipp}
	err := atdb.GetOneDoc[User](MongoConn, Colname, filter)
	users := err.Nipp
	return users != ""
}
