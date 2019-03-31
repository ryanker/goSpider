-- 第一步：分类（采集哪个网站）
-- 第二步：规则名称
-- 第三步：抓取规则细节参数

-- 分类表
CREATE TABLE RuleCate
(
  CateId     INTEGER PRIMARY KEY AUTOINCREMENT, -- 分类ID
  Name       VARCHAR(255) NOT NULL DEFAULT '',  -- 分类名称
  Brief      VARCHAR(255) NOT NULL DEFAULT '',  -- 分类备注
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
  Brief         VARCHAR(255) NOT NULL DEFAULT '',  -- 规则备注

  ListTable     VARCHAR(255) NOT NULL DEFAULT '',  -- 表名称
  ListUrl       VARCHAR(255) NOT NULL DEFAULT '',  -- 抓取列表网址
  ListPageStart INTEGER      NOT NULL DEFAULT '0', -- 列表开始页码
  ListPageEnd   INTEGER      NOT NULL DEFAULT '0', -- 列表结束页码
  ListPageSize  INTEGER      NOT NULL DEFAULT '0', -- 每页间隔，默认为1
  ListRule      TEXT                  DEFAULT '',  -- 列表规则

  ContentUrl    VARCHAR(255) NOT NULL DEFAULT '',  -- 内容测试网址

  UpdateDate    DATETIME              DEFAULT CURRENT_TIMESTAMP,
  CreateDate    DATETIME              DEFAULT CURRENT_TIMESTAMP
);

-- 规则参数表
CREATE TABLE RuleParam
(
  Pid           INTEGER PRIMARY KEY AUTOINCREMENT, -- 参数ID
  Rid           INTEGER      NOT NULL DEFAULT '0', -- 规则ID
  CateId        INTEGER      NOT NULL DEFAULT '0', -- 分类ID
  Type          VARCHAR(255) NOT NULL DEFAULT '',  -- 参数类型 值:List Content
  Field         VARCHAR(255) NOT NULL DEFAULT '',  -- 存放字段名称
  FieldType     VARCHAR(255) NOT NULL DEFAULT '',  -- 字段类型 值:INTEGER VARCHAR TEXT
  Brief         VARCHAR(255) NOT NULL DEFAULT '',  -- 存放字段备注
  Rule          TEXT                  DEFAULT '',  -- 匹配规则
  ValueType     VARCHAR(255) NOT NULL DEFAULT '',  -- 获取值类型， 值:Html Text Attr
  ValueAttr     VARCHAR(255) NOT NULL DEFAULT '',  -- 当为 Attr 时，需要指定具体哪个属性
  FilterType    VARCHAR(255) NOT NULL DEFAULT '',  -- 过滤规则，值:Trim(清理两端空白) Reg(正则替换)
  FilterRegexp  VARCHAR(255) NOT NULL DEFAULT '',  -- 过滤正则表达式
  FilterRepl    VARCHAR(255) NOT NULL DEFAULT '',  -- 处理结果
  Sort          INTEGER      NOT NULL DEFAULT '0', -- 排序
  IsSearch      INTEGER      NOT NULL DEFAULT '0', -- 是否参与搜索
  DownType      INTEGER      NOT NULL DEFAULT '0', -- 下载类型 0:不用下载 1:直接下载 2:规则匹配下载
  DownRule      TEXT                  DEFAULT '',  -- 匹配规则
  DownValueType VARCHAR(255) NOT NULL DEFAULT '',  -- 获取值类型， 值:Text Attr
  DownValueAttr VARCHAR(255) NOT NULL DEFAULT '',  -- 当为 Attr 时，需要指定具体哪个属性
  CreateDate    DATETIME              DEFAULT CURRENT_TIMESTAMP
);
