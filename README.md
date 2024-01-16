# Gf-Extend 基于v1.16.6版本

**GoFrame框架扩展工具包**

![go-kit](https://img.shields.io/badge/go-kit-ea7b99)
![GitHub last commit](https://img.shields.io/github/last-commit/unitsvc/go-kit?style=flat-square)
[![Go Doc](https://godoc.org/unitsvc/go-kit?status.svg)](https://pkg.go.dev/github.com/unitsvc/go-kit)
[![Go Report](https://goreportcard.com/badge/unitsvc/go-kit?v=1)](https://goreportcard.com/report/unitsvc/go-kit)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/unitsvc/go-kit)
[![License](https://img.shields.io/github/license/unitsvc/go-kit.svg?style=flat)](https://github.com/unitsvc/go-kit)

#### **安装**

> go get -u -v -d github.com/unitsvc/go-kit

`指定版本`
> go get github.com/unitsvc/go-kit@latest

**或**

> require github.com/unitsvc/go-kit latest

#### 实例化gf-casbin实例bean（推荐使用）

* **自动注册（无需关心数据源种类）**

```go
e, err := gfadapter.NewEnforcerBean()
```

* **手动注册**

```go
e, err := gfadapter.NewEnforcerBean(g.DB())
e, err := gfadapter.NewEnforcerBean(g.DB("sqlite"))
e, err := gfadapter.NewEnforcerBean(g.DB("mysql"))
e, err := gfadapter.NewEnforcerBean(g.DB("pgsql"))
```

#### 实例化gf-casbin执行器

* **自动注册**

```go
e, err := gfadapter.NewEnforcer()
```

* **手动注册**

```go
e, err := gfadapter.NewEnforcer(g.DB())
e, err := gfadapter.NewEnforcer(g.DB("mysql"))
e, err := gfadapter.NewEnforcer(g.DB("sqlite"))
e, err := gfadapter.NewEnforcer(g.DB("pgsql"))
```

#### 解压二进制中单文件到本地

```go
gfboot.SingleFileMemoryToLocal("./db", "sqlite3.db", "db/sqlite3.db")
```

#### 获取web响应对象

```go
hres.Ok()
```

#### 获取page分页对象

```go
hdto.NewPage()
```

#### **更多...**
