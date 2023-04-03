package eorm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type SqlInfo struct {
	selected  bool          // 查询
	insert    bool          // 插入
	delete    bool          // 删除
	update    bool          // 更新
	Table     string        // 表名
	query     []string      // 参数
	column    []string      //
	omit      []string      // 排除字段列表
	where     []string      // 条件语句
	other     []string      // 其它语句
	limit     uint          // 排序
	args      []interface{} // 参数
	whereLock bool          // 条件锁，判断是否已有where条件
	onLock    bool          // on语句锁，判断是否已有On语句
	err       error         // 错误
	//CacheMap  sync.Map      // 缓存
}

func Init() *SqlInfo {
	return &SqlInfo{
		query:  make([]string, 0),
		column: make([]string, 0),
		omit:   make([]string, 0),
		where:  make([]string, 0),
		other:  make([]string, 0),
		args:   make([]interface{}, 0),
	}
}

func (info *SqlInfo) Select(in interface{}, tag, tableName string) *SqlInfo {
	var colNameList = make([]string, 0)
	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Array && !v.IsZero() {
		v = reflect.ValueOf(in.([]interface{})[0])
	}

	// 判断in是否为指针类型， 如果v为指针类型这将v替换为指针对应的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 获取in的类型
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagVal := fi.Tag.Get(tag); tagVal != "" {
			if tagVal == "-" {
				continue
			}
			colNameList = append(colNameList, tagVal)
		}
	}
	info.column = colNameList
	info.Table = tableName
	info.selected = true

	info.column = RemoveRepeatedElement(info.column) // 去重

	return info
}

func (info *SqlInfo) Search(in interface{}, tag, tableName string) *SqlInfo {
	var colNameList = make([]string, 0)
	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Array && !v.IsZero() {
		v = reflect.ValueOf(in.([]interface{})[0])
	}

	// 判断in是否为指针类型， 如果v为指针类型这将v替换为指针对应的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 获取in的类型
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagVal := fi.Tag.Get(tag); tagVal != "" {
			if tagVal == "-" {
				continue
			}
			colNameList = append(colNameList, tagVal)
			value := v.Field(i)
			if value.IsZero() {
				info.Where(tagVal, value.Interface())
			}
		}
	}
	info.column = colNameList
	info.Table = tableName
	info.selected = true

	info.column = RemoveRepeatedElement(info.column) // 去重

	return info
}

func (info *SqlInfo) Insert(in interface{}, tag, tableName string) *SqlInfo {
	var colNameList = make([]string, 0)
	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Array && !v.IsZero() {
		v = reflect.ValueOf(in.([]interface{})[0])
	}

	// 判断in是否为指针类型， 如果v为指针类型这将v替换为指针对应的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 获取in的类型
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagVal := fi.Tag.Get(tag); tagVal != "" {
			if tagVal == "-" {
				continue
			}
			colNameList = append(colNameList, tagVal)
		}
	}

	info.insert = true
	info.column = colNameList
	info.Table = tableName
	return info
}

func (info *SqlInfo) Delete(tableName string) *SqlInfo {
	info.Table = tableName
	info.delete = true
	return info
}

func (info *SqlInfo) Update(in interface{}, tag, tableName string) *SqlInfo {
	var colNameList = make([]string, 0)
	v := reflect.ValueOf(in)

	if v.Kind() == reflect.Array && !v.IsZero() {
		v = reflect.ValueOf(in.([]interface{})[0])
	}

	// 判断in是否为指针类型， 如果v为指针类型这将v替换为指针对应的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 获取in的类型
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagVal := fi.Tag.Get(tag); tagVal != "" {
			field := v.Field(i)
			if tagVal == "-" || tagVal == "id" || field.IsZero() {
				continue
			}
			colNameList = append(colNameList, fmt.Sprintf("%s=:%s", tagVal, tagVal))
		}
	}

	info.update = true
	info.column = colNameList
	info.Table = tableName

	return info
}

func (info *SqlInfo) Where(query string, args ...interface{}) *SqlInfo {
	defer func() {
		if info.whereLock {
			info.query = append(info.query, "AND")
		} else {
			info.query = append(info.query, "WHERE")
		}
		if len(info.where) > 0 {
			info.query = append(info.query, strings.Join(info.where, " AND "))
		}
		info.where = make([]string, 0)
		info.whereLock = true
	}()
	if query == "" {
		info.err = errors.New("QUERY IS NOT NIL")
		return info
	}

	if !strings.Contains(query, "?") {
		info.where = append(info.where, query)
		return info
	}

	if len(args) == 0 || (len(args) > 0 && !strings.Contains(query, "?")) {
		info.err = errors.New("THE QUERY AND ARGS DO NOT MATCH")
		return info
	}
	info.where = append(info.where, query)
	info.args = append(info.args, args...)
	return info
}

//func (info *SqlInfo) Preload(query string, args ...interface{})*SqlInfo{
//
//}

