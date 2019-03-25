CREATE TABLE Com99reImgList
(
  Pid        INTEGER PRIMARY KEY AUTOINCREMENT,
  Url        VARCHAR(255) NOT NULL DEFAULT '',
  Title      VARCHAR(255) NOT NULL DEFAULT '',
  ImgUrl     VARCHAR(255) NOT NULL DEFAULT '',
  ImgUrlNew  VARCHAR(255) NOT NULL DEFAULT '',
  ImgNum     INTEGER      NOT NULL DEFAULT '0',
  Views      INTEGER      NOT NULL DEFAULT '0',
  Date       VARCHAR(255) NOT NULL DEFAULT '',
  CreateDate DATETIME              DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Com99reImgPost
(
  Pid         INTEGER PRIMARY KEY AUTOINCREMENT,
  Url         VARCHAR(255) NOT NULL DEFAULT '',
  Title       VARCHAR(255) NOT NULL DEFAULT '',
  Content     TEXT                  DEFAULT '',
  Description VARCHAR(255) NOT NULL DEFAULT '',
  ImgNum      INTEGER      NOT NULL DEFAULT '0',
  Views       INTEGER      NOT NULL DEFAULT '0',
  Date        VARCHAR(255) NOT NULL DEFAULT '',
  Author      VARCHAR(255) NOT NULL DEFAULT '',
  AuthorHtml  TEXT                  DEFAULT '',
  Cate        VARCHAR(255) NOT NULL DEFAULT '',
  CateHtml    TEXT                  DEFAULT '',
  Tags        VARCHAR(255) NOT NULL DEFAULT '',
  TagsHtml    TEXT                  DEFAULT '',
  CreateDate  DATETIME              DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Com99reImgPostData
(
  dataId     INTEGER PRIMARY KEY AUTOINCREMENT,
  Pid        INTEGER      NOT NULL DEFAULT '0',
  ImgUrl     VARCHAR(255) NOT NULL DEFAULT '',
  ImgUrlNew  VARCHAR(255) NOT NULL DEFAULT '',
  CreateDate DATETIME              DEFAULT CURRENT_TIMESTAMP
);
