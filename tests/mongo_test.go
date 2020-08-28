package test

import (
	"SrvControl/models/db/mgo"
	util "SrvControl/utils"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"os"
	"testing"
	"time"
)

// BSON（/ˈbiːsən/）是一种计算机数据交换格式,主要被用作MongoDB数据库中的数据存储和网络传输格式。
// BSON之名缘于JSON，含义为Binary JSON（二进制JSON）

// MongoDB每个集合都会一个名为_id的主键，这是一个24位的16进制字符串。
// 对应到mgo中就是bson.ObjectId。

type Data struct {
	Id      bson.ObjectId `bson:"_id"`
	Title   string        `bson:"title"`
	Des     string        `bson:"des"`
	Content string        `bson:"content"`
	Date    time.Time     `bson:"date"`
}

const (
	database   = "Test"
	collection = "TestModel"
)

// 测试方法1
func testNormal(t *testing.T) {
	util.InitLog()
	mgo.InitMongoDB()
	//insert one document
	data := &Data{
		Id:      bson.NewObjectId(), //returns a new unique ObjectId
		Title:   "博客的标题 1",
		Des:     "博客描述信息 1",
		Content: "博客的具体内容 1",
		Date:    time.Now(),
	}

	err := mgo.Insert(database, collection, data)
	if err != nil {
		fmt.Println("insert one doc", err)
	}

	// find one with all fields
	var result Data
	err = mgo.FindOne(
		database,
		collection,
		bson.M{"_id": bson.ObjectIdHex("5b3db2334d661ff46ee14b9c")},
		nil,
		&result)
	fmt.Println("find one with all fields", result)

	// find one without id field
	var result1 Data
	err = mgo.FindOne(
		database,
		collection,
		bson.M{"_id": bson.ObjectIdHex("5b3db2334d661ff46ee14b9c")},
		bson.M{"_id": 0},
		&result1)
	fmt.Println("find one without id field", result1)

	//find all documents
	var allResult []Data
	err = mgo.FindAll(
		database,
		collection,
		nil,
		nil,
		&allResult)
	fmt.Println("find all docs", allResult)

	// find all documents with query and selector
	var allResult1 []Data
	err = mgo.FindAll(database, collection, bson.M{"title": "博客的标题 1"}, bson.M{"_id": 0}, &allResult1)
	fmt.Println("find all docs with query and selector", allResult1)

	//find documents with page and limit
	var resultWithPage []Data
	err = mgo.FindPage(
		database,
		collection,
		0,
		4,
		nil,
		bson.M{"_id": 0},
		&resultWithPage)
	fmt.Println("find docs with page and limit", resultWithPage)

	//find the cursor
	var iterAll []Data
	iter := mgo.FindIter(database, collection, nil)
	err = iter.All(&iterAll)
	fmt.Println("find cursor ", iterAll)

	//update one document
	err = mgo.Update(
		database,
		collection,
		bson.M{"_id": bson.ObjectIdHex("5b3db2334d661ff46ee14b9c")},
		bson.M{"$set": bson.M{
			"title": "更新后的标题",
			"des":   "更新后的描述信息",
			"date":  time.Now(),
		}})

	if err != nil {
		fmt.Println("upate one error", err)
	}

	//update all docments
	/*err = db.UpdateAll(database, collection, nil, bson.M{"$set": bson.M{
		"title": "更新所有的标题",
		"date":  time.Now(),
	}})
	if err != nil {
		fmt.Println("update all docs error ", err)
	}*/

	//delete one docment
	err = mgo.Remove(database, collection, bson.M{"_id": bson.ObjectIdHex("5b3db2334d661ff46ee14b99")})
	if err != nil {
		fmt.Println("remove one doc error", err)
	}

	//upsert the docment
	err = mgo.Upsert(database, collection, bson.M{"title": "Title Upsert"}, bson.M{"$set": bson.M{
		"des":     "描述Upsert",
		"date":    time.Now(),
		"Content": "内容Upsert",
	}})
	if err != nil {
		fmt.Println("upsert docment error", err)
	}

}

// 测试方法2
func testBulk(t *testing.T) {
	//bulk insert docments
	d1 := &Data{
		Id:      bson.NewObjectId(),
		Title:   "bulk title",
		Des:     "bulk Des",
		Content: "bulk content",
		Date:    time.Now(),
	}
	d2 := &Data{
		Id:      bson.NewObjectId(),
		Title:   "bulk title",
		Des:     "bulk Des",
		Content: "bulk content",
		Date:    time.Now(),
	}

	insertResult, _ := mgo.BulkInsert(database, collection, d1, d2)
	fmt.Println("bulk insert docs", insertResult)

	//bulk update
	up1 := bson.M{"title": "bulk update title"}
	up2 := bson.M{"$set": bson.M{"title": "bulk update title"}}

	up3 := bson.M{"_id": bson.ObjectIdHex("5b3dbd7a9d5e3e314c93d150")}
	up4 := bson.M{"$set": bson.M{"des": "bulk update des"}}

	updateResult, _ := mgo.BulkUpdate(database, collection, up1, up2, up3, up4)
	fmt.Println("bulk update result", updateResult)

	//bulk update all
	up5 := bson.M{"title": "bulk title"}
	up6 := bson.M{"$set": bson.M{"title": "bulk update title"}}

	updateAllResult, _ := mgo.BulkUpdateAll(database, collection, up5, up6)
	fmt.Println("bulk update result", updateAllResult)
}

// 1. 测试函数T 在自定义测试函数TestOrder给定测试顺序
func TestOrder(t *testing.T) {
	t.Run("normal", testNormal) // T的Run方法，它用来执行子测试
	t.Run("bulk", testBulk)
}

// 2. 基准测试函数B
func BenchMongo(b *testing.B) {
	// b.Run("",)// B的Run方法，它用来执行子测试
}

// 3. 示例函数Example开头
func ExampleMongo() {

}

// 主函数 os.Exit(m.Run())
func TestMain(m *testing.M) {
	println("TestMain setup")
	util.InitLog()
	mgo.InitMongoDB()
	retCode := m.Run() // 执行测试，包括单元测试、性能测试和示例测试
	println("TestMain tear-down")
	os.Exit(retCode) //如果所有测试均通过测试，m.Run()返回0，否同m.Run()返回1，代表测试失败
}
