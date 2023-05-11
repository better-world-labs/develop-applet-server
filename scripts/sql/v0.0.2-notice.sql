CREATE TABLE `notice`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`     bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `type`        varchar(32)         NOT NULL COMMENT '消息类型 (mention: 被艾特, reference: 被引用)',
    `business_id` bigint(20) unsigned NOT NULL COMMENT '相关业务ID',
    `created_at`  datetime            NOT NULL default now() COMMENT '创建时间',
    `is_read`        tinyint                      default 0 COMMENT '是否已读',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX index_notice_created_at ON `notice` (created_at desc);
CREATE INDEX index_notice_read ON `notice` (`read`); -- 未读消息一般占很少一部分，索引有必要
CREATE INDEX notice__index_user_id
    ON notice (user_id);
