--- new
ALTER TABLE `user_setting` ADD COLUMN `site_settings` json NULL COMMENT '网页的一些设置，标题/icon'  AFTER `appearance_theme`;