#!/usr/bin/env node
/**
 * check-i18n-safety.mjs —— 汉化「禁改区」检查脚本
 *
 * 用途：
 *   扫描 git diff 中的「新增行」，判定其中的中文是否落入了禁改区（被程序读取的字符串）。
 *   命中任意 ERROR 规则即以非零码退出，作为 Phase 4–9 每个阶段的强制门禁（SAFE-03）。
 *
 * 契约来源：
 *   本脚本逐条实现 seanime-web/src/lib/i18n/SAFETY.md §2 定义的 10 条 ERROR 规则
 *   与 §5 约定的 CLI 用法。SAFETY.md 与脚本是同一份契约的两半，任何一方改动都必须同步另一方。
 *
 * 退出码约定（与 SAFETY.md §5 一致）：
 *   0  只有 WARN，或完全干净
 *   1  存在 ERROR（命中禁改区），必须修复
 *   2  脚本自身异常（参数错误、git 调用失败等）
 *
 * 纯 Node ESM，零第三方依赖。
 */

import fs from 'node:fs';
import path from 'node:path';
import { spawnSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import process from 'node:process';

// ---------------------------------------------------------------------------
// 规则表：13 条 = 11 ERROR + 2 WARN
// 正则一律用字面量书写，逐字对齐 SAFETY.md §2 的规则表，不得改动。
// ---------------------------------------------------------------------------
const RULES = [
  { id: 'value-cn',     level: 'ERROR', desc: 'value 字段/属性的值出现中文',
    re: /\bvalue\s*[:=]\s*["'`][^"'`]*[\u4e00-\u9fff]/ },
  { id: 'value-expr-cn', level: 'ERROR', desc: 'value={...} 表达式容器内出现中文字面量',
    re: /\bvalue\s*[:=]\s*\{[^}]*["'`][^"'`]*[\u4e00-\u9fff]/ },
  { id: 'case-cn',      level: 'ERROR', desc: 'switch case 标签出现中文',
    re: /\bcase\s*["'][^"']*[\u4e00-\u9fff]/ },
  { id: 'compare-cn',   level: 'ERROR', desc: '===/!== 或字符串方法比较出现中文',
    re: /([!=]==|\.(?:includes|startsWith|endsWith|indexOf|match|search)\s*\()\s*["'][^"']*[\u4e00-\u9fff]/ },
  { id: 'data-attr-cn', level: 'ERROR', desc: 'data-* 属性的值出现中文（选择器会被破坏）',
    re: /data-[a-zA-Z0-9-]+\s*=\s*["'`][^"'`]*[\u4e00-\u9fff]/ },
  { id: 'dom-id-cn',    level: 'ERROR', desc: 'JSX id 属性的值出现中文',
    re: /\bid\s*=\s*["'][^"']*[\u4e00-\u9fff]/ },
  { id: 'querykey-cn',  level: 'ERROR', desc: 'queryKey/mutationKey 出现中文',
    re: /(queryKey|mutationKey)\s*[:=]\s*\[[^\]]*[\u4e00-\u9fff]/ },
  { id: 'path-cn',      level: 'ERROR', desc: 'url/href/src/path/route/endpoint 出现中文',
    re: /\b(url|href|src|path|pathname|route|endpoint)\s*[:=]\s*["'`][^"'`]*[\u4e00-\u9fff]/ },
  { id: 'storage-cn',   level: 'ERROR', desc: 'localStorage/sessionStorage 键名出现中文',
    re: /(localStorage|sessionStorage)\s*\.\s*(getItem|setItem|removeItem)\s*\(\s*["'][^"']*[\u4e00-\u9fff]/ },
  { id: 'import-cn',    level: 'ERROR', desc: 'import/export 模块路径出现中文',
    re: /^\s*(import|export)\b.*\bfrom\s*["'][^"']*[\u4e00-\u9fff]/ },
  { id: 'const-cn',     level: 'ERROR', desc: '全大写常量被赋中文值',
    re: /\bconst\s+[A-Z][A-Z0-9_]{2,}\s*(?::[^=]+)?=\s*["'][^"']*[\u4e00-\u9fff]/ },
  { id: 'classname-cn', level: 'WARN',  desc: 'className 出现中文（tailwind 类名不可能含中文）',
    re: /className\s*=\s*["'`][^"'`]*[\u4e00-\u9fff]/ },
  { id: 'objkey-cn',    level: 'WARN',  desc: '对象字面量的键使用中文',
    re: /[{,]\s*["'][^"']*[\u4e00-\u9fff][^"']*["']\s*:/ },
];

// 三层范围过滤之第 1 层：扩展名白名单
const SCAN_EXT = new Set(['.ts', '.tsx', '.js', '.jsx', '.mjs', '.cjs']);
// 三层范围过滤之第 2 层：路径排除（一律用小写正斜杠路径做 includes 判断）
const EXCLUDE_SEGMENTS = ['node_modules/', '/out/', '/dist/', '.git/'];

// ---------------------------------------------------------------------------
// 内置自检样本
//   基础 37 条（21 BAD + 16 GOOD）逐字取自 03-02 plan，已被编排层 Python 原型实测验证。
//   后由编排层按 verifier 的绕过测试结果补充 6 条（3 BAD + 3 GOOD），覆盖 `value={...}`
//   表达式容器这一最大盲区（`value={ident}` 在 src 内出现 476 次，是最高频写法之一）。
//   改动样本时必须同步 SAFETY.md 的规则表与盲区表。
// ---------------------------------------------------------------------------
const SELFTEST_BAD = [
  '<SelectItem value="日本">',
  '<option value="已完成">done</option>',
  'const x = { value: "连载中", label: "连载中" }',
  'value="中文"',
  'value={"中文"}',
  '<SelectItem value={cond ? "动画" : "Anime"}>',
  '<SelectItem value={`中文`}>',
  'switch (s) { case "已完成": return 1 }',
  'if (status === "进行中") { }',
  'if (s !== "已停止") { }',
  'arr.includes("特别篇")',
  'name.startsWith("第")',
  '<div data-continue-watching-title="继续观看">',
  '<div data-testid="媒体卡片">',
  '<div id="主容器">',
  'queryKey: ["动漫列表", id]',
  'mutationKey: ["保存设置"]',
  'href="/中文路径"',
  'const path = "/媒体库"',
  'endpoint: "/接口/动漫"',
  'localStorage.setItem("用户偏好", v)',
  'sessionStorage.getItem("当前标签")',
  'import x from "@/中文模块"',
  'const SORT_OPTIONS = "中文排序"',
];

const SELFTEST_GOOD = [
  'label: "继续观看",',
  '<h2>继续观看</h2>',
  'toast.success("已保存")',
  '<SelectItem value="CURRENT">当前</SelectItem>',
  '{ value: "CURRENT", label: "当前" }',
  'value={ident}',
  '<SelectItem value={CURRENT}>当前</SelectItem>',
  'value={formatDate(createdAt)}',
  '<div data-testid="continue-watching">',
  'placeholder="搜索动漫"',
  '<span title="暂无内容">',
  '<img alt="封面图" />',
  'aria-label="关闭对话框"',
  'className="text-sm font-medium"',
  '// value: 保持英文枚举',
  '"library.filter.current": "当前",',
  '<DialogTitle>确认删除</DialogTitle>',
  '<Button variant="outline">取消</Button>',
  'description="选择要导入的文件"',
];

// 已知不覆盖：type 字段语义模糊（可能只是前端分类标签，也可能进 API），
// 做规则会大量误报。SAFETY.md §4 已把它列入「需人工判断」。断言它不命中，这是有意设计。
const SELFTEST_KNOWN_UNCOVERED = 'const y = { type: "动画" }';

// ---------------------------------------------------------------------------
// 工具函数
// ---------------------------------------------------------------------------

/** 统一转成正斜杠路径，便于跨平台做 includes / 展示 */
function toPosix(p) {
  return String(p).replace(/\\/g, '/');
}

/**
 * 检查器自身的路径（posix，相对 cwd），用于永久豁免。
 *
 * 为什么需要：本文件内建了 24 条故意违规的自检样本，它们必须是人可读的
 * 真实中文。但只要本文件处于「未跟踪」或「被修改」状态，git 就会把它交回来，
 * 24 条样本会全部误报成 ERROR（实测：20 errors）。语义上「检查器不检查自己」
 * 是准确且可解释的，因此按自身路径精确豁免，而不是放宽任何一条规则。
 *
 * 用 path.relative(cwd, ...) 归一化，保证与 git 输出的路径口径一致
 * （git 默认输出相对 cwd 的路径）。
 */
const SELF_FILE = toPosix(path.relative(process.cwd(), fileURLToPath(import.meta.url)));

/** 三层范围过滤之第 1、2 层：扩展名白名单 + 路径排除（外加自身豁免） */
function shouldScan(file) {
  const p = toPosix(file);
  if (!p) return false;
  if (p === SELF_FILE) return false;
  for (const seg of EXCLUDE_SEGMENTS) {
    if (p.includes(seg)) return false;
  }
  return SCAN_EXT.has(path.extname(p).toLowerCase());
}

/** 目录剪枝：命中排除段的目录不再递归 */
function isExcludedDir(dirPath) {
  const p = toPosix(dirPath);
  for (const seg of EXCLUDE_SEGMENTS) {
    if (p.includes(seg)) return true;
  }
  return false;
}

/** 三层范围过滤之第 3 层：注释行跳过（允许前导空白） */
function isCommentLine(text) {
  const t = text.replace(/^\s+/, '');
  return (
    t.startsWith('//') ||
    t.startsWith('*') ||
    t.startsWith('/*') ||
    t.startsWith('{/*')
  );
}

/** 对单行文本跑全部 12 条规则，返回命中列表（正则无 g 标志，exec 无状态） */
function runRules(text) {
  const hits = [];
  for (const rule of RULES) {
    const m = rule.re.exec(text);
    if (m) hits.push({ id: rule.id, level: rule.level, match: m[0] });
  }
  return hits;
}

// ---------------------------------------------------------------------------
// CLI 参数解析
// ---------------------------------------------------------------------------

function parseArgs(argv) {
  const opts = {
    files: [],
    commit: null,
    selftest: false,
    listRules: false,
    json: false,
    report: null,
    quiet: false,
    help: false,
    positional: [],
  };
  for (let i = 0; i < argv.length; i++) {
    const a = argv[i];
    switch (a) {
      case '--commit':
        opts.commit = argv[++i];
        break;
      case '--file':
        opts.files.push(argv[++i]);
        break;
      case '--selftest':
        opts.selftest = true;
        break;
      case '--list-rules':
        opts.listRules = true;
        break;
      case '--json':
        opts.json = true;
        break;
      case '--report':
        opts.report = argv[++i];
        break;
      case '--quiet':
        opts.quiet = true;
        break;
      case '--help':
      case '-h':
        opts.help = true;
        break;
      default:
        if (a.startsWith('-')) throw new Error(`未知选项: ${a}`);
        opts.positional.push(a);
    }
  }
  return opts;
}

// ---------------------------------------------------------------------------
// 采集：git diff 与直接读文件
// ---------------------------------------------------------------------------

/** 解析一段 unified diff 文本，提取「新增行」及其在新文件中的行号 */
function collectDiffLines(diffText) {
  const entries = [];
  if (!diffText) return entries;

  let file = null;
  let newLine = 0;
  let inHunk = false;

  for (const raw of diffText.split('\n')) {
    // 每个文件块以 "diff --git" 开始，重置状态
    if (raw.startsWith('diff --git')) {
      inHunk = false;
      file = null;
      continue;
    }
    // hunk header：@@ -a,b +c,d @@ —— 新文件起始行号取 c
    if (raw.startsWith('@@')) {
      const m = /^@@ -\d+(?:,\d+)? \+(\d+)(?:,\d+)? @@/.exec(raw);
      if (m) newLine = parseInt(m[1], 10);
      inHunk = true;
      continue;
    }
    // hunk 之外：只关心目标文件路径的 "+++ b/..." 行
    if (!inHunk) {
      if (raw.startsWith('+++ ')) {
        let p = raw.slice(4).trim();
        if (p === '/dev/null') {
          file = null;
        } else {
          if (p.startsWith('b/')) p = p.slice(2);
          file = p;
        }
      }
      continue;
    }
    // hunk 内：\ 表示 "No newline at end of file"，忽略
    if (raw.startsWith('\\')) continue;
    // 新增行：入结果并推进新文件行号
    if (raw.startsWith('+')) {
      if (file) entries.push({ file, line: newLine, text: raw.slice(1) });
      newLine++;
      continue;
    }
    // 删除行：不推进新文件行号
    if (raw.startsWith('-')) continue;
    // 上下文行：只推进新文件行号，不入结果
    if (raw.startsWith(' ')) {
      newLine++;
      continue;
    }
    // 其余（空串等）忽略
  }
  return entries;
}

/** 调用 git 并返回 stdout；失败即抛错（由顶层捕获转为退出码 2） */
function gitDiff(args) {
  const res = spawnSync('git', args, {
    encoding: 'utf8',
    maxBuffer: 128 * 1024 * 1024,
  });
  if (res.error) throw res.error;
  if (res.status !== 0) {
    const msg = (res.stderr || '').trim();
    throw new Error(`git ${args.join(' ')} 失败 (exit ${res.status})${msg ? ': ' + msg : ''}`);
  }
  return res.stdout || '';
}

/** 直接读文件（可传目录，递归）；不走 git */
function collectFileLines(inputs) {
  const files = [];
  const walk = (p) => {
    const st = fs.statSync(p);
    if (st.isDirectory()) {
      if (isExcludedDir(p)) return;
      for (const ent of fs.readdirSync(p, { withFileTypes: true })) {
        walk(path.join(p, ent.name));
      }
    } else if (st.isFile()) {
      files.push(p);
    }
  };
  for (const i of inputs) {
    if (!fs.existsSync(i)) throw new Error(`路径不存在: ${i}`);
    walk(i);
  }

  const entries = [];
  for (const f of files) {
    const content = fs.readFileSync(f, 'utf8');
    const lines = content.split(/\r?\n/);
    lines.forEach((text, idx) => entries.push({ file: f, line: idx + 1, text }));
  }
  return entries;
}

/**
 * 采集未跟踪的新文件（尊重 .gitignore，走 --exclude-standard）。
 *
 * 为什么必须单独做这一步：`git diff` 默认不显示未跟踪文件，而"新建文件"恰恰是
 * 本地化阶段最常见的改动形态之一 —— Phase 1/2 的提交几乎全是新增文件。
 * 少了这一步，工作区自检会在新建文件上出现静默盲区（实测：新建一个含
 * `value: "中文"` 的 .ts 文件，默认模式报 `0 added lines in 0 files`）。
 *
 * 每个未跟踪文件的全部行都视为"新增行"。扩展名过滤与路径排除在 scanEntries 中
 * 统一施加（这里先过一次 shouldScan，避免把二进制大文件读进内存）。
 */
function collectUntrackedLines() {
  const out = gitDiff(['ls-files', '--others', '--exclude-standard']);
  const entries = [];
  for (const raw of out.split('\n')) {
    const file = raw.trim();
    if (!file) continue;
    if (!shouldScan(file)) continue;
    let content;
    try {
      content = fs.readFileSync(file, 'utf8');
    } catch {
      continue; // 读不到（符号链接 / 权限 / 竞态删除）则跳过，不阻断
    }
    content.split(/\r?\n/).forEach((text, idx) => {
      entries.push({ file, line: idx + 1, text });
    });
  }
  return entries;
}

// ---------------------------------------------------------------------------
// 扫描与结果
// ---------------------------------------------------------------------------

/** 对采集到的行套用第 1、2、3 层过滤并跑规则，同时统计扫描行数/文件数 */
function scanEntries(entries) {
  const results = [];
  const files = new Set();
  let scannedLines = 0;

  for (const e of entries) {
    if (!shouldScan(e.file)) continue;
    files.add(toPosix(e.file));
    scannedLines++; // 通过扩展名 / 路径过滤的新增行计入扫描数
    if (isCommentLine(e.text)) continue; // 注释行跳过，不参与规则判定
    for (const h of runRules(e.text)) {
      results.push({
        rule: h.id,
        level: h.level,
        file: toPosix(e.file),
        line: e.line,
        text: e.text.trim(),
      });
    }
  }

  return { results, scannedLines, files: files.size };
}

/** 按「文件:行号:规则ID」去重（同一条改动可能同时出现在 git diff 与 --cached） */
function dedupe(results) {
  const seen = new Set();
  const out = [];
  for (const r of results) {
    const k = `${r.file}:${r.line}:${r.rule}`;
    if (seen.has(k)) continue;
    seen.add(k);
    out.push(r);
  }
  return out;
}

// ---------------------------------------------------------------------------
// 输出
// ---------------------------------------------------------------------------

function plural(n, word) {
  return `${n} ${word}${n === 1 ? '' : 's'}`;
}

/** 人类可读输出：命中明细 + 结论行（--quiet 只留结论行） */
function formatHuman(results, stats, quiet) {
  const errors = results.filter((r) => r.level === 'ERROR');
  const warnings = results.filter((r) => r.level === 'WARN');
  const out = [];

  if (!quiet) {
    for (const r of results) {
      const tag = r.level === 'ERROR' ? '[ERROR]' : '[WARN ]';
      out.push(`${tag} ${r.rule.padEnd(9)} ${r.file}:${r.line}`);
      out.push(`        ${r.text}`);
    }
    if (results.length) out.push('');
  }

  const summary = `scanned ${stats.scannedLines} added lines in ${stats.files} files (${stats.rangeLabel})`;
  if (errors.length === 0) {
    out.push(`✓ clean  —  ${summary}`);
  } else {
    out.push(`✗ ${plural(errors.length, 'error')}, ${plural(warnings.length, 'warning')}  —  ${summary}`);
  }
  return out.join('\n');
}

/** 构造 JSON 报告对象（--json 输出与 --report 落盘共用同一结构） */
function formatJson(results, stats, rangeLabel) {
  const errors = results.filter((r) => r.level === 'ERROR');
  const warnings = results.filter((r) => r.level === 'WARN');
  const pick = (r) => ({ rule: r.rule, level: r.level, file: r.file, line: r.line, text: r.text });
  return {
    version: 1,
    generatedAt: new Date().toISOString(),
    range: rangeLabel,
    scannedLines: stats.scannedLines,
    scannedFiles: stats.files,
    errorCount: errors.length,
    warningCount: warnings.length,
    errors: errors.map(pick),
    warnings: warnings.map(pick),
    rules: RULES.map((r) => ({ id: r.id, level: r.level, desc: r.desc })),
  };
}

/** --report <path>：写 JSON 报告，自动创建父目录 */
function writeReport(reportPath, payload) {
  const dir = path.dirname(path.resolve(reportPath));
  fs.mkdirSync(dir, { recursive: true });
  fs.writeFileSync(path.resolve(reportPath), JSON.stringify(payload, null, 2) + '\n', 'utf8');
}

function printRules() {
  const errCount = RULES.filter((r) => r.level === 'ERROR').length;
  const warnCount = RULES.length - errCount;
  console.log(`禁改区规则（${RULES.length} 条 = ${errCount} ERROR + ${warnCount} WARN）`);
  for (const r of RULES) {
    console.log(`${r.level === 'ERROR' ? '[ERROR]' : '[WARN ]'} ${r.id.padEnd(13)} ${r.desc}`);
  }
}

function printHelp() {
  console.log(`用法: node scripts/check-i18n-safety.mjs [选项] [rev-range]

位置参数:
  <rev1>..<rev2>        扫指定提交范围
  <rev>                 等价于 <rev>^..<rev>

选项:
  --commit <rev>        扫单个提交
  --file <path>         直接扫文件（可重复；传目录则递归；不走 git）
  --selftest            跑内置正反例自检
  --list-rules          列出全部规则后退出 0
  --json                JSON 输出
  --report <path>       同时把 JSON 报告写到指定路径（用于 SAFE-03 留存）
  --quiet               只输出结论行
  --help                帮助

默认（无参数）：扫工作区 = git diff（未暂存）+ git diff --cached（已暂存）+ 未跟踪新文件

范围过滤: 只扫 .ts/.tsx/.js/.jsx/.mjs/.cjs；跳过 node_modules/、out/、dist/；
         跳过纯注释行；豁免检查器自身（其内建自检样本为故意违规）
退出码: 0 干净或仅 WARN；1 存在 ERROR；2 脚本自身异常
契约:   seanime-web/src/lib/i18n/SAFETY.md`);
}

// ---------------------------------------------------------------------------
// 内置自检
// ---------------------------------------------------------------------------

function selftest() {
  const badMiss = [];
  for (const line of SELFTEST_BAD) {
    const hits = runRules(line).filter((h) => h.level === 'ERROR');
    if (hits.length === 0) badMiss.push(line);
  }

  const goodFalsePositives = [];
  for (const line of SELFTEST_GOOD) {
    const hits = runRules(line);
    if (hits.length > 0) goodFalsePositives.push({ line, hits });
  }

  const knownHits = runRules(SELFTEST_KNOWN_UNCOVERED).filter((h) => h.level === 'ERROR');

  const badOk = SELFTEST_BAD.length - badMiss.length;
  const goodOk = SELFTEST_GOOD.length - goodFalsePositives.length;
  const knownOk = knownHits.length === 0;
  const passed = badMiss.length === 0 && goodFalsePositives.length === 0 && knownOk;

  const out = [];
  out.push('=== 内置自检 (selftest) ===');
  out.push(`BAD 正例命中:      ${badOk}/${SELFTEST_BAD.length}`);
  out.push(`GOOD 反例零误报:   ${goodOk}/${SELFTEST_GOOD.length}`);
  out.push(
    `KNOWN_UNCOVERED:   已知不覆盖，需人工判断 —— ${knownOk ? '未命中 (符合预期)' : '意外命中: ' + knownHits.map((h) => h.id).join(', ')}`
  );
  out.push(`样本: ${SELFTEST_KNOWN_UNCOVERED}`);

  if (badMiss.length) {
    out.push('');
    out.push(`✗ 未命中的 BAD 正例（${badMiss.length}）:`);
    for (const l of badMiss) out.push(`    ${l}`);
  }
  if (goodFalsePositives.length) {
    out.push('');
    out.push(`✗ 误报的 GOOD 反例（${goodFalsePositives.length}）:`);
    for (const g of goodFalsePositives) {
      out.push(`    ${g.line}  ← ${g.hits.map((h) => h.id).join(', ')}`);
    }
  }

  out.push('');
  out.push(passed ? '✓ selftest passed' : '✗ selftest failed');
  console.log(out.join('\n'));
  return passed ? 0 : 1;
}

// ---------------------------------------------------------------------------
// 主流程
// ---------------------------------------------------------------------------

function main() {
  const opts = parseArgs(process.argv.slice(2));

  if (opts.help) {
    printHelp();
    return 0;
  }
  if (opts.listRules) {
    printRules();
    return 0;
  }
  if (opts.selftest) {
    return selftest();
  }

  let entries = [];
  let rangeLabel = '';

  if (opts.files.length) {
    entries = collectFileLines(opts.files);
    rangeLabel = '--file';
  } else if (opts.commit) {
    const r = `${opts.commit}^..${opts.commit}`;
    entries = collectDiffLines(gitDiff(['diff', '--unified=0', '--no-color', r]));
    rangeLabel = r;
  } else if (opts.positional.length) {
    let r = opts.positional[0];
    if (!r.includes('..')) r = `${r}^..${r}`;
    entries = collectDiffLines(gitDiff(['diff', '--unified=0', '--no-color', r]));
    rangeLabel = r;
  } else {
    // 工作区：未暂存 + 已暂存 + 未跟踪新文件，三者按 文件:行号 去重
    const unstaged = collectDiffLines(gitDiff(['diff', '--unified=0', '--no-color']));
    const staged = collectDiffLines(gitDiff(['diff', '--unified=0', '--no-color', '--cached']));
    const untracked = collectUntrackedLines();
    entries = unstaged.concat(staged, untracked);
    rangeLabel = 'worktree';
  }

  // 去重键必须带上**行内容**，不能只用 文件:行号。
  //
  // 原因：同一个 文件:行号 可能同时存在「已暂存版本」与「未暂存版本」，且两者内容
  // 不同 —— 典型场景是先 `git add` 了违规版本，随后又把工作区改回干净。此时
  // `git diff` 给出干净版本、`git diff --cached` 给出违规版本。若只按 文件:行号
  // 去重，先入者（干净版本）会挤掉后入者，脚本报 clean，而**真正会被提交的是索引里
  // 那份违规内容** —— 门禁被静默绕过。带上内容后两个版本都会被扫，任一出错即报错。
  const seen = new Set();
  const uniqueEntries = [];
  for (const e of entries) {
    const k = `${toPosix(e.file)}:${e.line}:${e.text}`;
    if (seen.has(k)) continue;
    seen.add(k);
    uniqueEntries.push(e);
  }

  const stats = scanEntries(uniqueEntries);
  stats.rangeLabel = rangeLabel;
  const results = dedupe(stats.results);
  const payload = formatJson(results, stats, rangeLabel);

  if (opts.report) writeReport(opts.report, payload);

  if (opts.json) {
    console.log(JSON.stringify(payload, null, 2));
  } else {
    console.log(formatHuman(results, stats, opts.quiet));
  }

  return payload.errorCount > 0 ? 1 : 0;
}

try {
  process.exitCode = main();
} catch (err) {
  console.error(`[check-i18n-safety] 内部错误: ${err && err.message ? err.message : String(err)}`);
  process.exitCode = 2;
}
