import './style.css';
import './app.css';

import { EventsOn } from '../wailsjs/runtime/runtime';
import { Bootstrap, ChooseDataDirectory, DataPath, SendChatMessage, UpdateDisplayName, UpdateLanguage } from '../wailsjs/go/main/App';

const dictionaries = {
  'zh-CN': {
    brandEyebrow: '局域网桌面通讯',
    nicknameLabel: '你的昵称',
    save: '保存',
    send: '发送',
    settings: '设置',
    language: '语言',
    change: '修改',
    onlineNow: '当前在线',
    storage: '存储位置',
    storageHint: '聊天记录和设置会保存在这个目录中。',
    storageUpdated: '存储位置已更新',
    loadingPath: '正在加载本地路径...',
    choosePeer: '选择一个联系人',
    noConversation: '未选择会话',
    waitingTitle: '正在等待局域网内的其他用户',
    waitingBody: '在同一网络中的另一台电脑上打开 LanTalk，它会自动出现在这里。',
    messagePlaceholder: '输入消息...',
    composerIdle: '选择在线联系人后即可开始聊天。',
    composerReady: '按 Enter 发送，Shift+Enter 换行。',
    noPeers: '暂时还没有发现其他用户。',
    emptyConversationTitle: '会话还是空的',
    emptyConversationBody: '和 {name} 打个招呼吧。',
    bootFailed: 'LanTalk 启动失败',
    idAndPort: 'ID: {id}  TCP: {port}',
    sent: '已发送',
    received: '已接收',
    failed: '发送失败',
    langZh: '简体中文',
    langEn: 'English',
  },
  'en-US': {
    brandEyebrow: 'Local Network Messenger',
    nicknameLabel: 'Your nickname',
    save: 'Save',
    send: 'Send',
    settings: 'Settings',
    language: 'Language',
    change: 'Change',
    onlineNow: 'Online now',
    storage: 'Storage',
    storageHint: 'Chat history and settings are stored in this directory.',
    storageUpdated: 'Storage location updated',
    loadingPath: 'Loading local path...',
    choosePeer: 'Choose a peer',
    noConversation: 'No conversation selected',
    waitingTitle: 'Waiting for teammates on the LAN',
    waitingBody: 'Open LanTalk on another computer in the same network and it will appear here automatically.',
    messagePlaceholder: 'Type a message...',
    composerIdle: 'Pick an online peer to start chatting.',
    composerReady: 'Enter sends. Shift+Enter inserts a new line.',
    noPeers: 'No peers discovered yet.',
    emptyConversationTitle: 'Conversation is empty',
    emptyConversationBody: 'Say hello to {name}.',
    bootFailed: 'LanTalk failed to boot',
    idAndPort: 'ID: {id}  TCP: {port}',
    sent: 'sent',
    received: 'received',
    failed: 'failed',
    langZh: 'Simplified Chinese',
    langEn: 'English',
  },
};

document.querySelector('#app').innerHTML = `
  <div class="shell">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-mark">L</div>
        <div>
          <div class="eyebrow" id="brandEyebrow"></div>
          <h1>LanTalk</h1>
        </div>
      </div>

      <section class="panel profile-panel">
        <label class="label" id="nicknameLabel" for="displayName"></label>
        <div class="name-row">
          <input id="displayName" class="text-input" maxlength="40" autocomplete="off" />
          <button id="saveName" class="primary-btn"></button>
        </div>
        <p class="hint" id="profileHint"></p>
      </section>

      <section class="panel settings-panel">
        <div class="section-head">
          <span id="settingsTitle"></span>
        </div>
        <label class="label" id="languageLabel"></label>
        <div class="language-menu">
          <button id="languageTrigger" class="select-button" type="button" aria-haspopup="listbox" aria-expanded="false">
            <span id="languageValue"></span>
          </button>
          <div id="languageMenu" class="menu-panel hidden" role="listbox" tabindex="-1">
            <button class="menu-option" type="button" data-language="zh-CN"></button>
            <button class="menu-option" type="button" data-language="en-US"></button>
          </div>
        </div>
      </section>

      <section class="panel peers-panel">
        <div class="section-head">
          <span id="onlineTitle"></span>
          <span class="badge" id="peerCount">0</span>
        </div>
        <div id="peerList" class="peer-list"></div>
      </section>

      <section class="panel data-panel">
        <div class="section-head">
          <span id="storageTitle"></span>
          <button id="changeDataPath" class="ghost-btn" type="button"></button>
        </div>
        <p class="hint data-path" id="dataPath"></p>
        <p class="hint" id="storageHint"></p>
      </section>
    </aside>

    <main class="chat-area">
      <div class="chat-header">
        <div>
          <div class="eyebrow" id="activeStatus"></div>
          <h2 id="activePeer"></h2>
        </div>
      </div>

      <section id="messageList" class="message-list empty"></section>

      <form id="composer" class="composer">
        <textarea id="messageInput" maxlength="2000" disabled></textarea>
        <div class="composer-actions">
          <span class="hint" id="composerHint"></span>
          <button type="submit" id="sendButton" class="primary-btn" disabled></button>
        </div>
      </form>
    </main>
  </div>
`;

