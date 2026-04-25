import './style.css';
import './app.css';

import { EventsOn } from '../wailsjs/runtime/runtime';
import { AddManualPeer, Bootstrap, ChooseDataDirectory, DataPath, EnsureDebugPeer, SendChatMessage, SendImageMessage, UpdateDisplayName, UpdateLanguage, UpdateTheme } from '../wailsjs/go/main/App';

const dictionaries = {
  'zh-CN': {
    brandEyebrow: '局域网桌面通讯',
    nicknameLabel: '你的昵称',
    save: '保存',
    send: '发送',
    settings: '设置',
    language: '语言',
    theme: '主题风格',
    change: '修改',
    debugTools: '本机调试',
    addBot: '添加回声机器人',
    debugHint: '只有一台电脑时，可以先用它验证消息发送、接收和本地存储。',
    manualConnect: '手动连接',
    manualName: '显示名称',
    manualAddress: '地址',
    manualPort: '端口',
    connect: '连接',
    debugTag: '调试',
    manualTag: '手动',
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
    emoji: '表情',
    image: '图片',
    emojiSmileys: '常用笑脸',
    emojiGestures: '手势互动',
    emojiHearts: '情绪氛围',
    emojiNature: '自然元素',
    imageTooLarge: '图片过大，请选择 4MB 以内的图片。',
    imageAlt: '图片消息',
    bootFailed: 'LanTalk 启动失败',
    idAndPort: 'ID: {id}  TCP: {port}',
    sent: '已发送',
    received: '已接收',
    failed: '发送失败',
    langZh: '简体中文',
    langEn: 'English',
    themeMidnight: '午夜蓝',
    themePaper: '暖白纸感',
    themeForest: '森林夜色',
  },
  'en-US': {
    brandEyebrow: 'Local Network Messenger',
    nicknameLabel: 'Your nickname',
    save: 'Save',
    send: 'Send',
    settings: 'Settings',
    language: 'Language',
    theme: 'Theme',
    change: 'Change',
    debugTools: 'Local testing',
    addBot: 'Add Echo Bot',
    debugHint: 'When you only have one computer, use this to verify sending, receiving, and local persistence first.',
    manualConnect: 'Manual connect',
    manualName: 'Display name',
    manualAddress: 'Address',
    manualPort: 'Port',
    connect: 'Connect',
    debugTag: 'Debug',
    manualTag: 'Manual',
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
    emoji: 'Emoji',
    image: 'Image',
    emojiSmileys: 'Smileys',
    emojiGestures: 'Gestures',
    emojiHearts: 'Mood',
    emojiNature: 'Nature',
    imageTooLarge: 'Image is too large. Please choose one under 4 MB.',
    imageAlt: 'Image message',
    bootFailed: 'LanTalk failed to boot',
    idAndPort: 'ID: {id}  TCP: {port}',
    sent: 'sent',
    received: 'received',
    failed: 'failed',
    langZh: 'Simplified Chinese',
    langEn: 'English',
    themeMidnight: 'Midnight Blue',
    themePaper: 'Paper Light',
    themeForest: 'Forest Night',
  },
};

const emojiGroups = [
  {
    labelKey: 'emojiSmileys',
    items: ['😀', '😄', '😁', '😂', '🙂', '😉', '😍', '🥳', '🤔', '😭', '😴', '😎'],
  },
  {
    labelKey: 'emojiGestures',
    items: ['👍', '👋', '👏', '🙏', '💪', '✌️', '👌', '🤝', '🙌', '🤞', '👀', '🔥'],
  },
  {
    labelKey: 'emojiHearts',
    items: ['❤️', '💙', '💚', '💛', '💜', '🧡', '💯', '✨', '⭐', '🎉', '🎈', '🎵'],
  },
  {
    labelKey: 'emojiNature',
    items: ['🌞', '🌙', '🌧️', '⛄', '🌸', '🍀', '🌲', '🌊', '🍎', '🍕', '☕', '🐱'],
  },
];

