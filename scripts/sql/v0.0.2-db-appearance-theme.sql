ALTER TABLE `user_setting` ADD COLUMN `appearance_theme` VARCHAR ( 32 ) NOT NULL DEFAULT 'dark' COMMENT '外观主题，默认是dark' AFTER `end_off_time`;

