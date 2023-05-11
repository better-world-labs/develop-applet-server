## 当前版本: v1.4
### 修改记录

- 1.2.1 读取应用列表,1.2.2 读取我创建的应用, 1.2.3 读取我收藏的应用, 1.2.6 读取某个应用 增加 `hot` 字段
- 1.2.10 读取某个应用的运行结果 增加 `nextCursor` 返回以支持分页
- 新增接口 1.2.11 读取某个运行结果
- 新增接口 1.3.6 批量读取我对应用的点赞状态
- 新增文档 1.9 用户
- 新增接口 1.9.3 读取用户个人统计
- [用户上报事件](#user-event-definition)增加事件类型 热度标记

[[_TOC_]]

# 域名

| 环境    | 域名  |
|-------|-----|
| 开发环境  | ai.moyu.dev.openviewtech.com|
| 测试环境  | ai.moyu.test.openviewtech.com|
| 生产环境  | ai.moyu.chat|

# API

## 0. 调用约定

### 0.1 响应体

对于 HTTP 状态码,有

| 状态码 | 说明            |
|-----|---------------|
| 200 | OK            |
| 400 | 客户端异常 (参数错误等) |
| 500 | 服务端异常         |
| 401 | 鉴权异常          |

对于每个 HTTP 请求，都会有以下格式的应答

| 字段   | 说明                                       |
|------|------------------------------------------|
| code | 业务状态码，0 代表正常|
| msg  | 消息描述，进一步描述具体的 code，若 code 为 0，不携带此字段     |
| data | 响应数据，格式为 Json，若 code 不为 0 或者无需返回业务数据，不携带此字段 |

- 例如

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      id: "xxx",
      name: "xxx"
    }
  };
  ```

  或

  ```js
  // HTTP/1.1 401 UNAUTHORIZED
  
  res = {
    "code": 401,
    "msg": "invalid token"
  };
  ```

### 0.2 鉴权约定

- 鉴权信息放在Cookie中，采用 HttpOnly（前端不用关心，js也无法读取）
- 请求头中 需要携带 csrf-token，防止 跨站攻击，csrf-token在登录时生成：

```http
GET /  HTTP/1.1
X-CSRF-TOKEN: {csrf-token}
```

### 0.3 链路追踪约定

为了方便问题排查时的链路追踪，需在每个请求的 Request Header 携带以下字段

```http
GET /  HTTP/1.1
X-Request-ID: {id}
```

并且会随着响应头原样返回；**注意**：若不携带此字段，服务端将自动生成一个

- 字段

| 字段  | 说明                         |
|-----|----------------------------|
| id  | 全局唯一的 id，唯一标识一个请求，可以是 UUID |

- 例子

    - 请求

      ```http
      GET /xxx HTTP/1.1
      X-Request-ID: 123,
      X-CSRF-TOKEN: {csrf-token}
      ```

    - 应答

      ```js
      // HTTP/1.1 200 OK
      // X-Request-ID: 123
      
      httpRes = {
        "code": 0,
        "data": {
            //...
        }
      }
      ```

### 0.4 调用域名

- 开发环境：ai.moyu.dev.openviewtech.com
- 测试环境：ai.moyu.test.openviewtech.com
- 生产环境：ai.moyu.chat

## 1. 接口列表

### 1.1 用户

复用摸鱼接口

### 1.2 应用

#### 1.2.1 读取应用列表

按照热度排序

- 请求

  ```http
  GET /api/apps?category=1 HTTP/1.1
  ```

- 其中

  | 字段   | 说明   |
    |------| --- |
  | category    | 应用类型ID,非必传 |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "list": [
        {
           "id": 1,
           "uuid": "uuid",
           "name": "模板名称",
           "price": 5, //花费积分
           "soldPoints": 100,
           "description": "模板描述",
           "results": [
             {
               "id": 1,
               "type": "text",
               "content": "xxxxxxxxxxxxxxxxxxx"
             },
             {
               "id": 1,
               "type": "text",
               "content": "xxxxxxxxxxxxxxxxxxx"
             }
           ],
           "createdBy": {
             "id": 12,
             "nickname": "xxxx",
             "avatar": "http://xxx/xxx"
           },
           "runTimes": 12,
           "useTimes": 11,
           "hot": 2345,
           "commentTimes": 10,
           "likeTimes": 2,
           "createdAt": "2023-03-22T07:08:02.851Z",
           "updatedAt": "2023-03-22T07:08:02.851Z",
           "status": 0 // 生命周期 (0.未发布 1.已发布)
        }
      ]
    }
  };
  ```

