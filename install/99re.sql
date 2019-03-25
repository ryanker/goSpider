CREATE TABLE imgList
(
  pid        INTEGER PRIMARY KEY AUTOINCREMENT,
  url        VARCHAR(255) NOT NULL DEFAULT '',
  title      VARCHAR(255) NOT NULL DEFAULT '',
  imgUrl     VARCHAR(255) NOT NULL DEFAULT '',
  imgUrlNew  VARCHAR(255) NOT NULL DEFAULT '',
  imgNum     INTEGER      NOT NULL DEFAULT '0',
  views      INTEGER      NOT NULL DEFAULT '0',
  date       VARCHAR(255) NOT NULL DEFAULT '',
  createDate DATETIME              DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE imgPost
(
  pid         INTEGER PRIMARY KEY AUTOINCREMENT,
  url         VARCHAR(255) NOT NULL DEFAULT '',
  title       VARCHAR(255) NOT NULL DEFAULT '',
  content     TEXT                  DEFAULT '',
  description VARCHAR(255) NOT NULL DEFAULT '',
  imgNum      INTEGER      NOT NULL DEFAULT '0',
  views       INTEGER      NOT NULL DEFAULT '0',
  date        VARCHAR(255) NOT NULL DEFAULT '',
  author      VARCHAR(255) NOT NULL DEFAULT '',
  authorHtml  TEXT                  DEFAULT '',
  cate        VARCHAR(255) NOT NULL DEFAULT '',
  cateHtml    TEXT                  DEFAULT '',
  tags        VARCHAR(255) NOT NULL DEFAULT '',
  tagsHtml    TEXT                  DEFAULT '',
  createDate  DATETIME              DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE imgPostData
(
  dataId     INTEGER PRIMARY KEY AUTOINCREMENT,
  pid        INTEGER      NOT NULL DEFAULT '0',
  imgUrl     VARCHAR(255) NOT NULL DEFAULT '',
  imgUrlNew  VARCHAR(255) NOT NULL DEFAULT '',
  createDate DATETIME              DEFAULT CURRENT_TIMESTAMP
);
