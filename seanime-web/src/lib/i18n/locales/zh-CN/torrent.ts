import type { Dictionary } from "../../types"

/**
 * torrent 词条表
 * 归属：种子客户端（torrent-client/torrent-list/qbittorrent）
 * 独占：09-01
 *
 * 只读复用：common.* / library.* 等同名文案按契约复用，不新建重复词条。
 * 已复用：common.action.retry / player.cast.pause / library.filter.all /
 *   library.filter.status / library.explorer.size / library.common.file /
 *   player.prefs.tab_general / media.field.progress / home.option.name /
 *   home.settings.add / entry.metadata.remove / entry.torrent_filter.seeders_unit /
 *   mpv.shader.open_folder / manga.manual_match.unknown / manga.settings.open /
 *   settings.field.category
 */
export const torrentDictionary = {
    // —— 内置客户端控制台：页面骨架（torrent-client/page.tsx）——
    "torrent.client.inactive_title": "Seanime 内置种子客户端未启用", // torrent-client/page.tsx LuffyError title
    "torrent.client.inactive_hint": "请先将 Seanime 设为默认种子客户端，即可使用此控制台。", // torrent-client/page.tsx 提示文案
    "torrent.client.load_error": "无法加载种子列表", // torrent-client/page.tsx LuffyError title
    "torrent.client.title": "种子客户端", // torrent-client/page.tsx h2
    "torrent.client.subtitle": "管理直接在 Seanime 中运行的下载任务。", // torrent-client/page.tsx 副标题
    "torrent.client.add_torrent": "添加种子", // torrent-client/page.tsx 按钮 + 弹窗 title
    "torrent.client.pause_all": "全部暂停", // torrent-client/page.tsx 按钮
    "torrent.client.resume_all": "全部恢复", // torrent-client/page.tsx 按钮
    "torrent.client.total_torrents": "{count} 个种子", // torrent-client/page.tsx 侧栏统计
    "torrent.client.filter_placeholder": "筛选种子", // torrent-client/page.tsx placeholder

    // —— 内置客户端控制台：工具栏 / 弹窗（torrent-client/page.tsx）——
    "torrent.client.global_speed_limits": "全局速度限制", // torrent-client/page.tsx Popover 标题
    "torrent.client.download_limit": "下载 (KB/s)", // torrent-client/page.tsx label
    "torrent.client.upload_limit": "上传 (KB/s)", // torrent-client/page.tsx label
    "torrent.client.global_download_speed": "全局下载速度", // torrent-client/page.tsx title
    "torrent.client.global_upload_speed": "全局上传速度", // torrent-client/page.tsx title
    "torrent.client.showing_range": "显示 {start}–{end}，共 {total} 个", // torrent-client/page.tsx 分页统计
    "torrent.client.per_page": "每页 {count} 条", // torrent-client/page.tsx 每页条数选项
    "torrent.client.no_match": "当前视图没有匹配的种子。", // torrent-client/page.tsx 空状态
    "torrent.client.add_description": "向 Seanime 种子客户端添加磁力链接。", // torrent-client/page.tsx 弹窗 description
    "torrent.client.start_download": "开始下载", // torrent-client/page.tsx 弹窗按钮
    "torrent.client.magnet_link": "磁力链接", // torrent-client/page.tsx label
    "torrent.client.move_description": "Seanime 将暂停该种子、移动文件并校验数据。", // torrent-client/page.tsx 弹窗 description
    "torrent.client.new_save_path": "新保存路径", // torrent-client/page.tsx label
    "torrent.client.rename_title": "重命名种子", // torrent-client/page.tsx 弹窗 title
    "torrent.client.display_name": "显示名称", // torrent-client/page.tsx label
    "torrent.client.remove_single": "移除种子", // torrent-client/page.tsx 确认弹窗 title
    "torrent.client.remove_multiple": "移除 {count} 个种子", // torrent-client/page.tsx 确认弹窗 title（多选）
    "torrent.client.remove_warning": "已下载的文件也将一并删除，此操作无法撤销。", // torrent-client/page.tsx 确认弹窗 description
    "torrent.client.select_hint": "选择一个种子以查看文件、Tracker 与连接。", // torrent-client/page.tsx Inspector 空状态

    // —— 内置客户端控制台：操作（torrent-client/page.tsx）——
    "torrent.action.resume": "恢复", // torrent-client/page.tsx Tooltip / 右键菜单
    "torrent.action.move_up": "上移", // torrent-client/page.tsx Tooltip / 右键菜单
    "torrent.action.move_down": "下移", // torrent-client/page.tsx Tooltip / 右键菜单
    "torrent.action.force_start": "强制开始", // torrent-client/page.tsx Tooltip / 右键菜单
    "torrent.action.change_save_path": "更改保存路径", // torrent-client/page.tsx Tooltip / 右键菜单 / 弹窗 title
    "torrent.action.recheck": "重新校验", // torrent-client/page.tsx Tooltip / 右键菜单
    "torrent.action.reannounce": "重新汇报", // torrent-client/page.tsx Tooltip / 右键菜单
    "torrent.action.speed_limits": "速度限制", // torrent-client/page.tsx Tooltip
    "torrent.action.apply_limits": "应用限制", // torrent-client/page.tsx 按钮
    "torrent.action.rename": "重命名", // torrent-client/page.tsx 右键菜单 / Inspector 按钮 / 弹窗按钮
    "torrent.action.move_files": "移动文件", // torrent-client/page.tsx 弹窗按钮
    "torrent.action.sequential": "顺序下载", // torrent-client/page.tsx Inspector 按钮
    "torrent.action.enable_sequential": "启用顺序下载", // torrent-client/page.tsx 右键菜单
    "torrent.action.disable_sequential": "关闭顺序下载", // torrent-client/page.tsx 右键菜单
    "torrent.action.end": "结束做种", // torrent-list/page.tsx Tooltip

    // —— 内置客户端控制台：筛选 / 表格（torrent-client/page.tsx）——
    "torrent.filter.downloading": "下载中", // torrent-client/page.tsx 筛选标签
    "torrent.filter.paused": "已暂停", // torrent-client/page.tsx 筛选标签
    "torrent.filter.active": "活跃", // torrent-client/page.tsx 筛选标签
    "torrent.filter.inactive": "非活跃", // torrent-client/page.tsx 筛选标签
    "torrent.table.seeds_peers": "做种 / 连接", // torrent-client/page.tsx 表头
    "torrent.table.speed": "速度", // torrent-client/page.tsx 表头
    "torrent.table.ratio": "分享率", // torrent-client/page.tsx 表头
    "torrent.table.date_added": "添加时间", // torrent-client/page.tsx 表头 / Inspector 字段
    "torrent.table.save_path": "保存路径", // torrent-client/page.tsx 表头 / Inspector 字段 / 弹窗 label

    // —— 内置客户端控制台：详情面板（torrent-client/page.tsx）——
    "torrent.details.tab_trackers": "Tracker", // torrent-client/page.tsx 标签页
    "torrent.details.tab_peers": "用户", // torrent-client/page.tsx 标签页
    "torrent.details.downloaded": "已下载", // torrent-client/page.tsx Inspector 字段
    "torrent.details.uploaded": "已上传", // torrent-client/page.tsx Inspector 字段
    "torrent.details.queue_position": "队列位置", // torrent-client/page.tsx Inspector 字段
    "torrent.details.priority": "优先级", // torrent-client/page.tsx 文件表头
    "torrent.details.waiting_metadata": "正在等待种子元数据。", // torrent-client/page.tsx 文件空状态
    "torrent.details.skip": "跳过", // torrent-client/page.tsx 文件优先级
    "torrent.details.normal": "正常", // torrent-client/page.tsx 文件优先级
    "torrent.details.high": "高", // torrent-client/page.tsx 文件优先级
    "torrent.details.remove_tracker": "移除 Tracker", // torrent-client/page.tsx Tooltip
    "torrent.details.no_trackers": "暂无 Tracker。", // torrent-client/page.tsx Tracker 空状态
    "torrent.details.address": "地址", // torrent-client/page.tsx 表头
    "torrent.details.client": "客户端", // torrent-client/page.tsx 表头
    "torrent.details.no_peers": "暂无已连接的用户。", // torrent-client/page.tsx Peers 空状态
    "torrent.status.forced": "强制", // torrent-client/page.tsx 状态徽标

    // —— 种子任务列表（torrent-list/page.tsx）——
    "torrent.list.title": "当前种子任务", // torrent-list/page.tsx h2
    "torrent.list.subtitle": "查看正在下载或做种的种子任务", // torrent-list/page.tsx 副标题
    "torrent.list.built_in_client": "内置客户端", // torrent-list/page.tsx 按钮
    "torrent.list.qbittorrent_webui": "qBittorrent Web UI", // torrent-list/page.tsx 按钮
    "torrent.list.console": "种子控制台", // torrent-list/page.tsx 按钮
    "torrent.list.connection_failed": "连接失败", // torrent-list/page.tsx LuffyError title
    "torrent.list.connection_error": "无法连接至种子客户端，请检查设置并确认客户端正在运行。", // torrent-list/page.tsx 错误提示
    "torrent.stats.downloading": "正在下载", // torrent-list/page.tsx 统计
    "torrent.stats.seeding": "正在做种", // torrent-list/page.tsx 统计
    "torrent.list.stop_seeding": "停止做种", // torrent-list/page.tsx 按钮
    "torrent.list.stop_all_seeding": "停止所有做种", // torrent-list/page.tsx 确认弹窗 title
    "torrent.list.stop_all_seeding_desc": "该操作将停止所有已完成种子的做种状态。", // torrent-list/page.tsx 确认弹窗 description
    "torrent.list.confirm": "确定", // torrent-list/page.tsx 按钮
    "torrent.list.filter_by_category": "按分类筛选", // torrent-list/page.tsx placeholder
    "torrent.list.empty": "暂无种子任务", // torrent-list/page.tsx 空状态
    "torrent.list.delete_torrent": "删除种子", // torrent-list/page.tsx 确认弹窗 title
    "torrent.list.delete_confirm": "此操作无法撤销，确定要删除吗？", // torrent-list/page.tsx 确认弹窗 description
    "torrent.sort.newest": "最新", // torrent-list/page.tsx 排序切换
    "torrent.sort.oldest": "最旧", // torrent-list/page.tsx 排序切换
    "torrent.sort.name_asc": "名称 (A-Z)", // torrent-list/page.tsx 排序切换
    "torrent.sort.name_desc": "名称 (Z-A)", // torrent-list/page.tsx 排序切换

    // —— qBittorrent Web UI（qbittorrent/page.tsx）——
    "torrent.qbittorrent.description": "访问内置 qBittorrent 客户端的 Web UI。", // qbittorrent/page.tsx 副标题
    "torrent.qbittorrent.open_in_browser": "在浏览器中打开", // qbittorrent/page.tsx 按钮
    "torrent.qbittorrent.open_new_tab_title": "在新标签页中打开", // qbittorrent/page.tsx h3
    "torrent.qbittorrent.browser_policy": "由于浏览器安全策略（COEP 与点击劫持防护），无法在此页面的 iframe 中加载嵌入式客户端。", // qbittorrent/page.tsx 说明
    "torrent.qbittorrent.open_webui": "打开 qBittorrent Web UI", // qbittorrent/page.tsx 按钮
} satisfies Dictionary
