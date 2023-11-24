package port

import (
	"time"
)

type User struct {
	Nipp     string `json:"nipp" bson:"nipp"`
	Nama     string `json:"nama" bson:"nama"`
	Jabatan  string `json:"jabatan" bson:"jabatan"`
	Divisi   string `json:"divisi" bson:"divisi"`
	Bidang   string `json:"bidang" bson:"bidang"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
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

type ResponseBack struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type Report struct {
	Reportid             string               `json:"reportid" bson:"reportid"`
	Date                 string               `json:"date" bson:"date"`
	Account              User                 `json:"user" bson:"user,omitempty"`
	IncidentLocation     string               `json:"incidentLocation" bson:"incidentLocation"`
	Description          string               `json:"description" bson:"description"`
	ObservationPhoto     string               `json:"observationPhoto" bson:"observationPhoto"`
	TypeDangerousActions TypeDangerousActions `json:"typeDangerousActions" bson:"typeDangerousActions,omitempty"`
	Area                 Area                 `json:"area" bson:"area"`
	ImmediateAction      string               `json:"immediateAction" bson:"immediateAction"`
	ImprovementPhoto     string               `json:"improvementPhoto" bson:"improvementPhoto"`
	CorrectiveAction     string               `json:"correctiveAction" bson:"correctiveAction"`
}

type TypeDangerousActions struct {
	TypeId   int      `json:"typeId" bson:"typeId"`
	TypeName string   `json:"typeName" bson:"typeName"`
	SubTypes []string `json:"subTypes" bson:"subTypes"`
}

type Area struct {
	AreaId   int    `json:"areaId" bson:"areaId"`
	AreaName string `json:"areaName" bson:"areaName"`
}

type ResponseReport struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Report `json:"data"`
}

type ResponseReportBanyak struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []Report `json:"data"`
}

type Cred struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ReqUsers struct {
	Username string `json:"username"`
}

type RequestReport struct {
	Reportid string `json:"reportid"`
}
