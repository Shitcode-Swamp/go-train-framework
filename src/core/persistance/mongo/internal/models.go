package internal

import (
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/store"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ToMongoWriteModels[T any](docs []store.BulkUpdateOp[T]) []mongo.WriteModel {
	res := make([]mongo.WriteModel, 0, len(docs))
	for _, doc := range docs {
		res = append(res, toMongoWriteModel(doc))
	}
	return res
}

func toMongoWriteModel[T any](doc store.BulkUpdateOp[T]) mongo.WriteModel {
	if doc.Insert {
		return mongo.NewInsertOneModel().SetDocument(doc)
	}
	return mongo.NewUpdateOneModel().
		SetFilter(ToMongoFilter(doc.Filter)).
		SetUpdate(doc.Update).
		SetUpsert(doc.Upsert)
}
