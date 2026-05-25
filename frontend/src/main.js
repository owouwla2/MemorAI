import './style.css';
import './app.css';
import Chart from 'chart.js/auto';

// ===== SVG 图标系统（Lucide 风格，统一 currentColor） =====
const ICONS = {
    cpu: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="16" height="16" x="4" y="4" rx="2"/><rect width="6" height="6" x="9" y="9" rx="1"/><path d="M15 2v2"/><path d="M15 20v2"/><path d="M2 15h2"/><path d="M2 9h2"/><path d="M20 15h2"/><path d="M20 9h2"/><path d="M9 2v2"/><path d="M9 20v2"/></svg>',
    refresh: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/><path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/><path d="M3 21v-5h5"/></svg>',
    moon: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>',
    sun: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="4"/><path d="M12 2v2"/><path d="M12 20v2"/><path d="m4.93 4.93 1.41 1.41"/><path d="m17.66 17.66 1.41 1.41"/><path d="M2 12h2"/><path d="M20 12h2"/><path d="m6.34 17.66-1.41 1.41"/><path d="m19.07 4.93-1.41 1.41"/></svg>',
    layout: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="3" rx="2"/><path d="M3 9h18"/><path d="M9 21V9"/></svg>',
    list: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="8" x2="21" y1="6" y2="6"/><line x1="8" x2="21" y1="12" y2="12"/><line x1="8" x2="21" y1="18" y2="18"/><line x1="3" x2="3.01" y1="6" y2="6"/><line x1="3" x2="3.01" y1="12" y2="12"/><line x1="3" x2="3.01" y1="18" y2="18"/></svg>',
    activity: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-2.48a2 2 0 0 0-1.93 1.46l-2.35 8.36a.5.5 0 0 1-.96 0L9.24 2.18a.5.5 0 0 0-.96 0l-2.35 8.36A2 2 0 0 1 4 12H2"/></svg>',
    rocket: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 0 0-2.91-.09z"/><path d="m12 15-3-3a22 22 0 0 1 2-3.95A12.88 12.88 0 0 1 22 2c0 2.72-.78 7.5-6 11a22.35 22.35 0 0 1-4 2z"/><path d="M9 12H4s.55-3.03 2-4c1.62-1.08 5 0 5 0"/><path d="M12 15v5s3.03-.55 4-2c1.08-1.62 0-5 0-5"/></svg>',
    bot: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 8V4H8"/><rect width="16" height="12" x="4" y="8" rx="2"/><path d="M2 14h2"/><path d="M20 14h2"/><path d="M15 13v2"/><path d="M9 13v2"/></svg>',
    settings: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/></svg>',
    chevronRight: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>',
    play: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"/></svg>',
    pieChart: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21.21 15.89A10 10 0 1 1 8 2.83"/><path d="M22 12A10 10 0 0 0 12 2v10z"/></svg>',
    lineChart: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 3v18h18"/><path d="m19 9-5 5-4-4-3 3"/></svg>',
    download: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" x2="12" y1="15" y2="3"/></svg>',
    trash: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>',
    send: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m22 2-7 20-4-9-9-4z"/><path d="M22 2 11 13"/></svg>',
    sparkles: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m12 3-1.91 5.79a2 2 0 0 1-1.3 1.3L3 12l5.79 1.91a2 2 0 0 1 1.3 1.3L12 21l1.91-5.79a2 2 0 0 1 1.3-1.3L21 12l-5.79-1.91a2 2 0 0 1-1.3-1.3z"/></svg>',
    fileText: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>',
    clipboard: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="8" height="4" x="8" y="2" rx="1"/><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/></svg>',
    alertTriangle: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg>',
    checkCircle: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>',
    box: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" x2="12" y1="22.08" y2="12"/></svg>',
};
function icon(name, cls = '') {
    const svg = ICONS[name] || '';
    if (!svg) return '';
    if (cls) return svg.replace('<svg ', `<svg class="${cls}" `);
    return svg;
}

// ===== 状态 =====
const state = {
    currentTab: 'overview', // overview | processes | monitor | startup | ai | settings
    theme: localStorage.getItem('theme') || 'dark',
    report: null,
    expandedGroups: new Set(),
    // AI 对话
    chat: {
        messages: [],     // [{role: 'user'|'assistant', content: string}]
        sending: false,
        snapshotLoaded: false,
        input: '',        // 当前输入框内容（保持渲染时不丢失）
    },
    // 监控
    monitor: {
        running: false,
        duration: 30,
        interval: 1,
        progress: 0,
        progressText: '',
        result: null,        // 完整结果
        liveSamples: [],     // 采样过程中的实时数据
    },
    // 自启项
    startup: {
        items: null,
        loading: false,
        isAdmin: false,
    },
};

// ===== 工具函数 =====
function fmtMB(mb) {
    if (mb >= 1024) return (mb / 1024).toFixed(2) + ' GB';
    return mb.toFixed(0) + ' MB';
}
function fmtPct(p) { return p.toFixed(1) + '%'; }

function escapeHtml(s) {
    return String(s == null ? '' : s)
        .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;').replace(/'/g, '&#39;');
}

