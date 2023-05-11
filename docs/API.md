## 当前版本: v0.0.1

[[_TOC_]]

# API

## 0. 调用约定

### 0.1 响应体

对于每个 HTTP 请求，都会有以下格式的应答

| 字段   | 说明                                                                    |
|------|-----------------------------------------------------------------------|
| code | 业务状态码，0 代表成功，对应的 HTTP Status 为 200，若不为 0 则代表 error，HTTP Status 不为 200 |
| msg  | 消息描述，进一步描述具体的 code，若 code 为 0，不携带此字段                                  |
| data | 响应数据，格式为 Json，若 code 不为 0 或者无需返回业务数据，不携带此字段                           |

- 例如

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
      id: "xxx",
      name: "xxx",
    },
  };
  ```

  或

  ```js
  // HTTP/1.1 401 UNAUTHORIZED
  
  res = {
    code: 401,
    msg: "invalid token",
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
          "msg": "ok",
          "data": {
              //...
          }
      }
      ```

### 0.4 调用域名

- 开发环境：moyu.dev.openviewtech.com
- 测试环境：moyu.test.openviewtech.com
- 生产环境：moyu.chat

## 1. 接口列表

### 1.1 登录相关

扫描流程设计：

1. 【PC端】通过接口 [1.1.1 获取扫码登录的二维码【PC端调用】](#111-获取扫码登录的二维码pc端调用)获取二维码的`qrToken`；
    1. 通过`qrToken`进入登录场景，完成登录场景和websocket通道的绑定：
       ```js
       
       socket.emit("enter-scene", {scene:"login", param: `${qrToken}` })
       
       ```  
    2. 使用`qrToken`拼接出完整的二维码内容（如`https://moyu.chat/login/mobile-auth?qrToken=${qrToken}`）并渲染
