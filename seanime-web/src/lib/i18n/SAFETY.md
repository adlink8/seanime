# 禁改区清单 —— 汉化硬约束契约

面向 **Phase 4–9 的所有 executor**。本文是**规定**，不是建议，也不是说明。

`03-02` 交付的检查脚本 `seanime-web/scripts/check-i18n-safety.mjs` **逐条实现**本文档定义的 10 条 ERROR 规则。本文档与检查脚本是同一份契约的两半，**任何一方改动都必须同步另一方**，二者不得漂移。动手前先读完本文，改完再按 §5 自检。

## 1. 这份文档解决什么问题

Seanime 前端里有大量字符串**看着像展示文案、实际是程序标识符**。把它们翻成中文不会报错、不会抛异常——但功能会**静默失效**。这是最危险的一类改动：类型检查绿、构建绿、测试可能也绿，用户侧却坏了。

两个真实后果：

1. **筛选器返回空结果**：筛选器的 `value` 直接进 GraphQL 查询参数。`value="CURRENT"` 改成 `value="当前"` 后，请求照发，后端认不出这个枚举，返回空列表。界面不报错，只是"没有数据"。
2. **`switch` 分支永不命中**：`switch (status) { case "CURRENT": ... }` 改成 `case "当前":` 后，后端仍返回 `"CURRENT"`，没有任何分支匹配，逻辑直接走空（或落到 `default` 兜底），界面显示原始枚举值。

Phase 1 已在 44 个文件上验证过：**只改"给人看"的字符串（`label:` / JSX 文本 / `toast`）零风险**。本文档把那个隐式约定显式化为硬约束。

| 代码里的样子 | 看着像什么 | 实际是什么 | 汉化后的后果 |
| --- | --- | --- | --- |
| `value="CURRENT"` | 展示文案 | GraphQL 查询参数 / 表单值 | 后端拒收，筛选返回空 |
| `case "CURRENT":` | 展示文案 | 分支判别常量 | 分支永不命中，功能静默失效 |
| `data-settings-anime-library="advanced-accordion-trigger"` | 无意义标记 | `querySelector` 选择器锚点 | 选择器失配，新手引导步骤定位失败 |
| `queryKey: ["anime", id]` | 无意义数组 | React Query 缓存键 | 缓存永不命中，重复请求 |

## 2. 禁改区（ERROR 级 —— 命中即阻断提交）

下表的 11 条规则中任意一条被命中，即为 **ERROR**：**必须修复**后重新自检。

**关于「绕过」**：「不得用变量提取、字符串拼接、Unicode 转义等方式把中文送进禁改区」是一条**规定，不是机器保证**。脚本只做逐行文本匹配，抓不到 `const v = "中文"; <C value={v} />` 这类间接写法（详见 §5 盲区表）。这条约束靠人守 —— 脚本只是第一道网，不是唯一防线。

| # | 规则 ID | 禁改区 | 为什么 | 检查脚本是否自动拦截 |
| --- | --- | --- | --- | --- |
| 1 | `value-cn` | `value:` / `value=` 的值 | 直接进 API 参数或表单值；改了后端拒收 | ✅ 是 |
| 2 | `case-cn` | `switch` 的 `case "..."` | 分支永不命中，功能静默失效 | ✅ 是 |
| 3 | `compare-cn` | `===` / `!==` / `.includes()` / `.startsWith()` / `.endsWith()` / `.indexOf()` 的字面量 | 比较永假，逻辑走错分支 | ✅ 是 |
| 4 | `data-attr-cn` | `data-*` 属性的**值**（不是属性名） | 被 e2e 测试与 DOM 查询当选择器 | ✅ 是 |
| 5 | `dom-id-cn` | JSX `id="..."` 的值 | 被 `getElementById` / 锚点 / 测试依赖 | ✅ 是 |
| 6 | `querykey-cn` | `queryKey` / `mutationKey` 数组内的字符串 | React Query 缓存键，变了缓存永不命中 | ✅ 是 |
| 7 | `path-cn` | `url` / `href` / `src` / `path` / `pathname` / `route` / `endpoint` 的值 | 路由与接口路径 | ✅ 是 |
| 8 | `storage-cn` | `localStorage` / `sessionStorage` 的键名 | 持久化键变了，旧数据读不到 | ✅ 是 |
| 9 | `import-cn` | `import` / `export ... from` 的模块路径 | 模块解析失败，构建报错 | ✅ 是 |
| 10 | `const-cn` | 全大写常量（`SORT_OPTIONS` 这类）被赋中文字面量 | 通常是枚举表，值要进 API | ✅ 是 |
| 11 | `value-expr-cn` | `value={...}` **表达式容器内**的中文字面量 | 与第 1 条同源。`value={ident}` 是本项目最高频的写法之一（`src/` 内 476 处），一旦写成 `value={"中文"}` 或 `value={cond ? "中文" : "CURRENT"}`，危害与第 1 条完全相同 | ✅ 是 |