#### 1.2.2 读取我创建的应用

按照创建时间倒序

- 请求

  ```http
  GET /api/apps/mine HTTP/1.1
  ```

- 其中

  | 字段   | 说明   |
    |------| --- |
  | category    | 应用类型ID |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "list": [
        {
          "id": 1,
          "uuid": "uuid",
          "name": "模板名称",
          "price": 5, //花费积分
          "soldPoints": 100,
          "description": "模板描述",
          "results": [
            {
                "id": 1,
                "type": "text",
                "content": "xxxxxxxxxxxxxxxxxxx"
            },
            {
                "id": 1,
                "type": "text",
                "content": "xxxxxxxxxxxxxxxxxxx"
            }
          ],
          "createdAt": "2023-03-22T07:08:02.851Z",
          "updatedAt": "2023-03-22T07:08:02.851Z",
          "createdBy": {
            "id": 12,
            "nickname": "xxxx",
            "avatar": "http://xxx/xxx"
          },
          "runTimes": 12,
          "hot": 2345,
          "useTimes": 11,
          "commentTimes": 10,
          "likeTimes": 2,
          "status": 0 // 生命周期 (0.未发布 1.已发布)
        }
      ]
    }
  };
  ```

#### 1.2.3 读取我收藏的应用

按照收藏时间倒序

- 请求

  ```http
  GET /api/apps/collected HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "list": [
        {
          "id": 1,
          "uuid": "uuid",
          "name": "模板名称",
          "price": 5, //花费积分
          "soldPoints": 100,
          "description": "模板描述",
          "results": [
            {
              "id": 1,
              "type": "text",
              "content": "xxxxxxxxxxxxxxxxxxx"
            },
            {
              "id": 1,
              "type": "text",
              "content": "xxxxxxxxxxxxxxxxxxx"
            }
          ],
          "createdAt": "2023-03-22T07:08:02.851Z",
          "updatedAt": "2023-03-22T07:08:02.851Z",
          "createdBy": {
            "id": 12,
            "nickname": "xxxx",
            "avatar": "http://xxx/xxx"
          },
          "runTimes": 12,
          "useTimes": 11,
          "hot": 2345,
          "commentTimes": 10,
          "likeTimes": 2,
          "status": 0 // 生命周期 (0.未发布 1.已发布)
        }
      ]
    }
  };
  ```

#### 1.2.4 读取某用户的应用

按照创建时间倒序

- 请求

  ```http
  GET /api/users/:userId/apps HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "nextCursor": "xxx",
      "list": [
        {
          "id": 1,
          "uuid": "uuid",
          "name": "模板名称",
          "price": 5, //花费积分
          "soldPoints": 100,
          "description": "模板描述",
          "results": [
            {
              "id": 1,
              "type": "text",
              "content": "xxxxxxxxxxxxxxxxxxx"
            },
            {
              "id": 1,
              "type": "text",
              "content": "xxxxxxxxxxxxxxxxxxx"
            }
          ],
          "createdAt": "2023-03-22T07:08:02.851Z",
          "updatedAt": "2023-03-22T07:08:02.851Z",
          "createdBy": {
            "id": 12,
            "nickname": "xxxx",
            "avatar": "http://xxx/xxx"
          },
          "runTimes": 12,
          "useTimes": 11,
          "hot": 2345,
          "commentTimes": 10,
          "likeTimes": 2,
          "status": 0 // 生命周期 (0.未发布 1.已发布)
        }
      ]
    }
  };
  ```

#### 1.2.4 对应用进行收藏/取消收藏

- 请求

  ```http
  POST /api/apps/:uuid/collect HTTP/1.1

  {
    "collected": true
  }
  ```

- 其中

  | 字段   | 说明   |
    |------| --- |
  | uuid  | APP uuid |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0
  };
  ```

#### 1.2.5 读取我对某些的收藏状态

- 请求

  ```http
  POST /api/apps/is-collected HTTP/1.1

  {
    "uuids": ["app1","app2","app3"]
  }
  ```