2. 【用户】使用"微信"扫描【PC端】的二维码
3. 【手机端】通过接口 [1.1.3 获取鉴权跳转参数【Mobile端调用】](#113-获取鉴权跳转参数mobile端调用)获取`qrToken`
   参数，校验过期时间，跳转到链接`authUrl`
4. 【手机端】页面`authUrl`鉴权完成，跳回到 手机登录业务页面
5. 【手机端】在 手机登录业务页面 ，取回鉴权`code`和`state`
   ，调用接口[1.1.4 移动端提交登录 `code`【Mobile端调用】](#114-移动端提交登录codemobile端调用) 将信息提交给后端
6. 【后端】通过 鉴权`code`获取用户唯一ID，通过`qrCode`标识的socket通道给【PC端】发送登录成功的`Trigger`:
7.
【PC端】在收到登录成功的Trigger时，调用接口 [1.1.2 使用 `loginToken` 完成登录【PC端调用】](#112-使用logintoken完成登录pc端调用)
完成登录:
 ```js
     socket.on('t', ({type, params})=>{
         if (type == "loginSuc"){
             api
                 .post(`/api/users/login`, {loginToken: params[0], nickname, avatar})
                 .then(res=>{
                     //...
                 })
         }
     })
 ```

#### 1.1.1 获取扫码登录的二维码【PC端调用】

- 请求
    ```http
    GET /api/users/login/qr HTTP/1.1
    ```

- 应答
    ```js
    httpStatus = 200;
    httpRes = {
        code: 0,
        data: {
            // 二维码token; 
            // 1. 通过websocket发给服务器:`socket.emit("enter-scene", {scene:"login", param: `${qrToken}` })`
            // 2. 拼接成为二维码内容，渲染显示
            qrToken: "xxx-xxwsxx-xxxs",  
            expired: "2022-09-20T05:28:51.132Z", //token 过期时间
        },
    }
    ```

#### 1.1.2 使用`loginToken`完成登录【PC端调用】

v0.2: 将字段 `nickname`、`avatar` 两个字段调整为非必填；在首次登录时，如果未填`nickname`、`avatar`
字段，接口随机分配；非首次登录，未填则不修改用户昵称和头像；

- 请求
    ```http
    POST /api/users/login HTTP/1.1
    
    {
        //登录token，来自websocket推送
        loginToken: `${loginToken}`,
        nickname: "昵称", //用户昵称，(v0.2 调整为 非必填)
        avatar: "https://xxx.com/image/xx.png", //头像链接 (v0.2 调整为 非必填)
        invitedBy: 1, //邀请用户,若存在则携带，非必填
        formApp: "xxx" // ['mini-app', 'moyu']
    }
    ```
- 应答
    - 正常登录
        ```js
        // SET-Cookie: token=xxxx
      
        httpStatus = 200;
        httpRes = {
            code: 0,
            data: {
                csrfToken: "xxxxx", //csrf Token
                user: {
                    id: 110001,
                    nickname: "大熊",
                    avatar: "https://xx.com/image/xx.png"
                },
            },
        };
        ```
    - 异常情况，在黑名单中-禁止登录
        ```js
      
        httpStatus = 400;
        httpRes = {
            code: 400,
            msg: "用户被禁用",
        };
        ```

#### 1.1.3 获取鉴权跳转参数【Mobile端调用】

> qrToken 二维码携带的token，从页面链接中获取

- 请求
    ```http
    GET /api/users/login/auth-params?qrToken=${qrToken}&redirectUrl=${redirectUrl}  HTTP/1.1
    ```
- 应答
    ```js
    httpStatus = 200;
    httpRes = {
        code: 0,
        data: {
            tokenExpired: "2022-09-20T05:28:51.132Z", //token 过期时间
            //页面授权地址
            authUrl: "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx520c15f417810387&redirect_uri=https%3A%2F%2Fchong.qq.com%2Fphp%2Findex.php%3Fd%3D%26c%3DwxAdapter%26m%3DmobileDeal%26showwxpaytitle%3D1%26vb2ctag%3D4_2030_5_1194_60&response_type=code&scope=snsapi_base&state=123#wechat_redirect", 
        },
    };
    ```

#### 1.1.4 移动端提交登录code【Mobile端调用】

- 请求
    ```http
    POST /api/users/login/code  HTTP/1.1
    Content-Type: application/json;charset=utf8
    
    {
        //qrToken 二维码携带的token，通过authUrl中的state传递
        qrToken: `${state}`,
        
        //微信鉴权code，authUrl页面回调获得
        code: `${code}`,
    }
    ```
- 应答
 
    ```js
    httpStatus = 200;
    httpRes = {
        code: 0,
    }
    ```
- 或者
    ```js
    httpStatus = 200;
    httpRes = {
        code: 401001  //用户未授权
    }
    ```

#### 1.1.5 测试登录接口 【仅用于测试】

- 请求
    ```http
    POST /api/users/login/test HTTP/1.1
    
    {
        openId: `${wechat-openId}`,
    }
    ```
- 应答
    ```js
    // SET-Cookie: token=xxxx
    
    httpStatus = 200;
    httpRes = {
        code: 0,
        data: {
            csrfToken: "xxxxx", //csrf Token
            user: {
                id: 110001,
                nickname: "大熊",
                avatar: "https://xx.com/image/xx.png",
            },
        },
    };
    ```

#### 1.1.6 注销登陆

- 请求
    ```http
    POST /api/users/logout HTTP/1.1
    ```
- 应答
    ```js
    // SET-Cookie: token=xxxx
    
    httpStatus = 200;
    httpRes = {
        code: 0
    };
    ```

### 1.2 星球管理

#### 1.2.1 修改星球名称，头像，封面 【管理员/超级管理员有权调用，目前icon字段未用上，传参可以不携带此字段】

- 请求

    ```http
    PUT /api/planets/1/msg HTTP/1.1
   
    {
        icon: "https://xx",
        frontCover: "https://xxa",
        name: "摸鱼猩球"
    }
    ```
- 应答gc
    ```js
    httpStatus = 200;
    httpRes = {
        code: 0,
        data: {}
    }
    ```

#### 1.2.2 读取星球成员列表 (分页)

- 请求

  ```http
  GET /api/planets/1/members?page=1&size=20 HTTP/1.1
  ```

- 其中

  | 字段  | 说明   |
    |------| --- |
  | page   | 页码 |
  | size   | 单页数据数 |

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          total: 100,
          list: [
              {
                  user: {
                      id: 1,
                      nickname: "沙发",
                      avatar: "https://xxx/xxx",
                      online: true
                  },
                  status: 1, // 成员状态 (1: 黑名单, 0: 正常)
                  role: 1  // 成员角色 (0. 普通成员 1. 管理员, 2. 超级管理员)
              },
              {
                  user: {
                      id: 2,
                      nickname: "沙发2",
                      avatar: "https://xxx/xxx",
                      online: true
                  },
                  status: 1, // 成员状态 (1: 黑名单, 0: 正常)
                  role: 1  // 成员角色 (0. 普通成员 1. 管理员, 2. 超级管理员)
              },
          ]  
      }
  }
  ```

#### 1.2.3 查询星球基本属性

- 请求

  ```http
  GET /api/planets/1/msg HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
        "id": 1,
        "name": "摸鱼猩球1",
        "icon": "https://xx1",
        "frontCover": "https://xxa1",
        "createdAt": "2022-09-22T11:11:27+08:00"
      }
  }
  ```

#### 1.2.4 读取用户在某个星球的信息

- 请求

  ```http
  GET /api/planets/1/members/me HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  {
     code: 0,
     data: {
        role: 1, // 角色
        status: 1 // 用户状态
     }
  }
  ```

#### 1.2.5 读取星球成员数

- 请求

  ```http
  GET /api/planets/1/members-count HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          "count": 10
      }
  }
  ```

### 1.3 频道管理

#### 1.3.1 读取频道组列表

- 请求

  ```http
  GET /api/channels/groups?planetId=1 HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          list: [
              {
                  id: 1,
                  name: "欢迎大厅",
                  icon: "https://xxx/xxx",
                  planetId: 1,
              },
              {
                  id: 2,
                  name: "信息处",
                  icon: "https://xxx/xxx",
                  planetId: 1,
              }
          ]  
      }
  }
  ```

#### 1.3.2 创建频道组【目前icon字段未用上，传参可以不携带此字段】

需要管理员权限

- 请求

  ```http
  POST /api/channels/groups HTTP/1.1
  
  {
      "planetId": 1,
      "name": "交流区",
      "icon": "https://xxx/xx",
  }
  ```
- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```

