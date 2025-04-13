package mivus_interface

import (
	"reflect"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"

	client "github.com/milvus-io/milvus/client/v2/milvusclient"
)

// 一列数据
type columnItem[T any | string] struct {
	columnName string
	columnData []T
}

func NewcolumnItem[T any]() *columnItem[T] {
	return &columnItem[T]{}
}

type columnVecItem[T any] struct {
	columnName string
	columnData [][]T
	dim        int
}

func NewcolumnVectItem[T any]() *columnVecItem[T] {
	return &columnVecItem[T]{}
}

type miluvsInsertOption struct {
	collectName string
	//
	columnScar map[int]any //key: ScalarColumnBoolType, value: []*columnItem[T], T: int64, int32, int16, int8, bool, string
	//
	columnVect map[int]any // []*columnVecItem[T]; T:float32, byte, int8
}
type InsertOptions func(*miluvsInsertOption)

func WithCollectName(collName string) InsertOptions {
	return func(item *miluvsInsertOption) {
		item.collectName = collName
	}
}

const (
	ScalarColumnBoolType = 1
	ScalarColumnI8Type   = 2
	ScalarColumnI16Type  = 3
	ScalarColumnI32Type  = 4
	ScalarColumnI64Type  = 5
	ScalarColumnStr      = 6
	//
	VectColumnF32  = 10 //data type float32
	VectColumnF16  = 11 // data type float32
	VectColumnBF16 = 12 // data type float32
	VectColumnByte = 13 //data type type
	// VectColumnI8 = 14
)

// 增加一列向量; dataType 值有：VectColumnF32/VectColumnF16/VectColumnBF16/VectColumnByte, data type 有：float32, type
func WithVectColumnValue[T any](columnName string, dim int, dataType int, data [][]T) InsertOptions {
	return func(item *miluvsInsertOption) {
		var colItem *columnVecItem[T] = &columnVecItem[T]{
			columnName: columnName,
			columnData: data,
			dim:        dim,
		}

		listColumns, hasItems := item.columnVect[dataType]
		switch dataType {
		case VectColumnF32:
			if hasItems {
				f32Column := listColumns.([]*columnVecItem[T])
				f32Column = append(f32Column, colItem)
				item.columnVect[VectColumnF32] = f32Column
			} else {
				f32Column := []*columnVecItem[T]{colItem}
				item.columnVect[VectColumnF32] = f32Column
			}

		case VectColumnF16:
			if hasItems {
				f16Column := listColumns.([]*columnVecItem[T])
				f16Column = append(f16Column, colItem)
				item.columnVect[VectColumnF16] = f16Column
			} else {
				f16Column := []*columnVecItem[T]{colItem}
				item.columnVect[VectColumnF16] = f16Column
			}

		case VectColumnBF16:
			if hasItems {
				bf16Column := listColumns.([]*columnVecItem[T])
				bf16Column = append(bf16Column, colItem)
				item.columnVect[VectColumnBF16] = bf16Column
			} else {
				bf16Column := []*columnVecItem[T]{colItem}
				item.columnVect[VectColumnBF16] = bf16Column
			}

		case VectColumnByte:
			if hasItems {
				bColumn := listColumns.([]*columnVecItem[T])
				bColumn = append(bColumn, colItem)
				item.columnVect[VectColumnByte] = bColumn
			} else {
				bColumn := []*columnVecItem[T]{colItem}
				item.columnVect[VectColumnByte] = bColumn
			}
		}
	}
}

// 增加一列标量, data 类型有：bool, uint8, int16, init32, int64, string
func WithScalarColumnNums[T any](columnName string, data []T) InsertOptions {
	return func(item *miluvsInsertOption) {
		var colItem *columnItem[T] = &columnItem[T]{
			columnName: columnName,
			columnData: data,
		}
		//
		tmp := new(T)
		vType := reflect.ValueOf(tmp).Elem().Kind()
		vTypeInt := ScalarColumnBoolType

		switch vType {
		case reflect.Bool:
			vTypeInt = ScalarColumnBoolType
		case reflect.Int8, reflect.Uint8:
			vTypeInt = ScalarColumnI8Type

		case reflect.Int16, reflect.Uint16:
			vTypeInt = ScalarColumnI16Type
		case reflect.Int32, reflect.Uint32, reflect.Int, reflect.Uint:
			vTypeInt = ScalarColumnI32Type

		case reflect.Uint64, reflect.Int64:
			vTypeInt = ScalarColumnI64Type
			logger.Debugf("is int64 type.")

		case reflect.String:
			vTypeInt = ScalarColumnStr
		default:
			vTypeInt = ScalarColumnBoolType
			logger.Infof("no support vType: %v", vType)
		}
		listColumns, hasItems := item.columnScar[vTypeInt]

		switch vType {
		case reflect.Bool:
			if hasItems {
				boolColumn := listColumns.([]*columnItem[T])
				boolColumn = append(boolColumn, colItem)
				item.columnScar[ScalarColumnBoolType] = boolColumn
			} else {
				boolColumn := []*columnItem[T]{colItem}
				item.columnScar[ScalarColumnBoolType] = boolColumn
			}

		case reflect.Int8, reflect.Uint8:
			if hasItems {
				i8Column := listColumns.([]*columnItem[T])
				i8Column = append(i8Column, colItem)
				item.columnScar[ScalarColumnI8Type] = i8Column
			} else {
				i8Column := []*columnItem[T]{colItem}
				item.columnScar[ScalarColumnI8Type] = i8Column
			}

		case reflect.Int16, reflect.Uint16:
			if hasItems {
				i16Column := listColumns.([]*columnItem[T])
				i16Column = append(i16Column, colItem)
				item.columnScar[ScalarColumnI16Type] = i16Column
			} else {
				i16Column := []*columnItem[T]{colItem}
				item.columnScar[ScalarColumnI16Type] = i16Column
			}

		case reflect.Int32, reflect.Uint32, reflect.Int, reflect.Uint:
			if hasItems {
				i32Column := listColumns.([]*columnItem[T])
				i32Column = append(i32Column, colItem)
				item.columnScar[ScalarColumnI32Type] = i32Column
			} else {
				i32Column := []*columnItem[T]{colItem}
				item.columnScar[ScalarColumnI32Type] = i32Column
			}

		case reflect.Int64, reflect.Uint64:
			item.columnScar[ScalarColumnI64Type] = colItem
			if hasItems {
				i64Column := listColumns.([]*columnItem[T])
				i64Column = append(i64Column, colItem)
				item.columnScar[ScalarColumnI64Type] = i64Column
			} else {
				i64Column := []*columnItem[T]{colItem}
				item.columnScar[ScalarColumnI64Type] = i64Column
			}

		case reflect.String:
			item.columnScar[ScalarColumnStr] = colItem
			if hasItems {
				strColumn := listColumns.([]*columnItem[T])
				strColumn = append(strColumn, colItem)
				item.columnScar[ScalarColumnStr] = strColumn
			} else {
				strColumn := []*columnItem[T]{colItem}
				item.columnScar[ScalarColumnStr] = strColumn
			}
		}
	}
}