**重要澄清**：第 4 条 `data-attr-cn` 管的是 `data-*` 的**值**，不是属性名。属性名是 `data-foo` 形式的 ASCII 标识符，本身不可能含中文，无需约束。ROADMAP 里「`data-*` 属性名」的表述有歧义，**本文档按「值」执行**，检查脚本亦按「值」实现。

以下每条规则给出**真实代码形态的正例与反例**，照此判断你的代码是否违规。

### R1 · `value-cn`

```tsx
违规：<SelectItem value="日本">
违规：<option value="已完成">done</option>
违规：const x = { value: "连载中", label: "连载中" }
合法：<SelectItem value="CURRENT">当前</SelectItem>
合法：{ value: "CURRENT", label: "当前" }
```

### R2 · `case-cn`

```tsx
违规：switch (s) { case "已完成": return 1 }
合法：switch (s) { case "CURRENT": return 1 }
```

### R3 · `compare-cn`

```tsx
违规：if (status === "进行中") { }
违规：arr.includes("特别篇")
违规：name.startsWith("第")
合法：if (status === "CURRENT") { }
```

### R4 · `data-attr-cn`

```tsx
违规：<div data-continue-watching-title="继续观看">
违规：<div data-testid="媒体卡片">
合法：<div data-testid="continue-watching">
```

**这不是假想风险 —— 本仓库有真实消费方。** `data-*` 属性的值被 `querySelector` / `closest` / `getAttribute` 在运行时读取：

| 属性值 | 定义处 | 消费处 |
| --- | --- | --- |
| `data-settings-anime-library="advanced-accordion-trigger"` | `src/app/(main)/settings/_containers/anime-library-settings.tsx:75` | `src/app/(main)/_features/tour/changelog-tour.tsx:64,70,71`（`waitForSelector` + `click`） |
| `data-home-toolbar-scan-button` | `src/app/(main)/_features/home/home-toolbar.tsx:151` | `src/app/(main)/_features/tour/changelog-tour.tsx:39` |
| `data-chapter-page-container` | `src/app/(main)/manga/_containers/chapter-reader/_components/chapter-page.tsx:75` | `src/app/(main)/manga/_lib/handle-chapter-reader.ts:485-486`（`closest`） |
| `data-css-scope`（动态值） | `src/app/(main)/_features/plugin/components/registry-components.tsx:78` | 同文件 `:85`（`[data-css-scope="${scopeId}"]`） |

把 `advanced-accordion-trigger` 改成中文，`changelog-tour` 的引导步骤会直接找不到目标元素。

**两点如实说明**：① 样本里的 `data-testid` 在当前仓库**不存在**，保留它是因为这类属性能被测试依赖，规则对新增代码同样适用；② `data-continue-watching-title`（`src/app/(main)/_features/anime-library/_containers/continue-watching.tsx:159`）目前**只有定义、没有消费方**，但它属于同一类约束，改动风险与上表同源。

### R5 · `dom-id-cn`

```tsx
违规：<div id="主容器">
合法：<div id="main-container">
```

### R6 · `querykey-cn`

```tsx
违规：queryKey: ["动漫列表", id]
违规：mutationKey: ["保存设置"]
合法：queryKey: ["anime", id]
```

### R7 · `path-cn`

