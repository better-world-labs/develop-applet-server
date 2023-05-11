alter table robot_config
    add env varchar(10) default '' not null comment '指定环境 (只处理来自指定环境的机器人消息)';
