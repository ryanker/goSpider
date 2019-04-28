-- 规则表
CREATE TABLE Rule
(
  Rid              INTEGER PRIMARY KEY AUTOINCREMENT,               -- 规则ID
  Status           INTEGER      NOT NULL DEFAULT '0',               -- 任务状态 1:关闭执行 2:执行一次 3:间隔执行
  IntervalHour     INTEGER      NOT NULL DEFAULT '0',               -- 间隔执行时间(小时)
  Name             VARCHAR(255) NOT NULL DEFAULT '',                -- 规则名称
  Brief            VARCHAR(255) NOT NULL DEFAULT '',                -- 规则备注

  Database         VARCHAR(255) NOT NULL DEFAULT '',                -- 数据库名称
  Cookie           TEXT                  DEFAULT '',                -- 登录抓取 Cookie
  Charset          VARCHAR(255) NOT NULL DEFAULT '',                -- 抓取网页的编码
  Timeout          INTEGER      NOT NULL DEFAULT '0',               -- 抓取超时时间(秒)

  ListSpecialUrl   TEXT                  DEFAULT '',                -- 抓取列表网址(特殊路径)
  ListUrl          VARCHAR(255) NOT NULL DEFAULT '',                -- 抓取列表网址
  ListPageStart    INTEGER      NOT NULL DEFAULT '0',               -- 列表开始页码
  ListPageEnd      INTEGER      NOT NULL DEFAULT '0',               -- 列表结束页码
  ListPageSize     INTEGER      NOT NULL DEFAULT '0',               -- 每页间隔，默认为1
  ListRule         TEXT                  DEFAULT '',                -- 列表规则

  ContentUrl       VARCHAR(255) NOT NULL DEFAULT '',                -- 内容测试网址

  IsList           INTEGER      NOT NULL DEFAULT '0',               -- 是否抓取列表页
  IsListDownAna    INTEGER      NOT NULL DEFAULT '0',               -- 是否分析列表页下载地址
  IsListDownRun    INTEGER      NOT NULL DEFAULT '0',               -- 是否下载列表页资源
  IsContent        INTEGER      NOT NULL DEFAULT '0',               -- 是否抓取内容页
  IsContentDownAna INTEGER      NOT NULL DEFAULT '0',               -- 是否分析内容页下载地址
  IsContentDownRun INTEGER      NOT NULL DEFAULT '0',               -- 是否抓取内容页资源

  SaveType         INTEGER      NOT NULL DEFAULT '0',               -- 存放文件类型 1:本机存放 2:OSS存放

  LastStartDate    DATETIME              DEFAULT CURRENT_TIMESTAMP, -- 上次采集开始时间
  LastEndDate      DATETIME              DEFAULT CURRENT_TIMESTAMP, -- 上次采集结束时间
  NextStartDate    DATETIME              DEFAULT CURRENT_TIMESTAMP, -- 下次采集开始时间

  UpdateDate       DATETIME              DEFAULT CURRENT_TIMESTAMP,
  CreateDate       DATETIME              DEFAULT CURRENT_TIMESTAMP
);

-- 规则参数表
CREATE TABLE RuleParam
(
  Pid           INTEGER PRIMARY KEY AUTOINCREMENT, -- 参数ID
  Rid           INTEGER      NOT NULL DEFAULT '0', -- 规则ID
  Type          VARCHAR(255) NOT NULL DEFAULT '',  -- 参数类型 值:List Content
  Field         VARCHAR(255) NOT NULL DEFAULT '',  -- 字段名称
  FieldType     VARCHAR(255) NOT NULL DEFAULT '',  -- 字段类型 值:INTEGER VARCHAR TEXT
  Brief         VARCHAR(255) NOT NULL DEFAULT '',  -- 字段备注
  Rule          TEXT                  DEFAULT '',  -- 匹配规则
  ValueType     VARCHAR(255) NOT NULL DEFAULT '',  -- 获取类型 值:Html Text Attr
  ValueAttr     VARCHAR(255) NOT NULL DEFAULT '',  -- 获取属性 (当为 Attr 时有效)
  FilterType    VARCHAR(255) NOT NULL DEFAULT '',  -- 过滤规则，值:Trim(清理两端空白) Reg(正则替换)
  FilterRegexp  VARCHAR(255) NOT NULL DEFAULT '',  -- 过滤正则
  FilterRepl    VARCHAR(255) NOT NULL DEFAULT '',  -- 正则结果
  Sort          INTEGER      NOT NULL DEFAULT '0', -- 排序
  IsSearch      INTEGER      NOT NULL DEFAULT '0', -- 是否参与搜索
  IsOrder       INTEGER      NOT NULL DEFAULT '0', -- 是否参与排序
  IsRequired    INTEGER      NOT NULL DEFAULT '0', -- 是否必须有值(采集入库时)
  DownType      INTEGER      NOT NULL DEFAULT '0', -- 下载类型 0:不用下载 1:直接下载 2:规则下载
  DownRule      TEXT                  DEFAULT '',  -- 下载地址匹配规则
  DownValueType VARCHAR(255) NOT NULL DEFAULT '',  -- 下载地址获取类型 值:Text Attr
  DownValueAttr VARCHAR(255) NOT NULL DEFAULT '',  -- 下载地址获取属性 (当为 Attr 时有效)
  DownFileType  VARCHAR(255) NOT NULL DEFAULT '',  -- 下载文件类型 值:Image Movie File
  DownTimeout   INTEGER      NOT NULL DEFAULT '0', -- 下载超时时间(秒)
  CreateDate    DATETIME              DEFAULT CURRENT_TIMESTAMP
);

-- 用户表
CREATE TABLE `User`
(
  `Uid`        INTEGER PRIMARY KEY AUTOINCREMENT,               -- 用户ID
  `Gid`        INTEGER      NOT NULL DEFAULT '0',               -- 用户组ID 1:管理员 2:普通用户
  `Name`       VARCHAR(255) NOT NULL DEFAULT '',                -- 用户名
  `Email`      VARCHAR(255) NOT NULL DEFAULT '',                -- 邮箱
  `Mobile`     VARCHAR(255) NOT NULL DEFAULT '',                -- 手机号
  `Password`   VARCHAR(32)  NOT NULL DEFAULT '',                -- 密码MD5
  `Salt`       VARCHAR(32)  NOT NULL DEFAULT '',                -- 密码盐
  `LoginNum`   INTEGER      NOT NULL DEFAULT '0',               -- 登录次数
  `LastIP`     VARCHAR(255) NOT NULL DEFAULT '',                -- 最后登陆IP
  `LastDate`   DATETIME              DEFAULT CURRENT_TIMESTAMP, -- 最后登陆时间
  `CreateIP`   VARCHAR(255) NOT NULL DEFAULT '',                -- 创建IP
  `CreateDate` DATETIME              DEFAULT CURRENT_TIMESTAMP, -- 创建时间
  `UpdateDate` DATETIME              DEFAULT CURRENT_TIMESTAMP  -- 更新时间
);
INSERT INTO `User` (`uid`, `gid`, `name`, `email`, `mobile`, `password`, `salt`)
VALUES (1, 1, 'dog', 'dog@dog.com', '18812345678', 'ea7ac79f95e032b0c73eec73ba6d1263', '12345678'); -- 密码: 88888888
