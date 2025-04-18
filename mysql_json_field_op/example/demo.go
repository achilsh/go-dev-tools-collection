package main

import (
	"encoding/json"
	"fmt"
	"time"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	. "github.com/achilsh/go-dev-tools-collection/mysql_json_field_op"
	"gorm.io/datatypes"
)

func main() {
	InitDB()
	//
	InitDemoDBTableInstance()

	//
	RunJsonFieldOp()

}

type UserItem struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

func RunJsonFieldOp() {
	// datatypes.JSON_ARRAY
	userItemStr, _ := json.Marshal(&UserItem{
		Name:    "achilsh",
		Age:     18,
		Address: "guangzhou.zh",
	})
	toInsertItems := []*DemoTestTabRecord{
		&DemoTestTabRecord{
			ID:           time.Now().UTC().UnixMilli(),
			SessionId:    200,
			Language:     "en",
			QuestionType: 1,
			// QuestionWordContent: datatypes.NewJSONType([]string{"sdfad"}),
			// QuestionOrigin:      datatypes.NewJSONType([]string{"dfadfa", "xxx"}),
			QuestionWordContent: datatypes.NewJSONSlice([]string{"sdfad"}),
			QuestionOrigin:      datatypes.NewJSONSlice([]string{"dfadfa", "xxx"}),
			//
			AnswerType:    2,
			AnswerContent: datatypes.NewJSONSlice([]string{"我是谁", "我是 achilsh"}),
			UserContent:   datatypes.JSON(userItemStr),
		},

		&DemoTestTabRecord{
			ID:                  time.Now().UTC().UnixMilli() + 100,
			SessionId:           201,
			Language:            "en",
			QuestionType:        1,
			QuestionWordContent: datatypes.NewJSONSlice([]string{"1111", "bbbb", "cccc"}),
			QuestionOrigin:      datatypes.NewJSONSlice([]string{"http://www.sina.com.cn", "http://baidu.com"}),
			//
			AnswerType:    2,
			AnswerContent: datatypes.NewJSONSlice([]string{"我是谁", "我是 achilsh"}),
		},
	}
	// nums err := DemoDbGetObj().Insert(toInsertItems)
	// insertSql := fmt.Sprintf("insert into history_question_answer_record(id,session_id, language, question_type,  question_words, question_origin, answer_type,answer_content)")
	// insertSql += fmt.Sprintf(" VALUES(?,?,?,?,JSON_ARRAY(?),JSON_ARRAY(?),?,JSON_ARRAY(?))")
	// DemoDbGetObj().GetDBHandle().Raw(insertSql, toInsertItems[0].ID, toInsertItems[0].SessionId, toInsertItems[0].Language, toInsertItems[0].QuestionType,
	// 	toInsertItems[0].QuestionWordContent, toInsertItems[0].QuestionOrigin, toInsertItems[0].AnswerType, toInsertItems[0].AnswerContent)

	logger.Infof("------- 插入 带有 json 字段的数据到 表中 ------------")
	//内置函数： json_object() , 插入时。
	//查询 jdoc 字段中 key  为 x 记录， 使用如下方式：
	//SELECT jdoc->'$.x' FROM t1;
	// SELECT JSON_EXTRACT(jdoc,'$.x') FROM t1;

	//返回值去掉 ""
	// SELECT JSON_UNQUOTE(json_extract('{"id": 14, "name": "Aztalan"}', '$.name'));

	n, err := DemoDbGetObj().Insert(toInsertItems)
	// DemoDbGetObj().GetDBHandle().Create(toInsertItems)
	if err != nil {
		logger.Errorf("insertg fail, err: %v", err)
		return
	}
	//query ....
	logger.Infof("insert item nums: %v", n)
	findRet, err := DemoDbGetObj().QueryItemOnCond("session_id = ? and language = ?", []any{201, "en"}, map[string]bool{"id": false}, 0, 4)
	if err != nil {
		logger.Errorf("query fail, err: %v", err)
		return
	}
	for _, item := range findRet {
		if item == nil {
			logger.Errorf("is nil for item.")
			continue
		}
		logger.Infof("item: %+v", *item)
	}

	// 查询 json 是数组， 从数组中查找 是否包含 某个值
	logger.Debugf("-----------------使用 json 字段查询 --------------------")
	sqlQuery := fmt.Sprintf("select *  from  %v where JSON_CONTAINS(question_origin, ?)", DemoDbGetObj().GetTabName())
	var retfind = []DemoTestTabRecord{}
	err = DemoDbGetObj().GetDBHandle().Raw(sqlQuery, "\"http://www.sina.com.cn\"").Scan(&retfind).Error
	if err != nil {
		logger.Errorf("query fail, in contails, err: %v", err)
		return
	}
	var idList = []int64{}
	for _, item := range retfind {
		logger.Debugf("item: %+v", item)
		if item.ID > 0 {
			idList = append(idList, item.ID)
		}
	}

	// 查询 json map 字段key 是否存在：
	logger.Infof("------- query json field value ..................")

	findMapRet := []DemoTestTabRecord{}
	// SELECT * FROM `users` WHERE JSON_EXTRACT(`user`, '$.age') IS NOT NULL
	err = DemoDbGetObj().GetDBHandle().Find(&findMapRet, datatypes.JSONQuery("user").HasKey("age")).Error
	if err != nil {
		logger.Errorf("find user info error: %+v", err)
	} else {
		logger.Infof("find user has age item lists: %+v", findMapRet)
	}

	// 查询 json map中key 对应的value 是否等于某个值
	logger.Debugf("------------find map key 's value is eq some value.........")
	findMapRetEqValue := []DemoTestTabRecord{}
	//SELECT * FROM `user` WHERE JSON_EXTRACT(`user`, '$.age') = 18
	//  Check JSON extract value from keys equal to value
	err = DemoDbGetObj().GetDBHandle().Find(&findMapRetEqValue, datatypes.JSONQuery("user").Equals(18, "age")).Error
	if err != nil {
		logger.Errorf("find user info error: %+v", err)
	} else {
		logger.Infof("find user has age item and value is 18, lists: %+v", findMapRet)
	}

	// 修改 json 中的某个 field 对应的value, if field 存在则更新，不存在则插入
	//更新一个不存map中的 值：
	err = DemoDbGetObj().GetDBHandle().
		Model(&DemoTestTabRecord{}).
		Where("id = ?", 1744445745841).
		UpdateColumn("user", datatypes.JSONSet("user").Set("age", 1000)).
		Error
	if err != nil {
		logger.Errorf("update id 1744445745841 user.age fail, err: %v", err)
		return
	} else {
		logger.Infof("update id 1744445745841 user.age succ.")
	}

}
