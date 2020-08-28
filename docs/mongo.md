
#### 

```go
c *mogo.Session

//单条件查询

// = ($eq)
c.Find(bson.M{"name": "Jimmy Kuu"}).All(&users)

// != ($ne)
c.Find(bson.M{"name": bson.M{"$ne": "Jimmy Kuu"}}).All(&users)

// > ($gt)
c.Find(bson.M{"age": bson.M{"$gt": 32}}).All(&users)

// <($lt)
c.Find(bson.M{"age": bson.M{"$lt": 32}}).All(&users)

// >=($gte)
c.Find(bson.M{"age": bson.M{"$gte": 33}}).All(&users)

// <=($lte)
c.Find(bson.M{"age": bson.M{"$lte": 31}}).All(&users)

// in($in)
c.Find(bson.M{"name": bson.M{"$in": []string{"Jimmy Kuu", "Tracy Yu"}}}).All(&users)


// 多条件查询

// and($and)
c.Find(bson.M{"name": "Jimmy Kuu", "age": 33}).All(&users)

// or($or)
c.Find(bson.M{"$or": []bson.M{bson.M{"name": "Jimmy Kuu"}, bson.M{"age": 31}}}).All(&users)


// 注意修改单个或多个字段需要通过$set操作符号，否则集合会被替换。 修改字段的值($set)
c.Update(bson.M{"_id": bson.ObjectIdHex("5204af979955496907000001")}, bson.M{"$set": bson.M{ "name": "Jimmy Gu", "age": 34, }})


// inc($inc) 字段增加值
c.Update(bson.M{"_id": bson.ObjectIdHex("5204af979955496907000001")}, bson.M{"$inc": bson.M{ "age": -1, }})


// push($push) 从数组中增加一个元素
c.Update(bson.M{"_id": bson.ObjectIdHex("5204af979955496907000001")}, bson.M{"$push": bson.M{ "interests": "Golang", }})


// pull($pull) 从数组中删除一个元素
c.Update(bson.M{"_id": bson.ObjectIdHex("5204af979955496907000001")}, bson.M{"$pull": bson.M{ "interests": "Golang", }})


// 删除
c.Remove(bson.M{"name": "Jimmy Kuu"})
```