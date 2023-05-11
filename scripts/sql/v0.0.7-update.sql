alter table channel
    add mute tinyint default false not null comment '禁言' after type;

create table points
(
    id          bigint auto_increment,
    user_id     bigint                 not null,
    points      bigint                 not null,
    description varchar(255)           not null,
    created_at  datetime default now() not null,
    constraint points_pk
        primary key (id)
)
    comment '积分';

create index points__index_user_id
    on points (user_id);

create table message_like
(
    id         bigint auto_increment,
    message_id bigint  not null comment '消息ID',
    user_id    bigint  not null comment '用户ID',
    is_like    tinyint not null default false comment '是否点赞',
    created_at datetime         default now() not null,

    constraint points_pk
        primary key (id),

    constraint unique_message_id_user_id
        unique key (message_id, user_id)
)
    comment '消息点赞';

alter table channel
    add notice varchar(2048) null comment '公告';

create table system_config
(
    id         bigint auto_increment primary key,
    `key`      varchar(128) unique    not null comment '键',
    `value`    json                   not null comment '值',
    created_at datetime default now() not null

)
    comment '系统配置';

