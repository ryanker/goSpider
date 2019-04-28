package model

import (
	"database/sql"

	"../lib/dbs"
)

type Setting struct {
	Key   string
	Value string
}

func SettingMap() (ptr *Setting, fields string, args *[]interface{}) {
	row := Setting{}
	fields, scanArr := dbs.GetSqlRead(dbs.H{
		"Key":   &row.Key,
		"Value": &row.Value,
	})
	ptr = &row
	args = &scanArr
	return
}

func SettingSet(Key string, Value string) (err error) {
	_, err = SettingRead(Key)
	if err == nil {
		_, err = db.Update("Setting", dbs.H{"Value": Value}, dbs.H{"Key": Key})
	} else if err == sql.ErrNoRows {
		_, err = db.Insert("Setting", dbs.H{"Key": Key, "Value": Value})
	}
	return
}

func SettingGet(Key string) (Value string, err error) {
	row, err := SettingRead(Key)
	if err == nil {
		return row.Value, nil
	} else if err == sql.ErrNoRows {
		return "", nil
	}
	return "", err
}

func SettingRead(Key string) (row Setting, err error) {
	data, fields, scanArr := SettingMap()
	err = db.Read("Setting", fields, *scanArr, dbs.H{"Key": Key})
	row = *data
	return
}

func SettingList() (list []Setting, err error) {
	data, fields, scanArr := SettingMap()
	err = db.Find("Setting", fields, *scanArr, dbs.H{}, "Sort ASC", 0, 1000, func() {
		list = append(list, *data)
	})
	if err != nil {
		return
	}
	return
}
