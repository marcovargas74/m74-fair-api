package cpfcnpj

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

func (q *MyQuery) saveQueryInMongoDB() error {

	_, err := collectionQuery.InsertOne(ctx, q)
	return err
}

func (q *MyQuery) showQueryAllMongoDB() (string, error) {

	filter := bson.M{}

	cursor, err := collectionQuery.Find(ctx, filter)
	if err != nil {
		return err.Error(), err
	}
	defer cursor.Close(ctx)

	var queryList []MyQuery
	for cursor.Next(ctx) {
		var aQuery MyQuery
		err := cursor.Decode(&aQuery)
		if err != nil {
			return err.Error(), err
		}
		queryList = append(queryList, aQuery)
	}

	if len(queryList) == 0 {
		errEmpty := errors.New("MONGODB: is Empty")
		return errEmpty.Error(), errEmpty
	}

	json, err := json.Marshal(queryList)
	if err != nil {
		return err.Error(), err
	}

	return string(json), nil
}

func (q *MyQuery) showQuerysByTypeMongoDB(isCPF bool) (string, error) {

	filter := bson.M{"is_cpf": bson.M{"$eq": isCPF}}

	cursor, err := collectionQuery.Find(ctx, filter)
	if err != nil {
		return err.Error(), err
	}
	defer cursor.Close(ctx)

	var queryList []MyQuery
	for cursor.Next(ctx) {
		var aQuery MyQuery
		err := cursor.Decode(&aQuery)
		if err != nil {
			return err.Error(), err
		}
		queryList = append(queryList, aQuery)
	}

	if len(queryList) == 0 {
		errEmpty := errors.New("MONGODB: Not Found elements to this Type")
		return errEmpty.Error(), errEmpty
	}

	json, err := json.Marshal(queryList)
	if err != nil {
		return err.Error(), err
	}

	return string(json), nil

}

func (q *MyQuery) deleteQuerysByNumMongoDB(findNum string) error {

	filter := bson.M{"cpf": bson.M{"$eq": findNum}}

	result, err := collectionQuery.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("num %q Not Found", findNum)
	}

	return nil
}

func (q *MyQuery) showQuerysByNumMongoDB(findNum string) (string, error) {

	filter := bson.M{"cpf": bson.M{"$eq": findNum}}

	cursor, err := collectionQuery.Find(ctx, filter)
	if err != nil {
		return err.Error(), err
	}
	defer cursor.Close(ctx)

	var queryList []MyQuery
	for cursor.Next(ctx) {
		var aQuery MyQuery
		err := cursor.Decode(&aQuery)
		if err != nil {
			return err.Error(), err
		}
		queryList = append(queryList, aQuery)
	}

	if len(queryList) == 0 {
		errEmpty := fmt.Errorf("num %q Not Found", findNum)
		return errEmpty.Error(), errEmpty
	}

	json, err := json.Marshal(queryList)
	if err != nil {
		return err.Error(), err
	}

	return string(json), nil
}
