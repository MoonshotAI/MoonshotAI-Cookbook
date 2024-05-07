### 简介

webman/openai 是一个异步非阻塞的openai客户端，配合webman可以做到同一时刻支持上万并发调用，使用简单，返回如丝般的顺滑，无卡顿。

> 开源地址：https://github.com/webman-php/openai

### 安装

```php
cd examples/php_demo

composer install
```
> 安装该插件之前记得先安装webman框架。

### 使用

#### 运行服务

启动 workerman 提供服务

```php
# php demo.php start
Workerman[chunk.php] start in DEBUG mode
------------------------------------------- WORKERMAN --------------------------------------------
Workerman version:4.1.15          PHP version:7.4.27           Event-Loop:\Workerman\Events\Event
-------------------------------------------- WORKERS ---------------------------------------------
proto   user            worker          listen                  processes    status
tcp     root            none            http://0.0.0.0:8782/    1             [OK] 
--------------------------------------------------------------------------------------------------
Press Ctrl+C to stop. Start success.
```
> 对外提供服务地址：`http://0.0.0.0:8782`