- 其中

  | 字段   | 说明   |
  |------| --- |
  | uuids  | APP uuid |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "app1": true,
      "app2": false,
      "app3": true
    }
  };
  ```
***结果中会包含参数给出的所有 uuid***

#### 1.2.6 读取某个应用

- 请求

  ```http
  GET /api/apps/:uuid HTTP/1.1
  ```

- 其中

  | 字段   | 说明   |
    |------| --- |
  | uuid  | APP uuid |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "id": 1,
      "uuid": "uuid",
      "price": 5, //花费积分
      "soldPoints": 100,
      "name": "模板名称",
      "category": 1,
      "description": "模板描述",
      "form": [
        {
          "id": "uuid",
          "label": "姓名",
          "type": "text",
          "properties": {
            "placeholder": "xxx"
          }
        },
        {
          "id": "uuid",
          "label": "性别",
          "type": "select",
          "properties": {
            "placeholder": "xxx",
            "values": "男\n女",
          }
        }
      ],
      "flow": [
        {
          "id": "uuid",
          "type": "chatgpt",
          "outputVisible": true,
          "prompt": [
            {
              "type": "text",
              "properties": {
                "value": "从"
              }
            },
            {
              "type": "tag",
              "properties": {
                "character": "uuid",
                "from": "result/form"
              }
            },
            {
              "type": "text",
              "properties": {
                "value": "选出最好的结果"
              }
            },
          ]
        }
      ],
      "createdBy": {
          "id": 12,
          "nickname": "xxxx",
          "avatar": "http://xxx/xxx"
      },
      "runTimes": 12,
      "useTimes": 11,
      "commentTimes": 10,
      "likeTimes": 2,
      "hot": 2345,
      "createdAt": "2023-03-22T07:08:02.851Z",
      "updatedAt": "2023-03-22T07:08:02.851Z",
      "status": 0 // 生命周期 (0.未发布 1.已发布)
    }
  };
  ```

#### 1.2.7 保存应用

- 请求

  ```http
  PUT /api/apps/:uuid HTTP/1.1

  {
    "name": "模板名称",
    "duplicateFrom": "uuid", //若从某个模板复制，此字段为源模板 uuid ，否则不传。此字段只有第一次保存生效，后续的修改无效
    "category": 1,
    "description": "模板描述",
    "form": [
        {
            "id": "uuid",
            "label": "姓名",
            "type": "text",
            "properties": {
                "placeholder": "xxx"
            }
        },
        {
            "id": "uuid",
            "label": "性别",
            "type": "select",
            "properties": {
                "placeholder": "xxx",
                "values": "男\n女"
            }
        }
    ],
    "flow": [
      {
        "type": "chatgpt",
        "outputVisible": true,
        "prompt": [
          {
            "id": "uuid",
            "type": "text",
            "properties": {
                "value": "从"
            }
          },
          {
            "id": "uuid",
            "type": "tag",
            "properties": {
                "from": "result", //form 或者 result
                "character": "uuid"
            }
          },
          {
            "id": "uuid",
            "type": "text",
            "properties": {
                "value": "选出最好的结果"
            }
          },
        ]
      }
    ],
    "createdBy": {
      "id": 12,
      "nickname": "xxxx",
      "avatar": "http://xxx/xxx"
    },
    "createdAt": "2023-03-22T07:08:02.851Z",
    "updatedAt": "2023-03-22T07:08:02.851Z",
    "status": 0 // 状态 (0.未发布 1.已发布)
  }

  ```

- 其中

  | 字段   | 说明   |
    |------| --- |
  | uuid    | 应用 UUID 如果是创建模板，则由前端生成一个 |


- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "id": 1,
      "uuid": "uuid",
      "name": "模板名称",
      "category": 1,
      "description": "模板描述",
      "form": [
        {
          "id": "uuid",
          "label": "姓名",
          "type": "text",
          "properties": {
            "placeholder": "xxx"
          }
        },
        {
          "id": "uuid",
          "label": "性别",
          "type": "select",
          "properties": {
            "placeholder": "xxx",
            "values": "男\n女",
          }
        }
      ],
      "flow": [
        {
          "id": "uuid",
          "type": "chatgpt",
          "outputVisible": true,
          "prompt": [
            {
              "type": "text",
              "properties": {
                "value": "从"
              }
            },
            {
              "type": "tag",
              "properties": {
                "charactor": "uuid"
                "from": "result/form"
              }
            },
            {
              "type": "text",
              "properties": {
                "value": "选出最好的结果"
              }
            },
          ]
        }
      ],
      "createdBy": {
        "id": 12,
        "nickname": "xxxx",
        "avatar": "http://xxx/xxx"
      },
      "runTimes": 12,
      "useTimes": 11,
      "commentTimes": 10,
      "likeTimes": 2,
      "createdAt": "2023-03-22T07:08:02.851Z",
      "updatedAt": "2023-03-22T07:08:02.851Z",
      "status": 0 // 生命周期 (0.未发布 1.已发布)
    }
  };
  ```

#### 1.2.8 删除应用

只能删除自己的应用

- 请求

  ```http
  DELETE /api/apps/:uuid HTTP/1.1
  ```

- 其中

  | 字段   | 说明   |
    |------| --- |
  | uuid    | 应用 UUID 如果是创建模板，则由前端生成一个 |


- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0
  }
  ```

#### 1.2.9 运行应用

- 请求

  ```http
  POST /api/apps/:uuid/run HTTP/1.1
  Accept: text/event-stream
  Content-Type: application/json;charset=utf-8

  {
    "values": ["text参数","男"],
    "open": true
  }
  ```

- 其中

  | 字段   | 说明   |
    |------| --- |
  | uuid   | app uuid |
  | values   | 表单参数，按顺序传递,均传递字符串 |
  | open   | 是否开放运行结果 |

- 应答

  ```
  HTTP/1.1 200 OK
  Content-Type: text/event-stream

  event: data
  data: {"flow":"uuid","type":"text","content":"这"}
  
  event: data
  data: {"flow":"uuid","type":"text","content":"是"}
  
  event: data
  data: {"flow":"uuid","type":"text","content":"数"}
  
  event: data
  data: {"flow":"uuid","type":"text","content":"据"}
  
  event: done
  data: {"code":0}
  ```

* 说明

  数据成组出现，每组数据流由 `event` 和 `data` 两种结构组成， 并用换行分割, 每组数据用两个换行(\n\n)进行分割

| 消息    | 说明                                                                        |
|-------|---------------------------------------------------------------------------|
| event | 消息头，描述事件类型。若为 `data` 则表示数据，若为 `done`则表示结束                                 | 
| data  | 消息载荷，若 `event` 为 `data`，则携带数据，若 `event` 为 `done` 则携带结束状态，与 http 业务响应体含义一致 |

- 对于 `data` 的特别说明

  事件流的数据结构与普通 `output` 实体结构一致，但 `content` 变为分片数据，其中 `flow` 为逻辑流的 `id` 当数据来自不同的逻辑流时，`flow` 字段用于区分

- 可能的业务 `code`

  | code | 说明     |
    |--------| --- | 
  | 500000 | 积分余额不足 | 

#### 1.2.10 读取某个应用的运行结果

按照最新排序分页返回，一页50条记录

- 请求

  ```http
  GET /api/apps/:uuid/outputs HTTP/1.1
  ```

- 其中

  | 字段   | 说明   |
    |------| --- |
  | uuid    | 应用ID |

- 应答

  ```js
  // HTTP/1.1 200 OK
  resp = {
    "code": 0,
    "data": {
      "nextCursor": "xxxxx",
      "list": [
        {
          "id": "1", 
          "type": "text",
          "inputArgs":["xxx","xxx"],
          "content": "哈哈哈",
          "likeTimes": 12, //点赞数
          "hateTimes": 12, //踩数
          "commentTimes": 12, //评论数
          "createdAt": "xxxx",
          "createdBy": {
              "id": 12,
              "nickname": "xxxx",
              "avatar": "http://xxx/xxx"
          }
        }
      ]
    }
  }
  ```

#### 1.2.11 读取某个运行结果

- 请求

  ```http
  GET /api/outputs/:outputId HTTP/1.1
  ```

- 其中

  | 字段   | 说明   |
    |------| --- |
  | outputId    | 运行结果ID |

