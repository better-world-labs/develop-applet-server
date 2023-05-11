-- auto-generated definition
create table notify_message
(
    id           bigint unsigned auto_increment comment 'ID'
        primary key,
    user_id      bigint unsigned                    not null comment '用户ID',
    type         varchar(32)                        not null comment '通知类型',
    title        varchar(127)                       not null comment '通知标题',
    content      varchar(1024)                      not null comment '通知正文',
    created_at   datetime default CURRENT_TIMESTAMP not null comment '创建时间',
    is_read      tinyint  default 0                 null comment '是否已读',
    operation_id varchar(64)                        not null unique comment '操作ID,全局唯一'
);

create index notify_message__index_user_id
    on notice (user_id);


-- auto-generated definition
create table points_recharge_goods_definition
(
    id          bigint auto_increment
        primary key,
    price       int          not null comment '价格 (分)',
    points      bigint       not null comment '积分',
    tag         varchar(16)  not null comment '标签',
    description varchar(255) not null comment '描述',
    coupon      int          null comment '优惠券'
)
    comment '积分商品';

-- auto-generated definition
create table points_recharge_order
(
    id                     bigint auto_increment
        primary key,
    order_id               varchar(32)  not null unique comment '订单编号',
    goods_id               bigint       not null comment '商品ID',
    goods                  json         not null comment '商品快照',
    price                  int          not null comment '价格 (分)',
    user_id                bigint       not null comment '用户ID',
    created_at             datetime     not null default now() comment '下单时间',
    state                  int                   default 0 not null comment '订单状态 0. 未支付， 1.已支付, 2.已关闭',
    payed_at               varchar(32) comment '支付时间',
    pay_expires_at         datetime     null,
    pay_transaction_id     varchar(64)  null,
    pay_trade_type         varchar(16)  null comment '交易类型',
    pay_trade_state        varchar(32)  null comment '交易状态',
    pay_bank_type          varchar(32)  null comment '银行类型',
    pay_openid             varchar(128) null,
    pay_amount_total       int          null comment '支付订单总金额(分)',
    pay_amount_payer_total int          null comment '用户支付总金额(分)'
)
    comment '积分订单';
-- auto-generated definition
create table sign_in_daily
(
    id         bigint unsigned auto_increment comment 'ID'
        primary key,
    user_id    bigint                             not null comment '用户ID',
    date       date                               not null comment '日期',
    created_at datetime default CURRENT_TIMESTAMP not null comment '创建时间'
);

create unique index sign_in_daily_user_id_date_uindex
    on sign_in_daily (user_id, date);

alter table user
    add from_app varchar(24) default '' not null;

alter table user
    add source varchar(255) default '' not null;

-- auto-generated definition
create table mini_app_collection
(
    id         bigint auto_increment
        primary key,
    app_id     varchar(64)            not null,
    user_id    bigint                 not null,
    created_at datetime default now() not null,
    constraint mini_app_like_unique
        unique (app_id, user_id)
);

create index mini_app_collection_app_id_index
    on mini_app_collection (app_id);

alter table user
    add login_at datetime default now() null;

alter table user
    add last_login_at datetime default now() null;

alter table points
    add type varchar(32) null;

create table retain_message
(
    id      bigint unsigned not null primary key auto_increment comment 'ID',
    user_id bigint unsigned not null comment '用户ID',
    type    varchar(32)     not null comment '消息类型',
    payload json            not null comment '消息内容'
);
create index idx_retain_message_user_id on retain_message (user_id);

create table retain_message_offset
(
    id        bigint unsigned not null primary key auto_increment comment 'ID',
    user_id   bigint unsigned not null unique comment '用户ID',
    offset_id bigint          not null comment '消息消费偏移量 (id)'
);

# mini_app 的 price 字段废弃