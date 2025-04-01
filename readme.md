# jenkinsdk


对 Jenkins [Remote Access API](https://www.jenkins.io/doc/book/using/remote-access-api/) 的简单封装。

* **jenkins**: 获取jenkins信息(测活)
* **view**: 创建，删除，获取配置，更改description，更改配置
* **folder**: 创建，删除，获取配置，更改description，更改配置
* **job**: 创建(config)，复制job，删除，获取配置，启用，禁用，更改description，更改配置
* **csrf**: 获取csrf，设置csrf
* **plugin**: 获取所有plugins，获取某个plugin


# 1. 说明

**版本**

代码兼容 `jenkins 2.479.3`

**核心思想**

通过 config.xml 来创建 folder, view, job 等。

# 2. 安装更新

```shell
go get -u "github.com/mrlaojia/go-jenkins"
```
 

