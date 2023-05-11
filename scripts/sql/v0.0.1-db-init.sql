use moyu_prod;
SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for anonymous_identity
-- ----------------------------
DROP TABLE IF EXISTS `anonymous_identity`;
CREATE TABLE `anonymous_identity`
(
    `id`       bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar`   varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for channel
-- ----------------------------
DROP TABLE IF EXISTS `channel`;
CREATE TABLE `channel`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name`       varchar(90)  NOT NULL COMMENT '频道名称',
    `icon`       varchar(128) NOT NULL DEFAULT '' COMMENT '频道图标',
    `planet_id`  bigint(20) unsigned DEFAULT NULL COMMENT '星球ID',
    `group_id`   bigint(20) unsigned DEFAULT NULL COMMENT '频道所属组ID',
    `created_by` bigint(20) unsigned DEFAULT NULL COMMENT '创建人ID',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `sort`       int(10) unsigned DEFAULT NULL COMMENT '排序',
    PRIMARY KEY (`id`),
    KEY          `idx_planet_id` (`planet_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for channel_group
-- ----------------------------
DROP TABLE IF EXISTS `channel_group`;
CREATE TABLE `channel_group`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name`       varchar(90)  NOT NULL COMMENT '组名称',
    `icon`       varchar(128) NOT NULL DEFAULT '' COMMENT '组图标',
    `planet_id`  bigint(20) unsigned DEFAULT NULL COMMENT '星球ID',
    `created_by` bigint(20) unsigned DEFAULT NULL COMMENT '创建人ID',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `sort`       int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY          `idx_planet_id` (`planet_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for channel_member
-- ----------------------------
DROP TABLE IF EXISTS `channel_member`;
CREATE TABLE `channel_member`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `channel_id` bigint(20) unsigned NOT NULL COMMENT '频道ID',
    `user_id`    bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `channel_id` (`channel_id`,`user_id`),
    KEY          `idx_channel_id` (`channel_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
CREATE INDEX index_channel_member_user_id ON `channel_member`(user_id)

-- ----------------------------
-- Table structure for hot_issue
-- ----------------------------
DROP TABLE IF EXISTS `hot_issue`;
CREATE TABLE `hot_issue`
(
    `id`      bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `title`   varchar(128) NOT NULL COMMENT '标题',
    `content` varchar(255) NOT NULL COMMENT '内容',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for message_record
-- ----------------------------
DROP TABLE IF EXISTS `message_record`;
CREATE TABLE `message_record`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `channel_id` bigint(20) unsigned NOT NULL COMMENT '频道ID',
    `user_id`    bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `content`    varchar(1024) NOT NULL DEFAULT '' COMMENT '消息内容',
    `created_at` datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `seq_id`     bigint(20) NOT NULL COMMENT '消息ID',
    `send_id`    varchar(64)   NOT NULL,
    `send_at`    datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for planet
-- ----------------------------
DROP TABLE IF EXISTS `planet`;
CREATE TABLE `planet`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name`        varchar(128) NOT NULL COMMENT '星球名称',
    `icon`        varchar(128) NOT NULL DEFAULT '' COMMENT '星球图标',
    `created_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `front_cover` varchar(256) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for planet_member
-- ----------------------------
DROP TABLE IF EXISTS `planet_member`;
CREATE TABLE `planet_member`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `planet_id`  bigint(20) unsigned NOT NULL COMMENT '星球ID',
    `user_id`    bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `role`       int(10) unsigned NOT NULL COMMENT '角色',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `status`     int(11) DEFAULT '0' COMMENT '状态 0:正常, 1: 黑名单用户',
    PRIMARY KEY (`id`),
    UNIQUE KEY `planet_member_unique` (`planet_id`,`user_id`),
    KEY          `idx_planet_id` (`planet_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for resign_template
-- ----------------------------
DROP TABLE IF EXISTS `resign_template`;
CREATE TABLE `resign_template`
(
    `id`      bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `title`   varchar(255) NOT NULL COMMENT '模板标题',
    `content` text         NOT NULL COMMENT '模版内容',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`                bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `nickname`          varchar(128) NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar`            varchar(128) NOT NULL DEFAULT '' COMMENT '头像URL',
    `wx_open_id`        varchar(64)  NOT NULL COMMENT '微信OpenId',
    `wx_union_id`       varchar(64)  NOT NULL DEFAULT '' COMMENT '微信UnionId',
    `status`            int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户状态: 0.正常 1.黑名单',
    `online`            tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否在线',
    `created_at`        datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `invited_by`        bigint(20) DEFAULT NULL COMMENT '邀请用户ID',
    `online_changed_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_wx_open_id_uindex` (`wx_open_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for user_setting
-- ----------------------------
DROP TABLE IF EXISTS `user_setting`;
CREATE TABLE `user_setting`
(
    `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`      bigint(20) unsigned NOT NULL COMMENT 'ID',
    `end_off_time` varchar(32) NOT NULL COMMENT '上班时间',
    `boss_key`     varchar(32) NOT NULL DEFAULT '' COMMENT '老板键',
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_setting_user_id_uindex` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET
FOREIGN_KEY_CHECKS = 1;

INSERT INTO `planet` (`id`, `name`, `icon`, `created_at`, `front_cover`)
VALUES (1, '摸鱼猩球', 'https://openview-oss.oss-cn-chengdu.aliyuncs.com/aed-test/avatar/93.png', '2022-09-22 11:11:27',
        'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/moyu-dev/0/2022-10-11/1665460397918246');


INSERT INTO `hot_issue` (`id`, `title`, `content`)
VALUES (7, '摸鱼八卦公会', '22人同时在线，90000条热门讨论');
INSERT INTO `hot_issue` (`id`, `title`, `content`)
VALUES (8, '闺蜜帮帮公会', '20次情感问题讨论和分析');
INSERT INTO `hot_issue` (`id`, `title`, `content`)
VALUES (9, '中午吃什么', '疯狂星期四疯狂推荐');
INSERT INTO `hot_issue` (`id`, `title`, `content`)
VALUES (10, '中午吃什么', '猪脚饭 汉堡薯条 麻辣烫 黄焖鸡米饭  煲仔饭 冒菜 砂锅粥');
INSERT INTO `hot_issue` (`id`, `title`, `content`)
VALUES (11, '下班倒计时', '距离下班仅剩 07:00:00');
INSERT INTO `hot_issue` (`id`, `title`, `content`)
VALUES (12, '辞职模板', '尊敬的老板，很遗憾地通知你，你被fire了');
INSERT INTO `hot_issue` (`id`, `title`, `content`)
VALUES (13, '随机请假理由', '身份证快过期了，请假去办理');
INSERT INTO `hot_issue` (`id`, `title`, `content`)
VALUES (14, '随机请假理由', '老板，我分手了');

INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (11, '%', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/1.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (12, 'AAAA野猩CFO', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/2.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (13, 'AAAA野猩CTO', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/3.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (14, 'AAAA野猩家', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/4.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (15, 'AAAA野猩猴王', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/5.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (16, 'AAAA猩球懂事长', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/6.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (17, 'AAAA野猩宇航员', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/7.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (18, 'AAAA野猩总锦鲤', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/8.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (19, 'AAAA野猩星星眼', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/9.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (20, 'AAAA野猩音乐家', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/10.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (21, 'AAAA野猩西部牛仔', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/11.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (22, 'Beryl', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/12.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (23, 'Big小新新新', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/13.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (24, 'Fe', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/14.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (25, 'ICE', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/15.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (26, 'LYnn', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/16.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (27, 'OR', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/17.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (28, 'Yuyu', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/18.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (29, 'format', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/19.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (30, 'ulion', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/20.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (31, 'zzz', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/21.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (32, '曾经', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/22.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (33, '朝阳', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/23.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (34, '比心', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/24.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (35, '泥泥', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/25.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (36, '西药', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/26.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (37, '霜降', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/27.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (38, '月初', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/28.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (39, '唤小醒', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/29.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (40, '大鲨鱼', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/30.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (41, '木头钟', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/31.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (42, '东土大堂', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/32.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (43, '别回头哦', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/33.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (44, '北冥有鱼', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/34.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (45, '《风云》', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/35.png');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (46, '吃瓜群的众', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/36.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (47, '飞机飞过天空', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/37.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (48, '渔业研究院前台', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/38.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (49, '绝情谷深情谷主', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/39.jpg');
INSERT INTO `anonymous_identity` (`id`, `nickname`, `avatar`)
VALUES (50, '国家级摸鱼研究员', 'https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/40.jpg');


INSERT INTO `resign_template` (`id`, `title`, `content`)
VALUES (1, '世界那么大，我想去看看',
        '尊敬的[经理姓名]，\n本人正式递交辞呈，通知您我已辞去[公司名称]的[职务]。我在公司的最后一天是[日期]。\n在[公司名称]任职期间，我逐渐意识到，这个职位的工作范围和我的预期存在较大距离。因此，我想探索其他工作机会。\n感谢您的理解。我也非常感谢您在我任职期间给予的支持和指导。\n对于在[通知期]的交接，如您需要我的任何帮助和支持，请不吝告知。\n诚挚地\n[你的名字]');
INSERT INTO `resign_template` (`id`, `title`, `content`)
VALUES (2, '正经辞职模板 除非你看到最后',
        '尊敬的XX领导： \n本人正式向您递交此份辞呈，以通知我辞去[公司名称]的[职务]的打算。按照我的通知期，我的最后一天将是[最后一天的日期]。\n我衷心感谢您给予这个机会，让我在过去的[工作时间]里在这个职位上效力。在这里我得到了很多经验，与同事们的合作亦很愉悦。在这里学到的很多东西对我整个职业生涯都会有所裨益，回顾以往，在这里渡过的时光是我一生宝贵的记忆。\n在接下来的（几周后的通知期内），我将尽我所能做好交接，将我的职责移交给同事们或继任者。在交接过程中，如需进一步帮助，请不吝告知。\n对不起，我家拆迁了。赔得太多了。告辞！');
INSERT INTO `resign_template` (`id`, `title`, `content`)
VALUES (3, '来上班？上不赚钱的班？    先去炒老板吧！',
        '尊敬的[经理姓名]，\n本人正式递交辞呈，通知您我将辞去[公司名称]的[职务]。按照我的通知期，我的最后一天将是[最后一天的日期]。\n我获得了另一个职位，新工作不仅能使我每天的通勤时间减半，也能让我在工作之余有更多的时间陪伴家人。\n在过去的[服务的岁月和月份]中，我在[公司]度过了非常愉快的时光，在此我对您表示衷心的感谢。\n在[截止日期]前的未来几周，我将尽我所能，确保交接工作平稳完成。\n诚挚地\n[你的名字]');
INSERT INTO `resign_template` (`id`, `title`, `content`)
VALUES (4, '工作是双向选择，你可以随时离开，只要想清楚了未来的路',
        '尊敬的公司领导：\n首先感谢公司近段时间对我的信任和关照，给予了我一个发展的平台，使我有了长足的进步。如今由于个人原因，无法为公司继续服务，现在我正式向公司提出辞职申请，将于20xx年xx月xx日离职，请公司做好相应的安排，在此期间我一定站好最后一班岗，做好交接工作。对此为公司带来的不便，我深感歉意。\n望公司批准，谢谢!\n祝公司业绩蒸蒸日上。\n此致\n敬礼!\n辞职人：xxx\n20xx年xx月xx日');
INSERT INTO `resign_template` (`id`, `title`, `content`)
VALUES (5, '不如好聚好散吧，挥挥手，不带走一片云彩',
        'xx：\n自xx年入职以来，我一直很喜欢这份工作，但因某些个人原因，我要重新确定自己未来的方向，最终选择了开始新的工作。\n希望公司能早日找到合适人手接替我的工作并希望能于今年5月底前正式辞职。如能给予我支配更多的时间来找工作我将感激不尽，希望公司理解!在我提交这份辞呈时，在未离开岗位之前，我一定会尽自己的职责，做好应该做的事。\n最后，衷心的说：“对不起”与“谢谢”!\n祝愿公司开创更美好的未来!\n望领导批准我的申请!并协助办理相关离职手续。\n此致\n敬礼!\n辞职人：xxx\n20xx年xx月xx日');