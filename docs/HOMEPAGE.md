# 下载仓库首页维护

README.md 是 GitHub 仓库面向用户的主页。原来的构建说明保留在 RELEASE-MAINTENANCE.md。

- 顶部宣传图：assets/hero.png；工具展示图：assets/agents.png。
- 视觉源文件：assets/homepage-artwork.html。修改后在浏览器中分别导出 .hero 和 .strip 区域；图片中文字使用本机字体渲染，GitHub 不需要加载字体服务。
- 品牌图标复用 EasyAI 直播素材；各工具名称和商标属于相应权利人。
- 下载按钮使用真实已发布安装包 URL。scripts/update-homepage.py 读取 GitHub 的最新正式 Release，只更新 release-downloads 标记内的区块，不改其他文案和二维码。
- Refresh download links 工作流在 Build EasyAI Desktop 成功后执行，也可手动运行。部分平台尚未发布时，显示查看可用版本入口，不生成不存在的下载地址。
- 按维护者要求，交流群先使用明确标注“示意 · 不可扫码”的占位图 assets/community-qr.png，源文件为 assets/community-qr-placeholder.html，不包含真实入群地址。
- 换成正式二维码：覆盖 assets/community-qr.png，删除 README 中“二维码为展示占位”的提示，并修改图片 alt。真实二维码建议保留完整白边；如果图片为正方形，把 README 图片的 height 改为 240（或删除 height）。发布版本不会覆盖二维码区。
- 工具和模型支持范围核对 EasyAI 生产目录 /v1/agents/catalog；最低系统要求同时核对 Desktop 打包配置。不要宣传所有工具都支持所有系统或所有模型。

本次核对日期：2026-11-05。可用安装包：Windows x64、macOS arm64、macOS x64。

运行 python scripts/test-update-homepage.py 检查版本、部分平台发布和内容保留；运行 python scripts/update-homepage.py 更新下载区。
