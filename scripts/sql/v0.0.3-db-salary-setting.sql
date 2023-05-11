ALTER TABLE `user_setting`
ADD COLUMN `monthly_salary` bigint UNSIGNED NOT NULL DEFAULT 10000 COMMENT '月薪，单位元' AFTER `boss_key`,
ADD COLUMN `monthly_working_days` bigint UNSIGNED NOT NULL DEFAULT 22 COMMENT '月工作时长，单位天，默认22天' AFTER `monthly_salary`;