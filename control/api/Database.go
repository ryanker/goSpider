package api

import (
	"errors"
	"os"
	"time"

	"github.com/xiuno/gin"

	"../../lib/dbs"
	"../../model"
)

// 数据库结构设置
func DatabaseSet(c *gin.Context) {
	m := struct {
		Rid  int64
		Type int64
	}{}
	err := c.ShouldBind(&m)
	if err != nil {
		c.Message("-1", "参数不正确: "+err.Error())
		return
	}

	// 1、读取规则
	row, err := model.RuleRead(m.Rid)
	if err != nil {
		c.Message("-1", "读取规则失败: "+err.Error())
		return
	}

	// 2、拼接创建SQL
	s, err := generateSql(m.Rid)
	if err != nil {
		c.Message("-1", "读取参数失败: "+err.Error())
		return
	}

	// 3、创建数据表
	err = createDatabase(m.Type, row.DateBase, s)
	if err != nil {
		c.Message("-1", "创建数据库失败: "+err.Error())
		return
	}

	c.Message("0", "success")

}

// 拼接创建SQL（4 张表：List Content ListDownload ContentDownload）
func generateSql(Rid int64) (s string, err error) {
	ListData, err := model.RuleParamList(dbs.H{"Rid": Rid, "Type": "List"}, 0, 0)
	if err != nil {
		return
	}
	ContentData, err := model.RuleParamList(dbs.H{"Rid": Rid, "Type": "Content"}, 0, 0)
	if err != nil {
		return
	}

	// List
	s += "CREATE TABLE List(\n"
	s += "  `Id` INTEGER PRIMARY KEY AUTOINCREMENT,\n"
	s += "  `Status` INTEGER NOT NULL DEFAULT '0',\n"
	for _, v := range ListData {
		Suffix := ""
		if v.FieldType == "INTEGER" {
			Suffix = "INTEGER NOT NULL DEFAULT '0',"
		} else if v.FieldType == "VARCHAR" {
			Suffix = "VARCHAR(255) NOT NULL DEFAULT '',"
		} else if v.FieldType == "TEXT" {
			Suffix = "TEXT DEFAULT '',"
		}
		s += "  `" + v.Field + "` " + Suffix + "\n"
	}
	s += "  `CreateDate` DATETIME DEFAULT CURRENT_TIMESTAMP\n"
	s += ");\n"

	// Content
	s += "CREATE TABLE Content(\n"
	s += "  `Id` INTEGER PRIMARY KEY AUTOINCREMENT,\n"
	for _, v := range ContentData {
		Suffix := ""
		if v.FieldType == "INTEGER" {
			Suffix = "INTEGER NOT NULL DEFAULT '0',"
		} else if v.FieldType == "VARCHAR" {
			Suffix = "VARCHAR(255) NOT NULL DEFAULT '',"
		} else if v.FieldType == "TEXT" {
			Suffix = "TEXT DEFAULT '',"
		}
		s += "  `" + v.Field + "`" + Suffix + "\n"
	}
	s += "  `CreateDate` DATETIME DEFAULT CURRENT_TIMESTAMP\n"
	s += ");\n"

	// Download
	s += `CREATE TABLE ListDownload
(
  Id           INTEGER PRIMARY KEY AUTOINCREMENT,
  Status       INTEGER      NOT NULL DEFAULT '0',
  Field        VARCHAR(255) NOT NULL DEFAULT '',
  OldUrl       VARCHAR(255) NOT NULL DEFAULT '',
  NewUrl       VARCHAR(255) NOT NULL DEFAULT '',
  FileSize     INTEGER      NOT NULL DEFAULT '0',
  Sort         INTEGER      NOT NULL DEFAULT '0',
  DownloadDate DATETIME              DEFAULT CURRENT_TIMESTAMP,
  CreateDate   DATETIME              DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ContentDownload
(
  Id           INTEGER PRIMARY KEY AUTOINCREMENT,
  Status       INTEGER      NOT NULL DEFAULT '0',
  Field        VARCHAR(255) NOT NULL DEFAULT '',
  OldUrl       VARCHAR(255) NOT NULL DEFAULT '',
  NewUrl       VARCHAR(255) NOT NULL DEFAULT '',
  FileSize     INTEGER      NOT NULL DEFAULT '0',
  Sort         INTEGER      NOT NULL DEFAULT '0',
  DownloadDate DATETIME              DEFAULT CURRENT_TIMESTAMP,
  CreateDate   DATETIME              DEFAULT CURRENT_TIMESTAMP
);`
	// fmt.Println(s)
	return
}

// 创建数据库
func createDatabase(Type int64, DateBase string, s string) (err error) {
	dbFile := "./db/" + DateBase + ".db"
	db, err := dbs.Open(dbFile)
	if err != nil {
		return
	}

	// 文件不存在，则创建表
	if Type == 1 {
		_, err = db.Exec(s)
		if err != nil {
			return
		}
	} else if Type == 2 {
		if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
			newPath := dbFile + "." + time.Now().Format("20060102_150405")
			if err := os.Rename(dbFile, newPath); err != nil {
				err = errors.New("备份数据库失败: " + err.Error())
				return
			}
		}

		_, err = db.Exec(s)
		if err != nil {
			return
		}
	} else if Type == 3 {
		if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
			if err = os.Remove(dbFile); err != nil {
				err = errors.New("删除文件失败: " + err.Error())
				return
			}
		}

		_, err = db.Exec(s)
		if err != nil {
			return
		}
	}
	return
}
