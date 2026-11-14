package internal

import (
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/filters"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToMongoFilter[T any](filtered filters.Filtered[T]) any {
	return filteredToBSON(filtered)
}

func filteredToBSON[T any](f filters.Filtered[T]) bson.M {
	switch f.Op {
	case filters.OpAnd, filters.OpOr:
		operator := "$and"
		if f.Op == filters.OpOr {
			operator = "$or"
		}
		subs := make(bson.A, len(f.Subs))
		for i, sub := range f.Subs {
			subs[i] = filteredToBSON(sub)
		}
		return bson.M{operator: subs}
	default:
		if f.Field == "" {
			return bson.M{}
		}
		return bson.M{f.Field: valueExprToBSON(f.Value)}
	}
}

func valueExprToBSON(v filters.ValueExpr) bson.M {
	switch v.Op {
	case filters.OpEq:
		return bson.M{"$eq": unaryValue(v.Data)}
	case filters.OpNe:
		return bson.M{"$ne": unaryValue(v.Data)}
	case filters.OpGt:
		return bson.M{"$gt": unaryValue(v.Data)}
	case filters.OpGte:
		return bson.M{"$gte": unaryValue(v.Data)}
	case filters.OpLt:
		return bson.M{"$lt": unaryValue(v.Data)}
	case filters.OpLte:
		return bson.M{"$lte": unaryValue(v.Data)}
	case filters.OpIn:
		return bson.M{"$in": multiValues(v.Data)}
	case filters.OpNot:
		return bson.M{"$not": valueExprToBSON(v.Subs[0])}
	default:
		return bson.M{}
	}
}

func unaryValue(data filters.Operable) any {
	return data.(filters.OperableUnary).Value
}

func multiValues(data filters.Operable) []any {
	return data.(filters.OperableMulti).Values
}