const state = {
  self: null,
  settings: {
    language: 'zh-CN',
  },
  peers: [],
  conversations: {},
  activePeerId: null,
};

const elements = {
  displayName: document.getElementById('displayName'),
  saveName: document.getElementById('saveName'),
  profileHint: document.getElementById('profileHint'),
  languageTrigger: document.getElementById('languageTrigger'),
  languageMenu: document.getElementById('languageMenu'),
  languageValue: document.getElementById('languageValue'),
  languageOptions: Array.from(document.querySelectorAll('[data-language]')),
  peerCount: document.getElementById('peerCount'),
  peerList: document.getElementById('peerList'),
  activePeer: document.getElementById('activePeer'),
  activeStatus: document.getElementById('activeStatus'),
  messageList: document.getElementById('messageList'),
  messageInput: document.getElementById('messageInput'),
  sendButton: document.getElementById('sendButton'),
  composerHint: document.getElementById('composerHint'),
  composer: document.getElementById('composer'),
  dataPath: document.getElementById('dataPath'),
  brandEyebrow: document.getElementById('brandEyebrow'),
  nicknameLabel: document.getElementById('nicknameLabel'),
  settingsTitle: document.getElementById('settingsTitle'),
  languageLabel: document.getElementById('languageLabel'),
  onlineTitle: document.getElementById('onlineTitle'),
  storageTitle: document.getElementById('storageTitle'),
  changeDataPath: document.getElementById('changeDataPath'),
  storageHint: document.getElementById('storageHint'),
};

function languageLabel(language) {
  return language === 'en-US' ? t('langEn') : t('langZh');
}

function t(key, vars = {}) {
  const language = state.settings.language || 'zh-CN';
  const dictionary = dictionaries[language] || dictionaries['zh-CN'];
  let value = dictionary[key] || key;

  Object.entries(vars).forEach(([name, replacement]) => {
    value = value.replaceAll(`{${name}}`, replacement);
  });

  return value;
}

function escapeHtml(value) {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');
}

function formatTime(timestamp) {
  return new Intl.DateTimeFormat(state.settings.language || 'zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    month: '2-digit',
    day: '2-digit',
  }).format(new Date(timestamp));
}

function setActivePeer(peerId) {
  state.activePeerId = peerId;
  renderPeers();
  renderConversation();
}

function activePeer() {
  return state.peers.find((peer) => peer.id === state.activePeerId) || null;
}

function renderStaticText() {
  elements.brandEyebrow.textContent = t('brandEyebrow');
  elements.nicknameLabel.textContent = t('nicknameLabel');
  elements.saveName.textContent = t('save');
  elements.settingsTitle.textContent = t('settings');
  elements.languageLabel.textContent = t('language');
  elements.languageValue.textContent = languageLabel(state.settings.language);
  elements.onlineTitle.textContent = t('onlineNow');
  elements.storageTitle.textContent = t('storage');
  elements.changeDataPath.textContent = t('change');
  elements.storageHint.textContent = t('storageHint');
  elements.sendButton.textContent = t('send');
  elements.messageInput.placeholder = t('messagePlaceholder');
  elements.languageOptions.forEach((option) => {
    const language = option.dataset.language;
    option.textContent = languageLabel(language);
    option.classList.toggle('active', language === state.settings.language);
  });

  if (!elements.dataPath.dataset.loaded) {
    elements.dataPath.textContent = t('loadingPath');
  }
}

