# i18n 模块使用说明

面向 P5–P9 的 executor。本文描述的都是**当前代码的真实行为**，动手前请先按此照抄，不要凭记忆自创调用方式。

## 1. 定位

`seanime-web/src/lib/i18n/` 是**全项目唯一的中文文案来源**。

- 页面侧（`src/app/`、`src/lib/` 下的组件与工具）**不得再写裸中文字符串**，一律经 `t()` 取值。
- 词条表只此一份；不存在第二套文案。
- `t()` 是纯函数、框架无关：不 import React、不依赖浏览器 API，服务端工具与数据转换模块同样可用。

## 2. 导入方式

```ts
import { t } from "@/lib/i18n"
```

- 路径别名 `@` 同时由 `tsconfig.json`（`"@/*": ["./src/*"]`）与 `rsbuild.config.ts`（`alias: { "@": path.resolve(__dirname, "./src") }`）配置，写法已与仓库一致。
- `@/lib/i18n` 指向 `src/lib/i18n/index.ts`（barrel）。barrel 只导出公共 API：`t`、`zhCN`、类型 `TranslationKey` / `Dictionary` / `DictionaryModule` / `TranslateParams`，以及运行时模块清单 `DICTIONARY_MODULES`。
- **页面侧只应导入 `t`**。`zhCN`、`DICTIONARY_MODULES` 供测试与工具使用，业务代码不要直接读字典。

## 3. key 命名规范

统一格式：

```
<module>.<section>.<name>
```

- 全小写 ASCII，以点号分段；例如 `navigation.sidebar.search`、`common.action.confirm`、`settings.field.auto_refresh_library`。
- 段内多个单词用下划线连接（`auto_refresh_library`），**不要**用连字符或驼峰。
- **严禁**在 key 中出现中文、空格或大写字母。
- 每个新增词条行尾加来源注释（`// <源文件路径>`），便于回溯与 P4 存量回填核对。

## 4. 六个模块的分工

模块顺序即 `types.ts` 中 `DICTIONARY_MODULES` 的顺序，也是 `locales/zh-CN/index.ts` 的合并顺序：

| 模块 | 负责的界面 | 对应 ROADMAP Phase |
|---|---|---|
| `common` | 跨界面通用状态名、空状态、通用操作按钮（首页 / Debrid / 服务端工具） | P5–P9 通用 |
| `navigation` | 顶部菜单、主侧边栏、离线侧边栏的导航项与导航相关操作 | Phase 5（导航骨架） |
| `library` | 媒体库详情页（统计卡、筛选器）与媒体库排序选项 | Phase 5（媒体库） |
| `search` | 探索发现（搜索）页、高级搜索筛选项、流派/季度/状态/格式/国家常量表 | Phase 6（搜索、发现与时间表） |
| `settings` | 系统设置页的标签页、设置卡片标题与字段标签 | Phase 7（设置页） |
| `player` | 播放偏好、桌面播放器 / 外部播放器关联、连播与流式播放设置 | Phase 7（条目详情与播放） |

漫画模块（Phase 8）与次要页面/动态文案（Phase 9）在 P5–P9 推进时按需新增模块或复用 `common`，届时同步更新本表与 `DICTIONARY_MODULES`。

## 5. 如何新增词条

1. 打开 `locales/zh-CN/<module>.ts`，找到对应模块的 `xxxDictionary` 对象。
2. 在对象内按 **key 字典序**插入新条目：`"<module>.<section>.<name>": "中文文案", // <源文件路径>`。
3. 不要新建「未在 `DICTIONARY_MODULES` 中登记」的模块文件；确需新模块时，先在 `types.ts` 扩展 `DictionaryModule` 与 `DICTIONARY_MODULES`，再在 `locales/zh-CN/index.ts` 的 `moduleDictionaries` 与 `zhCN` 合并处挂上。
4. 合并采用 `...commonDictionary, ...navigationDictionary, ...` 的对象展开，**后展开的同名 key 会覆盖先展开的**；`locales/zh-CN/index.ts` 内置开发期重复 key 检测，键数不一致会在控制台 `console.warn` 列出冲突项——新增后应确认无警告。

## 6. 未命中行为（重要）

`t.ts` 对不存在的 key 的处理是：

```ts
const template = dictionary[key]
if (template === undefined) return key
```

即 **原样返回 key 本身，不抛异常、不打日志**。

- 因此：**界面上出现 `navigation.sidebar.xxx` 这样的字面量，就等于该处漏了词条**。这是排查漏翻最快的手段——看到点号小写 ASCII 就当漏翻处理。
- 未命中属于设计内的正常降级（渐进式本地化），不要为了「消掉这个 key」而往词条表里塞占位词条。

## 7. 插值用法

词条值用 `{name}` 形式的占位符，调用时经第二个参数传入：

```ts
t("data.episodes.count", { count: 12 })   // 词条 "data.episodes.count": "{count} 集" → "12 集"
```

- 占位符正则：`/\{(\w+)\}/g`，`name` 为单词字符；`\w` 不含中文，故 key 与参数名都应是 ASCII。
- 参数类型 `TranslateParams = Record<string, string | number>`，数字会被 `String()` 字符串化。
- **缺少对应参数时占位符原样保留**（例：`t("data.episodes.count", { other: 1 })` 仍返回 `"{count} 集"`），同样不抛异常。
- 现成的真实插值词条示例：`player.external.scheme_placeholder` → `"例如: outplayer://{url} 或 iina://weblink?url={url}"`，用 `t(key, { url: "..." })` 渲染。

## 8. 禁改区（不得汉化、不得进词条表）

以下标识符一律**不汉化、不进词条表**：

- `value:` 枚举值与 `switch case` 分支常量
- `data-*` 属性名
- `queryKey` / `mutationKey`
- `URL` / `path` 等路由与接口路径

原因与完整口径见：

- `.planning/codebase/LOCALIZATION-INVENTORY.md` §4「下一步建议」中的指标口径更正（数据表 `advanced-search-constants.ts` 的 854 个 `label` 可译、`value` 必须冻结；`manga/_lib/language-map.ts` 的键与 `nativeName` 冻结）。
- `.planning/ROADMAP.md` 的「⚠ 指标口径更正：『残留英文串』被两张数据表严重放大」一节。

判断准则：**只要是运行时参与比较、请求、路由、存储或类型判别的字符串，都必须保持原样**；只有「给人看的展示文本」才走 `t()`。

## 9. 验证命令

在 `seanime-web/` 下依次执行：

```bash
npm run typecheck     # tsc --pretty false
npx vitest run        # 单测（含 src/lib/i18n/t.test.ts）
npm run build         # typecheck && rsbuild build
```

三者必须全部退出码 0。新增词条不改变 `t()` 的签名与行为契约，若 `t.test.ts` 失败，先在代码里找原因，**不得删测试或放宽断言**。

## 10. 文件行数上限

单个词条文件 **≤ 400 行**。当前六个模块文件最大 66 行（`settings.ts`），距上限尚有余量。某模块接近上限时按 `<module>.<section>` 再拆分子文件并在该模块的 `index` 中合并，拆分后同步更新本节与 §4 的分工表。
