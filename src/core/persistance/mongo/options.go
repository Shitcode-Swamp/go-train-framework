package mongo

import (
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/store"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func toInsertMongoOpt(opt []store.InsertOptions) []options.Lister[options.InsertOneOptions] {
	return nil
} 

func toInsertManyMongoOpt(opt []store.InsertOptions) []options.Lister[options.InsertManyOptions] {
	return nil
}

func toUpdateOneMongoOpt(opt []store.UpdateOptions) []options.Lister[options.UpdateOneOptions] {
	out := make([]options.Lister[options.UpdateOneOptions], 0, len(opt))
	for _, op := range opt {
		out = append(out, options.UpdateOne().SetUpsert(*op.Upsert))
	}
	return  out
}