#### 1.3.3 删除频道组

需要管理员权限

- 请求

  ```http
  DELETE /api/channels/groups/:id HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```

#### 1.3.4 修改频道在分组下的顺序

[1, 4, 3, 6] => id为1的频道顺序改为0，id为4的字段顺序改为1 ...

- 请求

    ```http
    PUT /api/channels/groups/sort HTTP/1.1
  
    {
        groupId: 1
        sortedChannelIds: [1,4,3,6]
    }
    ```

- 应答gc
    ```js
    httpStatus = 200;
    httpRes = {
        code: 0,
        data: {}
    }
    ```

#### 1.3.5 读取频道列表

后端按照sort字段返回给前端，以小到大顺序排序

- 请求

  ```http
  GET /api/channels?planetId=1 HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          list: [
              {
                  id: 1,
                  name: "新人接待",
                  icon: "https://xxx/xxx",
                  type: 1, // (1. 普通频道， 2. 私密频道)
                  mute: true, // 禁言
                  notice: "公告", // 公告
                  groupId: 1,
                  planetId: 1,
                  sort: 1,
                  status: 1, // (1. 有效， 2. 已过期),
                  expiresAt: "2022-09-22T11:11:27+08:00"
              },
              {
                  id: 2,
                  name: "社区守则",
                  icon: "https://xxx/xxx",
                  type: 1, // (1. 普通频道， 2. 私密频道)
                  mute: true, // 禁言
                  notice: "公告", // 公告
                  groupId: 1,
                  planetId: 1,
                  sort: 0,
                  status: 1, // (1. 有效， 2. 已过期),
                  expiresAt: "2022-09-22T11:11:27+08:00"
              }
          ]  
      }
  }
  ```

#### 1.3.6 读取某个频道信息

- 请求

  ```http
  GET /api/channels/:channelId HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          id: 1,
          name: "新人接待",
          icon: "https://xxx/xxx",
          type: 1, // (1. 普通频道， 2. 私密频道)
          mute: true, // 禁言
          notice: "公告", // 公告
          groupId: 1,
          planetId: 1,
          sort: 1,
          status: 1, // (1. 有效， 2. 已过期),
          expiresAt: "2022-09-22T11:11:27+08:00"
      }
  }
  ```

* 状态码

| 状态码    | 说明    |
|--------|-------|
| 404001 | 频道不存在 |

#### 1.3.7 批量某个频道的基本信息

- 请求

  ```http
  POST /api/channels/query-many HTTP/1.1
  
  {
      "ids": [45, 46, 47, 49, 50, 51]
  }
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          list: [
              {
                  id: 1,
                  name: "新人接待",
                  icon: "https://xxx/xxx",
                  type: 1, // (1. 普通频道， 2. 私密频道)
                  mute: true, // 禁言
                  notice: "公告", // 公告
                  groupId: 1,
                  planetId: 1,
                  sort: 1,
                  status: 1, // (1. 有效， 2. 已过期),
                  expiresAt: "2022-09-22T11:11:27+08:00"
              },
              {
                  id: 2,
                  name: "社区守则",
                  icon: "https://xxx/xxx",
                  type: 1, // (1. 普通频道， 2. 私密频道)
                  mute: true, // 禁言
                  notice: "公告", // 公告
                  groupId: 1,
                  planetId: 1,
                  sort: 0,
                  status: 1, // (1. 有效， 2. 已过期),
                  expiresAt: "2022-09-22T11:11:27+08:00"
              }
          ]
      }
  ```

#### 1.3.7 创建频道

需要管理员权限

- 请求

  ```http
  POST /api/channels HTTP/1.1
  
  {
      "planetId": 1,
      "name": "公告",
      "icon": "https://xxx/xx",
      "type": 1, // (1. 普通频道， 2. 私密频道)
      "mute": true, // 禁言
      "groupId": 3,
      "exipresIn": 1000000000   // jueduimiaoshu
  }
  ```
- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          channelId: 2  
      }
  }
  ```

#### 1.3.8 删除频道

需要管理员权限

- 请求

  ```http
  DELETE /api/channels/:id HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```

#### 1.3.9 读取频道成员

- 请求

  ```http
  GET /api/channels/:channelId/members HTTP/1.1
  ```

- 其中

  | 字段 | 说明   |
                |------| --- |
  | channelId | 频道ID |

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          total: 10,
          online: 10,
          list: [
              {
                  id: 1,
                  nickname: "沙发",
                  avatar: "https://xxx/xxx",
                  online: true,
                  role: 1,
              },
              {
                  id: 2,
                  nickname: "欧阳铁柱",
                  avatar: "https://xxx/xxx",
                  online: false,
                  role: 1,
              },
          ]  
      }
  }
  ```

#### 1.3.10 读取我在某个频道的授权状态

- 请求

  ```http
  GET /api/channels/:channelId/member-state HTTP/1.1
  ```

- 其中

  | 字段 | 说明   |
          |------| --- |
  | channelId | 频道ID |

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          state: 1 // (0. 未加入, 1.申请中, 2.已加入, 2.被移出)
      }j
  }
  ```

#### 1.3.11 申请加入私密频道

若已经存在申请没被处理,则不做任何事

- 请求

  ```http
  POST /api/channels/:channelId/apply HTTP/1.1
  
  {
      "reason": "申请理由" 
  }
  ```
- 其中

  | 字段 | 说明   |
          |------| --- |
  | channelId | 频道ID |

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```

