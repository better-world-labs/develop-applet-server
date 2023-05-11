# 做个小程序

设计文档

## APP 模板设计

### 1.基本信息

APP 的基本信息，例如名称，描述等

### 2.表单信息

表单定义了应用的输入参数与类型，定义在 APP 的 `form` 字段，使用 uuid 索引标识,支持以下类型

- 文本框

  ```Json
  {
      "id": "uuid",
      "label": "姓名",
      "type": "text",
      "properties": {
          "placeholder": "请输入姓名"
      }
  }
  ```

- 单选框

  ```Json
  {
      "id": "uuid",
      "label": "性别",
      "type": "select",
      "properties": {
          "placeholder": "请选择性别",
          "values": "男\n女"
      }
  }
  ```

- 多选框

  ```Json
  {
      "id": "uuid",
      "label": "分类",
      "type": "checkbox",
      "properties": {
          "placeholder": "请勾选分类",
          "values": "男\n女"
      }
  }
  ```

- 其中

  | 字段 | 说明 |
        | --- | --- |
  | uuid | 元素唯一标识 |
  | label | 表单元素 label|
  | type | 表单元素类型，支持 text 和 select |
  | properties | 不同表单类型的特有属性 |

- 运行传参

在运行应用时，需要按照表单类型和表单顺序分别传递参数，说明如下

| 类型 | 传参说明 | 传参示例|
  | --- | --- | --- |
| text | 文本 | 小明 |
| select | `values` 枚举的值列表中包含的值 | "男" |
| checkbox | `values` 枚举的值列表中包含的值集合，回车分隔 | "男\n女" |

### 3.逻辑流

逻辑流定义了程序执行逻辑，定义在 APP 的 `flow` 字段, 每个 flow 使用 `uuid` 唯一标识

#### 3.1 流程基本信息

定义了流程使用的模型，输出是否可见，提示词列表等信息

#### 3.2 提示词列表

流程的提示词列表相当于程序编写的逻辑部分,是由文本和标签(Tag)拼接而成的数组，其中标签可以是来自 `form` 的表单的`id`，也可以是来自 `flow` 的输出，类型定义如下

- 文本

  ```Json
  {
      "type": "text",
      "properties": {
          "value": "文本内容"
      }
  }
  ```

- 标签

  ```Json
  {
      "type": "tag",
      "properties": {
          "from": "result",
          "character": "uuid" //可以是来自 form 的id,也可以是来自 flow 的id
      }
  }
  ```
- 其中

  | 字段|说明 |
        | --- | --- |
  | from | result/form 表示来自表单或者执行结果 |
  | character| 引用 form 或者 flow 的ID|

