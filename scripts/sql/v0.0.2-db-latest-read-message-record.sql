--- 增加字段
ALTER TABLE `channel_member` ADD COLUMN `last_read_message_id` bigint NOT NULL DEFAULT 0 COMMENT '用户上次在该频道读取的最后一条消息' AFTER `user_id`;

--- 更新最新未读消息id
UPDATE channel_member a
    LEFT JOIN ( SELECT channel_id, max( id ) AS max_id FROM message_record GROUP BY channel_id ) b ON a.channel_id = b.channel_id
    SET a.last_read_message_id = b.max_id
WHERE
    b.max_id IS NOT NULL;



-- 验证，最新未读id是否变成频道下最新未读id
SELECT
    *
FROM
    channel_member a
        LEFT JOIN ( SELECT channel_id, max( id ) as max_id FROM message_record GROUP BY channel_id ) b ON a.channel_id = b.channel_id
WHERE
        a.last_read_message_id = b.max_id or (a.last_read_message_id = 0 and b.max_id is null);