function renderPeers() {
  elements.peerCount.textContent = String(state.peers.length);

  if (!state.peers.length) {
    elements.peerList.classList.add('is-empty');
    elements.peerList.innerHTML = `
      <div class="peer-empty">
        <span>${escapeHtml(t('noPeers'))}</span>
      </div>
    `;
    return;
  }

  elements.peerList.classList.remove('is-empty');
  elements.peerList.innerHTML = state.peers.map((peer) => {
    const isActive = peer.id === state.activePeerId;
    return `
      <button class="peer-card ${isActive ? 'active' : ''}" data-peer-id="${peer.id}">
        <div class="peer-avatar">${escapeHtml(peer.name.slice(0, 1).toUpperCase())}</div>
        <div class="peer-copy">
          <strong>${escapeHtml(peer.name)}</strong>
          <span>${escapeHtml(peer.address)}:${peer.listenPort}</span>
        </div>
      </button>
    `;
  }).join('');

  elements.peerList.querySelectorAll('[data-peer-id]').forEach((node) => {
    node.addEventListener('click', () => setActivePeer(node.dataset.peerId));
  });
}

function renderConversation() {
  const peer = activePeer();
  const messages = peer ? (state.conversations[peer.id] || []) : [];

  if (!peer) {
    elements.activePeer.textContent = t('noConversation');
    elements.activeStatus.textContent = t('choosePeer');
    elements.messageList.className = 'message-list empty';
    elements.messageList.innerHTML = `
      <div class="empty-state">
        <h3>${escapeHtml(t('waitingTitle'))}</h3>
        <p>${escapeHtml(t('waitingBody'))}</p>
      </div>
    `;
    elements.messageInput.disabled = true;
    elements.sendButton.disabled = true;
    elements.composerHint.textContent = t('composerIdle');
    return;
  }

  elements.activePeer.textContent = peer.name;
  elements.activeStatus.textContent = `${peer.address}:${peer.listenPort}`;
  elements.messageInput.disabled = false;
  elements.sendButton.disabled = false;
  elements.composerHint.textContent = t('composerReady');

  if (!messages.length) {
    elements.messageList.className = 'message-list empty';
    elements.messageList.innerHTML = `
      <div class="empty-state">
        <h3>${escapeHtml(t('emptyConversationTitle'))}</h3>
        <p>${escapeHtml(t('emptyConversationBody', { name: peer.name }))}</p>
      </div>
    `;
    return;
  }

  elements.messageList.className = 'message-list';
  elements.messageList.innerHTML = messages.map((message) => `
    <article class="message ${message.direction}">
      <div class="message-meta">
        <span>${escapeHtml(message.senderName)}</span>
        <span>${formatTime(message.timestamp)}</span>
        <span class="status status-${message.status}">${escapeHtml(t(message.status))}</span>
      </div>
      <div class="message-bubble">${escapeHtml(message.text).replaceAll('\n', '<br>')}</div>
    </article>
  `).join('');

  elements.messageList.scrollTop = elements.messageList.scrollHeight;
}

async function loadBootstrap() {
  const snapshot = await Bootstrap();
  state.self = snapshot.self;
  state.settings = snapshot.settings || state.settings;
  state.peers = snapshot.peers || [];
  state.conversations = snapshot.conversations || {};

  elements.displayName.value = snapshot.self.name || '';
  elements.profileHint.textContent = t('idAndPort', {
    id: snapshot.self.id.slice(0, 8),
    port: String(snapshot.self.listenPort),
  });

  if (!state.activePeerId && state.peers.length) {
    state.activePeerId = state.peers[0].id;
  }

  renderStaticText();
  renderPeers();
  renderConversation();
}

