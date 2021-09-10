package testingdb

import (
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFilter(t *testing.T) {
	stringValue := "abc"

	type exampleNestedField struct {
		Bar string
	}

	scenarios := []struct {
		name              string
		matchingFilter    bson.M
		nonMatchingFilter bson.M
		data              interface{}
	}{
		{
			"empty filter",
			bson.M{},
			bson.M{"a": true},
			struct{}{},
		},
		{
			"bool field match",
			bson.M{"foo": true},
			bson.M{"foo": false},
			struct{ Foo bool }{true},
		},
		{
			"int field match",
			bson.M{"foo": 123},
			bson.M{"foo": 1},
			struct{ Foo int16 }{123},
		},
		{
			"string field match",
			bson.M{"foo": "123"},
			bson.M{"foo": "abc"},
			struct{ Foo string }{"123"},
		},
		{
			"bson tag",
			bson.M{"bar": "123"},
			bson.M{"foo": "123"},
			struct {
				Foo string `bson:"bar"`
			}{"123"},
		},
		{
			"pointer value",
			bson.M{"foo": "abc"},
			bson.M{"foo": nil},
			struct {
				Foo *string
			}{&stringValue},
		},
		{
			"object id",
			bson.M{"foo": primitive.ObjectID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
			bson.M{"foo": primitive.ObjectID{}},
			struct {
				Foo primitive.ObjectID
			}{primitive.ObjectID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
		},
		{
			"inline test",
			bson.M{"bar": "abc"},
			bson.M{"foo": "abc"},
			struct {
				Foo exampleNestedField `bson:",inline"`
			}{exampleNestedField{"abc"}},
		},
		{
			"$gt with int",
			bson.M{"foo": bson.M{"$gt": 5}},
			bson.M{"foo": bson.M{"$gt": 10}},
			struct{ Foo int }{Foo: 7},
		},
		{
			"$lt with int",
			bson.M{"foo": bson.M{"$lt": 10}},
			bson.M{"foo": bson.M{"$lt": 5}},
			struct{ Foo int }{Foo: 7},
		},
		{
			"$gt with uint",
			bson.M{"foo": bson.M{"$gt": 5}},
			bson.M{"foo": bson.M{"$gt": 10}},
			struct{ Foo uint }{Foo: 7},
		},
		{
			"$lt with uint",
			bson.M{"foo": bson.M{"$lt": 10}},
			bson.M{"foo": bson.M{"$lt": 5}},
			struct{ Foo uint }{Foo: 7},
		},
		{
			"$gte",
			bson.M{"foo": bson.M{"$gte": 7}},
			bson.M{"foo": bson.M{"$gte": 10}},
			struct{ Foo int }{Foo: 7},
		},
		{
			"$lte",
			bson.M{"foo": bson.M{"$lte": 7}},
			bson.M{"foo": bson.M{"$lte": 5}},
			struct{ Foo int }{Foo: 7},
		},
		{
			"$gt with time",
			bson.M{"foo": bson.M{"$gt": time.Now()}},
			bson.M{"foo": bson.M{"$gt": time.Now().Add(time.Hour)}},
			struct{ Foo time.Time }{Foo: time.Now().Add(time.Minute * 30)},
		},
		{
			"$lt with time",
			bson.M{"foo": bson.M{"$lt": time.Now().Add(time.Hour)}},
			bson.M{"foo": bson.M{"$lt": time.Now()}},
			struct{ Foo time.Time }{Foo: time.Now().Add(time.Minute * 30)},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			True(t, newFilter(s.matchingFilter).matches(s.data))
			False(t, newFilter(s.nonMatchingFilter).matches(s.data))
		})
	}
}
