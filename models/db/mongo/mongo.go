package mongo

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var mongoClient *mongo.Client

// 使用mongo-driver

func InitMongo() {
	var host = beego.AppConfig.String("mongo::mongo_host")
	var uri = "mongodb://" + host

	var err error
	//client, err = mongo.NewClient( options.Client().ApplyURI(uri) )
	//if err != nil{
	//	fmt.Println("mongo.NewClient error: ", err)
	//}

	// 连接数据库
	//使用了options，可以设置连接数，连接时间，socket时间，超时时间
	var opts *options.ClientOptions = options.Client().ApplyURI(uri)
	mongoClient, err = mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 判断服务是不是可用
	if err = mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatal(err)
		return
	}
	logs.Info("Mongo DB Connect OK")
}