```tsx
违规：href="/中文路径"
违规：const path = "/媒体库"
违规：endpoint: "/接口/动漫"
合法：href="/anime"
合法：<img src="/images/cover.png" alt="封面图" />
```

### R8 · `storage-cn`

```tsx
违规：localStorage.setItem("用户偏好", v)
合法：localStorage.setItem("user-preferences", v)
```

### R9 · `import-cn`

```tsx
违规：import x from "@/中文模块"
合法：import x from "@/lib/helpers"
```

### R10 · `const-cn`

```tsx
违规：const SORT_OPTIONS = "中文排序"
合法：const SORT_OPTIONS = "TITLE_DESC"
```

### R11 · `value-expr-cn`

```tsx
违规：value={"中文"}
违规：<SelectItem value={cond ? "动画" : "Anime"}>
违规：<SelectItem value={`中文`}>
合法：value={ident}
合法：<SelectItem value={CURRENT}>当前</SelectItem>
合法：value={formatDate(createdAt)}
```

注意区分：`value={ident}`（变量）与 `value={formatDate(x)}`（函数调用）是**合法且高频**的写法，规则不碰；只有容器内部出现**中文字面量**才报。

## 3. 可安全汉化（白名单 —— 大胆改）

下表位置里的字符串**确认可以自由汉化**，这是日常最高频的改动面。注意：组件允许汉化的是其**子节点（children）**，不是它的 props——`<Button value="..." />` 仍受 R1 约束。

| 位置 | 样本 |
| --- | --- |
| 对象/数组里的 `label:` | `label: "继续观看",` |
| JSX 元素文本 | `<h2>继续观看</h2>` |
| `DialogTitle` / `CardTitle` 等标题组件子节点 | `<DialogTitle>确认删除</DialogTitle>` |
| `Button` / `Badge` 等组件子节点（**不含 `value` 等 props**） | `<Button variant="outline">取消</Button>` |
| toast / 通知消息 | `toast.success("已保存")` |
| `placeholder` | `placeholder="搜索动漫"` |
| `title` | `<span title="暂无内容">` |
| `alt` | `<img alt="封面图" />` |
| `aria-label` / `aria-description` | `aria-label="关闭对话框"` |
| `description` / `helperText` 类 props | `description="选择要导入的文件"` |
| 工具函数返回的展示文案 | `return "已跳过"` |
| 词条表条目本身 | `"library.filter.current": "当前",` |

> 实例核对：`DialogTitle`、`aria-description`、`helperText` 三个子项在当前仓库中暂无实例（同行的 `CardTitle` / `aria-label` / `description` 均有实例，故该行类别整体成立）。本仓库暂无实例，规则仍适用于新增代码。

**判定心智模型**（拿不准时只问这一句）：

> 问一句「这个字符串会被程序读，还是只会被人读？」
> 被人读 → 放心改。被程序读（当参数、当键、当比较对象）→ 一个字都别动。

## 4. 需人工判断（不在自动检查范围）

下列位置**故意不纳入** §2 的自动规则：它们的语义取决于上下文，写正则会大量误报，因此交给人工逐处判断。命中这些位置时，不要依赖检查脚本，**自己去看代码**。

| 位置 | 为什么模糊 | 建议 |
| --- | --- | --- |
| 对象字段 `type: "动画"` | 可能只是前端内部分类标签，也可能进 API。`type` 二字太宽，做规则会大量误报 | 改动前 grep 该字段名，看是否出现在 API 请求/响应构造里 |
| 对象字段 `name: "..."` | 可能是展示名，也可能是后端校验用的键 | 同上 |
| 对象字段 `key: "..."` | React 列表 key 用中文**不报错**但语义怪异 | 保持英文，除非确认只用于展示 |
| `console.log("中文")` | 开发日志，不影响功能 | 可保留中文，也可英文；不强制 |
| 注释里的中文 | 完全无害 | 鼓励写中文注释 |

## 5. 检查脚本用法

脚本路径 `seanime-web/scripts/check-i18n-safety.mjs`（由 `03-02` 交付）。在 `seanime-web/` 下执行：