#### 1.3.12 退出频道

- 请求

  ```http
  POST /api/channels/:channelId/exit HTTP/1.1
  ```
- 其中

  | 字段 | 说明   |
          |------| --- |
  | channelId | 频道ID |

- ~~应答~~

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```

#### 1.3.13 管理移除成员

管理移除成员后成员在此频道处于`被移除`状态，可再次通过申请加入

- 请求

  ```http
  DELETE /api/channels/:channelId/users/:userId HTTP/1.1
  ```

- 其中

  | 字段        | 说明   |
        |------| --- |
  | channelId | 频道ID |
  | userId     | 用户ID |

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```

#### 1.3.14 修改频道分组的顺序

[1, 4, 3, 6] => id为1的频道分组顺序改为0，id为4的字段顺序改为1 ...

- 请求

    ```http
    PUT /api/channel-groups/sort HTTP/1.1
  
    {
        planetId: 1,
        sortedGroupIds: [1,4,3,6]
    }
    ```

- 应答gc
    ```js
    httpStatus = 200;
    httpRes = {
        code: 0,
        data: {}
    }
    ```

#### 1.3.15 获取当前用户在频道下的最后读取消息id

- 请求

    ```http
    GET /api/channels/:channelId/last-read HTTP/1.1
    ```

- 应答gc
    ```js
    httpStatus = 200;
    httpRes = {
        code: 0,
        data: {
            lastReadMessageId: 0
        }
    }
    ```

#### 1.3.16 修改分组名

- 请求

  ```http
  PUT /api/channels/groups/group-name HTTP/1.1

  {
    planetId: 1,
    channelGroupId: 37,
    name: "生活测试"
  } 
  ```
- 应答

  ```js
  //HTTP/1.1 200 OK
  {
      code: 0,
      data: null
  }
  ```

#### 1.3.17 修改频道名

- 请求

  ```http
  PUT /api/channels/channel-name HTTP/1.1

  {
     planetId: 1,
     channelId: 41,
     name: "吃饭测试"
  }
  ```
- 应答

  ```js
  //HTTP/1.1 200 OK
  {
      code: 0,
      data: null
  }
  ```

#### 1.3.18 获取所有频道下用户消息未读数量

只返回用户加入了的频道下未读消息条数

- 请求

  ```http
  GET /api/channels/unread-msg-num?planetId=1 HTTP/1.1
  ```
    - 应答

      ```js
      //HTTP/1.1 200 OK
      {
         code: 0,
         data: {
             list: [
                {
                   channelId: 42,
                   unreadNum: 627
                },
                ...
             ]
         }
      }
      ```

#### 1.3.19 修改频道公告

修改频道公告

- 请求

  ```http
  PUT /api/channels/:channelId/notice HTTP/1.1
  
  {
    "notice": "xxxxxxxxxxxx"
  }
  ```
- 应答
- 
  ```js
  //HTTP/1.1 200 OK
  {
     code: 0
  }
  ```
### 1.4 用户管理

#### 1.4.1 批量更改用户在星球的角色

需要超级管理员权限,只能将角色小于当前用户用户角色改为小于当前用户角色

- 请求

  ```http
  PUT /api/planets/1/members/role HTTP/1.1
  
  {
      "userIds": [1, 2, 3],
      "role": 1
  }
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```

- 角色列表

| 角色标识 | 名称    | 说明                  |
|------|-------|---------------------|
| 0    | 普通成员  | 无特殊权限               |
| 1    | 管理员   | 频道管理，用户黑名单          |
| 2    | 超级管理员 | 频道管理，用户黑名单,提升用户为管理员 |

#### 1.4.2 修改用户状态

需要管理员权限

- 请求

  ```http
  PUT /api/planets/:planetId/members/status HTTP/1.1
  
  {
      "userIds" : [1,2,3]
      "status": 1
  }
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```
- 用户状态标识

| 状态标识 | 名称    | 说明                  |
|------|-------|---------------------|
| 0    | 正常    | 用户状态一切正常            |
| 1    | 黑名单   | 黑名单用户限制某些操作         |

#### 1.4.3 读取候选匿名身份列表

- 请求

  ```http
  GET /api/anonymous-identities
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          list: [
              {
                  nickname: "随机昵称1",
                  avatar: "https://xxx/xxx",
              },
              {
                  nickname: "随即昵称2",
                  avatar: "https://xxx/xxx",
              },
          ]  
      }
  }
  ```

#### 1.4.4 获取用户个人信息

- 请求

  ```http
  GET /api/users/me/info HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  {
     code: 0,
     data: {
        id: 1,
        nickname: "哈哈哈哈哈",
        avatar: "https://openview-oss.oss-cn-chengdu.aliyuncs.com/aed-test/avatar/93.png",
        lastLoginAt: "xxxxx" //为空表示首次登录
        invitedBy: 23, //为空表示自然流量
        points: 80
     }
  }
  ```

#### 1.4.5 修改个人资料信息

昵称长度：八个汉字

