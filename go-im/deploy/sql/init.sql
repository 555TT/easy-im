-- easy_chat 后端数据库初始化脚本
-- 适用服务：user-rpc、social-rpc
-- 字符集：utf8mb4 / utf8mb4_unicode_ci
-- 主键策略：业务侧使用 Snowflake ID（pkg/suid），全部 VARCHAR 主键，不使用 AUTO_INCREMENT

CREATE DATABASE IF NOT EXISTS `easy_chat`
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE `easy_chat`;

-- ============================================================
-- user-rpc
-- ============================================================

-- 用户表
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `id`         VARCHAR(24)  NOT NULL,
    `avatar`     VARCHAR(255) NOT NULL DEFAULT 'https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png',
    `nickname`   VARCHAR(24)  NOT NULL,
    `phone`      VARCHAR(24)  NOT NULL,
    `email`      VARCHAR(24)           DEFAULT NULL,
    `password`   VARCHAR(191)          DEFAULT NULL,
    `status`     TINYINT(1)            DEFAULT 0,
    `sex`        TINYINT(1)            DEFAULT 0,
    `created_at` TIMESTAMP    NULL     DEFAULT NULL,
    `updated_at` TIMESTAMP    NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_users_phone` (`phone`),
    KEY `idx_users_nickname` (`nickname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- social-rpc
-- ============================================================

-- 好友关系表
DROP TABLE IF EXISTS `friends`;
CREATE TABLE `friends` (
    `id`         VARCHAR(64) NOT NULL,
    `user_id`    VARCHAR(64) NOT NULL,
    `friend_uid` VARCHAR(64) NOT NULL,
    `remark`     VARCHAR(255)         DEFAULT NULL,
    `add_source` TINYINT              DEFAULT NULL,
    `created_at` TIMESTAMP   NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_friends_user_friend` (`user_id`, `friend_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 好友申请表
DROP TABLE IF EXISTS `friend_requests`;
CREATE TABLE `friend_requests` (
    `id`            VARCHAR(64)  NOT NULL,
    `user_id`       VARCHAR(64)  NOT NULL,
    `req_uid`       VARCHAR(64)  NOT NULL,
    `req_msg`       VARCHAR(255)          DEFAULT NULL,
    `req_time`      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `handle_result` TINYINT               DEFAULT NULL,
    `handle_msg`    VARCHAR(255)          DEFAULT NULL,
    `handled_at`    TIMESTAMP    NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_friend_requests_user` (`user_id`),
    KEY `idx_friend_requests_req` (`req_uid`, `user_id`, `handle_result`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 群组表
DROP TABLE IF EXISTS `groups`;
CREATE TABLE `groups` (
    `id`               VARCHAR(24)  NOT NULL,
    `name`             VARCHAR(255) NOT NULL,
    `icon`             VARCHAR(255) NOT NULL DEFAULT 'https://c-ssl.duitang.com/uploads/item/201802/24/20180224083913_yhrX2.jpeg',
    `status`           TINYINT               DEFAULT NULL,
    `creator_uid`      VARCHAR(64)  NOT NULL,
    `group_type`       TINYINT      NOT NULL,
    `is_verify`        TINYINT(1)   NOT NULL DEFAULT 0,
    `notification`     VARCHAR(255)          DEFAULT NULL,
    `notification_uid` VARCHAR(64)           DEFAULT NULL,
    `created_at`       TIMESTAMP    NULL     DEFAULT NULL,
    `updated_at`       TIMESTAMP    NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_groups_creator` (`creator_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 群成员表
DROP TABLE IF EXISTS `group_members`;
CREATE TABLE `group_members` (
    `id`           VARCHAR(64) NOT NULL,
    `group_id`     VARCHAR(64) NOT NULL,
    `user_id`      VARCHAR(64) NOT NULL,
    `role_level`   TINYINT     NOT NULL,
    `join_time`    TIMESTAMP   NULL     DEFAULT NULL,
    `join_source`  TINYINT              DEFAULT 0,
    `inviter_uid`  VARCHAR(64)          DEFAULT NULL,
    `operator_uid` VARCHAR(64)          DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_group_members_group_user` (`group_id`, `user_id`),
    KEY `idx_group_members_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 入群申请表
DROP TABLE IF EXISTS `group_requests`;
CREATE TABLE `group_requests` (
    `id`              VARCHAR(64)  NOT NULL,
    `req_id`          VARCHAR(64)  NOT NULL,
    `group_id`        VARCHAR(64)  NOT NULL,
    `req_msg`         VARCHAR(255)          DEFAULT NULL,
    `req_time`        TIMESTAMP    NULL     DEFAULT NULL,
    `join_source`     TINYINT               DEFAULT NULL,
    `inviter_user_id` VARCHAR(64)           DEFAULT NULL,
    `handle_user_id`  VARCHAR(64)           DEFAULT NULL,
    `handle_time`     TIMESTAMP    NULL     DEFAULT NULL,
    `handle_result`   TINYINT               DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_group_requests_group_result` (`group_id`, `handle_result`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
