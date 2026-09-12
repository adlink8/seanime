import type { Dictionary } from "../../types"

/**
 * autodownloader 词条表
 * 归属：自动下载（auto-downloader/）
 * 独占：09-02
 *
 * 只读复用：common.* / library.* / settings.* / entry.* 等同名文案按契约复用，不新建重复词条。
 */
export const autodownloaderDictionary = {
    // —— 页面标题（auto-downloader/page.tsx）——
    "autodownloader.page.subtitle": "有新剧集发布时自动下载。", // auto-downloader/page.tsx

    // —— 主页面：标签页与规则列表（_containers/autodownloader-page.tsx）——
    "autodownloader.tab.rules": "规则", // _containers/autodownloader-page.tsx
    "autodownloader.tab.profiles": "配置文件", // autodownloader-page.tsx, autodownloader-rule-item.tsx
    "autodownloader.alert.disabled": "自动下载器当前已停用。", // autodownloader-page.tsx
    "autodownloader.alert.enable_here": "在此启用。", // autodownloader-page.tsx
    "autodownloader.rules_desc": "规则允许你根据设定的参数自动下载新剧集。", // autodownloader-page.tsx
    "autodownloader.action.one_series": "单个番剧", // autodownloader-page.tsx
    "autodownloader.action.multiple_series": "一次添加多个番剧", // autodownloader-page.tsx
    "autodownloader.action.check_rss": "检查 RSS 订阅", // autodownloader-page.tsx
    "autodownloader.action.remove_finished": "移除已完结媒体", // autodownloader-page.tsx（菜单项与确认弹窗标题共用）
    "autodownloader.confirm_remove_finished_desc": "此操作将移除所有媒体已完结（不再放送）的规则。确定要继续吗？", // autodownloader-page.tsx
    "autodownloader.empty.no_rules": "暂无规则", // autodownloader-page.tsx
    "autodownloader.toast.interval_min": "检查间隔已调整为至少 15 分钟", // autodownloader-page.tsx settingsSchema

    // —— 主页面：设置表单（autodownloader-page.tsx）——
    "autodownloader.field.use_debrid": "使用 Debrid 服务", // autodownloader-page.tsx
    "autodownloader.alert.debrid_title": "自动下载器已停用", // autodownloader-page.tsx
    "autodownloader.alert.debrid_desc": "Debrid 服务未启用或未配置。请在设置中启用。", // autodownloader-page.tsx
    "autodownloader.empty.no_extensions": "未找到扩展", // autodownloader-page.tsx
    "autodownloader.field.download_immediately": "立即下载剧集", // autodownloader-page.tsx
    "autodownloader.help.download_immediately": "关闭后，种子将加入队列等待下载。", // autodownloader-page.tsx
    "autodownloader.field.interval": "检查间隔", // autodownloader-page.tsx
    "autodownloader.help.interval": "检查新剧集的频率。", // autodownloader-page.tsx
    "autodownloader.field.every": "每", // autodownloader-page.tsx NumberInput 左缀
    "autodownloader.field.strict_season_check": "严格季度校验", // autodownloader-page.tsx
    "autodownloader.help.strict_season_check": "启用后，种子名称与媒体标题需包含相同的季度编号。这可能导致漏检。", // autodownloader-page.tsx
    "autodownloader.modal.create_rules": "批量新建规则", // autodownloader-page.tsx
    "autodownloader.batch_intro": "一次创建多条规则。除目标文件夹外，每条规则都将使用相同的参数。默认剧集类型为「最近发布」。", // autodownloader-page.tsx

    // —— 规则表单（_containers/autodownloader-rule-form.tsx）——
    "autodownloader.toast.episode_number_required": "必须至少指定一个剧集编号", // autodownloader-rule-form.tsx
    "autodownloader.empty.no_media": "媒体库中未找到媒体", // autodownloader-rule-form.tsx, autodownloader-batch-rule-form.tsx
    "autodownloader.toast.check_fields": "发生错误，请检查各字段。", // autodownloader-rule-form.tsx, autodownloader-batch-rule-form.tsx
    "autodownloader.action.delete_rule": "删除此规则", // autodownloader-rule-form.tsx
    "autodownloader.filter.airing_only": "仅显示放送中", // autodownloader-rule-form.tsx
    "autodownloader.filter.airing_upcoming": "显示放送中与即将播出", // autodownloader-rule-form.tsx
    "autodownloader.filter.all": "显示全部", // autodownloader-rule-form.tsx
    "autodownloader.empty.media_not_in_library": "该媒体不在你的媒体库中", // autodownloader-rule-form.tsx
    "autodownloader.status.no_longer_airing": "已完结（不再放送）", // autodownloader-rule-form.tsx, autodownloader-rule-item.tsx
    "autodownloader.field.comparison_title": "比较标题", // autodownloader-rule-form.tsx, autodownloader-batch-rule-form.tsx
    "autodownloader.field.destination": "目标文件夹", // autodownloader-rule-form.tsx, autodownloader-batch-rule-form.tsx
    "autodownloader.help.destination": "文件将保存到本地媒体库中的此文件夹", // 同上
    "autodownloader.field.search_type": "搜索方式", // autodownloader-rule-form.tsx, autodownloader-batch-rule-form.tsx
    "autodownloader.option.most_likely": "可能性最高", // autodownloader-rule-form.tsx, autodownloader-batch-rule-form.tsx
    "autodownloader.option.most_likely_desc": "将解析种子名称并使用比较算法进行分析", // 同上
    "autodownloader.option.exact_match": "精确匹配", // 同上
    "autodownloader.option.exact_match_desc": "种子名称必须包含你设置的比较标题（不区分大小写）", // 同上
    "autodownloader.hint.also_titles": "还将使用以下标题：", // autodownloader-rule-form.tsx
    "autodownloader.field.episodes_to_look_for": "要查找的剧集", // autodownloader-rule-form.tsx
    "autodownloader.option.recent": "最近发布", // autodownloader-rule-form.tsx, autodownloader-rule-item.tsx
    "autodownloader.option.recent_desc": "你尚未观看的新剧集", // autodownloader-rule-form.tsx
    "autodownloader.option.selected_desc": "仅查找媒体库中尚未拥有的指定剧集", // autodownloader-rule-form.tsx
    "autodownloader.field.episode_numbers": "剧集编号", // autodownloader-rule-form.tsx
    "autodownloader.field.episode_offset": "剧集编号绝对偏移", // autodownloader-rule-form.tsx
    "autodownloader.help.episode_offset": "例如，若发布组从 13 而非 1 开始编号，则将此值设为 12。", // autodownloader-rule-form.tsx
    "autodownloader.section.constraints": "约束条件", // autodownloader-rule-form.tsx, autodownloader-batch-rule-form.tsx
    "autodownloader.action.run_simulation": "运行模拟", // autodownloader-rule-form.tsx
    "autodownloader.action.create": "创建", // autodownloader-rule-form.tsx, autodownloader-profile-form.tsx, autodownloader-batch-rule-form.tsx
    "autodownloader.action.update": "更新", // autodownloader-rule-form.tsx, autodownloader-profile-form.tsx
    "autodownloader.modal.result": "模拟结果", // autodownloader-rule-form.tsx
    "autodownloader.modal.simulation_prefix": "以下规则的模拟结果：", // autodownloader-rule-form.tsx
    "autodownloader.modal.simulation_suffix": "（ID: {id}）", // autodownloader-rule-form.tsx
    "autodownloader.modal.check_logs": "详情请查看服务器日志。", // autodownloader-rule-form.tsx

    // —— 共享字段（_containers/autodownloader-shared-fields.tsx）——
    "autodownloader.section.release_groups": "发布组", // autodownloader-shared-fields.tsx, autodownloader-profile-form.tsx
    "autodownloader.help.release_groups": "要查找的发布组列表。留空则接受任何发布组。规则可覆盖此设置。", // autodownloader-profile-form.tsx
    "autodownloader.help.release_groups_short": "要查找的发布组列表。留空则接受任何发布组。", // autodownloader-shared-fields.tsx
    "autodownloader.placeholder.release_group": "例如：SubsPlease", // autodownloader-shared-fields.tsx
    "autodownloader.help.resolutions": "要查找的分辨率列表。留空则接受最高分辨率。", // autodownloader-shared-fields.tsx
    "autodownloader.help.resolutions_sortable": "拖拽调整顺序。将选用第一个匹配的分辨率。规则可覆盖此设置。", // autodownloader-profile-form.tsx
    "autodownloader.placeholder.resolution": "例如：1080p", // autodownloader-shared-fields.tsx
    "autodownloader.section.terms": "视频、音频、来源、字幕", // autodownloader-shared-fields.tsx
    "autodownloader.help.additional_terms": "种子名称必须包含所有选项才会被接受。", // autodownloader-shared-fields.tsx 加粗段
    "autodownloader.help.additional_terms_more": "每个选项内可用逗号分隔多个变体。（不区分大小写）", // autodownloader-shared-fields.tsx
    "autodownloader.help.exclude_terms": "种子包含所有选项时将被拒绝。", // autodownloader-shared-fields.tsx 加粗段
    "autodownloader.help.exclude_terms_more": "包含其中任一关键词的种子都会被拒绝。（不区分大小写）", // autodownloader-shared-fields.tsx
    "autodownloader.placeholder.terms": "例如：H265,H.265,H 265,x265", // autodownloader-shared-fields.tsx
    "autodownloader.separator.and": "且", // autodownloader-shared-fields.tsx, autodownloader-batch-rule-form.tsx
    "autodownloader.section.profile": "配置文件", // autodownloader-shared-fields.tsx
    "autodownloader.help.profile": "选择配置文件以应用共享筛选。本地筛选将覆盖配置文件筛选。", // autodownloader-shared-fields.tsx
    "autodownloader.field.select_providers": "选择提供商", // autodownloader-shared-fields.tsx, autodownloader-profile-form.tsx
    "autodownloader.help.providers": "选择要查找的特定提供商。留空则使用默认提供商。", // 同上
    "autodownloader.field.select_profile": "选择配置文件", // autodownloader-shared-fields.tsx
    "autodownloader.empty.no_profile": "未找到配置文件", // autodownloader-shared-fields.tsx

    // —— 关键词建议标签（shared-fields / profile-form，技术专名保留原文）——
    "autodownloader.suggestion.hevc": "HEVC / H.265", // autodownloader-shared-fields.tsx
    "autodownloader.suggestion.avc": "AVC / H.264", // autodownloader-shared-fields.tsx
    "autodownloader.suggestion.dolby_vision": "Dolby Vision", // autodownloader-shared-fields.tsx（专有名词保留）
    "autodownloader.suggestion.remux": "Remux", // autodownloader-shared-fields.tsx（专有名词保留）
    "autodownloader.suggestion.opus": "Opus", // autodownloader-shared-fields.tsx（专有名词保留）
    "autodownloader.suggestion.dual_audio": "双语音轨", // autodownloader-shared-fields.tsx
    "autodownloader.suggestion.dubbed": "配音版", // autodownloader-shared-fields.tsx
    "autodownloader.suggestion.cam": "枪版", // autodownloader-shared-fields.tsx

    // —— 配置文件表单（_containers/autodownloader-profile-form.tsx）——
    "autodownloader.field.global": "全局", // autodownloader-profile-form.tsx, autodownloader-profiles.tsx
    "autodownloader.help.global": "自动将此配置文件应用到所有规则", // autodownloader-profile-form.tsx
    "autodownloader.field.minimum_score": "最低评分", // autodownloader-profile-form.tsx
    "autodownloader.help.minimum_score": "评分低于此值的种子将被拒绝", // autodownloader-profile-form.tsx
    "autodownloader.section.delay": "延迟", // autodownloader-profile-form.tsx
    "autodownloader.help.delay": "下载前等待更好的发布版本。延迟期将在找到首个匹配后开始计时。若规则分配了更长的延迟配置或全局应用了更长延迟，则忽略此处的设置。", // autodownloader-profile-form.tsx
    "autodownloader.field.delay_minutes": "延迟时间", // autodownloader-profile-form.tsx
    "autodownloader.help.delay_minutes": "下载前等待的分钟数", // autodownloader-profile-form.tsx
    "autodownloader.field.skip_delay_score": "跳过延迟评分", // autodownloader-profile-form.tsx
    "autodownloader.help.skip_delay_score": "种子评分超过此值时跳过延迟", // autodownloader-profile-form.tsx
    "autodownloader.section.conditions": "条件", // autodownloader-profile-form.tsx
    "autodownloader.help.conditions": "添加条件以筛选种子或调整其评分。", // autodownloader-profile-form.tsx
    "autodownloader.action.add_release_group": "添加发布组", // autodownloader-profile-form.tsx
    "autodownloader.action.add_resolution": "添加分辨率", // autodownloader-profile-form.tsx
    "autodownloader.action.add_condition": "添加条件", // autodownloader-profile-form.tsx
    "autodownloader.placeholder.condition_term": "例如：Blu-Ray、BluRay 或 \\b(group)\\bi", // autodownloader-profile-form.tsx
    "autodownloader.help.condition_term": "逗号分隔的不区分大小写的值，或正则表达式", // autodownloader-profile-form.tsx
    "autodownloader.field.action": "动作", // autodownloader-profile-form.tsx Select 标签
    "autodownloader.option.block": "阻止", // autodownloader-profile-form.tsx
    "autodownloader.option.require": "必须满足", // autodownloader-profile-form.tsx
    "autodownloader.field.regex": "正则表达式", // autodownloader-profile-form.tsx

    // —— 配置文件列表（_containers/autodownloader-profiles.tsx）——
    "autodownloader.profiles_desc": "配置文件允许你定义一组可应用到规则上的筛选条件。", // autodownloader-profiles.tsx
    "autodownloader.action.new_profile": "新建配置文件", // autodownloader-profiles.tsx（按钮与弹窗标题共用）
    "autodownloader.empty.no_profiles": "尚未创建任何配置文件。", // autodownloader-profiles.tsx
    "autodownloader.action.edit": "编辑", // autodownloader-profiles.tsx
    "autodownloader.modal.edit_profile": "编辑配置文件：{name}", // autodownloader-profiles.tsx
    "autodownloader.modal.delete_profile": "删除配置文件", // autodownloader-profiles.tsx
    "autodownloader.confirm_delete_profile": "确定要删除配置文件「{name}」吗？此操作无法撤销。", // autodownloader-profiles.tsx
    "autodownloader.status.conditions": "{count} 个条件", // autodownloader-profiles.tsx
    "autodownloader.status.min_delay": "{minutes} 分钟延迟", // autodownloader-profiles.tsx

    // —— 下载队列（_containers/autodownloader-queue.tsx）——
    "autodownloader.queue.file_downloaded": "文件已下载", // autodownloader-queue.tsx
    "autodownloader.queue.manual_action": "需要手动操作", // autodownloader-queue.tsx
    "autodownloader.queue.delayed": "延迟中", // autodownloader-queue.tsx
    "autodownloader.queue.delayed_for": "延迟 {time}。", // autodownloader-queue.tsx
    "autodownloader.queue.added": "添加于 {time}", // autodownloader-queue.tsx
    "autodownloader.queue.not_scanned": "尚未扫描", // autodownloader-queue.tsx
    "autodownloader.queue.desc": "队列显示等待下载或扫描的项目。", // autodownloader-queue.tsx
    "autodownloader.queue.empty": "队列为空", // autodownloader-queue.tsx

    // —— 规则条目（_components/autodownloader-rule-item.tsx）——
    "autodownloader.item.rule_for": "规则：", // autodownloader-rule-item.tsx
    "autodownloader.item.not_in_library": "此动漫不在你的媒体库中", // autodownloader-rule-item.tsx
    "autodownloader.item.edit_rule": "编辑规则", // autodownloader-rule-item.tsx
    "autodownloader.option.select_episodes": "选择剧集", // autodownloader-rule-item.tsx getEpisodeTypeName

    // —— 批量规则表单（_containers/autodownloader-batch-rule-form.tsx）——
    "autodownloader.action.add_currently_watching": "添加全部正在观看", // autodownloader-batch-rule-form.tsx
    "autodownloader.action.add_upcoming": "添加全部即将播出", // autodownloader-batch-rule-form.tsx
    "autodownloader.help.comparison_title_batch": "用于比较。使用「精确匹配」时，请填写最可能出现在种子名称中的标题。", // autodownloader-batch-rule-form.tsx
} satisfies Dictionary
