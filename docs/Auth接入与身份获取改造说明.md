# Auth 接入与身份获取改造说明

## 1. 背景

当前项目中，`userId` 和 `sitterId` 的获取方式仍然是开发态假实现：

- 创建订单相关接口继续使用前端传参
- 其余需要鉴权的接口，当前统一通过 `auth.GetUserID(c)` / `auth.GetSitterID(c)` 获取身份
- 这两个方法目前返回的是固定写死的测试 ID，并没有真正从 token 中解析

当前这样做是为了先把业务权限逻辑补齐，后续再统一替换为真实鉴权实现。

这份文档是给后续接手同学的改造说明，目标是把“假 auth”替换成“真实从 token 获取当前登录用户身份”。

## 2. 当前现状

### 2.1 当前 auth 实现

文件位置：

- [authorizer.go](file:///Users/bytedance/workplace/go/src/ph-order/util/auth/authorizer.go)

当前关键函数：

- `GetUserID(c)`：直接返回固定 user ID
- `GetSitterID(c)`：直接返回固定 sitter ID
- `GetUserIDFromToken(c)`：预留接口，但当前返回空字符串

也就是说，当前项目里并没有真正的 token 校验和身份注入逻辑。

### 2.2 当前路由现状

文件位置：

- [route.go](file:///Users/bytedance/workplace/go/src/ph-order/route.go)

当前 `gin` 路由没有统一挂载认证中间件，也没有把认证后的用户信息写入 `gin.Context`。

### 2.3 当前业务约定

这个约定非常重要，改造时不要改错：

- `创建订单`：
  - 暂时继续允许前端传 `userId` 和 `SitterID`
  - 不要求这次改造顺手改掉
- `其他需要鉴权的接口`：
  - 当前已经大量依赖 `auth.GetUserID(c)` / `auth.GetSitterID(c)`
  - 这次真正需要改造的是这两个方法背后的实现方式

换句话说，这次 auth 改造的核心目标不是“改业务接口签名”，而是“让现有接口里的 auth 调用真正可用”。

## 3. 改造目标

改造完成后，应满足以下目标：

1. 后端能从请求头中的 token 识别当前登录用户
2. `auth.GetUserID(c)` 能返回当前登录用户的真实 user ID
3. `auth.GetSitterID(c)` 能返回当前登录用户的真实 sitter ID
4. service 层不再依赖固定测试 ID
5. 尽量不改动已经写好的业务权限判断逻辑

## 4. 推荐实现原则

### 4.1 不要让业务层自己解析 token

不要在每个 service 里重复读取：

- `Authorization`
- `Bearer xxx`
- JWT Claims

正确做法是：

- 在中间件里统一解析 token
- 把解析后的身份写入 `gin.Context`
- `auth.GetUserID(c)` / `auth.GetSitterID(c)` 只负责从 `Context` 读取并返回

### 4.2 统一“当前登录身份”的来源

当前 owner/sitter 在业务上是两个角色，但本质上都应该来自同一个认证主体。

推荐统一使用 token 中的唯一身份字段，例如：

- `sub`
- `cognito:username`
- `user_id`
- 你们认证系统中约定的主键字段

最终使用哪个字段，要以认证平台的实际规范为准。

### 4.3 不要先改业务逻辑，再想 auth

当前很多接口已经写成了：

- owner 侧：`auth.GetUserID(c)` 对比 `order.OwnerID`
- sitter 侧：`auth.GetSitterID(c)` 对比 `order.SitterID`

最稳妥的方式是：

- 保持 service 层调用方式不变
- 只替换 auth 底层实现

这样改动面最小，回归风险最低。

## 5. 建议改造方案

## 5.1 第一步：新增统一认证中间件

建议新增一个中间件文件，例如：

- `util/auth/middleware.go`

中间件职责：

1. 读取请求头 `Authorization`
2. 校验格式是否为 `Bearer <token>`
3. 调用认证组件解析 token
4. 从 token claims 中提取当前登录主体 ID
5. 将认证结果写入 `gin.Context`

建议写入的上下文字段：

- `principal_id`
- `principal_claims`
- 如果你们有角色信息，也可以补：
  - `principal_role`

示例伪代码：

```go
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "missing authorization header"})
			return
		}

		token := parseBearerToken(authHeader)
		claims, err := verifyAndParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "invalid token"})
			return
		}

		principalID := extractPrincipalID(claims)
		if principalID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "principal id not found"})
			return
		}

		c.Set("principal_id", principalID)
		c.Set("principal_claims", claims)
		c.Next()
	}
}
```

## 5.2 第二步：把 auth 包改成从 Context 读取

改造文件：

- [authorizer.go](file:///Users/bytedance/workplace/go/src/ph-order/util/auth/authorizer.go)

建议目标：

- `GetUserID(c)`：从 `c.Get("principal_id")` 读取
- `GetSitterID(c)`：从 `c.Get("principal_id")` 读取
- 如果未来 user/sitter 真的是两套不同身份字段，再扩展；当前先保持同源即可

建议改造后语义：

- `GetUserID(c)`：返回“当前登录者的主体 ID，用于 owner 侧校验”
- `GetSitterID(c)`：返回“当前登录者的主体 ID，用于 sitter 侧校验”

注意：

- 当前这两个函数不应该继续返回写死值
- 也不建议在里面直接解析 token，应该只从 `Context` 读

示例伪代码：

```go
func GetUserID(c *gin.Context) (string, error) {
	v, ok := c.Get("principal_id")
	if !ok {
		return "", errors.New("principal_id not found in context")
	}
	id, ok := v.(string)
	if !ok || id == "" {
		return "", errors.New("invalid principal_id")
	}
	return id, nil
}

func GetSitterID(c *gin.Context) (string, error) {
	v, ok := c.Get("principal_id")
	if !ok {
		return "", errors.New("principal_id not found in context")
	}
	id, ok := v.(string)
	if !ok || id == "" {
		return "", errors.New("invalid principal_id")
	}
	return id, nil
}
```

## 5.3 第三步：在路由层挂载中间件

改造文件：

- [route.go](file:///Users/bytedance/workplace/go/src/ph-order/route.go)

建议做法：

- 对需要鉴权的接口组挂载 `AuthMiddleware`
- 如果短期内无法区分公共接口和私有接口，也可以先全局挂载，再为少量匿名接口放行

建议优先挂载到这些接口：

- 订单查询我的订单
- 订单修改
- 订单修改确认
- 评价创建 / 更新 / 删除 / 查询
- 平均评分查询

需要注意的例外：

- 创建订单相关接口，当前项目明确允许参数传入，不一定必须在这一轮强制改

## 5.4 第四步：确认 token claim 与数据库字段一致

这是最容易踩坑的地方。

必须确认：

- token 中解析出来的主体 ID
- `orders.owner_id`
- `orders.sitter_id`
- `sitter_rating.user_id`
- `sitter_rating.sitter_id`
- `owner_rating.owner_id`
- `owner_rating.sitter_id`

这些字段存的是不是同一套用户主键。

如果不是同一套，就会出现：

- token 解析没问题
- 但权限判断永远不匹配

如果字段体系不一致，需要先定义一层映射关系，而不是直接替换。

## 6. 具体文件改造清单

建议至少检查以下文件：

### 必改

- [authorizer.go](file:///Users/bytedance/workplace/go/src/ph-order/util/auth/authorizer.go)
- [route.go](file:///Users/bytedance/workplace/go/src/ph-order/route.go)

### 建议联调核验

- [order_modification_log_service.go](file:///Users/bytedance/workplace/go/src/ph-order/service/order_modification_log_service.go)
- [order_rating.go](file:///Users/bytedance/workplace/go/src/ph-order/service/order_rating.go)
- [order_service.go](file:///Users/bytedance/workplace/go/src/ph-order/service/order_service.go)
- [order_create_with_pricing.go](file:///Users/bytedance/workplace/go/src/ph-order/service/order_create_with_pricing.go)

## 7. 推荐改造顺序

建议按下面顺序做，避免一下子改太多难以定位问题：

1. 实现 token 解析函数
2. 实现 `AuthMiddleware`
3. 改造 `auth.GetUserID` / `auth.GetSitterID`
4. 在少量接口上挂中间件进行验证
5. 验证 owner / sitter 权限判断通过
6. 再逐步扩展到全部私有接口

## 8. 验证清单

改造完成后，至少验证这些场景：

### owner 场景

- owner 可以创建 sitter 评价
- owner 可以更新自己的 sitter 评价
- owner 可以删除自己的 sitter 评价
- owner 可以查看自己订单的评价详情
- owner 不能确认 sitter 角色专属接口

### sitter 场景

- sitter 可以创建 owner 评价
- sitter 可以查看自己参与订单的评价
- sitter 不能更新 owner 创建的 sitter 评价
- sitter 不能删除 owner 创建的 sitter 评价

### 订单修改场景

- owner 只能确认 sitter 发起的修改
- sitter 只能确认 owner 发起的修改
- 非订单参与方无法确认修改

### 评价状态场景

- 创建评价后，订单接口返回的 `UserRatingState` / `SitterRatingState` 正确变为 `1`
- 删除评价后，订单接口返回的 `UserRatingState` / `SitterRatingState` 正确恢复为 `0`

## 9. 已知注意事项

### 9.1 创建订单接口暂不统一改造

当前项目约定是：

- 创建订单暂时继续允许前端传 `userId` / `SitterID`

所以后续接手时，不要误以为这次 auth 改造必须连创建订单一起重构。

### 9.2 目前有些“我的订单”接口还在读 Query 参数

例如 `UserGetMyOrdersPage` / `SitterGetMyOrdersPage` 当前仍然通过 `Query("UserID")` 获取当前用户身份。

这部分建议后续一起收口到：

- `auth.GetUserID(c)`
- `auth.GetSitterID(c)`

否则整体权限模型会出现“部分接口走真实 auth，部分接口仍信任 query 参数”的不一致问题。

### 9.3 不要把真实鉴权逻辑散落在 service 层

如果把 token 解析写到每个 handler 里，后续很快会失控。

正确边界应当是：

- middleware：校验 token、解析 claims、写 context
- auth 包：从 context 取当前身份
- service：只做业务权限判断

## 10. 一句话总结

后续真正要做的不是“重写业务权限逻辑”，而是：

- 用认证中间件把当前登录用户身份放进 `gin.Context`
- 再把 `auth.GetUserID(c)` / `auth.GetSitterID(c)` 从“返回固定值”替换成“从 context 读取真实主体 ID”

这样就能在尽量不改 service 业务逻辑的前提下，把整套权限判断切换到真实认证之上。
