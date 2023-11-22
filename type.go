package port

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Properties struct {
	Name string `json:"name" bson:"name"`
}

type User struct {
	Nipp     string `json:"nipp" bson:"nipp"`
	Nama     string `json:"nama" bson:"nama"`
	Jabatan  string `json:"jabatan" bson:"jabatan"`
	Divisi   string `json:"divisi" bson:"divisi"`
	Bidang   string `json:"bidang" bson:"bidang"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
}

type ReqUsers struct {
	Nipp string `json:"nipp"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Cred struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseDataUser struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Data    []User `json:"data,omitempty" bson:"data,omitempty"`
}

type Response struct {
	Token string `json:"token,omitempty" bson:"token,omitempty"`
}

type ResponseEncode struct {
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
}

type Payload struct {
	User string    `json:"user"`
	Role string    `json:"role"`
	Exp  time.Time `json:"exp"`
	Iat  time.Time `json:"iat"`
	Nbf  time.Time `json:"nbf"`
}

type Report struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" `
	Reportid           string             `json:"reportid" bson:"reportid"`
	Date               string             `json:"date" bson:"date"`
	Supervisorid       int                `json:"supervisorid" bson:"supervisorid"`
	SupervisorName     string             `json:"supervisorname" bson:"supervisorname"`
	SupervisorPosition string             `json:"supervisorposition" bson:"supervisorposition"`
	IncidentLocation   string             `json:"incidentlocation" bson:"incidentlocation"`
	Description        string             `json:"description" bson:"description"`
	ObservationPhoto   string             `json:"observationphoto" bson:"observationphoto"`
	PeopleReactions    string             `json:"peoplereactions" bson:"peoplereactions"`
	PPE                string             `json:"ppe" bson:"ppe"`
	PersonPosition     string             `json:"personposition" bson:"personposition"`
	Equipment          string             `json:"equipment" bson:"equipment"`
	WorkProcedure      string             `json:"workprocedure" bson:"workprocedure"`
	Area               string             `json:"area" bson:"area"`
	ImmediateAction    string             `json:"immediateaction" bson:"immediateaction"`
	ImprovementPhoto   string             `json:"improvementphoto" bson:"improvementphoto"`
	CorrectiveAction   string             `json:"correctiveaction" bson:"correctiveaction"`
}