- 应答

  ```js
  // HTTP/1.1 200 OK
  resp = {
    "code": 0,
    "data": {
      "id": "1", 
      "type": "text",
      "inputArgs":["xxx","xxx"],
      "content": "哈哈哈",
      "likeTimes": 12, //点赞数
      "hateTimes": 12, //踩数
      "commentTimes": 12, //评论数
      "createdAt": "xxxx",
      "createdBy": {
          "id": 12,
          "nickname": "xxxx",
          "avatar": "http://xxx/xxx"
      }
    }
  }
  ```
#### 1.2.12 读取AI模型列表

按照热度排序

- 请求

  ```http
  GET /api/ai-models HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  resp = {
    "code": 0,
    "data": {
      "list": [
        {
          "category": "语言类模型",
          "models": [
            {
              "id": 1, 
              "name": "chatgpt",
              "description": "ChatGPT",
              "icon":"https://xxx/xx",
              "available": true
            }
          ]
        }
      ]
    }
  }
  ```

#### 1.2.13 读取 App 分类

- 请求

  ```http
  GET /api/app-categories HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  resp = {
    "code": 0,
    "data": {
      "list": [
        {
          "id": 1, 
          "text": "文学"
        },
        {
          "id": 2, 
          "text": "文学"
        }
      ]
    }
  }
  ```

#### 1.2.14 读取首页 Tab

- 请求

  ```http
  GET /api/app-tabs HTTP/1.1
  ```

- 应答

  ```js
  resp = {
    "code": 0,
    "data": {
      "tabs": [
        {
          "label": "热门",
          "category": 0,
          "order": 0
        },
        {
          "label": "生活",
          "category": 1,
          "order": 1
        },
        {
          "label": "文学",
          "category": 2,
          "order": 2
        },
        {
          "label": "办公神器",
          "category": 3,
          "order": 3
        }
      ]
    }
  }
   ```

#### 1.2.15 事件上报

前端上报事件

- 请求

  ```http
  POST /api/events HTTP/1.1
  
  {
    "type": "event-type"
    "args": ["args"]
  }
  ```

- 应答

  ```js
  resp = {
    "code": 0
  }
   ```
- 注意
 
 调用此接口可以不登录

