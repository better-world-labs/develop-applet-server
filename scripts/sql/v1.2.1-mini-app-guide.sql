create table mini_app_user_extra
(
    id         bigint auto_increment
        primary key,
    user_id    bigint unique                            not null comment '用户ID',
    complete_guidance  tinyint     default 0                    not null comment '是否完成新手引导',
    created_at datetime(3) default CURRENT_TIMESTAMP(3) not null comment '创建时间'
);

