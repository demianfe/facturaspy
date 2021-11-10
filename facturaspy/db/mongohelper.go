package db

import (
	"context"
	"errors"
	"log"
	"reflect"
	"time"

	helper "github.com/demianfe/facturaspy/facturaspy/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoHost = "mongodb://localhost:27017/"
	//MongoDBName is the mongo DB currently being used.
	MongoDBName = "facturaspy"
)

//GetMongoConnection returns a mongodb client
func GetMongoConnection() *mongo.Client {
	customValues := []interface{}{
		"",       // string
		int(0),   // int
		int32(0), // int32
	}

	rb := bson.NewRegistryBuilder()
	for _, v := range customValues {
		t := reflect.TypeOf(v)
		defDecoder, err := bson.DefaultRegistry.LookupDecoder(t)
		if err != nil {
			panic(err)
		}
		rb.RegisterDecoder(t, &nullawareDecoder{defDecoder, reflect.Zero(t)})
	}

	clientOpts := options.Client().
		ApplyURI(mongoHost).
		SetRegistry(rb.Build())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	return client

}

type nullawareDecoder struct {
	defDecoder bsoncodec.ValueDecoder
	zeroValue  reflect.Value
}

func (d *nullawareDecoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if vr.Type() != bsontype.Null {
		return d.defDecoder.DecodeValue(dctx, vr, val)
	}

	if !val.CanSet() {
		return errors.New("value not settable")
	}
	if err := vr.ReadNull(); err != nil {
		return err
	}
	// Set the zero value of val's type:
	val.Set(d.zeroValue)
	return nil
}

// TODO: move to another file
func GetLIEData(mongodb *mongo.Client, ruc string) LIE {
	collection := mongodb.Database(MongoDBName).Collection("lie")
	// filter := bson.D{primitive.E{Key: "informante.ruc", Value: ruc}}
	var lie LIE
	// err := collection.FindOne(context.TODO(), filter).Decode(&lie)
	err := collection.FindOne(context.TODO(), bson.M{}).Decode(&lie)
	helper.CheckError(err)
	return lie
}

func GetDetailsDataObjId(mongodb *mongo.Client, docID string) LIEDetalles {
	collection := mongodb.Database(MongoDBName).Collection("lie_details")

	filter := bson.M{"_id": docID} // bson.D{primitive.E{Key: "informante.ruc", Value: ruc}}
	var res LIEDetalles
	err := collection.FindOne(context.TODO(), filter).Decode(&res)
	helper.CheckError(err)
	return res
}

func GetDetailsDataRuc(mongodb *mongo.Client, ruc string) LIEDetalles {
	collection := mongodb.Database(MongoDBName).Collection("lie_details")
	filter := bson.D{primitive.E{Key: "informante.ruc", Value: ruc}}
	var res LIEDetalles
	err := collection.FindOne(context.TODO(), filter).Decode(&res)
	helper.CheckError(err)
	return res
}
