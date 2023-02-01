package scikits

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoClient struct {
	Label  string // yaml文件中的mongo标签
	client *mongo.Database
}

func (m *MongoClient) GetClient() *mongo.Database {
	return m.client
}

func (m *MongoClient) Init() {
	if m.Label == "" {
		m.Label = "mongo"
	}
	host := MyViper.GetString(fmt.Sprintf("%s.host", m.Label))
	port := MyViper.GetString(fmt.Sprintf("%s.port", m.Label))
	user := MyViper.GetString(fmt.Sprintf("%s.user", m.Label))
	pw := MyViper.GetString(fmt.Sprintf("%s.pass", m.Label))
	db := MyViper.GetString(fmt.Sprintf("%s.db", m.Label))

	// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%v/%s", user, pw, host, port, db)
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	mongoDB := client.Database(db)
	m.client = mongoDB
}

func (m *MongoClient) Update(colName string, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	collection := m.client.Collection(colName)
	updateBson := bson.M{"$set": update}
	res, err := collection.UpdateOne(context.TODO(), filter, updateBson)
	return res, err
}

// 根据唯一键有则更新无则添加
func (m *MongoClient) MongoUpdateOrInsert(colName string, filter bson.M, bMap bson.M) error {
	timeNow := time.Now().Unix()
	bMap["CreateTime"] = timeNow
	_, err := m.MongoInsertOne(colName, bMap)
	if err != nil {
		delete(bMap, "CreateTime")
		bMap["UpdateTime"] = timeNow
		err = m.MongoFindOneAndUpdate(colName, filter, bMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MongoClient) MongoFindOneAndUpdate(colName string, filter bson.M, update bson.M) error {
	updateBson := bson.M{"$set": update}
	collection := m.client.Collection(colName)
	res := collection.FindOneAndUpdate(context.TODO(), filter, updateBson)
	err := res.Err()
	if err != nil {
		//log.Fatal(err)
		return err
	} else {
		return nil
	}
}

func (m *MongoClient) MongoInsertOne(colName string, document interface{}) (*mongo.InsertOneResult, error) {
	collection := m.client.Collection(colName)
	res, err := collection.InsertOne(context.TODO(), document)
	return res, err
}

func (m *MongoClient) MongoJudgeExist(colName string, filter bson.M) bool {
	collection := m.client.Collection(colName)
	singleResult := collection.FindOne(context.TODO(), filter)
	err := singleResult.Err()
	if err != nil {
		return false
	} else {
		return true
	}
}

func (m *MongoClient) MongoFindOneLoadStruct(colName string, filter bson.M, model interface{}) error {
	collection := m.client.Collection(colName)
	singleResult := collection.FindOne(context.TODO(), filter)
	err := singleResult.Decode(model)
	return err
}

func (m *MongoClient) MongoFindAll(colName string, filter bson.M, opts ...*options.FindOptions) []map[string]interface{} {
	collection := m.client.Collection(colName)
	cur, _ := collection.Find(context.TODO(), filter, opts...)
	defer cur.Close(context.TODO())
	results := getMongoListDataByCur(cur)
	return results
}

func getMongoListStructByCur(cur *mongo.Cursor, responseStruct struct{}) []struct{} {
	var results []struct{}
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		tmpStruct := responseStruct
		err := cur.Decode(&tmpStruct)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, tmpStruct)
	}
	return results
}

func (m *MongoClient) GetMongoFindCur(colName string, filter bson.M, opts ...*options.FindOptions) *mongo.Cursor {
	collection := m.client.Collection(colName)
	cur, _ := collection.Find(context.TODO(), filter, opts...)
	return cur
}

func (m *MongoClient) MongoCount(colName string, filter bson.M) int64 {
	collection := m.client.Collection(colName)
	num, _ := collection.CountDocuments(context.TODO(), filter, nil)
	return num
}

func getMongoListDataByCur(cur *mongo.Cursor) []map[string]interface{} {
	var results []map[string]interface{}
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return results
}

func (m *MongoClient) MongoSql(colName string, filter bson.M, opts *options.FindOptions) []map[string]interface{} {
	collection := m.client.Collection(colName)
	cur, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.TODO())
	results := getMongoListDataByCur(cur)
	return results
}