- **事件类型参考 [2. 事件定义](#user-event-definition)**


### 1.3 点赞&评论

#### 1.3.1 为应用添加评论

- 请求

  ```http
  POST /api/apps/:uuid/comments HTTP/1.1

  {
    "content": "xxx"
  }
  ```

- 应答

  ```js
  resp = {
    "code": 0
  }
   ```

#### 1.3.2 为应用输出添评论 (暂不实现)

- 请求

  ```http
  POST /api/outputs/:id/comments HTTP/1.1

  {
    "content": "xxx"
  }
  ```

- 应答

  ```js
  resp = {
    "code": 0
  }
   ```

#### 1.3.3 对应用进行点赞/取消点赞

- 请求

  ```http
  POST /api/apps/:uuid/like HTTP/1.1

  {
    "like": true
  }
  ```
- 其中

  | 字段 | 说明 |
    | --- | --- |
  | like | 点赞/取消点赞 |

- 应答

  ```js
  resp = {
    "code": 0
  }
   ```

#### 1.3.4 对应用输出进行顶/踩/取消点赞

- 请求

  ```http
  POST /api/outputs/:id/like HTTP/1.1

  {
    "like": 1
  }
  ```
- 其中

  | 字段 | 说明                  |
    |---------------------| --- |
  | like | 1. 顶, 0. 取消点赞 -1. 踩 |

- 应答

  ```js
  resp = {
    "code": 0
  }
   ```

#### 1.3.5 读取我对应用的点赞状态

- 请求

  ```http
  GET /api/apps/:id/like HTTP/1.1
  ```
- 其中

  | 字段 | 说明 |
    | --- | --- |
  | like | 点赞/取消点赞 |

- 应答

  ```js
  resp = {
    "code": 0,
    "data": {
      "like": true,
    }
  }
   ```

#### 1.3.6 批量读取我对应用的点赞状态

- 请求

  ```http
  POST /api/apps/:id/get-likes HTTP/1.1
  
  {
    "appIds": ["xxx"]
  }
  ```

- 其中

  | 字段 | 说明 |
    | --- | --- |
  | like | 点赞/取消点赞 |

- 应答

  ```js
  resp = {
    "code": 0,
    "data": {
      "app1": true,
      "app2": true,
      "app3": false,
    }
  }
   ```
  ***输入参数均会出现在应答中***

#### 1.3.7 读取应用的评论列表

- 请求

  ```http
  GET /api/apps/:uuid/comments HTTP/1.1
  ```

- 应答

  ```js
  resp = {
    "code": 0,
    "data": {
      "list": [
        {
          "id": 12,
          "content": "评论",
          "createdAt": "2023-03-22T07:08:02.851Z",
          "createdBy": {
              "id": 12,
              "nickname": "xxxx",
              "avatar": "http://xxx/xxx"
          }
        }
      ]
    }
  }
   ```

#### 1.3.8 读取我对应用输出的点赞状态

- 请求

  ```http
  GET /api/outputs/likes?outputIds=xxx,xxx,xxx HTTP/1.1
  ```
- 其中

  | 字段  | 说明   |
    |------| --- |
  | outputIds   | 输出ID |

- 应答

  ```js
  resp = {
    "code": 0,
    "data": {
      "list": [
        {
          "outputId": "12",
          "like": 1
        },
        {
          "outputId": "12",
          "like": 0
        }
      ]
    }
  }
   ```

#### 1.3.9 读取应用输出的评论列表 (暂不实现)

- 请求

  ```http
  GET /api/outputs/:id/comments HTTP/1.1
  ```

- 应答

  ```js
  resp = {
    "code": 0,
    "data": {
      "list": [
        {
          "id": 12,
          "content": "评论",
          "createdAt": "2023-03-22T07:08:02.851Z",
          "createdBy": {
            "id": 12,
            "nickname": "xxxx",
            "avatar": "http://xxx/xxx"
          }
        }
      ]
    }
  }
   ```

### 1.4 新手引导

### 1.4.1 读取新手引导完成状态

- 请求

  ```http
  GET /api/users/me/guidance HTTP/1.1
  ```

- 应答

  ```js
  resp = {
    "code": 0,
    "data": {
      "completed": true
    }
  }
   ```

### 1.4.2 完成新手引导

- 请求

  ```http
  POST /api/users/me/guidance/completion HTTP/1.1
  ```

- 应答

  ```js
  resp = {
    "code": 0
  }
   ```

### 1.5 积分

#### 1.5.1 读取用户积分流水

- 请求

  ```http
  GET /api/points?page=xxx&size=xxx HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "total": 100, //结果总数
      "list": [
        {
          "userId": 31, // 用户ID
          "points": 100, // 积分数
          "description": "空投奖励", // 描述
          "createdAt": "2022-09-22T11:11:27+08:00"
        } 
      ] 
    }
  };
  ```
#### <span id="getPointTotal">1.5.2 积分总数查询</span>

查询用户当前积分总数

- 请求

  ```http
  GET /api/points/total HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      total: 200, //积分总数
      withdrawAmount: 2, //可提现金额 (元)
    }
  };
  ```

#### 1.5.3 读取积分充值档位

- 请求

  ```http
  GET /api/points-goods HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "list": [
        {
          "id": 12,
          "price": 9.9,
          "points": 300,
          "tag": ""
        },
        {
          "id": 13,
          "price": 0.99,
          "points": 50,
          "tag": "new-deal"
        }
      ]
    },
  };
  ```
- 其中

对于 `tag` 标识了商品的特殊属性，前端可以根据约定的 `tag` 标识做多种表现方式

 | TAG | 说明|
 | new-deal | 新人特惠，首次充值才可以使用|

#### 1.5.4 用户是否有充值记录

- 请求

  ```http
  GET /api/points-orders/exists HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "exists": true
    },
  };
  ```

#### 1.5.5 读取用户的充值订单

- 请求

  ```http
  GET /api/points-orders HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "list": [
         {
           "orderId": "7F50C7E74F884F94B72D348B1BF9492C",
           "goods": {
             "id": 1,
             "price": 0.01,
             "tag": "",
             "points": 50,
             "description": "测试"
           },
           "goodsId": 1,
           "price": 0.01,
           "userId": 10058,
           "createdAt": "2023-04-26T10:51:58.3775153+08:00",
           "state": 0,
           "codeUrl": "weixin://wxpay/bizpayurl?pr=8MOQXJ2zz",
           "PayExpiresAt": "2023-04-26T12:51:58.3775153+08:00",
           "payedAt": null
        }
      ]
    }
  };
  ```

#### 1.5.5 积分充值下单

生成积分充值订单,并检查权限,阻止不满足条件的下单

- 请求

  ```http
  POST /api/points-goods/:goodsId/points-orders HTTP/1.1
  ```
- 其中

  | 字段| 说明|
  | --- | --- |
  | goodsId | 充值档位商品ID|

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "orderId": "7F50C7E74F884F94B72D348B1BF9492C",
      "goods": {
        "id": 1,
        "price": 0.01,
        "tag": "",
        "points": 50,
        "description": "测试"
      },
      "goodsId": 1,
      "price": 0.01,
      "userId": 10058,
      "createdAt": "2023-04-26T10:51:58.3775153+08:00",
      "state": 0,
      "codeUrl": "weixin://wxpay/bizpayurl?pr=8MOQXJ2zz",
      "PayExpiresAt": "2023-04-26T12:51:58.3775153+08:00",
      "payedAt": null
    }
  }
  ```
- 可能返回的状态码

   | 状态码| 说明|
   | ---| ---|
   |  403010 | 非新人禁止使用新人档充值  |


#### 1.5.6 读取充值订单信息 

- 请求

  ```http
  GET /api/points-orders/:orderId HTTP/1.1
  ```
- 其中

  | 字段  | 说明      |
  |---------| --- |
  | orderId   | 充值订单ID |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "orderId": "7F50C7E74F884F94B72D348B1BF9492C",
      "goods": {
        "id": 1,
        "price": 0.01,
        "tag": "",
        "points": 50,
        "description": "测试"
      },
      "goodsId": 1,
      "price": 0.01,
      "userId": 10058,
      "createdAt": "2023-04-26T10:51:58.3775153+08:00",
      "state": 0,
      "codeUrl": "weixin://wxpay/bizpayurl?pr=8MOQXJ2zz",
      "PayExpiresAt": "2023-04-26T12:51:58.3775153+08:00",
      "payedAt": null
    }
  }
  ```
