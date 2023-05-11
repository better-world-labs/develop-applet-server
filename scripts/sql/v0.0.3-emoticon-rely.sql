ALTER TABLE `message_record`
    ADD COLUMN `i_msg_type` varchar(255) GENERATED ALWAYS AS (content->>'$.type') STORED COMMENT '虚拟字段；消息类型；' NULL AFTER `send_at`,
    ADD COLUMN `i_msg_ref` bigint UNSIGNED GENERATED ALWAYS AS (content->>'$.reference') STORED COMMENT '虚拟字段；消息引用；' NULL AFTER `i_msg_type`,
    ADD INDEX `idx_ref`(`i_msg_ref`) USING BTREE COMMENT '引用消息索引';


ALTER TABLE `emoticon`
    ADD COLUMN `group` int NULL DEFAULT 1 COMMENT '分组' AFTER `id`;


INSERT INTO `emoticon` (`group`, `name`, `url`, `keywords`, `ref_count`, `ref_stat_time`, `created_at`) VALUES (2, '点赞', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/static-web/images/perfect.gif', '点赞', 0, NULL, '2022-11-11 14:52:28');
INSERT INTO `emoticon` (`group`, `name`, `url`, `keywords`, `ref_count`, `ref_stat_time`, `created_at`) VALUES (2, '鼓掌', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/static-web/images/clap.gif\n', '鼓掌', 0, NULL, '2022-11-11 14:53:32');
INSERT INTO `emoticon` (`group`, `name`, `url`, `keywords`, `ref_count`, `ref_stat_time`, `created_at`) VALUES (2, '喜好', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/static-web/images/love.gif', '喜好', 0, NULL, '2022-11-11 14:54:07');
INSERT INTO `emoticon` (`group`, `name`, `url`, `keywords`, `ref_count`, `ref_stat_time`, `created_at`) VALUES (2, '送花', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/static-web/images/flower.gif', '送花', 1, '2022-11-12 01:00:00', '2022-11-11 14:54:35');
INSERT INTO `emoticon` (`group`, `name`, `url`, `keywords`, `ref_count`, `ref_stat_time`, `created_at`) VALUES (2, '+1', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/static-web/images/addone.gif', '+1', 0, NULL, '2022-11-11 14:54:56');