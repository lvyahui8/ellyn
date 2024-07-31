# ellyn

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/lvyahui8/ellyn)](https://goreportcard.com/report/github.com/lvyahui8/ellyn)
[![codecov](https://codecov.io/gh/lvyahui8/ellyn/graph/badge.svg?token=YBV3TH2HQU)](https://codecov.io/gh/lvyahui8/ellyn)


### Requires

- Go Version >= 1.18

### 通用层模型

- program 代表应用程序
- file 代表一个代码文件
- package 代表一个代码包
- method 代表一个方法
    - 开始
    - 结束
    - 出入参类型、变量名
- block 代表一个代码块
    - 开始
    - 结束
    - 块hash
- goStmt 代表一个go异步语句
- struct 代表一个go struct

暴露通用能力（handler）
- 插桩能力
    - 函数出入口植入代码
    - block块起点位置植入代码
    - go前后位置植入代码
- 对象迭代能力
- 文件迭代能力

### 应用场景

- ioc、aop 代码生成
    - 支持获取program上下文
    - 支持获取package上下文
    - 获取应用所有的package
    - 遍历package
    - 遍历文件
    - 遍历struct，判断是否内嵌的bean
    - 在package内生成一个bean init文件，统一对包内的bean实现
- 调用上下文采集：链路、参数、耗时、异常，应用在故障告警、定位，日常开发
    - 支持获取program上下文
        - 存放所有的method，统计methodId等等
    - 获取应用的所有的package
    - 遍历文件
    - 遍历方法
    - 获取方法的出入参名称、类型
    - 支持在方法入口植入代码
    - 支持修改方法签名
- 覆盖率、流量覆盖、单测覆盖
    - 支持获取program上下文
        - file、func、block全局id计数器
    - 获取应用的所有package
    - 遍历文件
    - 遍历block
    - 支持在block的入口注入代码

