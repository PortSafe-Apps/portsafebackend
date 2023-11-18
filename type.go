package port

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Properties struct {
	Name string `json:"name" bson:"name"`
}

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
	Token    string `json:"token,omitempty" bson:"token,omitempty"`
	Private  string `json:"private,omitempty" bson:"private,omitempty"`
	Public   string `json:"public,omitempty" bson:"public,omitempty"`
}

type UserRole struct {
	Role string `json:"role"`
}

type Payload struct {
	User string    `json:"user"`
	Role string    `json:"role"`
	Exp  time.Time `json:"exp"`
	Iat  time.Time `json:"iat"`
	Nbf  time.Time `json:"nbf"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type Response struct {
	Status  bool        `json:"status" bson:"status"`
	Message string      `json:"message" bson:"message"`
	Data    interface{} `json:"data" bson:"data"`
}

type Reporting struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" `
	Reportingid        int                `json:"reportingid" bson:"reportingid"`
	Date               string             `json:"date" bson:"date"`
	Title              string             `json:"title" bson:"title"`
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
	ViolatorName       string             `json:"violatorname" bson:"violatorname"`
	Status             bool               `json:"status" bson:"status"`
}