- Data 为空表示订单不存在 

#### 1.5.7 积分提现

积分申请提现

- 请求

  ```http
  POST /api/points-withdraw HTTP/1.1

  {
    "points": 2000
  }
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0
  };
  ```

- 可能返回的状态码

   | 状态码| 说明|
   | ---| ---|
   |  500010 | 积分余额未达到提现标准  |
   |  500000 | 积分余额不足  |

-  触发事件
  
  积分流水变化提醒


### 1.6 消息盒子

#### 1.6.1 读取消息列表

- 请求

  ```http
  GET /api/notify-messages?isRead=true&cursor=xxx HTTP/1.1
  ```
- 其中

 | 字段| 说明|
 | --- | --- |
 | isRead | 筛选未读消息，非必填, 若传递，按照指定状态筛选，若不传递,不筛选 |
 | cursor | 分页游标，来源于上一页的 data 数据, 不填返首页数据 |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "nextCursor": "xxxx", //下一页的游标
      "list": [
        {
          "id": 21,
          "icon": "https//xxx/xxx",
          "title": "空投奖励",
          "content": "嗨新伙伴，欢迎加入我们！送你 30 积分，和我们开始创造吧！",
          "isRead": false,
          "createdAt": "2023-03-22T07:08:02.851Z"
        },
        {
          "id": 21,
          "title": "好友邀请",
          "icon": "https//xxx/xxx",
          "content": "邀请用户（老用户）：哇塞！邀请好友成功，恭喜您获得 100 积分",
          "isRead": true,
          "createdAt": "2023-03-22T07:08:02.851Z"
        }
      ]
    }
  };
  ```

#### 1.6.2 读取未读消息数

- 请求

  ```http
  GET /api/notify-messages/unread-count HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "count": 12
    }
  };
  ```

#### 1.6.3 标记某条消息为已读

- 请求

  ```http
  PUT /api/notify-messages/:id/read HTTP/1.1
  ```
- 其中

 | 字段| 说明|
 | --- | --- |
 | id | 消息ID |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0
  };
  ```

#### 1.6.3 标记消息全部已读

- 请求

  ```http
  PUT /api/notify-messages/read-all HTTP/1.1
  ```