const maxImageBytes = 4 * 1024 * 1024;

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
        <label class="label" id="themeLabel"></label>
        <div class="theme-menu">
          <button id="themeTrigger" class="select-button" type="button" aria-haspopup="listbox" aria-expanded="false">
            <span id="themeValue"></span>
          </button>
          <div id="themeMenu" class="menu-panel hidden" role="listbox" tabindex="-1">
            <button class="menu-option" type="button" data-theme="midnight"></button>
            <button class="menu-option" type="button" data-theme="paper"></button>
            <button class="menu-option" type="button" data-theme="forest"></button>
          </div>
        </div>
      </section>

      <section class="panel debug-panel">
        <div class="section-head">
          <span id="debugTitle"></span>
        </div>
        <button id="addDebugBot" class="ghost-btn wide-btn" type="button"></button>
        <p class="hint" id="debugHint"></p>
        <div class="debug-grid">
          <label class="field">
            <span id="manualNameLabel" class="field-label"></span>
            <input id="manualName" class="text-input" maxlength="40" autocomplete="off" />
          </label>
          <label class="field">
            <span id="manualAddressLabel" class="field-label"></span>
            <input id="manualAddress" class="text-input" value="127.0.0.1" autocomplete="off" />
          </label>
          <label class="field">
            <span id="manualPortLabel" class="field-label"></span>
            <input id="manualPort" class="text-input" type="number" min="1" max="65535" autocomplete="off" />
          </label>
        </div>
        <button id="manualConnect" class="ghost-btn wide-btn" type="button"></button>
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
        <div class="composer-toolbar">
          <div class="composer-tools">
            <label class="ghost-btn image-trigger" for="imageInput" id="imageTrigger"></label>
            <input id="imageInput" class="visually-hidden" type="file" accept="image/*" />
            <div class="emoji-picker">
            <button id="emojiTrigger" class="ghost-btn emoji-trigger" type="button" aria-haspopup="dialog" aria-expanded="false"></button>
            <div id="emojiPanel" class="emoji-panel hidden" role="dialog" aria-label="emoji panel">
              ${emojiGroups.map((group) => `
                <section class="emoji-group">
                  <div class="emoji-group-title" data-emoji-group-label="${group.labelKey}"></div>
                  <div class="emoji-grid">
                    ${group.items.map((emoji) => `
                      <button class="emoji-option" type="button" data-emoji="${emoji}" aria-label="${emoji}">${emoji}</button>
                    `).join('')}
                  </div>
                </section>
              `).join('')}
            </div>
            </div>
          </div>
        </div>
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
    theme: 'midnight',
  },
  peers: [],
  conversations: {},
  activePeerId: null,
};

const elements = {
  displayName: document.getElementById('displayName'),
  saveName: document.getElementById('saveName'),
  profileHint: document.getElementById('profileHint'),
  languageWrapper: document.querySelector('.language-menu'),
  languageTrigger: document.getElementById('languageTrigger'),
  languageMenu: document.getElementById('languageMenu'),
  languageValue: document.getElementById('languageValue'),
  languageOptions: Array.from(document.querySelectorAll('[data-language]')),
  themeWrapper: document.querySelector('.theme-menu'),
  themeTrigger: document.getElementById('themeTrigger'),
  themeMenu: document.getElementById('themeMenu'),
  themeValue: document.getElementById('themeValue'),
  themeOptions: Array.from(document.querySelectorAll('[data-theme]')),
  peerCount: document.getElementById('peerCount'),
  peerList: document.getElementById('peerList'),
  activePeer: document.getElementById('activePeer'),
  activeStatus: document.getElementById('activeStatus'),
  messageList: document.getElementById('messageList'),
  messageInput: document.getElementById('messageInput'),
  sendButton: document.getElementById('sendButton'),
  imageTrigger: document.getElementById('imageTrigger'),
  imageInput: document.getElementById('imageInput'),
  emojiTrigger: document.getElementById('emojiTrigger'),
  emojiPanel: document.getElementById('emojiPanel'),
  emojiGroupLabels: Array.from(document.querySelectorAll('[data-emoji-group-label]')),
  emojiOptions: Array.from(document.querySelectorAll('[data-emoji]')),
  composerHint: document.getElementById('composerHint'),
  composer: document.getElementById('composer'),
  dataPath: document.getElementById('dataPath'),
  brandEyebrow: document.getElementById('brandEyebrow'),
  nicknameLabel: document.getElementById('nicknameLabel'),
  settingsTitle: document.getElementById('settingsTitle'),
  languageLabel: document.getElementById('languageLabel'),
  themeLabel: document.getElementById('themeLabel'),
  onlineTitle: document.getElementById('onlineTitle'),
  storageTitle: document.getElementById('storageTitle'),
  changeDataPath: document.getElementById('changeDataPath'),
  storageHint: document.getElementById('storageHint'),
  debugTitle: document.getElementById('debugTitle'),
  addDebugBot: document.getElementById('addDebugBot'),
  debugHint: document.getElementById('debugHint'),
  manualName: document.getElementById('manualName'),
  manualAddress: document.getElementById('manualAddress'),
  manualPort: document.getElementById('manualPort'),
  manualConnect: document.getElementById('manualConnect'),
  manualNameLabel: document.getElementById('manualNameLabel'),
  manualAddressLabel: document.getElementById('manualAddressLabel'),
  manualPortLabel: document.getElementById('manualPortLabel'),
};

