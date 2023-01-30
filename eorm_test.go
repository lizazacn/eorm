package eorm

import (
	"fmt"
	"testing"
)

type Student struct {
	Id   string `db:"id"`
	Name string `db:"name"`
	Sex  int    `db:"sex"`
	Age  int    `db:"age"`
}

func (s Student) Table() string {
	return "student"
}

var sqlInfo = Init()

func TestSqlInfo_Insert(t *testing.T) {
	var student Student
	sql, args, err := sqlInfo.Insert(student, "db", student.Table()).InArgs(student).ToBind()
	if err != nil {
		return
	}
	fmt.Printf("SQL语句：%s\n参数：%s\n", sql, args)
}

func TestSqlInfo_Delete(t *testing.T) {
	var student Student
	sql, args, err := sqlInfo.Delete(student.Table()).Where("id=?", "1").ToBind()
	if err != nil {
		return
	}
	fmt.Printf("SQL语句：%s\n参数：%s\n", sql, args)
}

func TestSqlInfo_Select(t *testing.T) {
	var student Student
	sql, args, err := sqlInfo.Select(student, "db", student.Table()).Where("id=?", "1").ToBind()
	if err != nil {
		return
	}
	fmt.Printf("SQL语句：%s\n参数：%s\n", sql, args)
}

func TestSqlInfo_Update(t *testing.T) {
	var student Student
	sql, args, err := sqlInfo.Update(student, "db", student.Table()).Where("id=?", "1").ToBind()
	if err != nil {
		return
	}
	fmt.Printf("SQL语句：%s\n参数：%s\n", sql, args)
}

func TestAll(t *testing.T) {
	var student Student
	student.Id = "1"
	student.Name = "xiaoming"
	student.Sex = 1
	student.Age = 18
	sql, args, err := sqlInfo.Update(student, "db", student.Table()).Where("id=?", "1").ToBind()
	if err != nil {
		return
	}
	fmt.Printf("SQL语句：%s\n参数：%s\n\n", sql, args)
	sql, args, err = sqlInfo.Delete(student.Table()).Where("id=?", "1").ToBind()
	if err != nil {
		return
	}
	fmt.Printf("SQL语句：%s\n参数：%s\n", sql, args)
}
