# 支付与订单核心流程接口文档（App 端专用）

本接口文档说明了用户在 App 端从“发起支付”到“Sitter 完成订单”的完整闭环流程。

---

## 一、接口详情

### 1. 发起支付授权 (App 唤起 Stripe 前调用)

当订单创建成功（状态为 `0`：初始化）后，用户点击支付时调用此接口，获取 Stripe 的 `clientSecret` 用于唤起客户端收银台。

- **接口地址**: `POST /Order/Pay`
- **请求参数**:

| 字段名    | 类型    | 必填 | 描述             |
| --------- | ------- | ---- | ---------------- |
| `OrderID` | integer | 是   | 订单 ID (例如 123) |

- **请求示例**:
```json
{
  "OrderID": 123
}
```

- **响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "clientSecret": "pi_3MtwBwLkdIwHu7ix28a3tqPc_secret_Yk3x...",
    "paymentIntentId": "pi_3MtwBwLkdIwHu7ix28a3tqPc"
  }
}
```
> **注意**：后端会在此接口自动将 `paymentIntentId` 记录到 Redis，方便订单完成时自动扣款。

---

### 2. 确认支付结果 (App 端支付成功后调用)

**【新增接口】** App 客户端通过 Stripe SDK 确认用户已完成本地授权/支付后，**必须**调用此接口通知服务器，将订单状态变更为已支付（状态 `1`）。

- **接口地址**: `POST /Order/ConfirmPayment`
- **请求参数**:

| 字段名    | 类型    | 必填 | 描述             |
| --------- | ------- | ---- | ---------------- |
| `OrderID` | integer | 是   | 订单 ID |

- **请求示例**:
```json
{
  "OrderID": 123
}
```

- **响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

---

### 3. Sitter 接受/拒绝订单

Sitter 只能在订单状态为已支付（`1`）时才能接受订单。

- **接口地址**: `POST /Order/SitterHandleInvite`
- **请求参数**:

| 字段名    | 类型    | 必填 | 描述             |
| --------- | ------- | ---- | ---------------- |
| `OrderID` | integer | 是   | 订单 ID |
| `State`   | integer | 是   | `2`: 接受, `-1`: 拒绝 |

- **请求示例**:
```json
{
  "OrderID": 123,
  "State": 2
}
```
> **注意**：如果订单未支付（比如未调用 `ConfirmPayment`），此处会返回 HTTP 400 `cannot accept unpaid order`。如果接受成功，非遛狗订单会在后台生成 StartCode (验证码)。

---

### 4. 订单完成与自动扣款 (Sitter 输入结束验证码)

当服务结束，Sitter 输入结束验证码。**此接口一旦验证通过，系统将在后台自动向 Stripe 发起 `Capture` (真实扣款) 请求。**

- **接口地址**: `POST /Order/SitterSetFinishCode`
- **请求参数**:

| 字段名    | 类型    | 必填 | 描述             |
| --------- | ------- | ---- | ---------------- |
| `OrderID` | integer | 是   | 订单 ID |
| `Code`    | string  | 是   | 4位结束验证码 |

- **请求示例**:
```json
{
  "OrderID": 123,
  "Code": "A1B2"
}
```

---

## 二、前端 (App 端) 调用流程规范指南

为了确保资金安全和用户体验，App 端的研发请严格按照以下时序逻辑进行开发：

### 阶段一：用户支付流程 (User 侧)

1. **获取待支付订单**: 确保你有一个状态为 `0`（`OrderInitialized`）的 `OrderID`。
2. **请求 Backend 获取凭证**: 
   - 调用 `POST /Order/Pay`。
   - 解析返回的 `clientSecret`。
3. **唤起 Stripe SDK**:
   - 使用上述的 `clientSecret` 初始化并唤起 Stripe Payment Sheet（客户端收银台）。
   - 用户在客户端完成绑卡/确认支付（此时资金只是被**授权冻结**，并未划走）。
4. **回传支付结果给 Backend**:
   - 当 Stripe SDK 返回 `PaymentSheetResult.Completed`（支付成功回调）时。
   - **必须立即调用** `POST /Order/ConfirmPayment`。
   - 收到后端 `code: 200` 后，才向用户展示“支付成功”页面，并将本地列表中的订单状态置为“等待保姆接单”。

### 阶段二：保姆接单流程 (Sitter 侧)

1. **接单校验**:
   - Sitter 在 App 看到新订单，点击“接受”。
   - 调用 `POST /Order/SitterHandleInvite`，传入 `State: 2`。
   - *（如果用户其实没付钱，该接口会拦截并报错，前端应提示“该订单用户尚未完成付款，无法接单”。）*
2. **正常服务流程**:
   - 输入开始验证码 (`POST /Order/SitterSetCreateCode`) -> 状态变为进行中。
   - 点击服务完成 (`POST /Order/SitterFinishOrder`) -> 状态变为待确认结束，生成结束验证码。
3. **核销与真实结算**:
   - Sitter 获取到用户的结束验证码，输入 App。
   - 调用 `POST /Order/SitterSetFinishCode`。
   - **注意**：由于该接口内部包含了与 Stripe 真实的 `Capture`（网络扣款请求），**响应时间可能会比平时长 1~2 秒**。App 端在此步骤**必须加上 Loading 动画 (菊花图)** 且**禁用按钮连点**，防止重复提交。
   - 接口返回成功后，提示保姆“订单已完成，资金即将入账”。