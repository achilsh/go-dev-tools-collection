package lib

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ModelTplConstraint 定义类型参数约束, 对多个数据模型，可以使用 或 | 方式, eg:  eg: bean.AAAA | bean.BBB 接口约束
type ModelTplConstraint interface {
	any
}

// DbOpsTplInterface 接口模板化， 通用 model 的数据操作接口。每个model实例化后就是对应的model 数据接口
type DbOpsTplInterface[T ModelTplConstraint] interface {
	// GetDBHandle 获取底层数据库的句柄
	GetDBHandle() *gorm.DB

	// Insert 插入记录
	Insert(items []*T) (int64, error)

	// UpdateItems 更新数据库表字段
	UpdateItems(whereCond string, whereValue []any, fieldName string, fieldNewValue any) (int64, error)

	// Updates 全量更新表字段
	Updates(whereCond string, whereValue []any, fieldNewValue any) (int64, error)

	// DelItems 根据条件删除记录（物理删除）
	DelItems(whereCond string, whereVal []any) error

	// CountItems 查询记录个数
	CountItems(whereCond string, whereVal []any) (int64, error)

	// DelItemOnID 根据主键删除记录
	DelItemOnID(ids []int64) error

	// QueryItemOnCond 查询记录
	QueryItemOnCond(whereCond string, whereVal []any, orderAsc map[string]bool, offset int, limit int) ([]*T, error)

	// QueryItems 根据主键查询记录
	QueryItems(ids []int64) ([]*T, error)

	// RawQueryRun 执行底层sql 语句
	RawQueryRun(whereCond string, whereVal []any) ([]*T, error)

	// GetTabName 获取当前表名
	GetTabName() string

	// SetTabName 设置表名
	SetTabName(name string) DbOpsTplInterface[T]
}

// DBModeTplWrapper struct 类型模板化，（可有函数模板化）实现 db 操作接口 DbOpsTplInterface[T] 的对象类型
type dBModeTplWrapper[T ModelTplConstraint] struct {
	tabName string
	dbOps   *gorm.DB
	isDebug bool
}

