# 简化版 mysqldump 工具设计规范 (MVP版本)

## 概述

本项目旨在实现一个简化版的 mysqldump 工具，作为最小可行产品(MVP)专注于将 MySQL 数据库表数据导出为 SQL 格式。该 MVP 版本只包含最基本但可用的功能，以便快速验证核心概念和价值。

## 功能需求

### 核心功能（MVP范围）
1. 连接到 MySQL 数据库
2. 选择特定的数据库和表
3. 通过 DDL 分析过滤掉虚拟列
4. 导出表数据为 INSERT 语句（包含具体列名）
5. 输出的 SQL 包含基本的环境设置
6. 输出到标准输出

### 非MVP功能（未来扩展）
1. 支持 WHERE 条件过滤数据
2. 输出到文件
3. 支持导出多个表
4. 支持压缩输出

## 技术选型

- 编程语言：Go
- 数据库驱动：[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- 配置管理：环境变量
- 日志记录：标准日志包

## MVP 架构设计

作为最小可行产品，我们将采用简化的设计方法，将核心功能组织在清晰的模块中：

### 核心组件

#### 配置管理模块
- 从环境变量读取所有配置信息
- 提供配置参数的默认值

#### 数据库连接模块
- 负责建立和管理数据库连接
- 执行基本的 SQL 查询获取表数据

#### 表结构分析模块
- 获取表的 DDL 信息
- 分析并过滤掉虚拟列（包含 VIRTUAL 关键字的列）
- 确定需要导出的实际列列表

#### 数据导出模块
- 将查询结果转换为 SQL INSERT 语句（包含具体列名）
- 处理数据转义和格式化
- 在输出中包含基本的环境设置（如 SET 语句）

#### 主程序模块
- 协调各模块工作流程
- 处理命令行参数解析
- 控制整体执行流程

## 详细设计

### 目录结构
```
mysqldump/
├── cmd/
│   └── dump/
│       └── main.go          # 应用入口
├── internal/
│   ├── config/
│   │   └── config.go        # 配置管理
│   ├── db/
│   │   └── connection.go    # 数据库连接管理
│   ├── schema/
│   │   └── analyzer.go      # 表结构分析
│   ├── exporter/
│   │   └── data.go          # 数据导出逻辑
│   └── core/
│       └── dump.go          # 核心导出流程
├── go.mod
└── go.sum
```

### 核心模块说明

#### 配置管理模块 (internal/config/config.go)
- 所有配置均通过环境变量加载
- 提供默认值处理（如 DB_HOST 默认为 localhost）
- 验证必需配置项（如 DB_NAME）

#### 数据库连接模块 (internal/db/connection.go)
- 建立和关闭数据库连接
- 执行 SQL 查询获取表数据和表结构信息

#### 表结构分析模块 (internal/schema/analyzer.go)
- 查询 information_schema.COLUMNS 表获取列信息，使用以下 SQL：
  SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS 
  WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
- 分析 EXTRA 字段，识别包含 'GENERATED' 关键字的列为生成列（包括 VIRTUAL 和 STORED 类型）
- 构建排除生成列的实际列列表，用于后续的数据查询和 INSERT 语句生成
- 提供获取实际列列表的接口方法

#### 数据导出模块 (internal/exporter/data.go)
- 根据分析得到的实际列列表查询数据
- 将查询结果转换为包含具体列名的 INSERT 语句
- 在输出中包含基本的环境设置（如 SET 语句）
- 处理数据转义和格式化

#### 主程序模块 (cmd/dump/main.go)
- 解析命令行参数
- 初始化各功能模块
- 协调执行导出流程
- 处理异常情况

## 配置管理

所有配置均通过环境变量加载，不使用命令行参数或配置文件：
- DB_HOST: 数据库主机地址（默认: localhost）
- DB_PORT: 数据库端口（默认: 3306）
- DB_USER: 数据库用户名（默认: root）
- DB_PASSWORD: 数据库密码（默认: 空）
- DB_NAME: 数据库名称（必需）

## 命令行参数

```
Usage: dump [options] table

Options:
  -h, --help     显示帮助信息

Arguments:
  table          要导出数据的表名
```

## 虚拟列处理

工具需要能够识别并排除虚拟列：
1. 通过查询 information_schema.COLUMNS 表获取表的列信息
2. 分析 EXTRA 字段，识别包含 'GENERATED' 关键字的列作为生成列（包括VIRTUAL和STORED类型）
3. 构建排除生成列的列列表
4. 在查询和 INSERT 语句中只包含实际存储的列

## SQL 输出格式

导出的 SQL 应包含基本的环境设置，参考 mysqldump 的格式：
```
-- MySQL dump 10.13  Distrib 8.0.x, for Linux (x86_64)
--
-- Host: localhost    Database: test
-- ------------------------------------------------------
-- Server version	8.0.x

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` (`id`, `name`) VALUES (1,'John Doe');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
```

## 数据流设计

1. 主程序解析命令行参数获取表名
2. 从环境变量读取数据库配置
3. 建立数据库连接
4. 获取表的 DDL 信息并分析列结构
5. 过滤掉虚拟列，确定实际导出列列表
6. 查询指定表的数据（仅包含实际列）
7. 将数据转换为包含具体列名的 INSERT 语句
8. 在输出中添加基本的环境设置
9. 输出到标准输出

## 错误处理

MVP 版本将包含基本的错误处理：
- 配置加载错误（缺少必需配置）
- 数据库连接错误
- SQL 执行错误
- DDL 解析错误
- 参数解析错误

## 测试策略

### 单元测试
- 配置管理模块测试（使用实际环境变量）
- 表结构分析模块测试（使用模拟 DDL 数据）
- 数据导出模块测试（使用模拟查询结果）

### 集成测试
- 完整导出流程测试

## 性能考虑

- 使用流式处理避免内存溢出
- 逐行处理和输出数据

## 安全考虑

- 密码等敏感信息不记录日志
- SQL 注入防护

## MVP交付标准

MVP版本应满足以下标准：
1. 能够成功连接到 MySQL 数据库
2. 能够正确识别并排除虚拟列
3. 能够导出指定表的数据为包含具体列名的 INSERT 语句
4. 输出的 SQL 包含基本的环境设置
5. 输出格式正确，可被 MySQL 直接导入
6. 提供基本的错误处理和提示信息
7. 有基本的单元测试覆盖

## 未来扩展

在MVP基础上，后续可以扩展以下功能：
- 支持 WHERE 条件过滤数据
- 支持输出到文件
- 支持导出多个表
- 支持其他数据库类型
- 支持压缩格式输出