- 其中

 | 字段| 说明|
 | --- | --- |
 | id | 消息ID |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0
  };
  ```

### 1.7 每日签到

#### 1.7.1 读取用户的签到状态

- 请求

  ```http
  GET /api/sign-in HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "signIn": false
    }
  };
  ```


#### 1.7.2 进行签到

- 请求

  ```http
  POST /api/sign-in HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0
  };
  ```

### 1.8 提醒

#### 1.8.3 读取保留消息

- 请求

  ```http
  GET /api/retain-messages HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    "code": 0,
    "data": {
      "list": [
         {
           "id": "xxx",
           "type": "friends-first-login", //邀请好友(合并消息)
           "payload": {
             "points": 80,
             "friends": [
               {
                  "id": 21, 
                  "nickname": "xx",
                  "avatar": "xx"
               } ,
               {
                  "id": 22, 
                  "nickname": "xx",
                  "avatar": "xx"
               } 
             ]
           }
         }
      ]
    }
  };
  ```

### 1.9 用户

#### 1.9.1 获取用户个人信息

复制自摸鱼

- 请求

  ```http
  GET /api/users/me/info HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  {
    "code": 0,
    "data": {
      "id": 0,
      "nickname": "哈哈哈哈哈",
      "avatar": "https://openview-oss.oss-cn-chengdu.aliyuncs.com/aed-test/avatar/93.png",
      "lastLoginAt": "xxxxx" //为空表示首次登录
      "invitedBy": 23, //为空表示自然流量
      "points": 80
    }
  }
  ```

#### 1.9.2 修改个人资料信息

复制自摸鱼

昵称长度：八个汉字

- 请求

  ```http
  PUT /api/users/me/info HTTP/1.1
  
  {
    nickname: "xxxx",
    avatar:"https://sxxx"
  }
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      "code": 0
  }
  ```

#### 1.9.3 读取用户个人统计

- 请求

  ```http
  GET /api/users/:userId/statistic HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      "code": 0,
      "data": {
        "apps": 390,
        "registeredDays": 30,
        "points": 30,
        "appUses": 30999,
        "appLikes": 30999,
      }
  }
  ```


### 1.10 服务端推送

通知接收类的服务端消息推送，采用 `Websocket` 的方式实现

#### 端点

```http
wss://ai.moyu.dev.openviewtech.com/push/endpoint?scrf=
```
- 其中

  | 字段   | 说明 |
  | --- | --- |
  | scrf    | CSRF Token |

### 客户端

[Websocket 客户端](https://github.com/shenweijiekdel/light-websocket-client-ts)

***事件定义见 [3. 服务端推送事件定义](#event-definition)***

## 2.<span id="user-event-definition">用户上报事件定义</span>

- 格式

   ```json
   {
     "type": "xxx",
     "args": []
   } 
   ```

- 其中

  | 字段   | 说明   |
  |------| --- |
  | type | 事件类型 |
  | args    | 参数 |

具体的类型有如下定义

### 2.1 APP

#### (1) 浏览 APP

  ```json
  {
    "type": "app-viewed",
    "args": ["uuid"]
  }
  ``` 

- 其中

  | 字段     |类型 | 说明   |
  | -------|------| --- |
  | args[0] | 字符串 | App uuid| 

#### (2) 热度标记

  ```json
  {
    "type": "app-hot-mark",
    "args": ["uuid"]
  }
  ``` 

- 其中

  | 字段     |类型 | 说明   |
  | -------|------| --- |
  | args[0] | 字符串 | App uuid| 

## 3.<span id="event-definition">服务端推送事件定义</span>

### 3.1 分享裂变

####（1）分享提示创建APP触发

  ```json
  {
    "type": "share-hint-create-app",
    "payload": {
      "createdApps": 1, // 创建小程序数量
      "earnPoints": 12 // 获得积分
    }
  } 
  ```

####（2）分享提示使用APP触发

  ```json
  {
    "type": "share-hint-use-app",
    "payload": {
      "usedApps": 3, // 使用小程序数量
      "costPoints": 12 // 花费积分
    }
  }
  ```

### 3.3 通知消息

#### (1) 通知消息变化提醒

  ```json
  {
    "type": "notify-message-changed",
    "payload": {
      "unread": 10 
    }
  } 
  ```

#### (2) 保留消息变化提醒

  ```json
  {
    "type": "retain-message-changed"
  } 
  ```

### 3.4 业务消息

#### (1) 用户积分变化

  ```json
  {
    "type": "user-points-changed"
  } 
  ```