// NewDBModelTplWrapper 创建数据模型模板
func NewDBModelTplWrapper[T ModelTplConstraint](db *gorm.DB) DbOpsTplInterface[T] {
	item := &dBModeTplWrapper[T]{}
	item.isDebug = true
	item.dbOps = db
	return item
}
func (a *dBModeTplWrapper[T]) GetDBHandle() *gorm.DB {
	if a == nil {
		return nil
	}
	return a.dbOps
}
func (a *dBModeTplWrapper[T]) checkInParam() error {
	if a == nil || a.dbOps == nil {
		return fmt.Errorf("obj is nil or db handle is nil")
	}
	return nil
}
func (a *dBModeTplWrapper[T]) GetTabName() string {
	if a == nil {
		return ""
	}
	return a.tabName
}
func (a *dBModeTplWrapper[T]) SetTabName(tabName string) DbOpsTplInterface[T] {
	if a == nil {
		return nil
	}
	a.tabName = tabName
	return a
}
func (a *dBModeTplWrapper[T]) CountItems(whereCond string, whereVal []any) (int64, error) {
	if err := a.checkInParam(); err != nil {
		return 0, err
	}
	if len(whereVal) == 0 {
		return 0, fmt.Errorf("where cond is empty")
	}
	if len(whereVal) == 0 {
		return 0, fmt.Errorf("where cond value is empty")
	}

	var execuateDb *gorm.DB = nil
	var retItemNums int64 = 0
	// eg: db.Where("name = ? AND age >= ?", "jinzhu", "22")
	execuateDb = a.dbOps.Table(a.tabName).Where(whereCond, whereVal...).Select("id").Count(&retItemNums)
	if execuateDb.Error != nil {
		if errors.Is(execuateDb.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("query  item nums  fail, e: %v", execuateDb.Error)
	}
	return retItemNums, nil
}
func (a *dBModeTplWrapper[T]) Insert(items []*T) (int64, error) {
	if e := a.checkInParam(); e != nil {
		return 0, e
	}
	//
	if len(items) == 0 {
		return 0, nil
	}

	insertRet := a.dbOps.Table(a.tabName).Create(items)
	if insertRet.Error != nil {
		fmt.Printf("insert items to tab: %s,  err: %v, item: %+v", a.tabName, insertRet.Error, items)
		return 0, insertRet.Error
	}

	//logger.Infof("insert item to tab: %s, success insert nums: %v, to insert nums: %v", a.tabName, insertRet.RowsAffected, len(items))
	return insertRet.RowsAffected, nil
}

// Updates  fieldValue may be map[]any = {field_name: field_value}
func (a *dBModeTplWrapper[T]) Updates(whereCond string, whereValue []any, fieldNewValue any) (int64, error) {
	if a == nil {
		return 0, errors.New("is empty for object")
	}
	if len(whereCond) == 0 || len(whereValue) == 0 {
		return 0, errors.New("where cond or where value is empty")
	}

	var execuateDb *gorm.DB = nil
	execuateDb = a.dbOps.Table(a.tabName).Where(whereCond, whereValue...).Updates(fieldNewValue)
	if execuateDb == nil || execuateDb.Error != nil {
		return 0, fmt.Errorf("update ,fail, err: %v", execuateDb.Error)
	}
	return execuateDb.RowsAffected, nil
}
func (a *dBModeTplWrapper[T]) UpdateItems(whereCond string, whereValue []any, fieldName string, fieldNewValue any) (int64, error) {
	if a == nil {
		return 0, errors.New("is empty for object")
	}
	if len(whereCond) == 0 || len(whereValue) == 0 {
		return 0, errors.New("where cond or where value is empty")
	}

	if len(fieldName) == 0 {
		return 0, errors.New("update field name is empty")
	}
	var execuateDb *gorm.DB = nil
	execuateDb = a.dbOps.Table(a.tabName).Debug().Where(whereCond, whereValue...).Update(fieldName, fieldNewValue)
	if execuateDb == nil || execuateDb.Error != nil {
		return 0, fmt.Errorf("update on feildName: %v, fail, err: %v", fieldName, execuateDb.Error)
	}
	return execuateDb.RowsAffected, nil
}

// DelItems 根据条件删除记录（物理删除）
func (a *dBModeTplWrapper[T]) DelItems(whereCond string, whereVal []any) error {
	if a == nil || a.dbOps == nil {
		return fmt.Errorf("input param is nil")
	}
	if len(whereVal) == 0 {
		return fmt.Errorf("where cond is empty")
	}
	if len(whereVal) == 0 {
		return fmt.Errorf("where cond value is empty")
	}

	var toDelItem []T
	delRet := a.dbOps.Table(a.tabName).Where(whereCond, whereVal...).Delete(&toDelItem)
	if delRet.Error != nil {
		return fmt.Errorf("del item fail, e: %w, whereCond: %v, whereVal: %v", delRet.Error, whereCond, whereVal)
	}

	fmt.Printf("del item nums: %v", delRet.RowsAffected)
	return nil
}

// DelItemOnID 根据主键删除记录
func (a *dBModeTplWrapper[T]) DelItemOnID(ids []int64) error {
	if a == nil || a.dbOps == nil {
		return fmt.Errorf("input param is nil")
	}

	if len(ids) == 0 {
		return nil
	}
	var (
		delRet     *gorm.DB = nil
		toDelItems []T
	)
	delRet = a.dbOps.Table(a.tabName).Delete(&toDelItems, ids)

	if delRet == nil || delRet.Error != nil {
		return fmt.Errorf("del item by key id: %v, fail, e: %v", ids, delRet.Error)
	}
	fmt.Printf("del item nums: %v", delRet.RowsAffected)
	return nil
}

func (a *dBModeTplWrapper[T]) QueryItemOnCond(whereCond string, whereVal []any, orderAsc map[string]bool, offset int, limit int) ([]*T, error) {
	if err := a.checkInParam(); err != nil {
		return nil, err
	}
	if len(whereVal) == 0 {
		return nil, fmt.Errorf("where cond is empty")
	}
	if len(whereVal) == 0 {
		return nil, fmt.Errorf("where cond value is empty")
	}

	var execuateDb *gorm.DB = nil
	// eg: db.Where("name = ? AND age >= ?", "jinzhu", "22")
	execuateDb = a.dbOps.Table(a.tabName).Where(whereCond, whereVal...)
	for fieldName, direct := range orderAsc {
		if len(fieldName) == 0 {
			continue
		}

		if direct == false { // desc
			execuateDb = execuateDb.Order(fmt.Sprintf("%s desc", fieldName))
		} else {
			execuateDb = execuateDb.Order(fmt.Sprintf("%s asc", fieldName))
		}
	}

	if offset > 0 {
		execuateDb = execuateDb.Offset(offset)
	}
	if limit > 0 {
		execuateDb = execuateDb.Limit(limit)
	}

	var findRet []T
	e := execuateDb.Find(&findRet).Error
	if e != nil {
		return nil, fmt.Errorf("find fail, e: %v", e)
	}

	var ret []*T
	for _, v := range findRet {
		vv := v
		ret = append(ret, &vv)
	}
	return ret, nil
}

func (a *dBModeTplWrapper[T]) QueryItems(ids []int64) ([]*T, error) {
	if err := a.checkInParam(); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}

	var queryRet []T
	e := a.dbOps.Table(a.tabName).Find(&queryRet, ids).Error
	if e != nil {
		return nil, fmt.Errorf("query by id: %v, fail, e: %v", ids, e)
	}

	var ret []*T = nil
	for _, v := range queryRet {
		vv := v
		ret = append(ret, &vv)
	}
	return ret, nil
}
func (a *dBModeTplWrapper[T]) RawQueryRun(whereCond string, whereVal []any) ([]*T, error) {
	if err := a.checkInParam(); err != nil {
		return nil, err
	}
	if len(whereVal) == 0 {
		return nil, fmt.Errorf("where cond is empty")
	}
	if len(whereVal) == 0 {
		return nil, fmt.Errorf("where cond value is empty")
	}

	var findRet []T
	e := a.dbOps.Raw(whereCond, whereVal...).Scan(&findRet).Error
	if e != nil {
		return nil, fmt.Errorf("find fail, e: %v", e)
	}

	var ret []*T
	for _, v := range findRet {
		vv := v
		ret = append(ret, &vv)
	}
	return ret, nil
}
