CREATE TABLE IF NOT EXISTS `user` (
    id bigint AUTO_INCREMENT,
    name varchar(255) NOT NULL DEFAULT '' COMMENT 'The username',
    type tinyint(1) NULL DEFAULT 0 COMMENT 'The user type, 0:normal,1:vip, for test golang keyword',
    create_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE = InnoDB COLLATE utf8mb4_general_ci COMMENT 'user table';

--  默认用户在建表时会创建createTime、updateTime字段(忽略大小写、下划线命名风格)且默认值均为CURRENT_TIMESTAMP，而updateTime支持ON UPDATE CURRENT_TIMESTAMP