- 请求

  ```http
  PUT /api/users/me/info HTTP/1.1
  
  {
      nickname: "xxxx",
      avatar:"https://sxxx",
  }
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```

#### 1.4.6 设置老板键

- 请求

  ```http
  PUT /api/users/me/boss-key HTTP/1.1
  
  {
      "bossKey": "Alt + F4"
  }
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```

#### 1.4.7 设置下班时间

- 请求

  ```http
  PUT /api/users/me/work-off-time HTTP/1.1
  
  {
      "time": "18:00:00",
  }
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
  }
  ```

#### 1.4.8 读取个人配置

siteSettings.type: default / office / custom, 取值非custom时，customTitle和customIcon值忽略

- 请求

  ```http
  GET /api/users/me/user-settings HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          bossKey: "xxx",
          endOffTime: "18:00:00",
          appearanceTheme: "dark",
          siteSettings: {
              type: "default",  // office、custom
              customTitle: "",
              customIcon: ""
          },
          monthlySalary: 10000,  // 月薪
          monthlyWorkingDays: 22 // 月工作日
      }
  }
  ```

#### 1.4.9 批量读取用户信息

- 请求

  ```http
  POST /api/users/list HTTP/1.1
  
  {
    "ids" :[1,3,4]
  }
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {
          list: [
              {
                  "id": 1,
                  "nickname": "xxx",
                  "avatar": "https://xxx/xxx",
              } 
          ]
      }
  }
  ```

#### 1.4.10 获取设置时间早于xx%

- 请求

  ```http
  GET /api/users/off-time-earlier?offTime=18:00:00  HTTP/1.1
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  resp = {
      code: 0,
      data: {   
          earlierThan: 83.12     // 83.12%
      }
  }
  ```

#### 1.4.11 上报在线状态

用户状态保活,需定时上报，15s/次

- 请求

  ```http
  GET /api/users/keepalive

  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  {
     code: 0
  }
  ```

#### 1.4.12 修改用户当前的用户配置

允许传多个字段，按需传递参数, 只会更新传入的参数配置，例如：只传appearanceTheme只会更新外观主题设置

- 请求

  ```http
  PUT /api/users/me/user-settings

  {
     appearanceTheme: "dark",
     siteSettings: {
         type: "custom",
         customIcon: "x",
         customTitle: "1"
     }
  }
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  { 
     code: 0
  }
  ```

#### 1.4.13 获取自定义多项配置

说明: 目前传入的参数只支持三类：appearanceTheme，siteSettings，分别获取不同配置的值，接口使用时根据返回数据中componentName判断使用，避免依赖传参数组中的顺序

- 请求

  ```http
  POST /api/users/me/simple-settings

  {
     componentNames: ["appearanceTheme", "siteSettings"]
  }
  ```

- 应答

  ```js
  //HTTP/1.1 200 OK
  {
    code: 0,
    data: {
      list: [
        {
          componentName: "appearanceTheme",
          componentSettings: "bright"
        },
        {
          componentName: "siteSettings",
          componentSettings: {
             type: "custom",
             customTitle: "1",
             customIcon: "x"
          }
        }
      ]
    }
  }
  ```

#### 1.4.14 设置用户平均月薪，月工作天数，下班时间

- 请求

    ```http
    PUT /api/users/me/work-settings
  
    {
      offWorkTime: "19:10:21",
      monthlySalary: 10000,
      monthlyWorkingDays: 22
    }
    ```

- 应答
    ```js
    res = {
        code: 0
    }
    ```

### 1.5 摸鱼功能

#### 1.5.1 获取所有辞职模板

- 请求

    ```http
    GET /api/resign/templates
    ```

- 应答
    ```js
    httpStatus = 200;
    httpRes = {
        code: 0,
        data: {
            list: [
                {
                   id: 3,
                   title: "辞职模板三",
                   content: "xxxxxxxxsaxx"
                }
            ]
        }
    }
    ```

#### 1.5.2 读取热门话题列表

- 请求

    ```http
    GET /api/issues/hot
    ```

- 应答
    ```js
    httpStatus = 200;
    httpRes = {
        code: 0,
        data: {
            list: [
                {
                   title: "中午吃什么",
                   content: "疯狂星期四疯狂推荐"
                },
                {
                   title: "摸鱼八卦公会",
                   content: "21人同时在线，21846条热门评论"
                }
            ]
        }
    }
    ```

### 1.6 系统接口

#### 1.6.1 获取 OSS Token

