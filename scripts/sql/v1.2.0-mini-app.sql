create table mini_app_comment
(
    id         bigint auto_increment,
    app_id     varchar(64)              not null,
    content    varchar(2048) default '' not null comment '评论内容',
    created_by bigint                   null,
    created_at datetime                 null,
    constraint mini_app_comment_pk
        primary key (id)
);

create index idx_mini_app_comment_app_id on mini_app_comment (app_id);
create index idx_mini_app_comment_created_by on mini_app_comment (created_by);

create table mini_app_like

(
    id         bigint primary key auto_increment,
    app_id     varchar(64) not null,
    `like`     tinyint     not null comment '评论内容',
    created_by bigint      null,
    updated_at bigint      not null,
    unique (app_id, created_by)
);

create index idx_mini_app_like_app_id on mini_app_comment (app_id);
create index idx_mini_app_like_created_by on mini_app_like (created_by);

alter table statistic_mini_app
    add view_times int default 0 not null;

alter table statistic_mini_app
    add view_times_updated_at datetime default now() null;

alter table statistic_mini_app
    add degree_of_heat float default 0 not null;

alter table statistic_mini_app
    add degree_of_heat_updated_at datetime default now() null;


-- auto-generated definition
create table mini_app_output_like
(
    id         bigint auto_increment
        primary key,
    output_id  varchar(64) unique not null,
    `like`     tinyint            not null comment '1. 顶, 0, 取消点赞, -1. 踩',
    created_by bigint             null,
    updated_at bigint             not null,
    constraint output_id
        unique (output_id, created_by)
);

create index idx_mini_app_output_like_created_by
    on mini_app_output_like (created_by);


-- auto-generated definition
create table statistic_mini_app_output
(
    id                       bigint auto_increment
        primary key,
    output_id                varchar(64) unique                       not null comment '应用输出 ID',
    like_times               int         default 0                    not null comment '点赞次数',
    like_times_updated_at    datetime(3) default CURRENT_TIMESTAMP(3) not null comment '点赞次数修改时间',
    hate_times               int         default 0                    not null comment '踩次数',
    hate_times_updated_at    datetime(3) default CURRENT_TIMESTAMP(3) not null comment '踩次数修改时间',
    comment_times            int         default 0                    not null comment '评论次数',
    comment_times_updated_at datetime(3) default CURRENT_TIMESTAMP(3) not null comment '评论次数修改时间',
    constraint app_id
        unique (output_id)
);

-- auto-generated definition
create table mini_app_user_extra
(
    id         bigint auto_increment
        primary key,
    user_id    bigint unique                            not null comment '用户ID',
    complete_guidance  tinyint     default 0                    not null comment '是否完成新手引导',
    created_at datetime(3) default CURRENT_TIMESTAMP(3) not null comment '创建时间'
);

