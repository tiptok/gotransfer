package comm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
)

var (
	ColumnNotExist = errors.New("DataTable not exsits colunm")
)

type SqlHelper struct {
}

type DataTable struct {
	Rows      *sql.Rows
	values    map[string]interface{}
	columns   []string //列名
	valuesPtr []interface{}
	Index     int
	mux       sync.RWMutex
}

//NewDataTable  新建一个DataTable实例.
func NewDataTable(rows *sql.Rows) *DataTable {
	dt := &DataTable{}
	dt.columns, _ = rows.Columns()
	dt.values = make(map[string]interface{}, len(dt.columns))
	dt.valuesPtr = make([]interface{}, len(dt.columns))
	for i := range dt.valuesPtr {
		var ptr interface{}
		dt.valuesPtr[i] = &ptr
	}
	dt.Rows = rows
	return dt
}

//ReadRow DataTable 读取当前行.
func (dt DataTable) ReadRow() {
	dt.mux.Lock()
	defer dt.mux.Unlock()
	err := dt.Rows.Scan(dt.valuesPtr...)
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(dt.columns); i++ {
		dt.values[dt.columns[i]] = reflect.ValueOf(dt.valuesPtr[i]).Elem().Interface()
	}
	dt.Index++
}

func (dt DataTable) NextRow() bool {
	return dt.Rows.Next()
}

func (dt DataTable) ReadNextRow() bool {
	if dt.NextRow() {
		dt.ReadRow()
		return true
	}
	return false
}

//GetColumn  获取列数据
//sColumnName 列名
//val         返回值，传递指针
func (dt DataTable) GetColumn(sColumnName string, val interface{}) error {
	dt.mux.RLock()
	defer dt.mux.RUnlock()
	if v, ok := dt.values[sColumnName]; ok {
		return ConvertAssign(val, v)
	}
	return fmt.Errorf("dataTable not exsits colunm:%s", sColumnName)
}

//GetColumnDefault  获取列数据，未取到值的按默认值返回
//sColumnName 列名
//val         返回值，传递指针
//def		  默认值
func (dt DataTable) GetColumnDefault(sColumnName string, val, def interface{}) error {
	dt.mux.RLock()
	defer dt.mux.RUnlock()
	if _, ok := dt.values[sColumnName]; ok {
		if def != nil {
			ConvertAssign(val, def)
			return ColumnNotExist
		}
	}
	return dt.GetColumn(sColumnName, val)
}

// func (sqlHelper) GetColumn() {

// }