- 请求

  ```http
  GET /api/system/oss-token?ext=pdf HTTP/1.1
  ```
    - 其中

      | 字段  | 类型   | 说明                         |
                              |------|------|------------------------------------------| 
      | ext   | string | 上传文件的后缀；不指定则支持图片上传；v0.2增加； | 

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
      host: "https://moyu-chat.oss-cn-hangzhou.aliyuncs.com",
      accessId: "LTAI5tADMa68F6QymxHJ5zKq",
      signature: "/xVKjQFVuDue/iK/X066+RWNklY=",
      policy: "eyJleHBpcmF0aW9uIjoiMjAyMS0xMi0wNlQwODo0NDoyNloiLCJjb25kaXRpb25zIjpbWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJhZWQvIl1dfQ==",
      key: "moyu-dev/100/111112222",
    },
  };
  ```

    - 其中

      | 字段        | 类型     | 说明                                                                  |
                              |-----------|--------|---------------------------------------------------------------------|
      | host      | string | 阿里云 oss 域名                                                          |
      | accessId  | string | 阿里云 accessId                                                        |
      | signature | string | 签名                                                                  |
      | policy    | string | 用户表单上传的策略（Policy),用于验证请求的合法性。Policy 为一段经过 UTF-8 和 Base64 编码 JSON 文本 |
      | key       | string | 上传文件key                                                             |

#### 1.6.2 创建短链接

- 请求

  ```http
  POST /api/l/short-link? HTTP/1.1
  
  {
      "url": "https://xxx/index?invided=1",
  }
  ```

    - 其中

      | 字段 | 说明 | 
                                              | --- | --- |
      | url | 原始链接,可携带 query 参数 |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
      "link": "https://xxx/xxx"
    },
  };
  ```

#### 1.6.3 获取摸鱼表情

- 请求
    ```http
    GET /api/system/emoticons?group&sort=none HTTP/1.1
    ```
    - 其中
        - group： 分组，不传表示获取所有；
            - 1：图片表情；
            - 2：回复表情；
        - sort:
            - none 不另外排序，默认排序
- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
        "list": [{
            id: 1, //表情Id
            name: "大笑", //表情名称
            url: "https://xxx...", //表情图片地址
            keywords: "关键词1;关键词2",
        }]
    },
  };
  ```

### 1.7 内容审核

#### 1.7.1 文本内容审核

- 请求

  ```http
  POST /api/audit/text HTTP/1.1
  
  {
      "text": "xxxxx"
  } 
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
        "valid": true // 通过 
    },
  };
  ```

### 1.8 通知消息

#### 1.8.1 分页读取用户通知

- 请求

  ```http
  GET /api/notices?page=1&size=20 HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
        "list": [
            {
                id: 1,
                user: {  // 来源用户
                    id: 10038,
                    nickname: "AAAA野猩CTO",
                    avatar: "https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/3.png"
                },
                business: {  // 来源消息
                    id: 21,
                    channelId: 21, 
                    userId: 22, // 来源用户 ID
                    createdAt: "2022-09-22T11:11:27+08:00",
                    content: {
                        "type": "text",
                        "text": "@[22]北极熊出来玩"
                    }
                },
                type: "mention", // 消息类型 (mention: 提及, reference: 引用)
                createdAt: "2022-09-22T11:11:27+08:00",
                read: false // 是否已读
            },
            {
                id: 1,
                user: {  // 来源用户
                    id: 10038,
                    nickname: "AAAA野猩CTO",
                    avatar: "https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/3.png"
                },
                business: {  // 申请单信息
                    id: 2,
                    approvalType: "channel-join",   // 审批类型
                    businessId: 3343, // channel-join 类型的审批，businessId 为 channelId
                    createdAt: "",  // 创建时间
                    reason: "申请理由", 
                    userId: 24, 
                    state: 1, // (审批状态, 0. 待审核, 1. 审核通过, 2. 审核驳回)
                },
                type: "approval", // 消息类型 (mention: 提及, reference: 引用, 审批: 频道申请)
                createdAt: "2022-09-22T11:11:27+08:00",
                read: false // 是否已读
            },
            {
                id: 2,
                user: {  // 来源用户
                    id: 10038,
                    nickname: "AAAA野猩CTO",
                    avatar: "https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/3.png"
                },
                business: {  // 来源消息
                    id: 21,
                    userId 22, // 来源用户 ID
                    channelId: 21, 
                    createdAt: "2022-09-22T11:11:27+08:00",
                    content: {
                        "type": "text",
                        "reference": 145, // 引用消息
                        "text": "哈哈哈"
                    }
                },     
                type: "reference", // 消息类型 (mention: 提及, reference: 引用)
                createdAt: "2022-09-22T11:11:27+08:00",
                read: false // 是否已读
            }     
        ]
    },
  };
  ```

* 其中

| 消息类型             | 说明                     |
|------------------|------------------------|
| mention          | `business` 表示源消息实体     |
| reference        | `business` 表示源消息实体     |
| approval | `business` 表示 `审批单` 实体 |

* 其中

#### 1.8.2 读取未读通知数

- 请求

  ```http
  GET /api/notices/unread-count HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
        "count": 33
    },
  };
  ```

#### 1.8.3 读取单条通知

- 请求

  ```http
  GET /api/notices/:id HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
        id: 1,
        user: {  // 来源用户
            id: 10038,
            nickname: "AAAA野猩CTO",
            avatar: "https://moyu-chat.oss-cn-hangzhou.aliyuncs.com/default/avatar/3.png"
        },
        business: {
            id: 21,
            userId: 22,  
            channelId: 21, 
            createdAt: "2022-09-22T11:11:27+08:00",
            content: {
                "type": "text",
                "reference": 145, // 引用消息
                "text": "哈哈哈"
            }
        },
        type: "mention", // 消息类型 (mention: 提及, reference: 引用)
        createdAt: "2022-09-22T11:11:27+08:00",
        read: false // 是否已读
    }     
  };
  ```

* 状态码

| 状态码    | 说明    |
|--------|-------|
| 404003 | 通知不存在 |

#### 1.8.4 批量标记通知已读

- 请求

  ```http
  POST /api/notices/read HTTP/1.1
  
  {
      "ids": [1,3,4]
  }
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
  };
  ```

### 1.9 IM 消息

#### 1.9.1 批量读取 IM 消息

- 请求

  ```http
  POST /api/messages/get HTTP/1.1
  
  {
      "ids:" [1,2,3,4,5]   // 消息ID

  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
        list: [
            {
                id: 1,
                userId: 22,  
                channelId: 21, 
                createdAt: "2022-09-22T11:11:27+08:00",
                content: {
                    "type": "messageType",
                    "reference": 12    // 字段不存在则没有引用
                }
            },
            {
                id: 2,
                userId: 22,  
                channelId: 21, 
                createdAt: "2022-09-22T11:11:27+08:00",
                content: {
                    "type": "messageType",
                    "reference": 12    // 字段不存在则没有引用
                }
            }     
        ]
    },
  };
  ```

#### 1.9.2 点赞/取消点赞精华消息

若不是精华消息则无视

- 请求

  ```http
  POST /api/messages/:id/like HTTP/1.1
  
  {
    "isLike": true // 点赞/取消点赞
  }
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
  ```

#### 1.9.3 读取消息的点赞情况

- 请求

  ```http
  POST /api/messages/like HTTP/1.1
  
  {
    ids: [1,2,3,4,5] 
  }
  ```
- 其中
 
  | 字段 | 说明 |
  | ---  | ---  |
  | ids  | 消息 ID |

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: { 
      list: [
        {
          id: 123, //消息 ID
          like: 30, //点赞数
          isLike: true //是否被我点赞
        }
      ]
    }
  }
  ```
  
