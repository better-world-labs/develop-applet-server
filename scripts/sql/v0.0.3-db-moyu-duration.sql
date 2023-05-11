ALTER TABLE `user`
ADD COLUMN `connect_time` datetime NULL COMMENT 'ws最新一次连接时间' AFTER `online_changed_at`,
ADD COLUMN `total_access_duration` bigint NOT NULL DEFAULT 0 COMMENT '用户在摸鱼星球挂机时长，单位为秒' AFTER `connect_time`,
ADD COLUMN `total_browse_duration` bigint NOT NULL DEFAULT 0 COMMENT '用户访问摸鱼星球网页时长，单位为秒，前端按时上报' AFTER `total_access_duration`;

CREATE TABLE `user_access_record` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` bigint(20) NOT NULL COMMENT '用户主键id',
  `start_time` datetime NOT NULL COMMENT '摸鱼开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '摸鱼结束时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_browse_record` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `browse_date` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '摸鱼浏览业务日期，yyyy-MM-dd',
  `browse_duration` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '摸鱼时长，单位s',
  `last_browse_time` datetime DEFAULT NULL COMMENT '上次摸鱼结束时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
