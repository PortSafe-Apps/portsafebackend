package port

import (
	"context"
	"fmt"
	"os"

	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetConnection(MongoString, dbname string) *mongo.Database {
	MongoInfo := atdb.DBInfo{
		DBString: os.Getenv(MongoString),
		DBName:   dbname,
	}
	conn := atdb.MongoConnect(MongoInfo)
	return conn
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

func DeleteUser(Mongoconn *mongo.Database, colname, nipp string) (deleted interface{}, err error) {
	filter := bson.M{"nipp": nipp}
	data := atdb.DeleteOneDoc(Mongoconn, colname, filter)
	return data, err
}

func PasswordValidator(MongoConn *mongo.Database, colname string, userdata User) bool {
	filter := bson.M{"nipp": userdata.Nipp}
	data := atdb.GetOneDoc[User](MongoConn, colname, filter)
	hashChecker := CheckPasswordHash(userdata.Password, data.Password)
	return hashChecker
}

func UpdatePassword(mongoconn *mongo.Database, user User) (Updatedid interface{}) {
	filter := bson.D{{Key: "nipp", Value: user.Nipp}}
	pass, _ := HashPassword(user.Password)
	update := bson.D{{Key: "$Set", Value: bson.D{
		{Key: "password", Value: pass},
	}}}
	res, err := mongoconn.Collection("user").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return "gagal update data"
	}
	return res
}

func CompareNipp(MongoConn *mongo.Database, Colname, nipp string) bool {
	filter := bson.M{"nipp": nipp}
	err := atdb.GetOneDoc[User](MongoConn, Colname, filter)
	user := err.Nipp
	return user != ""
}

func GetAreaByName(MongoConn *mongo.Database, areaName string) *Area {
	collection := MongoConn.Collection("area")
	filter := bson.D{{Key: "areaName", Value: areaName}}

	var area Area
	err := collection.FindOne(context.Background(), filter).Decode(&area)
	if err != nil {
		// Handle error, misalnya return nil atau tindakan lain yang sesuai
		return nil
	}

	return &area
}

func GetLocationByName(MongoConn *mongo.Database, locationName string) *Location {
	collection := MongoConn.Collection("location")
	filter := bson.D{{Key: "locationName", Value: locationName}}

	var location Location
	err := collection.FindOne(context.Background(), filter).Decode(&location)
	if err != nil {
		// Handle error, misalnya return nil atau tindakan lain yang sesuai
		return nil
	}

	return &location
}

func InsertReport(MongoConn *mongo.Database, colname string, rpt Report) (InsertedID interface{}) {
	req := new(Report)
	req.Reportid = rpt.Reportid
	req.Date = rpt.Date
	req.User = rpt.User
	req.Location = rpt.Location
	req.Description = rpt.Description
	req.ObservationPhoto = rpt.ObservationPhoto
	req.TypeDangerousActions = rpt.TypeDangerousActions
	req.Area = rpt.Area
	req.ImmediateAction = rpt.ImmediateAction
	req.ImprovementPhoto = rpt.ImprovementPhoto
	req.CorrectiveAction = rpt.CorrectiveAction
	return InsertOneDoc(MongoConn, colname, req)
}

func UpdateReport(Mongoconn *mongo.Database, ctx context.Context, emp Report) (UpdateId interface{}, err error) {
	filter := bson.D{{Key: "reportid", Value: emp.Reportid}}
	res, err := Mongoconn.Collection("report").ReplaceOne(ctx, filter, emp)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteReportData(mongoconn *mongo.Database, colname, Repid string) (deletedid interface{}, err error) {
	filter := bson.M{"reportid": Repid}
	data := atdb.DeleteOneDoc(mongoconn, colname, filter)
	return data, err
}

func GetOneReportData(mongoconn *mongo.Database, colname, Repid string) (dest Report) {
	filter := bson.M{"reportid": Repid}
	dest = atdb.GetOneDoc[Report](mongoconn, colname, filter)
	return
}

func GetAllReportData(Mongoconn *mongo.Database, colname string) []Report {
	data := atdb.GetAllDoc[[]Report](Mongoconn, colname)
	return data
}
