package internal

import (
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/store"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IDOnly struct {
	ID any `bson:"_id"`
}

func ToInsertMongoOpt(opt []store.InsertOptions) []options.Lister[options.InsertOneOptions] {
	return nil
}

func ToInsertManyMongoOpt(opt []store.InsertOptions) []options.Lister[options.InsertManyOptions] {
	return nil
}

func UpsertRequested(opt []store.UpdateOptions) bool {
	return len(opt) > 0 && opt[0].Upsert != nil && *opt[0].Upsert
}

func ToReplaceMongoOpt(upsert bool) options.Lister[options.FindOneAndReplaceOptions] {
	return options.FindOneAndReplace().
		SetUpsert(upsert).
		SetReturnDocument(options.After).
		SetProjection(bson.M{"_id": 1})
}
