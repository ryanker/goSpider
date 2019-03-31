-- 第一步：分类（采集哪个网站）
-- 第二步：规则名称
-- 第三步：抓取规则细节参数

-- 分类表
CREATE TABLE RuleCate
(
  CateId     INTEGER PRIMARY KEY AUTOINCREMENT, -- 分类ID
  Name       VARCHAR(255) NOT NULL DEFAULT '',  -- 分类名称
  Brief      VARCHAR(255) NOT NULL DEFAULT '',  -- 分类简述
  Url        VARCHAR(255) NOT NULL DEFAULT '',  -- 目标网址
  DateBase   VARCHAR(255) NOT NULL DEFAULT '',  -- 数据库名（一个分类，一个库）
  CreateDate DATETIME              DEFAULT CURRENT_TIMESTAMP
);

-- 规则表
CREATE TABLE Rule
(
  Rid           INTEGER PRIMARY KEY AUTOINCREMENT, -- 规则ID
  CateId        INTEGER      NOT NULL DEFAULT '0', -- 分类ID
  Name          VARCHAR(255) NOT NULL DEFAULT '',  -- 规则名称
  Brief         VARCHAR(255) NOT NULL DEFAULT '',  -- 规则简述

  ListTable     VARCHAR(255) NOT NULL DEFAULT '',  -- 列表表名称
  ListUrl       VARCHAR(255) NOT NULL DEFAULT '',  -- 抓取列表网址
  ListPageStart INTEGER      NOT NULL DEFAULT '0', -- 列表开始页码
  ListPageEnd   INTEGER      NOT NULL DEFAULT '0', -- 列表结束页码
  ListPageSize  INTEGER      NOT NULL DEFAULT '0', -- 每页间隔，默认为1
  ListRange     TEXT                  DEFAULT '',  -- 列表范围规则
  ListRule      TEXT                  DEFAULT '',  -- 列表规则

  ContentUrl    VARCHAR(255) NOT NULL DEFAULT '',  -- 内容测试网址

  UpdateDate    DATETIME              DEFAULT CURRENT_TIMESTAMP,
  CreateDate    DATETIME              DEFAULT CURRENT_TIMESTAMP
);

-- 规则参数表
CREATE TABLE RuleParam
(
  Pid        INTEGER PRIMARY KEY AUTOINCREMENT, -- 参数ID
  Rid        INTEGER      NOT NULL DEFAULT '0', -- 规则ID
  Type       INTEGER      NOT NULL DEFAULT '0', -- 参数类型 1:列表 2:内容
  Field      VARCHAR(255) NOT NULL DEFAULT '',  -- 存放字段名称
  FieldType  INTEGER      NOT NULL DEFAULT '0', -- 字段类型 1:INTEGER 2:VARCHAR 3:TEXT
  Rule       TEXT                  DEFAULT '',  -- 匹配规则
  ValueType  INTEGER      NOT NULL DEFAULT '0', -- 获取值类型， 1:Html 2:Text 3:Attr
  ValueAttr  VARCHAR(255) NOT NULL DEFAULT '',  -- 当为 Attr 时，需要指定具体哪个属性
  FilterType INTEGER      NOT NULL DEFAULT '0', -- 过滤规则，1:清理两端空白 2:正则替换
  FilterReg  VARCHAR(255) NOT NULL DEFAULT '',  -- 过滤正则
  Sort       INTEGER      NOT NULL DEFAULT '0', -- 排序
  IsDown     INTEGER      NOT NULL DEFAULT '0', -- 此字段是否需要下载
  CreateDate DATETIME              DEFAULT CURRENT_TIMESTAMP
);