function languageLabel(language) {
  return language === 'en-US' ? t('langEn') : t('langZh');
}

function themeLabel(theme) {
  switch (theme) {
    case 'paper':
      return t('themePaper');
    case 'forest':
      return t('themeForest');
    default:
      return t('themeMidnight');
  }
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

function renderMessageContent(message) {
  if (message.kind === 'image') {
    return `
      <div class="image-message">
        <img class="message-image" src="${escapeHtml(message.text)}" alt="${escapeHtml(message.mediaName || t('imageAlt'))}" />
        <div class="image-caption">${escapeHtml(message.mediaName || t('imageAlt'))}</div>
      </div>
    `;
  }

  return escapeHtml(message.text).replaceAll('\n', '<br>');
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
  document.documentElement.dataset.theme = state.settings.theme || 'midnight';
  elements.brandEyebrow.textContent = t('brandEyebrow');
  elements.nicknameLabel.textContent = t('nicknameLabel');
  elements.saveName.textContent = t('save');
  elements.settingsTitle.textContent = t('settings');
  elements.languageLabel.textContent = t('language');
  elements.languageValue.textContent = languageLabel(state.settings.language);
  elements.themeLabel.textContent = t('theme');
  elements.themeValue.textContent = themeLabel(state.settings.theme);
  elements.debugTitle.textContent = t('debugTools');
  elements.addDebugBot.textContent = t('addBot');
  elements.debugHint.textContent = t('debugHint');
  elements.manualConnect.textContent = t('connect');
  elements.manualNameLabel.textContent = t('manualName');
  elements.manualAddressLabel.textContent = t('manualAddress');
  elements.manualPortLabel.textContent = t('manualPort');
  elements.onlineTitle.textContent = t('onlineNow');
  elements.storageTitle.textContent = t('storage');
  elements.changeDataPath.textContent = t('change');
  elements.storageHint.textContent = t('storageHint');
  elements.sendButton.textContent = t('send');
  elements.imageTrigger.textContent = t('image');
  elements.emojiTrigger.textContent = t('emoji');
  elements.messageInput.placeholder = t('messagePlaceholder');
  elements.emojiGroupLabels.forEach((node) => {
    node.textContent = t(node.dataset.emojiGroupLabel);
  });
  elements.languageOptions.forEach((option) => {
    const language = option.dataset.language;
    option.textContent = languageLabel(language);
    option.classList.toggle('active', language === state.settings.language);
  });
  elements.themeOptions.forEach((option) => {
    const theme = option.dataset.theme;
    option.textContent = themeLabel(theme);
    option.classList.toggle('active', theme === state.settings.theme);
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
    const tag = peer.source === 'debug'
      ? `<span class="peer-tag">${escapeHtml(t('debugTag'))}</span>`
      : peer.source === 'manual'
        ? `<span class="peer-tag">${escapeHtml(t('manualTag'))}</span>`
        : '';
    return `
      <button class="peer-card ${isActive ? 'active' : ''}" data-peer-id="${peer.id}">
        <div class="peer-avatar">${escapeHtml(peer.name.slice(0, 1).toUpperCase())}</div>
        <div class="peer-copy">
          <strong>${escapeHtml(peer.name)}${tag}</strong>
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
    elements.imageInput.disabled = true;
    elements.imageTrigger.classList.add('disabled');
    elements.emojiTrigger.disabled = true;
    elements.composerHint.textContent = t('composerIdle');
    closeEmojiPanel();
    return;
  }

  elements.activePeer.textContent = peer.name;
  elements.activeStatus.textContent = `${peer.address}:${peer.listenPort}`;
  elements.messageInput.disabled = false;
  elements.sendButton.disabled = false;
  elements.imageInput.disabled = false;
  elements.imageTrigger.classList.remove('disabled');
  elements.emojiTrigger.disabled = false;
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
      <div class="message-bubble ${message.kind === 'image' ? 'message-bubble-image' : ''}">${renderMessageContent(message)}</div>
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
  elements.languageWrapper.classList.remove('open');
  if (elements.themeMenu.classList.contains('hidden')) {
    document.body.classList.remove('menu-open');
  }
}

function closeEmojiPanel() {
  elements.emojiPanel.classList.add('hidden');
  elements.emojiTrigger.setAttribute('aria-expanded', 'false');
}

function openEmojiPanel() {
  elements.emojiPanel.classList.remove('hidden');
  elements.emojiTrigger.setAttribute('aria-expanded', 'true');
}

function openLanguageMenu() {
  closeThemeMenu();
  elements.languageMenu.classList.remove('hidden');
  elements.languageTrigger.setAttribute('aria-expanded', 'true');
  elements.languageWrapper.classList.add('open');
  document.body.classList.add('menu-open');
}

function closeThemeMenu() {
  elements.themeMenu.classList.add('hidden');
  elements.themeTrigger.setAttribute('aria-expanded', 'false');
  elements.themeWrapper.classList.remove('open');
  if (elements.languageMenu.classList.contains('hidden')) {
    document.body.classList.remove('menu-open');
  }
}

function openThemeMenu() {
  closeLanguageMenu();
  elements.themeMenu.classList.remove('hidden');
  elements.themeTrigger.setAttribute('aria-expanded', 'true');
  elements.themeWrapper.classList.add('open');
  document.body.classList.add('menu-open');
}

elements.languageTrigger.addEventListener('click', () => {
  if (elements.languageMenu.classList.contains('hidden')) {
    openLanguageMenu();
    return;
  }
  closeLanguageMenu();
});

elements.themeTrigger.addEventListener('click', () => {
  if (elements.themeMenu.classList.contains('hidden')) {
    openThemeMenu();
    return;
  }
  closeThemeMenu();
});

elements.emojiTrigger.addEventListener('click', () => {
  if (elements.emojiTrigger.disabled) {
    return;
  }
  if (elements.emojiPanel.classList.contains('hidden')) {
    openEmojiPanel();
    elements.messageInput.focus();
    return;
  }
  closeEmojiPanel();
});

function insertEmoji(emoji) {
  const input = elements.messageInput;
  const start = input.selectionStart ?? input.value.length;
  const end = input.selectionEnd ?? input.value.length;
  input.value = `${input.value.slice(0, start)}${emoji}${input.value.slice(end)}`;
  const cursor = start + emoji.length;
  input.focus();
  input.setSelectionRange(cursor, cursor);
}

elements.emojiOptions.forEach((option) => {
  option.addEventListener('click', () => {
    insertEmoji(option.dataset.emoji);
  });
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

elements.themeOptions.forEach((option) => {
  option.addEventListener('click', async () => {
    closeThemeMenu();
    const settings = await UpdateTheme(option.dataset.theme);
    state.settings = settings;
    renderStaticText();
    renderPeers();
    renderConversation();
  });
});

document.addEventListener('click', (event) => {
  if (!event.target.closest('.language-menu')) {
    closeLanguageMenu();
  }
  if (!event.target.closest('.theme-menu')) {
    closeThemeMenu();
  }
  if (!event.target.closest('.emoji-picker')) {
    closeEmojiPanel();
  }
});

elements.changeDataPath.addEventListener('click', async () => {
  const path = await ChooseDataDirectory();
  elements.dataPath.textContent = path;
  elements.dataPath.title = path;
  elements.dataPath.dataset.loaded = 'true';
});

elements.addDebugBot.addEventListener('click', async () => {
  const peer = await EnsureDebugPeer();
  state.activePeerId = peer.id;
  renderPeers();
  renderConversation();
});

elements.imageInput.addEventListener('change', async (event) => {
  const [file] = event.target.files || [];
  if (!file || !state.activePeerId) {
    return;
  }
  if (file.size > maxImageBytes) {
    window.alert(t('imageTooLarge'));
    elements.imageInput.value = '';
    return;
  }

  const dataURL = await new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result);
    reader.onerror = () => reject(reader.error || new Error('failed to read image'));
    reader.readAsDataURL(file);
  });

  await SendImageMessage(state.activePeerId, dataURL, file.name);
  elements.imageInput.value = '';
});

elements.manualConnect.addEventListener('click', async () => {
  const port = Number(elements.manualPort.value);
  if (!port) {
    elements.manualPort.focus();
    return;
  }

  const peer = await AddManualPeer(
    elements.manualName.value,
    elements.manualAddress.value,
    port,
  );
  state.activePeerId = peer.id;
  renderPeers();
  renderConversation();
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
  closeThemeMenu();
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
  closeEmojiPanel();
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