### 1.10 版本通知

#### 1.10.1 读取有效版本通知

读取到时间的最新通知

- 请求

  ```http
  GET /api/alerts/version HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
        id: 2,
        title: "新版本上线啦",   // 标题
        headImage: "https://xxx/xxx",  // 头图
        content: "<p>xxx<p>",   // 富文本内容
    },
  };
  ```

* 或者

不存在

  ```js
  // HTTP/1.1 200 OK

res = {
    code: 0,
    data: null
};
  ```
### 1.10 版本通知

#### 1.10.1 读取有效版本通知

读取到时间的最新通知

- 请求

  ```http
  GET /api/alerts/version HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
        id: 2,
        title: "新版本上线啦",   // 标题
        headImage: "https://xxx/xxx",  // 头图
        content: "<p>xxx<p>",   // 富文本内容
    },
  };
  ```

* 或者

不存在

  ```js
  // HTTP/1.1 200 OK

res = {
    code: 0,
    data: null
};
  ```

### 1.11 审批

#### 1.11.1 读取某个审批单

- 请求

  ```http
  GET /api/approvals/:id HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
        id: 2,
        approvalType: "channel-join",   // 审批类型
        businessId: 3343,
        createdAt: "",  // 创建时间
        reason: "申请理由", 
        userId: 24, 
        state: 1, // (审批状态, 0. 待审核, 1. 审核通过, 2. 审核驳回)
    },
  };
  ```

* 状态码

| 状态码    | 说明    |
|--------|-------|
| 404002 | 审批单不存在 |

#### 1.11.2 审批

- 请求

  ```http
  POST /api/approvals/:id/audit HTTP/1.1
  
  {
      "pass": true  // 是否通过
  }
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0
  };
  ```

### 1.11 摸鱼数据统计

#### 1.11.1 上报摸鱼时长

- 请求

  ```http
  POST /api/users/me/browse-duration  HTTP/1.1
  
  {
      timeQuantum: [
           {
               "startTime": "2022-11-09T07:40:16.787Z",
               "endTime": "2022-11-09T07:40:16.787Z"
           },
           ...
      ]
  }
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0
  }
  ```

#### 1.11.2 获取摸鱼数据详情

- 请求

  ```http
  GET /api/users/moyu-detail?userId=10047 HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK

  res = {
      code: 0,
      data: {
         joinDate: "2022-11-14T07:40:16.787Z", // date
         todayBrowseDuration: 678, // s
         totalBrowseDuration: 18223, // s
         moreThan: 88.3,  // %
         secondSalary: 0.5,   // 元， 小数点后端不做保留，前端使用时判断保留几位，有的秒薪算下来太小
         lastReportBrowseTime: "2022-11-03 12:44:22",   // 前段考虑是否需要这个字段
         accumulateMsgCnt: 872,
         user: {
            id: 10047,
            avatar: "https://xxa",
            nickname: "alice"
         }
      }
   }
  ```

#### 1.11.3 获取摸鱼时长排行榜

- 请求

  ```http
  GET /api/users/moyu-time-ranking?top=5&period=1  HTTP/1.1
  ```
  | 参数 | 说明                            | 
    |-------------------------------| --- |
  | top | 前n名                           |
  | period | 查询周期， 0/默认-total, 1-daily（今日） |

- 应答

  ```js
  // HTTP/1.1 200 OK

  res = {
       code: 0,
       data: {
           list: [{
              rank: 1      // 排名       
              user: {
                 id: 10023,
                 nickname: "xxa",
                 avatar: "https://ssda"
              },
              browseDuration: 566 // 秒   
           }, {
              rank: 2 // 排名  
              user: {
                 id: 10033,
                 nickame: "xxa",
                 avatar: "https://ssda"
              },
              browseDuration: 461
           }]
       }
    } 
  ```

