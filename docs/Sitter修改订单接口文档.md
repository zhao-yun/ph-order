# Sitter 订单修改接口文档

本文档说明了 Sitter（保姆）端发起订单修改（针对单个宠物调整价格或修改订单结束日期）的接口调用规范。

---

## 1. Sitter 申请修改订单

Sitter 可以通过此接口提交订单的修改申请。目前支持修改订单的结束日期 (`ToDate`) 以及调整每个宠物的服务价格 (`PetPrice`)。

> **业务逻辑说明**：
> Sitter 提交修改后，订单状态不会立刻生效。系统会生成一条 `State=0`（初始化待确认）的订单修改记录。需要用户端确认后（通过另外的确认接口），修改才会真正生效。

### 1.1 基本信息

- **接口路径**: `POST /Order/SitterUpdateOrder`
- **请求方式**: `POST`
- **Content-Type**: `application/json`

### 1.2 请求参数 (Request Body)

| 字段名    | 类型      | 必填 | 说明                                                              |
| :-------- | :-------- | :--- | :---------------------------------------------------------------- |
| `OrderID` | integer   | 是   | 需要修改的订单 ID                                                 |
| `ToDate`  | string    | 否   | 新的订单结束日期。格式为 `YYYY-MM-DDTHH:mm:ssZ`。如果不修改日期，请传入原订单的结束日期。 |
| `PetList` | array     | 否   | 宠物列表。支持传入修改后的宠物价格 `PetPrice`。                   |

#### `PetList` 数组元素说明:
| 字段名     | 类型    | 必填 | 说明                                           |
| :--------- | :------ | :--- | :--------------------------------------------- |
| `PetID`    | string  | 是   | 宠物 ID (唯一标识)                             |
| `PetPrice` | float64 | 是   | 该宠物修改后的最新价格 (如 `50.00`)            |
| 其他字段   | -       | 否   | 宠物原本的其他信息可以一并带上，目前主要校验价格。 |

### 1.3 请求示例

```json
{
  "OrderID": 123,
  "ToDate": "2026-05-01T00:00:00Z",
  "PetList": [
    {
      "PetID": "pet_001",
      "PetPrice": 50.00
    },
    {
      "PetID": "pet_002",
      "PetPrice": 40.00
    }
  ]
}
```

### 1.4 响应参数

成功响应（HTTP 200）的 `data` 字段中会包含两部分：
1. `log`: 数据库生成的修改记录明细。
2. `priceChanges`: 后端自动计算的本次价格变动详情。

#### `data.priceChanges` 字段说明:
| 字段名          | 类型    | 说明                                                    |
| :-------------- | :------ | :------------------------------------------------------ |
| `previousPrice` | float64 | 修改前的订单总价（包含服务费、税费等）                  |
| `newPrice`      | float64 | 根据新宠物价格重新计算后的订单总价（包含服务费、税费等）|
| `difference`    | float64 | 差价 = `newPrice` - `previousPrice`                     |

### 1.5 响应示例

**成功响应:**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "log": {
      "ID": 10,
      "OrderID": 123,
      "OwnerID": "user_888",
      "SitterID": "sitter_999",
      "PreviousDate": "2026-04-30T00:00:00Z",
      "NewDate": "2026-05-01T00:00:00Z",
      "PreviousPetList": "[{\"PetID\":\"pet_001\",\"PetPrice\":40},{\"PetID\":\"pet_002\",\"PetPrice\":30}]",
      "NewPetList": "[{\"PetID\":\"pet_001\",\"PetPrice\":50},{\"PetID\":\"pet_002\",\"PetPrice\":40}]",
      "PreviousPrice": 80.00,
      "NewPrice": 100.00,
      "State": 0,
      "Type": 2,
      "CreatedAt": "2026-04-05T12:00:00Z",
      "UpdatedAt": "2026-04-05T12:00:00Z"
    },
    "priceChanges": {
      "previousPrice": 80.00,
      "newPrice": 100.00,
      "difference": 20.00
    }
  }
}
```

**失败响应 (订单不存在):**
```json
{
  "code": 500,
  "msg": "get order by id failed, err = record not found",
  "data": null
}
```

---

## 2. 前端调用建议

1. **界面交互**: 当 Sitter 想要修改订单价格时，界面应当列出订单内包含的每一个宠物，并允许对每个宠物单独输入新的价格。
2. **确认弹窗**: 前端调用本接口后，可以通过读取 `data.priceChanges` 拿到实时的价格变动情况。App 可以在收到成功响应后，直接弹窗提示 Sitter：*“您已成功发起价格修改申请。原订单总价 $80.00，新订单总价 $100.00，差额 $20.00。请等待用户确认。”*