type InsertOps struct {
	options      *miluvsInsertOption
	insertOption client.InsertOption
}

func NewInsertOpts() *InsertOps {
	return &InsertOps{
		options: &miluvsInsertOption{
			columnScar: make(map[int]any),
			columnVect: make(map[int]any),
		},
	}
}

// Register 构造insert options.
func (iops *InsertOps) Register(collectName string, ops ...InsertOptions) {
	WithCollectName(collectName)(iops.options)
	for _, op := range ops {
		if op == nil {
			continue
		}
		op(iops.options)
	}
	columnDataOption := client.NewColumnBasedInsertOption(iops.options.collectName)

	for scalarKey, scalarValue := range iops.options.columnScar {
		switch scalarKey {
		case ScalarColumnBoolType:
			columnList, ok := scalarValue.([]*columnItem[bool])
			if ok {
				for _, columnItem := range columnList {
					if columnItem == nil {
						continue
					}
					columnDataOption.WithBoolColumn(columnItem.columnName, columnItem.columnData)
				}
			}
		case ScalarColumnI8Type:
			columnList, ok := scalarValue.([]*columnItem[int8])
			if ok {
				for _, columnItem := range columnList {
					if columnItem == nil {
						continue
					}
					columnDataOption.WithInt8Column(columnItem.columnName, columnItem.columnData)
				}
			}
		case ScalarColumnI16Type:
			columnList, ok := scalarValue.([]*columnItem[int16])
			if ok {
				for _, columnItem := range columnList {
					if columnItem == nil {
						continue
					}
					columnDataOption.WithInt16Column(columnItem.columnName, columnItem.columnData)
				}
			}
		case ScalarColumnI32Type:
			columnList, ok := scalarValue.([]*columnItem[int32])
			if ok {
				for _, columnItem := range columnList {
					if columnItem == nil {
						continue
					}
					columnDataOption.WithInt32Column(columnItem.columnName, columnItem.columnData)
				}
			}
		case ScalarColumnI64Type:
			columnList, ok := scalarValue.([]*columnItem[int64])
			if ok {
				for _, columnItem := range columnList {
					if columnItem == nil {
						continue
					}
					columnDataOption.WithInt64Column(columnItem.columnName, columnItem.columnData)
				}
			}
		case ScalarColumnStr:
			columnList, ok := scalarValue.([]*columnItem[string])
			if ok {
				for _, columnItem := range columnList {
					if columnItem == nil {
						continue
					}
					columnDataOption.WithVarcharColumn(columnItem.columnName, columnItem.columnData)
				}
			}
		}
	}

	for vectKey, vectValue := range iops.options.columnVect {
		switch vectKey {
		case VectColumnF32:
			columnList, ok := vectValue.([]*columnVecItem[float32])
			if ok {
				for _, columnItem := range columnList {
					if columnItem == nil {
						continue
					}
					columnDataOption.WithFloatVectorColumn(columnItem.columnName, columnItem.dim, columnItem.columnData)
				}
			}
		case VectColumnF16:
			columnList, ok := vectValue.([]*columnVecItem[float32])
			if ok {
				for _, columnItem := range columnList {
					if columnItem == nil {
						continue
					}
					columnDataOption.WithFloat16VectorColumn(columnItem.columnName, columnItem.dim, columnItem.columnData)
				}
			}
		case VectColumnBF16:
			columnList, ok := vectValue.([]*columnVecItem[float32])
			if ok {
				for _, columnItem := range columnList {
					if columnItem == nil {
						continue
					}
					columnDataOption.WithBFloat16VectorColumn(columnItem.columnName, columnItem.dim, columnItem.columnData)
				}
			}

		case VectColumnByte:
			columnList, ok := vectValue.([]*columnVecItem[byte])
			if ok {
				for _, columnItem := range columnList {
					if columnItem == nil {
						continue
					}
					columnDataOption.WithBinaryVectorColumn(columnItem.columnName, columnItem.dim, columnItem.columnData)
				}
			}
		}
	}
	logger.Infof("collumn data options: %+v", columnDataOption)
	iops.insertOption = columnDataOption
}

func (iops *InsertOps) BuildInsertOption() (any, bool) {
	return iops.insertOption, true
}
