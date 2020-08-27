package mgo

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/globalsign/mgo"
	"log"
	"time"
)

var mgoSession *mgo.Session

func InitMongoDB() {
	var host = beego.AppConfig.String("mongo::mongo_host")
	//var username = beego.AppConfig.String("mongo::mongo_username")
	//var password = beego.AppConfig.String("mongo::mongo_password")
	var poollimit = beego.AppConfig.DefaultInt("mongo_poollimit", 4096)
	var timeout = beego.AppConfig.DefaultInt("mongo_timeout", 60)
	var authdb = beego.AppConfig.String("mongo::mongo_autudb")

	dialInfo := &mgo.DialInfo{
		Addrs:   []string{host},
		Timeout: time.Duration(timeout) * time.Second,
		Source:  authdb,
		//Username:  username,
		//Password:  password,
		PoolLimit: poollimit,
	}

	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalf("Create Session: %s\n", err)
	}
	mgoSession = s

	logs.Info("Create Mongo Session OK")
}

/**
Strong
session 的读写一直向主服务器发起并使用一个唯一的连接，因此所有的读写操作完全的一致。
Monotonic
session 的读操作开始是向其他服务器发起（且通过一个唯一的连接），只要出现了一次写操作，session 的连接就会切换至主服务器。
由此可见此模式下，能够分散一些读操作到其他服务器，但是读操作不一定能够获得最新的数据。
Eventual
session 的读操作会向任意的其他服务器发起，多次读操作并不一定使用相同的连接，也就是读操作不一定有序。
session 的写操作总是向主服务器发起，但是可能使用不同的连接，也就是写操作也不一定有序。

*/
/*
Session.DB()来切换相应的数据库
func (s *Session) DB(name string) *Database

通过Database.C()方法切换集合（Collection）
func (db *Database) C(name string) *Collection
*/

func connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	ms := mgoSession.Copy()         //每一次操作都copy一份 Session,避免每次创建Session,导致连接数量超过设置的最大值
	c := ms.DB(db).C(collection)    //获取文档对象 c := Session.DB(db).C(collection)
	ms.SetMode(mgo.Monotonic, true) //session设置模式
	return ms, c
}

func getDb(db string) (*mgo.Session, *mgo.Database) {
	ms := mgoSession.Copy()
	return ms, ms.DB(db)
}

/*
func connect(db, collection string) (*mgo.Session, *mgo.Collection)
func (c *Collection) Insert(docs ...interface{}) error
*/
// @title    插入数据
// @description   插入数据
// @auth      Berlin             时间（2019/6/18   10:57 ）
// @param     db        	string         "操作的数据库"
// @param     collection    string         "操作的文档(表)"
// @param     doc       	interface{}    "要插入的数据"
// @return    返回参数名     	error         "解释"
func Insert(db, collection string, doc interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close() //每次操作之后都要主动关闭
	return c.Insert(doc)
}

/*
func (c *Collection) Find(query interface{}) *Query 来进行查询，
返回的Query struct可以有附加各种条件来进行过滤

func (q *Query) Select(selector interface{}) *Query
对结果检索哪些字段

func (q *Query) All(result interface{}) error 可以获得所有结果
func (q *Query) One(result interface{}) (err error) 可以获得一个结果
*/
// 查询数据
func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}

func FindAll(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).All(result) //结果放到了result
}

// 只能用来修改单条记录，即使条件能匹配多条记录，也只会修改第一条匹配的记录
func Update(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Update(selector, update)
}

// 更新，如果不存在就插入一个新的数据 `upsert:true`
func Upsert(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.Upsert(selector, update)
	return err
}

// `multi:true` 批量更新
func UpdateAll(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	_, err := c.UpdateAll(selector, update)
	return err
}

// 删除数据
func Remove(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Remove(selector)
}

func RemoveAll(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	_, err := c.RemoveAll(selector)
	return err
}

func FindPage(db, collection string, page, limit int, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Skip(page * limit).Limit(limit).All(result)
}

func FindIter(db, collection string, query interface{}) *mgo.Iter {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Iter()
}

func IsEmpty(db, collection string) bool {
	ms, c := connect(db, collection)
	defer ms.Close()
	count, err := c.Count()
	if err != nil {
		log.Fatal(err)
	}
	return count == 0
}

func Count(db, collection string, query interface{}) (int, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Count()
}

//insert one or multi documents
func BulkInsert(db, collection string, docs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Insert(docs...)
	return bulk.Run()
}

func BulkRemove(db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()

	bulk := c.Bulk()
	bulk.Remove(selector...)
	return bulk.Run()
}

func BulkRemoveAll(db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.RemoveAll(selector...)
	return bulk.Run()
}

func BulkUpdate(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Update(pairs...)
	return bulk.Run()
}

func BulkUpdateAll(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.UpdateAll(pairs...)
	return bulk.Run()
}

func BulkUpsert(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Upsert(pairs...)
	return bulk.Run()
}

func PipeAll(db, collection string, pipeline, result interface{}, allowDiskUse bool) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}
	return pipe.All(result)
}

func PipeOne(db, collection string, pipeline, result interface{}, allowDiskUse bool) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}
	return pipe.One(result)
}

func PipeIter(db, collection string, pipeline interface{}, allowDiskUse bool) *mgo.Iter {
	ms, c := connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}

	return pipe.Iter()

}

func Explain(db, collection string, pipeline, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	pipe := c.Pipe(pipeline)
	return pipe.Explain(result)
}
func GridFSCreate(db, prefix, name string) (*mgo.GridFile, error) {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Create(name)
}

func GridFSFindOne(db, prefix string, query, result interface{}) error {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Find(query).One(result)
}

func GridFSFindAll(db, prefix string, query, result interface{}) error {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Find(query).All(result)
}

func GridFSOpen(db, prefix, name string) (*mgo.GridFile, error) {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Open(name)
}

func GridFSRemove(db, prefix, name string) error {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Remove(name)
}
