alter table channel_member
    add state int default 2 not null comment '0. 未加入, 1. 申请中, 2. 已加入, 3. 被移除';

alter table channel_member
    add apply_id bigint default 0 not null comment '申请单 ID';

alter table channel
    add expires_at datetime null comment '过期时间，私密频道具备
';

alter table channel
    add type int not null default 1 comment '1. 普通频道, 2. 私密频道
';

create table approval
(
    id            bigint unsigned auto_increment comment 'ID'
        primary key,
    business_id   bigint unsigned                        not null comment '业务ID',
    approval_type varchar(32)                            not null comment '审批类型 (channel-join: 频道加入审批)',
    user_id       bigint unsigned                        not null comment '用户ID',
    created_at    datetime     default CURRENT_TIMESTAMP not null comment '创建时间',
    state         int                                    not null comment '0. 待审核, 1. 通过, 2. 驳回',
    reason        varchar(255) default ''                not null comment '申请理由',
    audit_by      bigint unsigned                        null comment '审批人',
    audit_at      datetime                               null comment '审批时间'
);

create index approval__index_user_id
    on approval (user_id);

create index notice__index_type_business_id
    on notice (business_id, type);