func (info *SqlInfo) LeftJoin(tableName string) *SqlInfo {
	if tableName == "" {
		info.err = errors.New("TABLE IS NOT NULL")
		return info
	}
	info.query = append(info.query, fmt.Sprintf("LEFT JOIN %s", tableName))
	return info
}

func (info *SqlInfo) RightJoin(tableName string) *SqlInfo {
	if tableName == "" {
		info.err = errors.New("TABLE IS NOT NULL")
		return info
	}
	info.query = append(info.query, fmt.Sprintf("RIGHT JOIN %s", tableName))
	return info
}

func (info *SqlInfo) InnerJoin(tableName string) *SqlInfo {
	if tableName == "" {
		info.err = errors.New("TABLE IS NOT NULL")
		return info
	}
	info.query = append(info.query, fmt.Sprintf("INNER JOIN %s", tableName))
	return info
}

func (info *SqlInfo) On(query string, args ...interface{}) *SqlInfo {
	if query == "" {
		info.err = errors.New("QUERY IS NOT NULL")
		return info
	}
	return info
}

func (info *SqlInfo) Others(query string, args ...interface{}) *SqlInfo {
	defer func() {
		if len(info.other) > 0 {
			info.query = append(info.query, info.other...)
		}
		info.other = make([]string, 0)
		info.whereLock = false
	}()
	if query == "" {
		info.err = errors.New("QUERY IS NOT NIL")
		return info
	}
	if !strings.Contains(query, "?") {
		info.other = append(info.other, query)
		return info
	}
	if len(args) == 0 || (len(args) > 0 && !strings.Contains(query, "?")) {
		info.err = errors.New("THE QUERY AND ARGS DO NOT MATCH")
		return info
	}
	info.other = append(info.other, query)
	info.args = append(info.args, args...)
	return info
}

func (info *SqlInfo) InArgs(args ...interface{}) *SqlInfo {
	info.args = append(info.args, args...)
	return info
}

func (info *SqlInfo) Omit(columns ...string) *SqlInfo {
	info.omit = append(info.omit, columns...)
	info.column = SliceDiff(info.omit, info.column) //屏蔽敏感字段
	return info
}

func (info *SqlInfo) ToBind() (string, []interface{}, error) {
	defer info.cleanAll()
	var count = 0
	var args = info.args
	var sql string
	var querys = make([]string, 0)
	if info.err != nil {
		return "", nil, info.err
	}
	defer func() {
		info.err = nil
	}()
	if info.Table == "" {
		return "", nil, errors.New("TABLE NAME IS NOT NIL")
	}
	switch true {
	case info.selected:
		querys = append(querys, "SELECT")
		count++
	case info.insert:
		querys = append(querys, "INSERT INTO")
		count++
	case info.delete:
		querys = append(querys, "DELETE")
		count++
	case info.update:
		querys = append(querys, "UPDATE")
		count++
	case count == 1:
		break
	default:
		return "", nil, errors.New("MULTIPLE ACTIONS CANNOT BE PERFORMED SIMULTANEOUSLY")
	}

	if len(info.column) > 0 && info.selected {
		for idx, _ := range info.column {
			info.column[idx] = fmt.Sprintf("%s.%s", info.Table, info.column[idx])
		}
		querys = append(querys, strings.Join(info.column, ","))
		querys = append(querys, "FROM")
		querys = append(querys, info.Table)
	}
	if info.insert && len(info.column) > 0 {
		columString := strings.Join(info.column, ",")
		valString := strings.Join(info.column, ",:")
		querys = append(querys, fmt.Sprintf("%s(%s) VALUES (:%s)", info.Table, columString, valString))
	}
	if info.update && len(info.column) > 0 {
		valString := strings.Join(info.column, ",")
		querys = append(querys, fmt.Sprintf("%s SET %s", info.Table, valString))
	}
	if info.delete {
		querys = append(querys, info.Table)
	}
	//else {
	//	return "", nil, errors.New("MULTIPLE ACTIONS CANNOT BE PERFORMED SIMULTANEOUSLY")
	//}

	querys = append(querys, info.query...)
	sql = strings.Join(querys, " ")

	return sql, args, nil
}

func (info *SqlInfo) cleanAll() {
	info.query = make([]string, 0)
	info.args = make([]interface{}, 0)
	info.column = make([]string, 0)
	info.omit = make([]string, 0)
	info.other = make([]string, 0)
	info.where = make([]string, 0)
	info.delete = false
	info.update = false
	info.insert = false
	info.whereLock = false
	info.selected = false
	info.onLock = false
}

func RemoveRepeatedElement(in []string) []string {
	var result = make([]string, 0)
	var m = make(map[string]bool)
	for _, v := range in {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

func SliceDiff(slice1, slice2 []string) []string {
	var result = make([]string, 0)
	var m = make(map[string]int, 0)
	if len(slice1) > len(slice2) {
		var tmp = slice1
		slice1 = slice2
		slice2 = tmp
	}
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		if times, _ := m[v]; times < 1 {
			result = append(result, v)
		}
	}
	return result
}
