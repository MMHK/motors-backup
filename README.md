# motors-backup

A simple MySQL database backup tool written in Go.

一个用 Go 语言编写的简单 MySQL 数据库备份工具。

[![GitHub](https://img.shields.io/github/license/MMHK/motors-backup)](https://github.com/MMHK/motors-backup)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/MMHK/motors-backup)](https://github.com/MMHK/motors-backup)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/MMHK/motors-backup)](https://github.com/MMHK/motors-backup/releases)
[![Docker Image](https://img.shields.io/docker/v/mmhk/motors-backup)](https://hub.docker.com/r/mmhk/motors-backup)
[![Docker Pulls](https://img.shields.io/docker/pulls/mmhk/motors-backup)](https://hub.docker.com/r/mmhk/motors-backup)
[![GitHub stars](https://img.shields.io/github/stars/MMHK/motors-backup)](https://github.com/MMHK/motors-backup)

## Features | 功能特性

- Export all tables or specific tables from a MySQL database 导出 MySQL 数据库中的所有表或特定表
- Option to include or exclude `CREATE DATABASE` statement 可选择包含或排除 `CREATE DATABASE` 语句
- Ignore specific tables completely or ignore only their data 完全忽略特定表或仅忽略其数据
- Apply WHERE conditions to table exports 对表导出应用 WHERE 条件
- Clean and readable SQL output 清晰易读的 SQL 输出
- Requires MySQL version >= 8.0 要求 MySQL 版本 >= 8.0
- Automatically excludes generated column data from exports 自动排除导出中的生成列数据


## Usage | 使用方法

```shell
Usage: motors-backup [options] table

Options:
  --create-database             Include CREATE DATABASE statement (default true)
  -h, --help                    Display help information
  --ignore-table string         Tables to ignore completely
  --ignore-table-data string    Tables to ignore data only
  -w, --where string            WHERE conditions for tables (format: id>100) (default "")

Arguments:
  table          Table name(s) to export data from, multiple names separated by commas
```

#### Environment Variables

> Database connection configuration can only be set through environment variables, not supported via command line arguments.
>
> 数据库连接配置只能通过环境变量进行设置，不支持通过命令行参数指定。

```dotenv
# Database connection 数据库连接
DB_HOST=
# Database port 数据库端口
DB_PORT=3306
# Database user 数据库用户
DB_USER=
# Database password 数据库密码
DB_PASSWORD=
# Database name 数据库名称
DB_NAME=
```

###  Examples

```shell
  motors-backup                            Export all tables in the database
  motors-backup users                      Export users table
  motors-backup users,orders               Export users and orders tables
  motors-backup --ignore-table=logs --ignore-table=users
                                           Export all tables except logs and users
  motors-backup --ignore-table-data=logs --ignore-table-data=users users,orders
                                           Export orders table structure and data,
                                           but only structure for logs,users table
  motors-backup --create-database=false users
                                           Export users table without create database statement
  motors-backup --where='id>100' users
                                           Export users table with condition id>100
                                           
                                           
  motors-backup                            导出数据库中的所有表
  motors-backup users                      导出 users 表
  motors-backup users,orders               导出 users 和 orders 表
  motors-backup --ignore-table=logs --ignore-table=users
                                           导出除 logs 和 users 外的所有表
  motors-backup --ignore-table-data=logs --ignore-table-data=users users,orders
                                           导出 orders 表结构和数据，
                                           但 logs,logs 表仅导出结构
  motors-backup --create-database=false users
                                           导出 users 表但不包含 create database 语句
  motors-backup --where='id>100' users
                                           导出 users 表并应用条件 id>100
```