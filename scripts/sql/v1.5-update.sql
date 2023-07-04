create table ai_tool
(
    id          bigint auto_increment
        primary key,
    name        varchar(64)   not null comment '名称',
    description varchar(1024) not null default '' comment '描述',
    icon        varchar(512)  not null default '' comment '图标',
    category    int           not null,
    target      varchar(512)  not null
);

create table ai_tool_category
(
    id          bigint auto_increment
        primary key,
    category    int           not null comment '分类',
    description varchar(1024) not null default '' comment '描述',
    sort        int           not null default 0
);

INSERT INTO ai_tool_category (id, category, description, sort)
VALUES (1, 1, '提示词类', 1);
INSERT INTO ai_tool_category (id, category, description, sort)
VALUES (2, 2, '写作类', 2);
INSERT INTO ai_tool_category (id, category, description, sort)
VALUES (3, 3, '互动类', 3);
INSERT INTO ai_tool_category (id, category, description, sort)
VALUES (4, 4, '绘画类', 4);

create table gpt_chat_message
(
    id         bigint auto_increment
        primary key,
    role       varchar(20)   not null comment '角色: user: 用户, assistant: gpt',
    message_id       varchar(64)   not null unique comment '唯一标识',
    user_id     bigint        not null comment '用户ID',
    content    varchar(4096) not null default '' comment '消息内容',
    created_at datetime      not null default now(),
    is_like    tinyint       not null default 0 comment '点赞状态',
    INDEX (user_id)
);