#### 1.11.5 热门聊天榜(按频道)

- 请求
  ```http request
  GET /api/stat/hot-messages?top=5&channel=5
  ```

- 应答
     ```js
        // HTTP/1.1 200 OK
        res = {
            code: 0,
            data: {
                list: [{
                    msgId: 10, //消息Id
                    userId: 20, //用户Id
                    content: {
                        type: "text",
                        // ... 其他字段 
                        reply:[{ //表情回复内容
                            userId, 
                            emoticonId,
                        }]
                    },
                    replyCount: 100, //回复数
                }]
            },
        } 
     ```

### 1.13 机器人接入

#### 1.13.1 发送消息

- 请求

  ```http
  POST /api/robots/session-messages HTTP/1.1
  App-Id: {appId}
  Channel-Id: {channelId}
  Env: {Environment}
  
  {
      "type": "messageType",
      "reference": 12    // 字段不存在则没有引用
  }
  ```
- 消息类型参考 [概要设计](概要设计.md) 的 3.3 消息结构 部分

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0
  };
  ```
- 其中

  | 字段    | 说明                                       |
  |------------------------------------------| --- | 
  | appId | Moyu Server 提供给也业务方的 appId               |
  | Env     | 环境标识，moyu-server 会根据机器人配置选择处理或者忽略 |

#### 1.13.2 消息接收回调接口规范

此接口在业务方提供，由 Moyu Server 进行调用

- 请求

  ```http
  POST https://xxx/xxx HTTP/1.1
  App-Id: {appId}
  Token: {token}
  
  {
      "id": 1234,
      "userId": 21,
      "channelId": 21,
      "content": {
          "type": "messageType",
          "reference": 12    // 字段不存在则没有引用
      }
  }
  ```
- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0
  };
  ```

- 其中

  | 字段 | 说明 |
    | --- | --- | 
  | appId     | Moyu Server 提供给也业务方的 appId                 |
  | token | 使用 secret 对 appId 进行 SHA256 散列 后转换为 base64 |
 
### 1.14 积分流水

#### 1.14.1 创建积分流水 (管理员)

`points` > 0 表示增加，< 0 表示扣除 

- 请求

  ```http
  POST /api/points HTTP/1.1
  
  {
      "userId": 31, // 用户ID
      "points": 100, // 积分数
      "description": "空投奖励" // 描述
  }
  ```
  
- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0
  };
  ```
  
#### 1.14.2 读取用户积分流水

- 请求

  ```http
  GET /api/points?page=xxx&size=xxx HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0
    data: {
      total: 100, //结果总数
      list: [
        {
          userId: 31, // 用户ID
          points: 100, // 积分数
          description: "空投奖励" // 描述
          createdAt: "2022-09-22T11:11:27+08:00"
        } 
      ] 
    }
  };
  ```
#### 1.14.3 积分总数查询

查询用户当前积分总数

- 请求

  ```http
  GET /api/points/total HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
      total: 200, //积分总数
    },
  };
  ```

#### 1.14.4 用户积分排行榜 (日榜)

- 请求

  ```http
  GET /api/users/points-ranking/daily  HTTP/1.1
  ```
  
- 其中

  | 字段      | 说明                            | 
  |-------------------------------| --- |
  | top    | 前n名                           |

- 应答

  ```js
  // HTTP/1.1 200 OK

  res = {
    "code": 0,
    "data": {
      "list": [
        {
          "user": {
            "id": 234,
            "nickname": "xxx",
            "avatar": "http://xxx"
          },
          "points": 200
        } 
      ]
    }
  } 
  ```

#### 1.14.5 用户积分排行榜 (总榜)

- 请求

  ```http
  GET /api/users/points-ranking?top=xxx HTTP/1.1
  ```
  
- 其中
 
  | 字段      | 说明                            | 
  |-------------------------------| --- |
  | top    | 前n名                           |

- 应答

  ```js
  // HTTP/1.1 200 OK

  res = {
    "code": 0,
    "data": {
      "list": [
        {
          "user": {
            "id": 234,
            "nickname": "xxx",
            "avatar": "http://xxx"
          },
          "points": 200
        } 
      ]
    }
  } 
  ```
#### 1.15 系统配置

#### 1.15.1 读取某个键值对

- 请求

  ```http
  GET /api/system-configs?key=xxx HTTP/1.1
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0,
    data: {
      value: xxx
    },
  };
  ```

- 其中

  | 字段   | 说明   |
  |------| --- |
  | value    | 值，可以是任何数据类型 |

#### 1.15.2 设置键值对

内部 API， 不暴露公网，需要挂 vpn 访问

- 请求

  ```http
  PUT /inner/system-configs HTTP/1.1
  
  {
    "key": "xxx",
    "value": xxx
  }
  ```

- 应答

  ```js
  // HTTP/1.1 200 OK
  
  res = {
    code: 0
  };
  ```

- 其中

  | 字段   | 说明   |
  |------| --- |
  | value    | 值，可以是任何数据类型 |

