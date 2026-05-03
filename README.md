# Pet Services Backend API (Demo)

本项目是一个基于 Go 语言 (Gin + GORM) 开发的宠物服务（代遛、寄养等）后端 API 平台。

## 环境依赖

在启动项目之前，请确保您的本地或服务器环境中已安装并启动以下服务：

- **Go**: 1.18 或更高版本
- **Redis**: 默认连接 `localhost:6379`（用于验证码、支付凭证缓存等）
- **MySQL**: 默认数据库（可通过配置文件修改）

## 如何启动项目

1. **获取代码并进入项目目录**
   在终端中进入项目所在的根目录（包含 `main.go` 的目录）。

2. **配置数据库**
   打开 `conf/db.yml` 文件，配置您的数据库连接信息：
   ```yaml
   host: 您的数据库地址
   port: 3306
   username: 您的数据库用户名
   password: 您的数据库密码
   dbName: 您的数据库名称
   ```

3. **下载依赖**
   ```bash
   go mod tidy
   ```

4. **启动服务**
   项目包含多个入口文件，推荐直接运行整个目录：
   ```bash
   go run .
   ```
   *或者指定文件运行：*
   ```bash
   go run main.go route.go
   ```
   
   启动成功后，服务默认监听在 `http://localhost:8000`。

---

## 如何将数据库从 MySQL 切换为 PostgreSQL

项目中已经预置了 PostgreSQL 的支持代码。由于 GORM 不能同时激活同名的初始化逻辑，目前项目中激活的是 MySQL (`util/postgres/mysql.go`)，而 PostgreSQL 的代码 (`util/postgres/postgres.go`) 被注释掉了。

> **⚠️ 数据库初始化脚本说明**
> 项目根目录提供了两份初始化表结构的 SQL 文件：
> - **如果您使用 MySQL**，请在数据库中执行：[`ph_orders.sql`](./ph_orders.sql)
> - **如果您使用 PostgreSQL**，请在数据库中执行：[`ph_orders_pg.sql`](./ph_orders_pg.sql)

若要切换到 PostgreSQL，请按照以下步骤操作：

### 第一步：切换代码

1. **禁用 MySQL 初始化代码**：
   打开 `util/postgres/mysql.go` 文件，将该文件中的**所有代码注释掉**（或者直接将该文件重命名为 `mysql.go.bak`）。
   
2. **启用 PostgreSQL 初始化代码**：
   打开 `util/postgres/postgres.go` 文件，**取消该文件中所有代码的注释**。
   *(注意：取消注释后，请检查代码顶部的 `package postgres` 和 import 是否正确)*。

### 第二步：下载 PostgreSQL 驱动

确保您的 `go.mod` 中包含 GORM 的 Postgres 驱动：
```bash
go get gorm.io/driver/postgres
go mod tidy
```

### 第三步：修改配置文件

打开 `conf/db.yml`，将配置修改为您 PostgreSQL 的连接信息（注意端口号的变化）：
```yaml
host: 您的PG数据库地址 (如 localhost)
port: 5432
username: postgres
password: 您的PG数据库密码
dbName: 您的PG数据库名称
```

### 第四步：重新启动项目

配置完成后，重新运行项目即可连接到 PostgreSQL 数据库：
```bash
go run .
```

> **注意**：由于 PostgreSQL 的 DSN 格式与 MySQL 不同，启用 `util/postgres/postgres.go` 后，系统会自动使用 `host=%s port=%s user=%s password=%s dbname=%s sslmode=disable` 的格式进行连接，无需您手动修改拼接逻辑。
