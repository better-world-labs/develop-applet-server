create table alert
(
    id          bigint unsigned auto_increment comment 'ID'
        primary key,
    title       varchar(128) default ''                not null comment '标题',
    head_image  varchar(128)                           not null comment '头图',
    content     varchar(1024)                          not null comment '富文本内容',
    created_at  datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    target_time datetime     default CURRENT_TIMESTAMP not null comment '生效时间'
);
