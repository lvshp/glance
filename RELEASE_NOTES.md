# Release Notes

## Release Title Rule

后续发版统一使用纯版本号作为 Release 标题，例如：

* `v0.1.1`
* `v0.2.0`
* `v1.0.0`

不再额外追加“维护版发布”“首次发布”“ReadCLI 发布”等描述。

## Release Naming

建议保持以下约定：

* Git tag：`vX.Y.Z`
* GitHub Release 标题：`vX.Y.Z`
* 二进制和压缩包文件名：保留版本号，例如 `readcli-darwin-arm64-v0.2.0.tar.gz`

## Fork Note

仓库仍然基于 [TimothyYe/glance](https://github.com/TimothyYe/glance) 持续演进，Release 标题简化不会影响 fork 说明、许可证和致谢信息的保留。