async function loadDataPath() {
  const path = await DataPath();
  elements.dataPath.textContent = path;
  elements.dataPath.title = path;
  elements.dataPath.dataset.loaded = 'true';
}

elements.saveName.addEventListener('click', async () => {
  const profile = await UpdateDisplayName(elements.displayName.value);
  state.self = profile;
  elements.profileHint.textContent = t('idAndPort', {
    id: profile.id.slice(0, 8),
    port: String(profile.listenPort),
  });
  elements.displayName.value = profile.name;
});

function closeLanguageMenu() {
  elements.languageMenu.classList.add('hidden');
  elements.languageTrigger.setAttribute('aria-expanded', 'false');
  document.body.classList.remove('menu-open');
}

function openLanguageMenu() {
  elements.languageMenu.classList.remove('hidden');
  elements.languageTrigger.setAttribute('aria-expanded', 'true');
  document.body.classList.add('menu-open');
}

elements.languageTrigger.addEventListener('click', () => {
  if (elements.languageMenu.classList.contains('hidden')) {
    openLanguageMenu();
    return;
  }
  closeLanguageMenu();
});

elements.languageOptions.forEach((option) => {
  option.addEventListener('click', async () => {
    closeLanguageMenu();
    const settings = await UpdateLanguage(option.dataset.language);
    state.settings = settings;
    renderStaticText();
    renderPeers();
    renderConversation();
    if (state.self) {
      elements.profileHint.textContent = t('idAndPort', {
        id: state.self.id.slice(0, 8),
        port: String(state.self.listenPort),
      });
    }
  });
});

document.addEventListener('click', (event) => {
  if (!event.target.closest('.language-menu')) {
    closeLanguageMenu();
  }
});

elements.changeDataPath.addEventListener('click', async () => {
  const path = await ChooseDataDirectory();
  elements.dataPath.textContent = path;
  elements.dataPath.title = path;
  elements.dataPath.dataset.loaded = 'true';
});

EventsOn('data:pathUpdated', (path) => {
  elements.dataPath.textContent = path;
  elements.dataPath.title = path;
  elements.dataPath.dataset.loaded = 'true';
});

EventsOn('settings:updated', (settings) => {
  state.settings = settings;
  renderStaticText();
  renderPeers();
  renderConversation();
  if (state.self) {
    elements.profileHint.textContent = t('idAndPort', {
      id: state.self.id.slice(0, 8),
      port: String(state.self.listenPort),
    });
  }
  closeLanguageMenu();
});
elements.composer.addEventListener('submit', async (event) => {
  event.preventDefault();

  if (!state.activePeerId) {
    return;
  }

  const text = elements.messageInput.value.trim();
  if (!text) {
    return;
  }

  await SendChatMessage(state.activePeerId, text);
  elements.messageInput.value = '';
});

elements.messageInput.addEventListener('keydown', async (event) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault();
    elements.composer.requestSubmit();
  }
});

EventsOn('profile:updated', (profile) => {
  state.self = profile;
  elements.displayName.value = profile.name;
  elements.profileHint.textContent = t('idAndPort', {
    id: profile.id.slice(0, 8),
    port: String(profile.listenPort),
  });
});


EventsOn('peers:updated', (peers) => {
  state.peers = peers || [];

  if (state.activePeerId && !state.peers.some((peer) => peer.id === state.activePeerId)) {
    state.activePeerId = state.peers[0]?.id || null;
  }
  if (!state.activePeerId && state.peers.length) {
    state.activePeerId = state.peers[0].id;
  }

  renderPeers();
  renderConversation();
});

EventsOn('conversation:updated', (payload) => {
  state.conversations[payload.peerId] = payload.messages || [];

  if (!state.activePeerId) {
    state.activePeerId = payload.peerId;
  }

  renderConversation();
});

Promise.all([loadBootstrap(), loadDataPath()]).catch((error) => {
  console.error(error);
  elements.messageList.className = 'message-list empty';
  elements.messageList.innerHTML = `
    <div class="empty-state">
      <h3>${escapeHtml(t('bootFailed'))}</h3>
      <p>${escapeHtml(String(error))}</p>
    </div>
  `;
});
