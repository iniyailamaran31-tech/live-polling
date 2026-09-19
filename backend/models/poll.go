package models

import "go.mongodb.org/mongo-driver/v2/bson"

type PollOption struct {
	ID    string `bson:"id" json:"id"`
	Text  string `bson:"text" json:"text"`
	Count int64  `bson:"count" json:"count"`
}

type Poll struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Question  string        `bson:"question" json:"question"`
	Options   []PollOption  `bson:"options" json:"options"`
	CreatedBy bson.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt int64         `bson:"createdAt" json:"createdAt"`
}