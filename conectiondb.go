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

// CRUD User
func GetAllUser(MongoConn *mongo.Database, colname string) []User {
	data := atdb.GetAllDoc[[]User](MongoConn, colname)
	return data
}

func GetOneUser(MongoConn *mongo.Database, colname string, userdata User) User {
	filter := bson.M{"nipp": userdata.Nipp}
	data := atdb.GetOneDoc[User](MongoConn, colname, filter)
	return data
}

func GetUserByNipp(MongoConn *mongo.Database, nipp string) (*User, error) {
	collection := MongoConn.Collection("user")
	filter := bson.D{{Key: "nipp", Value: nipp}}

	var user User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		// Handle error, misalnya return nil atau tindakan lain yang sesuai
		return nil, err
	}

	return &user, nil
}

func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

func InsertUserdata(MongoConn *mongo.Database, nipp, nama, jabatan, password, role string) (InsertedID interface{}) {
	req := new(User)
	req.Nipp = nipp
	req.Nama = nama
	req.Jabatan = jabatan
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

// CRUD Unsafe Action
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

func GetOneReportData(mongoconn *mongo.Database, colname, Reportid string) (dest ReportUnsafeAction) {
	filter := bson.M{"reportid": Reportid}
	dest = atdb.GetOneDoc[ReportUnsafeAction](mongoconn, colname, filter)
	return
}

func GetAllReportData(Mongoconn *mongo.Database, colname string) []ReportUnsafeAction {
	data := atdb.GetAllDoc[[]ReportUnsafeAction](Mongoconn, colname)
	return data
}

func GetAllReportDataByUser(conn *mongo.Database, colname, nipp string) ([]ReportUnsafeAction, error) {
	var reports []ReportUnsafeAction

	filter := bson.D{{Key: "user.nipp", Value: nipp}} // Sesuaikan dengan struktur data yang digunakan
	cur, err := conn.Collection(colname).Find(context.Background(), filter)
	if err != nil {
		// Handle error
		return reports, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var report ReportUnsafeAction
		if err := cur.Decode(&report); err != nil {
			// Handle decoding error
			continue
		}
		reports = append(reports, report)
	}

	if err := cur.Err(); err != nil {
		// Handle error
		return reports, err
	}

	return reports, nil
}

func InsertReportAction(MongoConn *mongo.Database, colname string, rpt ReportUnsafeAction) (InsertedID interface{}) {
	req := new(ReportUnsafeAction)
	req.Reportid = rpt.Reportid
	req.Date = rpt.Date
	req.User = rpt.User
	req.Location = rpt.Location
	req.Area = rpt.Area
	req.Description = rpt.Description
	req.ObservationPhoto = rpt.ObservationPhoto
	req.TypeDangerousActions = rpt.TypeDangerousActions
	req.ImmediateAction = rpt.ImmediateAction
	req.ImprovementPhoto = rpt.ImprovementPhoto
	return InsertOneDoc(MongoConn, colname, req)
}

func UpdateReportAction(Mongoconn *mongo.Database, ctx context.Context, rpt ReportUnsafeAction) (UpdateId interface{}, err error) {
	filter := bson.D{{Key: "reportid", Value: rpt.Reportid}}
	res, err := Mongoconn.Collection("reportingUnsafe").ReplaceOne(ctx, filter, rpt)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteReportData(mongoconn *mongo.Database, colname, Reportid string) (deletedid interface{}, err error) {
	filter := bson.M{"reportid": Reportid}
	data := atdb.DeleteOneDoc(mongoconn, colname, filter)
	return data, err
}

// CRUD Compromised Action
func GetOneCompromisedActionData(mongoconn *mongo.Database, colname, Reportid string) (dest ReportCompromisedAction) {
	filter := bson.M{"reportid": Reportid}
	dest = atdb.GetOneDoc[ReportCompromisedAction](mongoconn, colname, filter)
	return
}

func GetAllCompromisedActionData(Mongoconn *mongo.Database, colname string) []ReportCompromisedAction {
	data := atdb.GetAllDoc[[]ReportCompromisedAction](Mongoconn, colname)
	return data
}

func GetAllCompromisedActionDataByUser(conn *mongo.Database, colname, nipp string) ([]ReportCompromisedAction, error) {
	var repotscompromised []ReportCompromisedAction

	filter := bson.D{{Key: "user.nipp", Value: nipp}} // Sesuaikan dengan struktur data yang digunakan
	cur, err := conn.Collection(colname).Find(context.Background(), filter)
	if err != nil {
		// Handle error
		return repotscompromised, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var report ReportCompromisedAction
		if err := cur.Decode(&report); err != nil {
			// Handle decoding error
			continue
		}
		repotscompromised = append(repotscompromised, report)
	}

	if err := cur.Err(); err != nil {
		// Handle error
		return repotscompromised, err
	}

	return repotscompromised, nil
}

func GetReportCompromisedByID(MongoConn *mongo.Database, reportID string) *ReportCompromisedAction {
	collection := MongoConn.Collection("reportingUnsafeCompromised")
	filter := bson.M{"reportid": reportID}

	var report ReportCompromisedAction
	err := collection.FindOne(context.Background(), filter).Decode(&report)
	if err != nil {
		return nil
	}
	return &report
}

func InsertReportCompromised(MongoConn *mongo.Database, colname string, rpt ReportCompromisedAction) (InsertedID interface{}) {
	req := new(ReportCompromisedAction)
	req.Reportid = rpt.Reportid
	req.Date = rpt.Date
	req.User = rpt.User
	req.Location = rpt.Location
	req.Area = rpt.Area
	req.Description = rpt.Description
	req.ObservationPhoto = rpt.ObservationPhoto
	req.TypeDangerousActions = rpt.TypeDangerousActions
	req.ImmediateAction = rpt.ImmediateAction
	req.Recomendation = rpt.Recomendation
	req.Status = rpt.Status
	return InsertOneDoc(MongoConn, colname, req)
}

// Fungsi untuk mengubah status laporan kondisi berbahaya dan menyimpan tindak lanjut
func UpdateReportCompromised(Mongoconn *mongo.Database, ctx context.Context, rpt ReportCompromisedAction) (UpdateId interface{}, err error) {
	filter := bson.D{{Key: "reportid", Value: rpt.Reportid}}
	res, err := Mongoconn.Collection("reportingCompromised").ReplaceOne(ctx, filter, rpt)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteCompromisedAction(mongoconn *mongo.Database, colname, Reportid string) (deletedid interface{}, err error) {
	filter := bson.M{"reportid": Reportid}
	data := atdb.DeleteOneDoc(mongoconn, colname, filter)
	return data, err
}
