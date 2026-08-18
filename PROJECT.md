# Pkgforge · APT 仓库 by-hash / InRelease 发布队列

## 这不是什么

不是通用网盘、不是代码片段管理、不是电商 SKU。记录是 **Debian 包或发行元数据**：要进入 dist 的 `.deb` 或 `InRelease`。

## 谁在用

内部 Linux 包镜像管理员。要把构建产物发布到符合 Debian 仓库布局的 by-hash 池，并排队生成 InRelease（本产品做队列与校验，不在仓库里实现 GnuPG 签名机）。

## 核心业务

1. 标题必须是 `name_version_arch.deb` 或恰好为 `InRelease`。
2. 正文必须含 `sha256=` 后接 64 位十六进制。
3. 标签为套件 `stable` / `testing` 与组件 `main`。
4. 附件是 `.deb` 实体，路径不得逃出数据根。
5. 备份 zip 写完后立刻 `ReadManifest` 回读，空包视为失败。

## 运行与验收

- 标题 `foo.tgz` 拒绝；`foo_1.0.0_amd64.deb` + 合法 sha256 接受。
- 不做积分商城、不做 Excel 商品导入。