function renderMarkdown(md) {
    if (!md) return '';
    const lines = md.split('\n');
    let html = '', inList = false, listType = '';
    function closeList() {
        if (inList) { html += listType === 'ol' ? '</ol>' : '</ul>'; inList = false; listType = ''; }
    }
    function inline(s) {
        s = escapeHtml(s);
        s = s.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
        s = s.replace(/`([^`]+)`/g, '<code>$1</code>');
        return s;
    }
    for (let raw of lines) {
        const line = raw.trimEnd();
        if (line === '') { closeList(); continue; }
        let m = line.match(/^(#{1,6})\s+(.*)$/);
        if (m) {
            closeList();
            const tag = m[1].length <= 2 ? 'h2' : 'h3';
            html += `<${tag}>${inline(m[2])}</${tag}>`;
            continue;
        }
        if (/^\d+[.)]\s+/.test(line)) {
            if (!inList || listType !== 'ol') { closeList(); html += '<ol>'; inList = true; listType = 'ol'; }
            html += `<li>${inline(line.replace(/^\d+[.)]\s+/, ''))}</li>`;
            continue;
        }
        if (/^[-*]\s+/.test(line)) {
            if (!inList || listType !== 'ul') { closeList(); html += '<ul>'; inList = true; listType = 'ul'; }
            html += `<li>${inline(line.replace(/^[-*]\s+/, ''))}</li>`;
            continue;
        }
        closeList();
        html += `<p>${inline(line)}</p>`;
    }
    closeList();
    return html;
}

function callBackend(method, ...args) {
    if (!window.go || !window.go.main || !window.go.main.App) {
        return Promise.reject(new Error('Wails 后端未就绪'));
    }
    const fn = window.go.main.App[method];
    if (!fn) return Promise.reject(new Error(`后端方法不存在: ${method}`));
    return fn(...args);
}

// Wails 事件监听
function listenEvent(eventName, handler) {
    if (window.runtime && window.runtime.EventsOn) {
        window.runtime.EventsOn(eventName, handler);
    }
}

function toast(msg, type = '') {
    const t = document.createElement('div');
    t.className = 'toast ' + type;
    t.textContent = msg;
    document.body.appendChild(t);
    setTimeout(() => t.remove(), 3500);
}

// ===== 主渲染 =====
function applyTheme() {
    document.documentElement.dataset.theme = state.theme;
}

function render() {
    applyTheme();
    const root = document.querySelector('#app');

    // HUD：从 report 里取核心指标
    const sys = state.report ? state.report.system : null;
    const cpu = state.report ? state.report.cpu : null;
    const totalProc = state.report ? state.report.total : 0;
    const pcCls = sys ? pctClass(sys.usedPercent) : '';
    const pcText = pcCls === 'danger' ? 'danger' : pcCls === 'warn' ? 'warn' : 'good';
    const cpuPct = cpu ? cpu.usedPercent : 0;
    const cpuCls = pctClass(cpuPct);
    const cpuText = cpuCls === 'danger' ? 'danger' : cpuCls === 'warn' ? 'warn' : 'good';

    const hud = sys ? `
        <div class="hud">
            <div class="hud-item">
                <span class="hud-label">CPU</span>
                <span class="hud-value ${cpuText}">${cpuPct.toFixed(1)}%</span>
            </div>
            <div class="hud-item">
                <span class="hud-label">内存</span>
                <span class="hud-value ${pcText}">${fmtPct(sys.usedPercent)}</span>
            </div>
            <div class="hud-item">
                <span class="hud-label">已用</span>
                <span class="hud-value">${sys.usedGB.toFixed(2)} GB</span>
            </div>
            <div class="hud-item">
                <span class="hud-label">可用</span>
                <span class="hud-value accent">${sys.availableGB.toFixed(2)} GB</span>
            </div>
            <div class="hud-item">
                <span class="hud-label">总计</span>
                <span class="hud-value">${sys.totalGB.toFixed(2)} GB</span>
            </div>
            <div class="hud-item">
                <span class="hud-label">进程</span>
                <span class="hud-value">${totalProc}</span>
            </div>
            <div class="hud-item">
                <span class="hud-label">已开机</span>
                <span class="hud-value">${(sys.uptimeSeconds / 3600).toFixed(1)} h</span>
            </div>
        </div>
    ` : `<div class="hud"><span style="color:var(--text-muted);font-size:12px;">尚未采集数据</span></div>`;

    root.innerHTML = `
        <div class="header">
            <div class="header-brand">
                <span class="brand-icon">${icon('cpu')}</span>
                <h1>MemorAI</h1>
                <span class="brand-sub">AI Memory Analyzer</span>
            </div>
            ${hud}
            <div class="header-actions">
                <button class="btn icon-only" id="btn-theme" title="切换主题">${icon(state.theme === 'dark' ? 'sun' : 'moon')}</button>
                <button class="btn primary" id="btn-refresh">${icon('refresh')}<span>刷新</span></button>
            </div>
        </div>
        <div class="tabs">
            <div class="tab ${cls('overview')}" data-tab="overview">${icon('layout')}<span>概览</span></div>
            <div class="tab ${cls('processes')}" data-tab="processes">${icon('list')}<span>进程</span></div>
            <div class="tab ${cls('monitor')}" data-tab="monitor">${icon('lineChart')}<span>监控</span></div>
            <div class="tab ${cls('startup')}" data-tab="startup">${icon('rocket')}<span>自启项</span></div>
            <div class="tab ${cls('ai')}" data-tab="ai">${icon('sparkles')}<span>AI 助手</span></div>
            <div class="tab ${cls('settings')}" data-tab="settings">${icon('settings')}<span>设置</span></div>
        </div>
        <div class="content" id="content"></div>
        <div class="footer">
            <span id="footer-status">${state.report ? `${state.report.total} processes · ${state.report.groups.length} apps` : 'no data yet'}</span>
            <span>MemorAI v0.4</span>
        </div>
    `;
    document.getElementById('btn-refresh').onclick = refresh;
    document.getElementById('btn-theme').onclick = () => {
        state.theme = state.theme === 'dark' ? 'light' : 'dark';
        localStorage.setItem('theme', state.theme);
        render();
    };
    document.querySelectorAll('.tab').forEach(t => {
        t.onclick = () => { state.currentTab = t.dataset.tab; render(); };
    });
    renderContent();
}

function cls(tab) { return state.currentTab === tab ? 'active' : ''; }

function renderContent() {
    const c = document.getElementById('content');
    switch (state.currentTab) {
        case 'overview': c.innerHTML = renderOverview(); bindOverview(); break;
        case 'processes': c.innerHTML = renderProcesses(); break;
        case 'monitor': c.innerHTML = renderMonitor(); bindMonitor(); break;
        case 'startup': renderStartup(c); break;
        case 'ai': c.innerHTML = renderAI(); bindAI(); break;
        case 'settings': renderSettings(c); break;
    }
}

function pctClass(p) {
    if (p >= 85) return 'danger';
    if (p >= 65) return 'warn';
    return '';
}

// 渲染 SVG 圆环图
function gaugeRing(percent, size = 110, stroke = 9, label = '已用内存') {
    const r = (size - stroke) / 2;
    const c = 2 * Math.PI * r;
    const dash = c * (Math.min(percent, 100) / 100);
    let color = 'var(--accent-cyan)';
    let textCls = '';
    if (percent >= 85) { color = 'var(--accent-red)'; textCls = 'danger'; }
    else if (percent >= 65) { color = 'var(--accent-yellow)'; textCls = 'warn'; }
    return `
        <div class="gauge-ring">
            <svg width="${size}" height="${size}">
                <circle cx="${size/2}" cy="${size/2}" r="${r}" fill="none" stroke="var(--bg-elev-2)" stroke-width="${stroke}"/>
                <circle cx="${size/2}" cy="${size/2}" r="${r}" fill="none"
                    stroke="${color}" stroke-width="${stroke}"
                    stroke-dasharray="${dash} ${c}"
                    stroke-linecap="round"
                    transform="rotate(-90 ${size/2} ${size/2})"/>
            </svg>
            <div class="gauge-text">
                <span class="num ${textCls}">${percent.toFixed(1)}<span style="font-size:13px;">%</span></span>
                <span class="lab">${label}</span>
            </div>
        </div>
    `;
}

// ===== 概览页 =====
function renderOverview() {
    if (!state.report) return `<div class="placeholder">${icon('refresh', 'ph-icon')}<p>正在采集...</p></div>`;
    const sys = state.report.system;
    const cpu = state.report.cpu;

    let html = `
        <div class="dashboard">
            <div class="card">
                <div class="card-title">${icon('cpu')}<span>内存</span></div>
                <div class="gauge">
                    ${gaugeRing(sys.usedPercent, 110, 9, '已用内存')}
                    <div class="gauge-info">
                        <div class="row"><span class="k">总内存</span><span class="v">${sys.totalGB.toFixed(2)} GB</span></div>
                        <div class="row"><span class="k">已使用</span><span class="v">${sys.usedGB.toFixed(2)} GB</span></div>
                        <div class="row"><span class="k">可用</span><span class="v">${sys.availableGB.toFixed(2)} GB</span></div>
                    </div>
                </div>
            </div>

            <div class="card">
                <div class="card-title">${icon('activity')}<span>CPU</span></div>
                <div class="gauge">
                    ${gaugeRing(cpu ? cpu.usedPercent : 0, 110, 9, 'CPU 占用')}
                    <div class="gauge-info">
                        <div class="row"><span class="k">物理核心</span><span class="v">${cpu ? cpu.cores : '-'}</span></div>
                        <div class="row"><span class="k">逻辑核心</span><span class="v">${cpu ? cpu.logicalCores : '-'}</span></div>
                        <div class="row"><span class="k">主频</span><span class="v">${cpu && cpu.mhz ? (cpu.mhz/1000).toFixed(2) + ' GHz' : '-'}</span></div>
                    </div>
                </div>
                ${cpu && cpu.modelName ? `<div style="margin-top:8px; font-size:11px; color:var(--text-muted); white-space:nowrap; overflow:hidden; text-overflow:ellipsis;" title="${escapeHtml(cpu.modelName)}">${escapeHtml(cpu.modelName)}</div>` : ''}
            </div>

            <div class="card">
                <div class="card-title">${icon('box')}<span>系统</span></div>
                <div class="stat-grid">
                    <div class="stat-cell">
                        <div class="l">操作系统</div>
                        <div class="v" style="font-size:13px;">${escapeHtml(sys.platform || '-')}</div>
                    </div>
                    <div class="stat-cell">
                        <div class="l">版本</div>
                        <div class="v" style="font-size:13px;">${escapeHtml(sys.platformVersion || '-')}</div>
                    </div>
                    <div class="stat-cell">
                        <div class="l">主机名</div>
                        <div class="v" style="font-size:13px;">${escapeHtml(sys.hostname || '-')}</div>
                    </div>
                    <div class="stat-cell">
                        <div class="l">已开机</div>
                        <div class="v">${(sys.uptimeSeconds / 3600).toFixed(1)}<span style="font-size:11px;color:var(--text-muted);"> 小时</span></div>
                    </div>
                </div>
            </div>

            <div class="card">
                <div class="card-title">${icon('list')}<span>统计</span></div>
                <div class="stat-grid">
                    <div class="stat-cell green">
                        <div class="l">总进程数</div>
                        <div class="v">${state.report.total}</div>
                    </div>
                    <div class="stat-cell">
                        <div class="l">软件分组</div>
                        <div class="v">${state.report.groups.length}</div>
                    </div>
                    <div class="stat-cell yellow">
                        <div class="l">最大占用</div>
                        <div class="v" style="font-size:12.5px;">${state.report.groups[0] ? escapeHtml(state.report.groups[0].name) : '-'}</div>
                    </div>
                    <div class="stat-cell red">
                        <div class="l">该项 MB</div>
                        <div class="v">${state.report.groups[0] ? fmtMB(state.report.groups[0].totalMB) : '-'}</div>
                    </div>
                </div>
            </div>
        </div>

        <div class="toolbar">
            <span style="font-size:11px; color:var(--text-muted); text-transform:uppercase; letter-spacing:0.6px; font-weight:600;">软件分组</span>
            <span class="meta">点击行展开查看子进程</span>
        </div>
    `;

    html += `<div class="list">
        <div class="list-header"><span></span><span>软件 / 类别</span><span style="text-align:right">内存</span><span style="text-align:right">占总内存</span><span style="text-align:center">进程数</span></div>`;
    for (const g of state.report.groups) {
        const expanded = state.expandedGroups.has(g.name);
        html += `
            <div class="list-row ${expanded ? 'expanded' : ''}" data-group="${escapeHtml(g.name)}">
                <span>${icon('chevronRight', 'chev')}</span>
                <span class="name">${escapeHtml(g.name)}<span class="sub">${escapeHtml(g.category)}</span></span>
                <span class="mem">${fmtMB(g.totalMB)}</span>
                <span class="pct">${fmtPct(g.totalPct)}</span>
                <span class="pcount">× ${g.processCount}</span>
            </div>
        `;
        if (expanded) {
            html += `<div class="sub-procs">`;
            for (const p of g.processes) {
                html += `<div class="sp-row"><span>${escapeHtml(p.name)} (PID ${p.pid})</span><span>${fmtMB(p.memoryMB)}</span></div>`;
            }
            html += `</div>`;
        }
    }
    html += `</div>`;
    return html;
}
function bindOverview() {
    document.querySelectorAll('.list-row[data-group]').forEach(row => {
        row.onclick = () => {
            const name = row.dataset.group;
            if (state.expandedGroups.has(name)) state.expandedGroups.delete(name);
            else state.expandedGroups.add(name);
            renderContent();
        };
    });
}

// ===== 进程详情 =====
function renderProcesses() {
    if (!state.report) return `<div class="placeholder">${icon('refresh', 'ph-icon')}<p>正在采集...</p></div>`;
    let html = `
        <div class="toolbar">
            <span style="font-size:11px; color:var(--text-muted); text-transform:uppercase; letter-spacing:0.6px; font-weight:600;">进程列表</span>
            <span class="meta">按内存占用降序排列 · 共 ${state.report.topProcs.length} 项</span>
        </div>
        <div class="list">
            <div class="list-header"><span></span><span>进程名 / 路径</span><span style="text-align:right">内存</span><span style="text-align:right">占比</span><span style="text-align:center">PID</span></div>`;
    for (const p of state.report.topProcs) {
        html += `<div class="list-row">
            <span></span>
            <span class="name">${escapeHtml(p.name)}${p.exePath ? `<span class="sub">${escapeHtml(p.exePath)}</span>` : ''}</span>
            <span class="mem">${fmtMB(p.memoryMB)}</span>
            <span class="pct">${fmtPct(p.memoryPct)}</span>
            <span class="pcount">${p.pid}</span>
        </div>`;
    }
    html += `</div>`;
    return html;
}

// ===== 监控采样 =====
let lineChart = null, pieChart = null;

function renderMonitor() {
    const m = state.monitor;
    let body;
    if (m.running) {
        body = `
            <div class="monitor-status">
                <h3>采样中... ${m.progressText}</h3>
                <div class="progress" style="height:14px;"><div class="progress-bar" style="width:${m.progress}%"></div></div>
                <p style="color:#8b949e; font-size:12px; margin-top:8px;">已采集 ${m.liveSamples.length} 个样本，请稍候</p>
            </div>
        `;
    } else if (m.result) {
        body = `
            <div class="monitor-status">
                <p style="color:var(--accent-cyan); margin:0; font-size:12.5px;">✓ 采样完成 — ${m.result.startedAt} · ${m.result.durationSec}秒 · ${m.result.samples.length} 个样本</p>
            </div>
            <div class="charts">
                <div class="chart-card"><h3>${icon('lineChart')}<span>内存占用变化</span></h3><canvas id="line-chart"></canvas></div>
                <div class="chart-card"><h3>${icon('pieChart')}<span>软件内存占比 (Top 10)</span></h3><canvas id="pie-chart"></canvas></div>
            </div>
        `;
    } else {
        body = `<div class="placeholder">${icon('lineChart', 'ph-icon')}<p>点击下方按钮开始采样</p></div>`;
    }

    return `
        <div class="toolbar">
            <button class="btn primary" id="btn-monitor-start" ${m.running ? 'disabled' : ''}>
                ${m.running ? `<span class="spinner"></span><span>采样中...</span>` : `${icon('play')}<span>开始采样</span>`}
            </button>
            <select id="monitor-duration" ${m.running ? 'disabled' : ''}>
                <option value="15" ${m.duration===15?'selected':''}>15 秒</option>
                <option value="30" ${m.duration===30?'selected':''}>30 秒</option>
                <option value="60" ${m.duration===60?'selected':''}>60 秒</option>
                <option value="120" ${m.duration===120?'selected':''}>2 分钟</option>
            </select>
            ${m.result && !m.running ? `
                <span class="gap"></span>
                <button class="btn" id="btn-monitor-export-html">${icon('fileText')}<span>导出 HTML</span></button>
                <button class="btn" id="btn-monitor-export-json">${icon('download')}<span>导出 JSON</span></button>
            ` : ''}
        </div>
        ${body}
    `;
}

function bindMonitor() {
    const startBtn = document.getElementById('btn-monitor-start');
    if (startBtn) startBtn.onclick = startMonitoring;
    const sel = document.getElementById('monitor-duration');
    if (sel) sel.onchange = e => { state.monitor.duration = parseInt(e.target.value, 10); };
    const expHtml = document.getElementById('btn-monitor-export-html');
    if (expHtml) expHtml.onclick = exportMonitorHTML;
    const expJson = document.getElementById('btn-monitor-export-json');
    if (expJson) expJson.onclick = exportMonitorJSON;

    // 渲染图表
    if (state.monitor.result && !state.monitor.running) {
        drawCharts(state.monitor.result);
    }
}

const COLOR_PALETTE = [
    '#58a6ff', '#f85149', '#3fb950', '#f0ad4e', '#a371f7',
    '#ff7b72', '#79c0ff', '#ffa657', '#56d364', '#d2a8ff',
];

function drawCharts(result) {
    if (lineChart) { lineChart.destroy(); lineChart = null; }
    if (pieChart) { pieChart.destroy(); pieChart = null; }

    // 折线图：x = 时间，多条线 = 总占用率 + CPU + Top 软件
    const lineCanvas = document.getElementById('line-chart');
    if (lineCanvas) {
        const labels = result.samples.map(s => s.wallTime);
        const datasets = [{
            label: '内存占用率(%)',
            data: result.samples.map(s => s.usedPercent),
            borderColor: '#7aa2f7',
            backgroundColor: 'rgba(122,162,247,0.1)',
            yAxisID: 'y',
            tension: 0.3,
            borderWidth: 2,
        }, {
            label: 'CPU 占用率(%)',
            data: result.samples.map(s => s.cpuPercent || 0),
            borderColor: '#f7768e',
            backgroundColor: 'rgba(247,118,142,0.08)',
            yAxisID: 'y',
            tension: 0.3,
            borderWidth: 2,
        }];
        // 每个 tracked 软件画一条线（MB）
        const tracked = result.trackedSoftware || [];
        tracked.slice(0, 6).forEach((sw, idx) => {
            const data = result.samples.map(s => {
                const found = (s.topGroups || []).find(g => g.name === sw);
                return found ? found.mb : null;
            });
            datasets.push({
                label: sw,
                data,
                borderColor: COLOR_PALETTE[(idx + 2) % COLOR_PALETTE.length],
                backgroundColor: 'transparent',
                yAxisID: 'y2',
                tension: 0.3,
                borderWidth: 1.5,
                spanGaps: true,
            });
        });

        lineChart = new Chart(lineCanvas, {
            type: 'line',
            data: { labels, datasets },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                interaction: { intersect: false, mode: 'index' },
                plugins: {
                    legend: { labels: { color: '#c9d1d9', font: { size: 11 } } },
                    tooltip: { backgroundColor: '#161b22', titleColor: '#c9d1d9', bodyColor: '#c9d1d9' },
                },
                scales: {
                    x: { ticks: { color: '#8b949e', maxTicksLimit: 10 }, grid: { color: '#21262d' } },
                    y: {
                        type: 'linear', position: 'left',
                        title: { display: true, text: '占用率 (%) — 内存/CPU', color: '#7aa2f7' },
                        ticks: { color: '#8b949e' }, grid: { color: '#21262d' }, min: 0, max: 100,
                    },
                    y2: {
                        type: 'linear', position: 'right',
                        title: { display: true, text: '软件内存 (MB)', color: '#f0ad4e' },
                        ticks: { color: '#8b949e' }, grid: { drawOnChartArea: false },
                    },
                },
            },
        });
    }

    // 饼图：Top 10 + 其他
    const pieCanvas = document.getElementById('pie-chart');
    if (pieCanvas) {
        const groups = result.finalGroups || [];
        const top = groups.slice(0, 10);
        const otherSum = groups.slice(10).reduce((sum, g) => sum + g.totalMB, 0);
        const labels = top.map(g => g.name);
        const data = top.map(g => g.totalMB);
        if (otherSum > 0) { labels.push('其他'); data.push(otherSum); }

        pieChart = new Chart(pieCanvas, {
            type: 'doughnut',
            data: {
                labels,
                datasets: [{
                    data,
                    backgroundColor: labels.map((_, i) => COLOR_PALETTE[i % COLOR_PALETTE.length]),
                    borderColor: '#0f1419',
                    borderWidth: 2,
                }],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { position: 'right', labels: { color: '#c9d1d9', font: { size: 11 }, padding: 8 } },
                    tooltip: {
                        backgroundColor: '#161b22', titleColor: '#c9d1d9', bodyColor: '#c9d1d9',
                        callbacks: {
                            label: ctx => {
                                const mb = ctx.parsed;
                                const total = ctx.dataset.data.reduce((a,b)=>a+b,0);
                                const pct = total > 0 ? (mb/total*100).toFixed(1) : '0';
                                return `${ctx.label}: ${fmtMB(mb)} (${pct}%)`;
                            }
                        }
                    },
                },
            },
        });
    }
}

async function startMonitoring() {
    state.monitor.running = true;
    state.monitor.progress = 0;
    state.monitor.progressText = '准备中...';
    state.monitor.result = null;
    state.monitor.liveSamples = [];
    renderContent();

    try {
        const result = await callBackend('StartMonitoring', state.monitor.duration, state.monitor.interval);
        state.monitor.result = result;
        toast('采样完成', 'success');
    } catch (e) {
        toast('采样失败: ' + e, 'error');
    } finally {
        state.monitor.running = false;
        renderContent();
    }
}

function exportMonitorJSON() {
    const data = JSON.stringify(state.monitor.result, null, 2);
    downloadFile('memory-monitor.json', data, 'application/json');
    toast('已导出 JSON', 'success');
}

function exportMonitorHTML() {
    const r = state.monitor.result;
    if (!r) return;
    const samples = r.samples;
    const tracked = r.trackedSoftware || [];

    const lineData = {
        labels: samples.map(s => s.wallTime),
        datasets: [{
            label: '内存占用率(%)',
            data: samples.map(s => s.usedPercent),
            borderColor: '#7aa2f7',
            backgroundColor: 'rgba(122,162,247,0.1)',
            yAxisID: 'y', tension: 0.3, borderWidth: 2,
        }, {
            label: 'CPU 占用率(%)',
            data: samples.map(s => s.cpuPercent || 0),
            borderColor: '#f7768e',
            backgroundColor: 'rgba(247,118,142,0.08)',
            yAxisID: 'y', tension: 0.3, borderWidth: 2,
        }, ...tracked.slice(0, 6).map((sw, idx) => ({
            label: sw,
            data: samples.map(s => {
                const f = (s.topGroups || []).find(g => g.name === sw);
                return f ? f.mb : null;
            }),
            borderColor: COLOR_PALETTE[(idx + 2) % COLOR_PALETTE.length],
            backgroundColor: 'transparent',
            yAxisID: 'y2', tension: 0.3, borderWidth: 1.5, spanGaps: true,
        }))],
    };

    const top = (r.finalGroups || []).slice(0, 10);
    const otherSum = (r.finalGroups || []).slice(10).reduce((s, g) => s + g.totalMB, 0);
    const pieLabels = top.map(g => g.name); const pieData = top.map(g => g.totalMB);
    if (otherSum > 0) { pieLabels.push('其他'); pieData.push(otherSum); }

    // 表格行
    const groupRows = (r.finalGroups || []).map(g => `
        <tr><td>${escapeHtml(g.name)}</td><td>${escapeHtml(g.category)}</td>
        <td style="text-align:right">${fmtMB(g.totalMB)}</td>
        <td style="text-align:right">${fmtPct(g.totalPct)}</td>
        <td style="text-align:center">${g.processCount}</td></tr>`).join('');

    const html = `<!DOCTYPE html>
