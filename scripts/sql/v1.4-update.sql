alter table mini_app
    add top int default 0 not null;

create table mini_app_recommend
(
    id         bigint auto_increment
        primary key,
    app_id     varchar(64) not null,
    recommend     tinyint     not null comment '是否推荐',
    created_by bigint      not null,
    updated_at bigint      not null,
    constraint mini_app_recommend_unique
        unique (app_id, created_by)
);

create index idx_mini_app_recommend_created_by
    on mini_app_like (created_by);

alter table statistic_mini_app
    modify comment_times int default 0 not null comment '评论次数' after comment_times_updated_at;

alter table statistic_mini_app
    add recommend_times int default 0 not null after view_times_updated_at;

alter table statistic_mini_app
    modify degree_of_heat_updated_at datetime default CURRENT_TIMESTAMP not null after recommend_times;

alter table statistic_mini_app
    add recommend_times_updated_at datetime default now() not null;