```bash
# 扫工作区未提交改动（日常用这个）
npm run check:i18n-safety

# 扫指定提交范围
node scripts/check-i18n-safety.mjs 9bdd052..77722aa

# 扫单个提交
node scripts/check-i18n-safety.mjs --commit 58096b0

# 扫指定文件（不走 git）
node scripts/check-i18n-safety.mjs --file src/app/page.tsx

# 列出全部规则
node scripts/check-i18n-safety.mjs --list-rules

# 内置正反例自检
node scripts/check-i18n-safety.mjs --selftest

# 输出 JSON 存档
node scripts/check-i18n-safety.mjs --json --report ../.planning/codebase/safety-checks/phase-4.json
```

**退出码约定**：

| 退出码 | 含义 | 处置 |
| --- | --- | --- |
| `0` | 只有 WARN，或完全干净 | 可提交 |
| `1` | 存在 ERROR | 必须修复，不得绕过 |

**范围过滤**（脚本行为）：只扫 `.ts` / `.tsx` / `.js` / `.jsx` / `.mjs` / `.cjs`；跳过 `node_modules/`、`out/`、`dist/`；跳过纯注释行（行首为 `//`、`*`、`/*`、`{/*`）；豁免检查器自身（其内建自检样本是故意违规，不豁免会自报 21 个假错）。

**默认扫描范围**（不带参数时）= `git diff`（未暂存）+ `git diff --cached`（已暂存）+ **未跟踪的新文件**。

最后一项不可省：`git diff` 默认**不显示未跟踪文件**，而「新建文件」恰是汉化阶段最常见的改动形态（Phase 1/2 的提交几乎全是新增文件）。缺了它，工作区自检会在新建文件上出现静默盲区 —— 实测：新建一个含 `value: "中文"` 的 `.ts` 文件，不加此项时脚本报 `0 added lines in 0 files` 并放行。

**已知不覆盖的盲区**（需要人工兜底，不要指望脚本）：

**漏检方向** —— 违规存在但脚本不报：

| # | 盲区 | 说明 |
| --- | --- | --- |
| 1 | **变量间接** | `const v = "中文"; <C value={v} />` —— 文本匹配看不到赋值链路 |
| 2 | **字符串拼接 / 表达式求值** | `value={"中" + "文"}`、`value={["中文"].join("")}` |
| 3 | **Unicode 转义** | `value="\u4e2d\u6587"` —— 源码里没有中文字符，运行时却是中文。**这属于故意绕过，明文禁止** |
| 4 | **非 `value` 位置的表达式容器** | R11 只覆盖 `value={...}`。`sort={"中文"}`、`filter={"中文"}` 这类其他 prop 抓不到 |
| 5 | **数组 / 对象字面量里的裸中文** | `const STATUS = ["已完结", "连载中"]` —— 没有 `value:` 前缀可依附 |
| 6 | **TS `enum` 成员** | `enum S { A = "中文" }` —— 与普通赋值语法无法逐行区分，做规则会大量误报 |
| 7 | **跨多行的禁改区写法** | 规则逐行匹配，`value:` 与值分处两行时抓不到 |
| 8 | **动态 `import()`** | `import("中文模块")` —— R9 只匹配静态 `import ... from` 形式 |
| 9 | **模板字面量比较** | `` x === `中文` `` —— R3 的引号类只含 `'` 与 `"` |
| 10 | **非白名单扩展名** | `.json` / `.md` / `.css` / `.graphql` 不在扫描范围 |
| 11 | **`type:` / `name:` / `key:` 字段** | §4 已列为需人工判断，规则故意不覆盖 |

**误报方向** —— 脚本报了但实际无害：

| 情况 | 说明 |
| --- | --- |
| 行内注释 | `foo(x) // value: 保持英文` —— 注释不在行首，第 3 层过滤抓不到。若退出码为 0 可忽略 |
| 字符串字面量里恰好含 `value:` 文本 | 如 `console.log("check value: 中文")`。罕见 |

> **怎么用这张表**：脚本返回 `0` **不等于**「没有违规」，只等于「已覆盖的 13 条规则里没有命中」。上表漏检方向第 1–11 项必须靠人看代码。反过来，脚本报 ERROR 时先去本节核对是不是误报，再决定改代码还是收紧规则。