<html lang="zh-CN"><head><meta charset="UTF-8"><title>内存监控报告 - ${r.startedAt}</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js"><\/script>
<style>
body{font-family:-apple-system,Segoe UI,sans-serif;background:#0f1419;color:#e6e6e6;margin:0;padding:24px;}
h1{color:#58a6ff;margin-top:0;}h2{color:#79c0ff;border-bottom:1px solid #21262d;padding-bottom:8px;}
.meta{background:#161b22;padding:14px 18px;border-radius:8px;margin-bottom:16px;}
.charts{display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-bottom:24px;}
.chart-box{background:#161b22;padding:16px;border-radius:8px;height:380px;}
table{width:100%;border-collapse:collapse;background:#161b22;border-radius:8px;overflow:hidden;}
th,td{padding:8px 12px;border-bottom:1px solid #21262d;text-align:left;}
th{background:#1c222a;color:#8b949e;font-size:12px;}
@media(max-width:900px){.charts{grid-template-columns:1fr;}}
</style></head><body>
<h1>💻 内存监控报告</h1>
<div class="meta">
<p><b>开始时间:</b> ${r.startedAt} &nbsp;|&nbsp; <b>持续:</b> ${r.durationSec}s &nbsp;|&nbsp; <b>采样数:</b> ${samples.length}</p>
<p><b>系统:</b> ${escapeHtml(r.system?.platform||'')} ${escapeHtml(r.system?.platformVersion||'')} &nbsp;|&nbsp;
<b>总内存:</b> ${(r.system?.totalGB||0).toFixed(2)} GB &nbsp;|&nbsp;
<b>采样末已用:</b> ${(r.system?.usedGB||0).toFixed(2)} GB (${(r.system?.usedPercent||0).toFixed(1)}%)</p>
</div>
<h2>📈 内存占用变化</h2>
<div class="chart-box"><canvas id="lc"></canvas></div>
<h2>🥧 最终软件占比（Top 10）</h2>
<div class="chart-box"><canvas id="pc"></canvas></div>
<h2>📋 最终软件分组明细</h2>
<table><thead><tr><th>软件</th><th>类别</th><th style="text-align:right">内存</th><th style="text-align:right">占比</th><th style="text-align:center">进程数</th></tr></thead><tbody>${groupRows}</tbody></table>
<script>
const palette=${JSON.stringify(COLOR_PALETTE)};
new Chart(document.getElementById('lc'),{type:'line',data:${JSON.stringify(lineData)},options:{responsive:true,maintainAspectRatio:false,plugins:{legend:{labels:{color:'#c9d1d9'}}},scales:{x:{ticks:{color:'#8b949e'},grid:{color:'#21262d'}},y:{type:'linear',position:'left',title:{display:true,text:'系统占用率 (%)',color:'#58a6ff'},ticks:{color:'#8b949e'},grid:{color:'#21262d'},min:0,max:100},y2:{type:'linear',position:'right',title:{display:true,text:'软件内存 (MB)',color:'#f0ad4e'},ticks:{color:'#8b949e'},grid:{drawOnChartArea:false}}}}});
new Chart(document.getElementById('pc'),{type:'doughnut',data:{labels:${JSON.stringify(pieLabels)},datasets:[{data:${JSON.stringify(pieData)},backgroundColor:${JSON.stringify(pieLabels)}.map((_,i)=>palette[i%palette.length]),borderColor:'#0f1419',borderWidth:2}]},options:{responsive:true,maintainAspectRatio:false,plugins:{legend:{position:'right',labels:{color:'#c9d1d9'}}}}});
<\/script></body></html>`;

    downloadFile(`memory-report-${r.startedAt.replace(/[: ]/g,'-')}.html`, html, 'text/html');
    toast('已导出 HTML 报告', 'success');
}

function downloadFile(filename, content, mime) {
    const blob = new Blob([content], { type: mime });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url; a.download = filename;
    document.body.appendChild(a); a.click();
    setTimeout(() => { URL.revokeObjectURL(url); a.remove(); }, 200);
}

// ===== 自启项 =====
function renderStartup(container) {
    container.innerHTML = `<div class="ai-placeholder">加载中...</div>`;
    Promise.all([
        callBackend('ListStartupItems'),
        callBackend('IsAdmin'),
    ]).then(([items, isAdmin]) => {
        state.startup.items = items || [];
        state.startup.isAdmin = isAdmin;
        renderStartupList(container);
    }).catch(e => {
        container.innerHTML = `<div class="ai-placeholder">加载失败: ${escapeHtml(String(e))}</div>`;
    });
}

function renderStartupList(container) {
    const items = state.startup.items;
    const isAdmin = state.startup.isAdmin;

    let banner = '';
    if (!isAdmin) {
        banner = `<div class="banner-warn">${icon('alertTriangle')}<span>当前不是管理员模式。可以管理"用户级"启动项，但"系统级"启动项需要以管理员身份重新启动本程序。</span></div>`;
    } else {
        banner = `<div class="banner-info">${icon('checkCircle')}<span>已以管理员身份运行，所有自启项均可修改。</span></div>`;
    }

    let html = `
        <div class="toolbar">
            <button class="btn primary" id="btn-startup-refresh">${icon('refresh')}<span>重新扫描</span></button>
            <span class="meta">共 ${items.length} 项 · 启用 ${items.filter(i=>i.enabled).length} · 已禁用 ${items.filter(i=>!i.enabled).length}</span>
        </div>
        ${banner}
        <div class="list">
            <div class="list-header startup-header">
                <span>名称</span><span>命令 / 路径</span><span style="text-align:center">来源</span><span style="text-align:center">范围</span><span style="text-align:center">操作</span>
            </div>
    `;

    if (items.length === 0) {
        html += `<div class="placeholder" style="padding:30px;">未发现任何自启项</div>`;
    }

    for (const it of items) {
        const cantModify = it.needAdmin && !isAdmin;
        const sourceText = it.source === 'registry' ? '注册表' : '启动文件夹';
        const scopeText = it.scope === 'user' ? '用户' : '系统';
        const scopeColor = it.scope === 'system' ? 'var(--accent-yellow)' : 'var(--accent-cyan)';
        html += `
            <div class="list-row startup-row" style="${it.enabled ? '' : 'opacity:0.55;'}">
                <span class="name">${escapeHtml(it.name)}<span class="sub">${escapeHtml(it.location)}</span></span>
                <span class="name" style="font-size:11.5px;color:var(--text-dim);font-family:'JetBrains Mono',monospace;">${escapeHtml(it.command)}</span>
                <span style="text-align:center"><span class="cat-tag">${sourceText}</span></span>
                <span style="text-align:center"><span class="cat-tag" style="color:${scopeColor}">${scopeText}</span></span>
                <span>
                    <button class="btn ${it.enabled ? 'danger' : 'primary'}" data-id="${it.id}" data-enable="${!it.enabled}" ${cantModify?'disabled title="需要管理员权限"':''}>
                        ${it.enabled ? '禁用' : '启用'}
                    </button>
                </span>
            </div>
        `;
    }
    html += `</div>`;
    container.innerHTML = html;

    document.getElementById('btn-startup-refresh').onclick = () => renderStartup(container);
    container.querySelectorAll('button[data-id]').forEach(b => {
        b.onclick = async () => {
            const id = b.dataset.id;
            const enable = b.dataset.enable === 'true';
            b.disabled = true;
            b.innerHTML = '<span class="spinner"></span>处理中';
            try {
                await callBackend('ToggleStartupItem', id, enable);
                toast(`已${enable?'启用':'禁用'}`, 'success');
                renderStartup(container);
            } catch (e) {
                toast('操作失败: ' + e, 'error');
                renderStartup(container);
            }
        };
    });
}

// ===== AI 对话 =====
function renderAI() {
    const c = state.chat;

    // 顶部工具栏
    let toolbar = `
        <div class="toolbar">
            <button class="btn primary" id="btn-ai-quick" ${c.sending ? 'disabled' : ''}>${icon('sparkles')}<span>一键分析当前内存</span></button>
            <button class="btn" id="btn-ai-snapshot" ${c.sending || c.snapshotLoaded ? 'disabled' : ''} title="把当前内存快照作为上下文加入对话">
                ${icon('clipboard')}<span>${c.snapshotLoaded ? '✓ 已载入快照' : '载入内存快照'}</span>
            </button>
            <button class="btn danger" id="btn-ai-clear" ${c.sending || c.messages.length === 0 ? 'disabled' : ''}>${icon('trash')}<span>清空对话</span></button>
        </div>
    `;

    // 消息列表
    let msgsHtml = '';
    if (c.messages.length === 0) {
        msgsHtml = `<div class="chat-placeholder">
            ${icon('bot', 'ph-icon')}
            <p style="font-size:14px; color:var(--text-dim); margin: 8px 0 4px;">开始与 AI 对话</p>
            <small style="color:var(--text-muted); line-height:1.8;">
                · 直接输入问题聊天<br>
                · 或点击"一键分析"让 AI 解读内存占用<br>
                · 或先"载入快照"再追问具体问题
            </small>
        </div>`;
    } else {
        msgsHtml = c.messages.map((m, idx) => {
            // 系统注入的快照消息特殊处理（折叠显示）
            const isSnapshot = m.role === 'user' && m._snapshot;
            if (isSnapshot) {
                return `<div class="msg msg-snapshot">
                    <div class="msg-meta">${icon('clipboard')} 内存快照 (已作为上下文发送给 AI)</div>
                    <details><summary>查看快照内容</summary><pre>${escapeHtml(m.content)}</pre></details>
                </div>`;
            }
            const cls = m.role === 'user' ? 'msg-user' : 'msg-assistant';
            const label = m.role === 'user' ? '你' : 'AI';
            const body = m.role === 'assistant' ? renderMarkdown(m.content) : `<p>${escapeHtml(m.content).replace(/\n/g, '<br>')}</p>`;
            return `<div class="msg ${cls}">
                <div class="msg-meta">${label}</div>
                <div class="msg-body">${body}</div>
            </div>`;
        }).join('');
    }

    if (c.sending) {
        msgsHtml += `<div class="msg msg-assistant">
            <div class="msg-meta">AI</div>
            <div class="msg-body"><span class="spinner"></span> 思考中...</div>
        </div>`;
    }

    // 输入框
    const inputDisabled = c.sending ? 'disabled' : '';
    const input = `
        <div class="chat-input">
            <textarea id="chat-input" placeholder="向 AI 提问，按 Enter 发送，Shift+Enter 换行..." ${inputDisabled} rows="2">${escapeHtml(c.input)}</textarea>
            <button class="btn primary" id="btn-chat-send" ${inputDisabled}>${icon('send')}<span>发送</span></button>
        </div>
    `;

    return `
        ${toolbar}
        <div class="chat-area">
            <div class="chat-messages" id="chat-messages">${msgsHtml}</div>
            ${input}
        </div>
    `;
}

function bindAI() {
    const quick = document.getElementById('btn-ai-quick');
    if (quick) quick.onclick = aiQuickAnalyze;
    const snap = document.getElementById('btn-ai-snapshot');
    if (snap) snap.onclick = aiLoadSnapshot;
    const clr = document.getElementById('btn-ai-clear');
    if (clr) clr.onclick = aiClearChat;
    const send = document.getElementById('btn-chat-send');
    if (send) send.onclick = aiSendMessage;

    const ta = document.getElementById('chat-input');
    if (ta) {
        ta.addEventListener('input', e => { state.chat.input = e.target.value; });
        ta.addEventListener('keydown', e => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                aiSendMessage();
            }
        });
        // 聚焦但不打断
        if (!state.chat.sending) ta.focus();
    }

    // 滚到底部
    const ml = document.getElementById('chat-messages');
    if (ml) ml.scrollTop = ml.scrollHeight;
}

async function aiQuickAnalyze() {
    if (state.chat.sending) return;
    // 重置对话，载入快照，发"请分析"
    state.chat.messages = [];
    state.chat.snapshotLoaded = false;
    state.chat.sending = true;
    renderContent();

    try {
        const snapshot = await callBackend('GetMemorySnapshotPrompt');
        state.chat.messages.push({ role: 'user', content: snapshot, _snapshot: true });
        state.chat.snapshotLoaded = true;
        // 添加一条明确的"请分析"消息
        state.chat.messages.push({ role: 'user', content: '请基于以上数据进行完整分析。' });
        renderContent();

        // 发送给 AI（去掉 _snapshot 标记，role/content 保留）
        const apiMsgs = state.chat.messages.map(m => ({ role: m.role, content: m.content }));
        const reply = await callBackend('AIChat', apiMsgs);
        state.chat.messages.push({ role: 'assistant', content: reply });
    } catch (e) {
        state.chat.messages.push({ role: 'assistant', content: `## ❌ 分析失败\n\n${String(e)}\n\n请检查「设置」中的 API 配置。` });
    } finally {
        state.chat.sending = false;
        state.chat.input = '';
        renderContent();
    }
}

async function aiLoadSnapshot() {
    if (state.chat.sending || state.chat.snapshotLoaded) return;
    try {
        const snapshot = await callBackend('GetMemorySnapshotPrompt');
        state.chat.messages.push({ role: 'user', content: snapshot, _snapshot: true });
        state.chat.snapshotLoaded = true;
        toast('快照已载入到对话上下文', 'success');
        renderContent();
    } catch (e) {
        toast('载入快照失败: ' + e, 'error');
    }
}

function aiClearChat() {
    if (state.chat.sending) return;
    if (state.chat.messages.length > 0 && !confirm('确定清空当前对话吗？')) return;
    state.chat.messages = [];
    state.chat.snapshotLoaded = false;
    state.chat.input = '';
    renderContent();
}

async function aiSendMessage() {
    if (state.chat.sending) return;
    const text = state.chat.input.trim();
    if (!text) return;

    state.chat.messages.push({ role: 'user', content: text });
    state.chat.input = '';
    state.chat.sending = true;
    renderContent();

    try {
        const apiMsgs = state.chat.messages.map(m => ({ role: m.role, content: m.content }));
        const reply = await callBackend('AIChat', apiMsgs);
        state.chat.messages.push({ role: 'assistant', content: reply });
    } catch (e) {
        state.chat.messages.push({ role: 'assistant', content: `❌ 错误: ${String(e)}` });
    } finally {
        state.chat.sending = false;
        renderContent();
    }
}

// ===== 设置 =====
function renderSettings(container) {
    callBackend('GetConfig').then(cfg => {
        cfg = cfg || {};
        container.innerHTML = `
            <div class="settings">
                <h3 class="settings-section">AI 接口</h3>
                <div class="field">
                    <label>API Base URL</label>
                    <input type="text" id="cfg-url" value="${escapeHtml(cfg.aiBaseURL)}" placeholder="https://api.deepseek.com/v1"/>
                    <div class="hint">兼容 OpenAI 协议的接口地址。常用：<br>
                        · DeepSeek: <code>https://api.deepseek.com/v1</code><br>
                        · OpenAI: <code>https://api.openai.com/v1</code><br>
                        · Kimi: <code>https://api.moonshot.cn/v1</code><br>
                        · 通义千问: <code>https://dashscope.aliyuncs.com/compatible-mode/v1</code>
                    </div>
                </div>
                <div class="field"><label>API Key</label>
                    <input type="password" id="cfg-key" value="${escapeHtml(cfg.aiApiKey)}" placeholder="sk-..."/>
                    <div class="hint">保存到本地 config.json，不会上传到任何地方</div>
                </div>
                <div class="field"><label>模型名称</label>
                    <input type="text" id="cfg-model" value="${escapeHtml(cfg.aiModel)}" placeholder="deepseek-chat"/>
                    <div class="hint">如 <code>deepseek-chat</code>、<code>gpt-4o-mini</code>、<code>moonshot-v1-8k</code></div>
                </div>

                <h3 class="settings-section">提示词注入</h3>
                <div class="field">
                    <label>自定义系统提示词（可选）</label>
                    <textarea id="cfg-extra" rows="6" placeholder="在这里写你想追加给 AI 的指令，例如：&#10;- 用更简洁的语气&#10;- 给出英文回答&#10;- 不要建议禁用 XX 软件，这是我必须用的&#10;- 假装你是一个毒舌系统优化大师">${escapeHtml(cfg.aiExtraPrompt || '')}</textarea>
                    <div class="hint">这段内容会追加到内置系统提示词之后，影响 AI 的所有回复（含一键分析和对话）。
                        <a href="#" id="link-show-default">查看默认提示词</a>
                    </div>
                </div>

                <h3 class="settings-section">采集</h3>
                <div class="field"><label>进程列表 Top N</label>
                    <input type="number" id="cfg-topn" value="${cfg.topN || 30}" min="10" max="200"/>
                    <div class="hint">"进程详情"页只显示前 N 个最占内存的进程</div>
                </div>
                <div class="btn-row">
                    <button class="btn primary" id="btn-cfg-save">保存</button>
                    <button class="btn" id="btn-cfg-test">测试连接</button>
                </div>
            </div>

            <div id="default-prompt-modal" class="modal hidden">
                <div class="modal-content">
                    <div class="modal-header">
                        <h3>默认系统提示词（只读）</h3>
                        <button class="btn" id="btn-modal-close">关闭</button>
                    </div>
                    <pre id="default-prompt-text">加载中...</pre>
                </div>
            </div>
        `;
        const collect = () => ({
            aiBaseURL: document.getElementById('cfg-url').value.trim(),
            aiApiKey: document.getElementById('cfg-key').value.trim(),
            aiModel: document.getElementById('cfg-model').value.trim(),
            aiExtraPrompt: document.getElementById('cfg-extra').value,
            topN: parseInt(document.getElementById('cfg-topn').value, 10) || 30,
        });
        document.getElementById('btn-cfg-save').onclick = async () => {
            try { await callBackend('SaveConfig', collect()); toast('配置已保存', 'success'); }
            catch (e) { toast('保存失败: ' + e, 'error'); }
        };
        document.getElementById('btn-cfg-test').onclick = async () => {
            try {
                await callBackend('SaveConfig', collect());
                toast('正在测试连接...');
                const r = await callBackend('TestAIConnection');
                toast('连接成功: ' + r, 'success');
            } catch (e) { toast('连接失败: ' + e, 'error'); }
        };
        document.getElementById('link-show-default').onclick = async (e) => {
            e.preventDefault();
            const modal = document.getElementById('default-prompt-modal');
            modal.classList.remove('hidden');
            try {
                const txt = await callBackend('GetDefaultSystemPrompt');
                document.getElementById('default-prompt-text').textContent = txt;
            } catch (err) {
                document.getElementById('default-prompt-text').textContent = '加载失败: ' + err;
            }
        };
        document.getElementById('btn-modal-close').onclick = () => {
            document.getElementById('default-prompt-modal').classList.add('hidden');
        };
    }).catch(e => {
        container.innerHTML = `<div class="ai-placeholder">加载配置失败: ${escapeHtml(String(e))}</div>`;
    });
}

// ===== 操作 =====
async function refresh() {
    const btn = document.getElementById('btn-refresh');
    if (btn) { btn.disabled = true; btn.innerHTML = `<span class="spinner"></span><span>采集中</span>`; }
    try {
        state.report = await callBackend('AnalyzeMemory');
        toast('采集完成', 'success');
    } catch (e) {
        toast('采集失败: ' + e, 'error');
    } finally {
        render();
    }
}

async function runAIAnalysis() {
    // 兼容老代码：转到对话页快速分析
    state.currentTab = 'ai';
    render();
    setTimeout(aiQuickAnalyze, 100);
}

// ===== 监控进度事件 =====
listenEvent('monitor:progress', payload => {
    if (!state.monitor.running) state.monitor.running = true;
    state.monitor.progress = payload.percent || 0;
    state.monitor.progressText = `${payload.index}/${payload.total} (${(payload.percent||0).toFixed(0)}%)`;
    if (payload.sample) state.monitor.liveSamples.push(payload.sample);
    if (state.currentTab === 'monitor') renderContent();
});

// ===== 启动 =====
render();
setTimeout(refresh, 200);
