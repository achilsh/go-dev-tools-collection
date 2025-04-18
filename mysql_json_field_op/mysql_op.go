package mysql_json_field_op

import (
	"github.com/achilsh/go-dev-tools-collection/mysql-wrapper/lib"
	"gorm.io/datatypes"
)

type DemoTestTabRecord struct {
	ID        int64  `gorm:"column:id;primaryKey"` // 每次对话id
	SessionId int64  `gorm:"column:session_id"`
	Language  string `gorm:"column:language"`
	// 请求数据
	QuestionType int `gorm:"column:question_type"` //提问内容类型
	// QuestionWordContent datatypes.JSONType[[]string] `gorm:"column:question_words; type:json"`  //存的文字，包括其他格式转义后的文字
	// QuestionOrigin      datatypes.JSONType[[]string] `gorm:"column:question_origin; type:json"` //语音url，文件url等

	QuestionWordContent datatypes.JSONSlice[string] `gorm:"column:question_words;type:json"`  //存的文字，包括其他格式转义后的文字
	QuestionOrigin      datatypes.JSONSlice[string] `gorm:"column:question_origin;type:json"` //语音url，文件url等

	// `gorm:"配置项" json:"字段名" xml:"字段名"`
	// `gorm:"column:field_name;type:varchar(100);not null"`
	AnswerType int `gorm:"column:answer_type"` //回答内容类型
	// AnswerContent datatypes.JSONType[[]string] `gorm:"column:answer_content;type:json"` //url
	AnswerContent datatypes.JSONSlice[string] `gorm:"column:answer_content;type:json"` //url

	UserContent datatypes.JSON `gorm:"column:user;type:json"` // user object
}

func (DemoTestTabRecord) TableName() string {
	return "history_question_answer_record"
}

type demoTestRecordType = lib.DbOpsTplInterface[DemoTestTabRecord]

// 定义 db 表操作的实例，后续对该表的操作使用该实例来操作
var (
	demoTestDbHandle demoTestRecordType = nil
)

// 直接使用该函数就可以操作
func DemoDbGetObj() demoTestRecordType {
	return demoTestDbHandle
}
func InitDemoDBTableInstance() {
	demoTestDbHandle = lib.NewDBModelTplWrapper[DemoTestTabRecord](GetDB()).SetTabName(DemoTestTabRecord{}.TableName())
}
