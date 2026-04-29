import './style.css';
import './app.css';

import { CanResolveFilePaths, ClipboardSetText, EventsOn, InitializeNotifications, IsNotificationAvailable, OnFileDrop, RequestNotificationAuthorization, SendNotification, WindowSetTitle } from '../wailsjs/runtime/runtime';
import { AddManualPeer, Bootstrap, CheckFileMessageAvailable, ChooseDataDirectory, CleanupUnusedMedia, ClearConversation, CopyImageToClipboard, DataPath, DeleteMessage, EnsureDebugPeer, LoadConversationPage, LoadFileSource, LoadImageSource, LoadStorageStats, OpenFileMessage, RetryMessage, RevealFileMessage, RevealImageMessage, SaveFileMessage, SaveImageMessage, SendChatMessage, SendFileMessage, SendLocalFileMessage, SetUnreadCount, UpdateDisplayName, UpdateLanguage, UpdateTheme } from '../wailsjs/go/main/App';

const dictionaries = {
  'zh-CN': {
    brandEyebrow: '局域网桌面通讯',
    moreActions: '更多操作',
    currentNickname: '当前昵称',
    profileSettings: '个人信息',
    nicknameLabel: '你的昵称',
    profileIdLabel: 'ID',
    save: '保存',
    send: '发送',
    settings: '设置',
    appearanceSettings: '界面偏好',
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
    contacts: '联系人',
    onlineNow: '当前在线',
    offline: '已离线',
    storage: '本地存储',
    storagePathLabel: '保存目录',
    storageHint: '聊天记录、媒体文件和程序设置都会保存在这里。',
    storageUpdated: '保存目录已更新',
    storageMessages: '消息总数',
    storageMedia: '媒体占用',
    cleanupMediaLabel: '无引用媒体',
    cleanupMedia: '清理媒体',
    cleanupMediaTitle: '清理无引用媒体？',
    cleanupMediaBody: '这会扫描本地存储目录，移除已经不再被任何聊天记录引用的图片和文件。当前会话内容不会受影响。',
    confirmCleanupMedia: '开始清理',
    cleanupMediaFailed: '清理媒体失败，请稍后再试。',
    cleanupMediaDone: '已清理 {count} 个无引用文件，释放 {size}。',
    cleanupMediaEmpty: '没有发现需要清理的无引用媒体。',
    clearConversation: '清空会话',
    clearConversationRecords: '清空聊天记录',
    clearConversationTitle: '清空当前会话？',
    clearConversationBody: '这会删除与 {name} 的全部本地聊天记录，并尝试清理不再被引用的图片和文件。此操作无法恢复。',
    confirmClearConversation: '清空会话',
    clearConversationFailed: '清空会话失败，请稍后再试。',
    conversationCleared: '已清空 {count} 条消息。',
    messagesStillSending: '仍有消息正在发送，请稍后再试。',
    loadingPath: '正在加载本地路径...',
    choosePeer: '选择一个联系人',
    noConversation: '未选择会话',
    waitingTitle: '正在等待局域网内的其他用户',
    waitingBody: '在同一网络中的另一台电脑上打开 LanTalk，它会自动出现在这里。',
    loadingConversation: '正在加载会话...',
    search: '搜索',
    findConversationContent: '查找聊天内容',
    searchPlaceholder: '搜索当前聊天记录',
    searchLoading: '正在搜索...',
    searchNoResult: '未找到匹配内容',
    searchCount: '{current} / {total}',
    searchPrev: '上一个',
    searchNext: '下一个',
    searchClose: '关闭搜索',
    loadOlderMessages: '加载更早消息',
    loadingOlderMessages: '正在加载更早消息...',
    loadHistoryFailed: '加载聊天记录失败，请稍后再试。',
    messagePlaceholder: '输入消息...',
    composerIdle: '选择在线联系人后即可开始聊天。',
    composerReady: '按 Enter 发送，Shift+Enter 换行。',
    composerPeerOffline: '对方当前已离线，你仍然可以查看历史消息。',
    noPeers: '暂时还没有联系人或历史会话。',
    emptyConversationTitle: '会话还是空的',
    emptyConversationBody: '和 {name} 打个招呼吧。',
    emoji: '表情',
    file: '文件/图片',
    image: '图片',
    emojiSmileys: '常用笑脸',
    emojiGestures: '手势互动',
    emojiHearts: '情绪氛围',
    emojiNature: '自然元素',
    fileTooLarge: '文件过大，请选择 20MB 以内的文件。',
    imageTooLarge: '图片过大，请选择 4MB 以内的图片。',
    fileAlt: '文件消息',
    imageAlt: '图片消息',
    copy: '复制',
    copyFileName: '复制文件名',
    copyImage: '复制图片',
    forward: '转发',
    saveAs: '另存为',
    openFolder: '打开所在位置',
    openFileFailed: '打开文件失败，请稍后再试。',
    openFolderFailed: '打开所在位置失败，请稍后再试。',
    saveFileFailed: '保存文件失败，请稍后再试。',
    sendMessageFailed: '发送消息失败，请稍后再试。',
    copyImageFailed: '当前环境下复制图片失败，请稍后再试。',
    sendFileFailed: '发送文件失败，请稍后再试。',
    sendImageFailed: '发送图片失败，请稍后再试。',
    saveImageFailed: '保存图片失败，请稍后再试。',
    forwardFailed: '转发消息失败，请稍后再试。',
    fileUnavailable: '文件已不存在或已被移动。',
    fileUnavailableStatus: '文件不可用',
    mediaStorageLimitReached: '本地媒体存储空间已达上限，请清理历史媒体后再试。',
    forwarded: '消息已转发',
    peerUnavailable: '该联系人当前不可达，请稍后再试。',
    dismiss: '关闭',
    forwardPreviewLabel: '转发内容',
    forwardTargetLabel: '发送给',
    fileMessageSummary: '文件消息',
    imageMessageSummary: '图片消息',
    forwarding: '转发中...',
    deleteMessage: '删除',
    cancel: '取消',
    forwardDialogTitle: '转发消息',
    forwardDialogBody: '选择一个在线联系人，将这条消息转发过去。',
    forwardEmpty: '当前没有可转发的在线联系人。',
    confirmForward: '发送转发',
    deleteConfirmTitle: '删除这条消息？',
    deleteConfirmBody: '删除后将同时从当前会话和本地存储中移除，且无法恢复。',
    confirmDelete: '删除',
    retryMessage: '重新发送',
    imagePreviewTitle: '图片预览',
    imageViewerHint: '滚轮可在鼠标位置缩放，双击快速切换倍率，左右方向键可切换图片，单击空白处关闭',
    zoomIn: '放大',
    zoomOut: '缩小',
    zoomReset: '原始比例',
    previousImage: '上一张',
    nextImage: '下一张',
    loadImageFailed: '加载图片失败，请稍后再试。',
    closePreview: '关闭预览',
    bootFailed: 'LanTalk 启动失败',
    idAndPort: 'ID: {id}  TCP: {port}',
    sent: '已发送',
    received: '已接收',
    sending: '发送中',
    failed: '发送失败',
    langZh: '简体中文',
    langEn: 'English',
    themeMidnight: '微信简约',
    themePaper: '暖白纸感',
    themeForest: '森林夜色',
  },
  'en-US': {
    brandEyebrow: 'Local Network Messenger',
    moreActions: 'More actions',
    currentNickname: 'Current nickname',
    profileSettings: 'Profile',
    nicknameLabel: 'Your nickname',
    profileIdLabel: 'ID',
    save: 'Save',
    send: 'Send',
    settings: 'Settings',
    appearanceSettings: 'Appearance',
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
    contacts: 'Contacts',
    onlineNow: 'Online now',
    offline: 'Offline',
    storage: 'Local Storage',
    storagePathLabel: 'Storage directory',
    storageHint: 'Chat history, media files, and app settings are stored here.',
    storageUpdated: 'Storage directory updated',
    storageMessages: 'Messages',
    storageMedia: 'Media usage',
    cleanupMediaLabel: 'Unreferenced media',
    cleanupMedia: 'Clean media',
    cleanupMediaTitle: 'Remove unreferenced media?',
    cleanupMediaBody: 'This scans local storage and removes images and files that are no longer referenced by any chat record. Current conversations will not be affected.',
    confirmCleanupMedia: 'Clean media',
    cleanupMediaFailed: 'Failed to clean unreferenced media.',
    cleanupMediaDone: 'Removed {count} unreferenced files and freed {size}.',
    cleanupMediaEmpty: 'No unreferenced media was found.',
    clearConversation: 'Clear chat',
    clearConversationRecords: 'Clear chat history',
    clearConversationTitle: 'Clear this conversation?',
    clearConversationBody: 'This deletes all local chat history with {name} and attempts to clean images and files that are no longer referenced. This cannot be undone.',
    confirmClearConversation: 'Clear chat',
    clearConversationFailed: 'Failed to clear the conversation.',
    conversationCleared: 'Cleared {count} messages.',
    messagesStillSending: 'Some messages are still sending. Please try again shortly.',
    loadingPath: 'Loading local path...',
    choosePeer: 'Choose a peer',
    noConversation: 'No conversation selected',
    waitingTitle: 'Waiting for teammates on the LAN',
    waitingBody: 'Open LanTalk on another computer in the same network and it will appear here automatically.',
    loadingConversation: 'Loading conversation...',
    search: 'Search',
    findConversationContent: 'Find chat content',
    searchPlaceholder: 'Search this conversation',
    searchLoading: 'Searching...',
    searchNoResult: 'No matches found',
    searchCount: '{current} / {total}',
    searchPrev: 'Previous',
    searchNext: 'Next',
    searchClose: 'Close search',
    loadOlderMessages: 'Load earlier messages',
    loadingOlderMessages: 'Loading earlier messages...',
    loadHistoryFailed: 'Failed to load chat history.',
    messagePlaceholder: 'Type a message...',
    composerIdle: 'Pick an online peer to start chatting.',
    composerReady: 'Enter sends. Shift+Enter inserts a new line.',
    composerPeerOffline: 'This contact is offline right now, but you can still browse the chat history.',
    noPeers: 'No contacts or conversation history yet.',
    emptyConversationTitle: 'Conversation is empty',
    emptyConversationBody: 'Say hello to {name}.',
    emoji: 'Emoji',
    file: 'Files / Images',
    image: 'Image',
    emojiSmileys: 'Smileys',
    emojiGestures: 'Gestures',
    emojiHearts: 'Mood',
    emojiNature: 'Nature',
    fileTooLarge: 'File is too large. Please choose one under 20 MB.',
    imageTooLarge: 'Image is too large. Please choose one under 4 MB.',
    fileAlt: 'File message',
    imageAlt: 'Image message',
    copy: 'Copy',
    copyFileName: 'Copy file name',
    copyImage: 'Copy image',
    forward: 'Forward',
    saveAs: 'Save as',
    openFolder: 'Show in folder',
    openFileFailed: 'Failed to open the file.',
    openFolderFailed: 'Failed to show the file in its folder.',
    saveFileFailed: 'Failed to save the file.',
    sendMessageFailed: 'Failed to send the message.',
    copyImageFailed: 'Failed to copy the image to the system clipboard.',
    sendFileFailed: 'Failed to send the file.',
    sendImageFailed: 'Failed to send the image.',
    saveImageFailed: 'Failed to save the image.',
    forwardFailed: 'Failed to forward the message.',
    fileUnavailable: 'This file is no longer available.',
    fileUnavailableStatus: 'Unavailable',
    mediaStorageLimitReached: 'Local media storage limit has been reached. Remove old media and try again.',
    forwarded: 'Message forwarded',
    peerUnavailable: 'This contact is currently unreachable.',
    dismiss: 'Dismiss',
    forwardPreviewLabel: 'Forward content',
    forwardTargetLabel: 'Send to',
    fileMessageSummary: 'File message',
    imageMessageSummary: 'Image message',
    forwarding: 'Forwarding...',
    deleteMessage: 'Delete',
    cancel: 'Cancel',
    forwardDialogTitle: 'Forward message',
    forwardDialogBody: 'Choose an available contact to forward this message to.',
    forwardEmpty: 'No available contacts can receive this message right now.',
    confirmForward: 'Forward',
    deleteConfirmTitle: 'Delete this message?',
    deleteConfirmBody: 'This will remove it from the current conversation and local storage. This action cannot be undone.',
    confirmDelete: 'Delete',
    retryMessage: 'Retry send',
    imagePreviewTitle: 'Image Preview',
    imageViewerHint: 'Use the mouse wheel to zoom at the pointer, double-click to toggle zoom, use the arrow keys to switch images, and click outside to close.',
    zoomIn: 'Zoom In',
    zoomOut: 'Zoom Out',
    zoomReset: 'Reset',
    previousImage: 'Previous',
    nextImage: 'Next',
    loadImageFailed: 'Failed to load the image.',
    closePreview: 'Close preview',
    bootFailed: 'LanTalk failed to boot',
    idAndPort: 'ID: {id}  TCP: {port}',
    sent: 'sent',
    received: 'received',
    sending: 'sending',
    failed: 'failed',
    langZh: 'Simplified Chinese',
    langEn: 'English',
    themeMidnight: 'WeChat Green',
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
const maxFileBytes = 20 * 1024 * 1024;
const conversationPageSize = 40;
const historyAutoLoadThreshold = 36;
const conversationStickToBottomThreshold = 24;
const peerOfflineThresholdMs = 8 * 1000;
const viewerMinScale = 1;
const viewerMaxScale = 4;
const viewerWheelIntensity = 0.0018;
const viewerDoubleClickScale = 2;

document.querySelector('#app').innerHTML = `
  <div class="shell">
    <aside class="sidebar">
      <div class="brand">
        <div id="brandMark" class="brand-mark">L</div>
        <h1 id="brandName" class="brand-title">LanTalk</h1>
        <button id="openSettings" class="ghost-btn brand-settings-btn" type="button"></button>
      </div>

      <section class="panel peers-panel">
        <div class="section-head">
          <span id="onlineTitle"></span>
          <span class="badge" id="peerCount">0</span>
        </div>
        <div id="peerList" class="peer-list"></div>
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

    </aside>

    <main class="chat-area">
      <div class="chat-header">
        <div class="chat-header-main">
          <div class="chat-header-copy">
            <div class="eyebrow" id="activeStatus"></div>
            <h2 id="activePeer"></h2>
          </div>
          <div class="chat-header-actions">
            <div class="chat-more">
              <button id="chatMenuTrigger" class="ghost-btn chat-more-trigger" type="button" aria-haspopup="menu" aria-expanded="false"></button>
              <div id="chatMenuPanel" class="chat-more-panel hidden" role="menu">
                <button id="searchToggle" class="chat-more-item" type="button" role="menuitem"></button>
                <button id="clearConversation" class="chat-more-item danger" type="button" role="menuitem"></button>
              </div>
            </div>
          </div>
        </div>
        <div id="chatSearchBar" class="chat-search-bar hidden">
          <input id="chatSearchInput" class="text-input chat-search-input" maxlength="120" autocomplete="off" />
          <span id="chatSearchMeta" class="chat-search-meta"></span>
          <button id="chatSearchPrev" class="ghost-btn chat-search-nav" type="button"></button>
          <button id="chatSearchNext" class="ghost-btn chat-search-nav" type="button"></button>
          <button id="chatSearchClose" class="ghost-btn chat-search-close" type="button"></button>
        </div>
      </div>

      <section class="message-list-shell">
        <section id="messageList" class="message-list empty"></section>
      </section>

      <form id="composer" class="composer">
        <div class="composer-toolbar">
          <div class="composer-tools">
            <label class="ghost-btn file-trigger" for="fileInput" id="fileTrigger"></label>
            <input id="fileInput" class="visually-hidden" type="file" multiple />
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
        <div class="composer-input-shell">
          <textarea id="messageInput" maxlength="2000" disabled></textarea>
        </div>
        <div class="composer-actions">
          <span class="hint" id="composerHint"></span>
          <button type="submit" id="sendButton" class="primary-btn" disabled></button>
        </div>
      </form>
    </main>
  </div>
  <div id="settingsDialog" class="dialog-overlay hidden" aria-hidden="true">
    <div class="dialog-backdrop"></div>
    <div class="dialog-card dialog-card-settings" role="dialog" aria-modal="true" aria-labelledby="settingsDialogTitle">
      <div class="dialog-header-row">
        <h3 id="settingsDialogTitle" class="dialog-title"></h3>
        <button id="settingsDialogClose" class="ghost-btn dialog-btn" type="button"></button>
      </div>
      <div class="settings-dialog-layout">
        <section class="settings-dialog-section">
          <div class="section-head">
            <span id="profileSettingsTitle"></span>
          </div>
          <div class="settings-form-row settings-form-row-edit">
            <label class="settings-form-label" id="nicknameLabel" for="displayName"></label>
            <div class="settings-form-control settings-form-control-inline">
              <input id="displayName" class="text-input" maxlength="40" autocomplete="off" />
              <button id="saveName" class="primary-btn"></button>
            </div>
          </div>
          <div class="settings-form-row">
            <span class="settings-form-label" id="profileIdLabel"></span>
            <div id="profileId" class="text-input profile-readonly settings-form-value" aria-readonly="true"></div>
          </div>
        </section>
        <section class="settings-dialog-section">
          <div class="section-head">
            <span id="settingsTitle"></span>
          </div>
          <div class="settings-form-row">
            <span class="settings-form-label" id="languageLabel"></span>
            <div class="settings-form-control">
              <div class="language-menu">
                <button id="languageTrigger" class="select-button" type="button" aria-haspopup="listbox" aria-expanded="false">
                  <span id="languageValue"></span>
                </button>
                <div id="languageMenu" class="menu-panel hidden" role="listbox" tabindex="-1">
                  <button class="menu-option" type="button" data-language="zh-CN"></button>
                  <button class="menu-option" type="button" data-language="en-US"></button>
                </div>
              </div>
            </div>
          </div>
          <div class="settings-form-row">
            <span class="settings-form-label" id="themeLabel"></span>
            <div class="settings-form-control">
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
            </div>
          </div>
        </section>
        <section class="settings-dialog-section">
          <div class="section-head">
            <div class="section-copy">
              <span id="storageTitle"></span>
              <p class="hint section-note" id="storageHint"></p>
            </div>
          </div>
          <div class="settings-form-row">
            <span class="settings-form-label" id="storagePathLabel"></span>
            <div class="storage-path-layout">
              <div class="storage-path-card">
                <p class="hint data-path" id="dataPath"></p>
              </div>
              <button id="changeDataPath" class="ghost-btn storage-manage-btn" type="button"></button>
            </div>
          </div>
          <div class="settings-form-row">
            <span class="settings-form-label" id="storageMessageCountLabel"></span>
            <div id="storageMessageCount" class="text-input profile-readonly settings-form-value" aria-readonly="true"></div>
          </div>
          <div class="settings-form-row">
            <span class="settings-form-label" id="storageMediaUsageLabel"></span>
            <div id="storageMediaUsage" class="text-input profile-readonly settings-form-value" aria-readonly="true"></div>
          </div>
          <div class="settings-form-row">
            <span class="settings-form-label" id="cleanupMediaLabel"></span>
            <button id="cleanupMedia" class="ghost-btn storage-clean-btn" type="button"></button>
          </div>
        </section>
      </div>
    </div>
  </div>
  <div id="imageViewer" class="image-viewer hidden" aria-hidden="true">
    <div class="image-viewer-backdrop"></div>
    <div class="image-viewer-dialog">
      <div class="image-viewer-toolbar">
        <div class="image-viewer-toolbar-copy">
          <div id="imageViewerTitle" class="image-viewer-title"></div>
          <span id="imageViewerCounter" class="image-viewer-counter hidden"></span>
          <p id="imageViewerHint" class="hint image-viewer-hint"></p>
        </div>
        <div class="image-viewer-actions">
          <button id="imageViewerZoomOut" class="ghost-btn image-viewer-action" type="button"></button>
          <button id="imageViewerReset" class="ghost-btn image-viewer-action" type="button"></button>
          <span id="imageViewerScale" class="image-viewer-scale">100%</span>
          <button id="imageViewerZoomIn" class="ghost-btn image-viewer-action" type="button"></button>
          <button id="imageViewerClose" class="ghost-btn image-viewer-close" type="button"></button>
        </div>
      </div>
      <div class="image-viewer-body">
        <div class="image-viewer-rail image-viewer-rail-start">
          <button id="imageViewerPrev" class="image-viewer-nav hidden" type="button">
            <span aria-hidden="true">‹</span>
          </button>
        </div>
        <div id="imageViewerFrame" class="image-viewer-frame">
          <div id="imageViewerStage" class="image-viewer-stage">
            <img id="imageViewerImg" class="image-viewer-img" alt="" />
          </div>
        </div>
        <div class="image-viewer-rail image-viewer-rail-end">
          <button id="imageViewerNext" class="image-viewer-nav hidden" type="button">
            <span aria-hidden="true">›</span>
          </button>
        </div>
      </div>
    </div>
  </div>
  <div id="messageContextMenu" class="message-context-menu hidden" aria-hidden="true">
    <button id="messageContextCopy" class="message-context-item" type="button"></button>
    <button id="messageContextSave" class="message-context-item hidden" type="button"></button>
    <button id="messageContextReveal" class="message-context-item hidden" type="button"></button>
    <button id="messageContextRetry" class="message-context-item hidden" type="button"></button>
    <button id="messageContextForward" class="message-context-item" type="button"></button>
    <button id="messageContextDelete" class="message-context-item danger" type="button"></button>
  </div>
  <div id="deleteConfirmDialog" class="dialog-overlay hidden" aria-hidden="true">
    <div class="dialog-backdrop"></div>
    <div class="dialog-card" role="dialog" aria-modal="true" aria-labelledby="deleteConfirmTitle">
      <h3 id="deleteConfirmTitle" class="dialog-title"></h3>
      <p id="deleteConfirmBody" class="dialog-body"></p>
      <div class="dialog-actions">
        <button id="deleteConfirmCancel" class="ghost-btn dialog-btn" type="button"></button>
        <button id="deleteConfirmSubmit" class="primary-btn dialog-btn danger-btn" type="button"></button>
      </div>
    </div>
  </div>
  <div id="actionConfirmDialog" class="dialog-overlay hidden" aria-hidden="true">
    <div class="dialog-backdrop"></div>
    <div class="dialog-card" role="dialog" aria-modal="true" aria-labelledby="actionConfirmTitle">
      <h3 id="actionConfirmTitle" class="dialog-title"></h3>
      <p id="actionConfirmBody" class="dialog-body"></p>
      <div class="dialog-actions">
        <button id="actionConfirmCancel" class="ghost-btn dialog-btn" type="button"></button>
        <button id="actionConfirmSubmit" class="primary-btn dialog-btn" type="button"></button>
      </div>
    </div>
  </div>
  <div id="forwardDialog" class="dialog-overlay hidden" aria-hidden="true">
    <div class="dialog-backdrop"></div>
    <div class="dialog-card dialog-card-wide" role="dialog" aria-modal="true" aria-labelledby="forwardDialogTitle">
      <h3 id="forwardDialogTitle" class="dialog-title"></h3>
      <p id="forwardDialogBody" class="dialog-body"></p>
      <div class="forward-preview-block">
        <span id="forwardPreviewLabel" class="forward-section-label"></span>
        <div id="forwardPreview" class="forward-preview-card"></div>
      </div>
      <span id="forwardTargetLabel" class="forward-section-label"></span>
      <div id="forwardPeerList" class="forward-peer-list"></div>
      <p id="forwardDialogEmpty" class="hint forward-empty hidden"></p>
      <div class="dialog-actions">
        <button id="forwardDialogCancel" class="ghost-btn dialog-btn" type="button"></button>
        <button id="forwardDialogSubmit" class="primary-btn dialog-btn" type="button"></button>
      </div>
    </div>
  </div>
  <div id="toastRegion" class="toast-region" aria-live="polite" aria-atomic="true"></div>
`;

const state = {
  self: null,
  settings: {
    language: 'zh-CN',
    theme: 'midnight',
  },
  peers: [],
  conversations: {},
  conversationMeta: {},
  imageCache: {},
  imageLoads: {},
  fileAvailability: {},
  fileChecks: {},
  transferProgress: {},
  storageStats: {
    messageCount: 0,
    mediaBytes: 0,
    loaded: false,
    refreshTimer: 0,
  },
  settingsDialog: {
    open: false,
  },
  chatMenu: {
    open: false,
  },
  preserveConversationScroll: false,
  activePeerId: null,
  viewer: {
    open: false,
    peerId: '',
    messageId: '',
    src: '',
    name: '',
    loading: false,
    scale: viewerMinScale,
    naturalWidth: 0,
    naturalHeight: 0,
    baseWidth: 0,
    baseHeight: 0,
  },
  messageMenu: {
    open: false,
    peerId: '',
    messageId: '',
    kind: 'text',
  },
  deleteDialog: {
    open: false,
    peerId: '',
    messageId: '',
  },
  actionDialog: {
    open: false,
    mode: '',
    peerId: '',
    busy: false,
  },
  forwardDialog: {
    open: false,
    peerId: '',
    messageId: '',
    targetPeerId: '',
  },
  search: {
    open: false,
    query: '',
    results: [],
    activeIndex: -1,
    loading: false,
    token: 0,
    peerId: '',
  },
  unreadCounts: {},
  windowFocused: document.hasFocus(),
  windowVisible: true,
  notifications: {
    available: false,
  },
  lastUnreadTotal: -1,
};

const elements = {
  openSettings: document.getElementById('openSettings'),
  brandMark: document.getElementById('brandMark'),
  brandName: document.getElementById('brandName'),
  displayName: document.getElementById('displayName'),
  saveName: document.getElementById('saveName'),
  profileIdLabel: document.getElementById('profileIdLabel'),
  profileId: document.getElementById('profileId'),
  settingsDialog: document.getElementById('settingsDialog'),
  settingsDialogTitle: document.getElementById('settingsDialogTitle'),
  settingsDialogClose: document.getElementById('settingsDialogClose'),
  chatMenuTrigger: document.getElementById('chatMenuTrigger'),
  chatMenuPanel: document.getElementById('chatMenuPanel'),
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
  fileTrigger: document.getElementById('fileTrigger'),
  fileInput: document.getElementById('fileInput'),
  emojiTrigger: document.getElementById('emojiTrigger'),
  emojiPanel: document.getElementById('emojiPanel'),
  clearConversation: document.getElementById('clearConversation'),
  searchToggle: document.getElementById('searchToggle'),
  chatSearchBar: document.getElementById('chatSearchBar'),
  chatSearchInput: document.getElementById('chatSearchInput'),
  chatSearchMeta: document.getElementById('chatSearchMeta'),
  chatSearchPrev: document.getElementById('chatSearchPrev'),
  chatSearchNext: document.getElementById('chatSearchNext'),
  chatSearchClose: document.getElementById('chatSearchClose'),
  emojiGroupLabels: Array.from(document.querySelectorAll('[data-emoji-group-label]')),
  emojiOptions: Array.from(document.querySelectorAll('[data-emoji]')),
  composerHint: document.getElementById('composerHint'),
  composer: document.getElementById('composer'),
  imageViewer: document.getElementById('imageViewer'),
  imageViewerFrame: document.getElementById('imageViewerFrame'),
  imageViewerStage: document.getElementById('imageViewerStage'),
  imageViewerImg: document.getElementById('imageViewerImg'),
  imageViewerTitle: document.getElementById('imageViewerTitle'),
  imageViewerCounter: document.getElementById('imageViewerCounter'),
  imageViewerHint: document.getElementById('imageViewerHint'),
  imageViewerScale: document.getElementById('imageViewerScale'),
  imageViewerZoomOut: document.getElementById('imageViewerZoomOut'),
  imageViewerZoomIn: document.getElementById('imageViewerZoomIn'),
  imageViewerReset: document.getElementById('imageViewerReset'),
  imageViewerPrev: document.getElementById('imageViewerPrev'),
  imageViewerNext: document.getElementById('imageViewerNext'),
  imageViewerClose: document.getElementById('imageViewerClose'),
  messageContextMenu: document.getElementById('messageContextMenu'),
  messageContextCopy: document.getElementById('messageContextCopy'),
  messageContextSave: document.getElementById('messageContextSave'),
  messageContextReveal: document.getElementById('messageContextReveal'),
  messageContextRetry: document.getElementById('messageContextRetry'),
  messageContextForward: document.getElementById('messageContextForward'),
  messageContextDelete: document.getElementById('messageContextDelete'),
  deleteConfirmDialog: document.getElementById('deleteConfirmDialog'),
  deleteConfirmTitle: document.getElementById('deleteConfirmTitle'),
  deleteConfirmBody: document.getElementById('deleteConfirmBody'),
  deleteConfirmCancel: document.getElementById('deleteConfirmCancel'),
  deleteConfirmSubmit: document.getElementById('deleteConfirmSubmit'),
  actionConfirmDialog: document.getElementById('actionConfirmDialog'),
  actionConfirmTitle: document.getElementById('actionConfirmTitle'),
  actionConfirmBody: document.getElementById('actionConfirmBody'),
  actionConfirmCancel: document.getElementById('actionConfirmCancel'),
  actionConfirmSubmit: document.getElementById('actionConfirmSubmit'),
  forwardDialog: document.getElementById('forwardDialog'),
  forwardDialogTitle: document.getElementById('forwardDialogTitle'),
  forwardDialogBody: document.getElementById('forwardDialogBody'),
  forwardPreviewLabel: document.getElementById('forwardPreviewLabel'),
  forwardPreview: document.getElementById('forwardPreview'),
  forwardTargetLabel: document.getElementById('forwardTargetLabel'),
  forwardPeerList: document.getElementById('forwardPeerList'),
  forwardDialogEmpty: document.getElementById('forwardDialogEmpty'),
  forwardDialogCancel: document.getElementById('forwardDialogCancel'),
  forwardDialogSubmit: document.getElementById('forwardDialogSubmit'),
  toastRegion: document.getElementById('toastRegion'),
  dataPath: document.getElementById('dataPath'),
  storageMessageCount: document.getElementById('storageMessageCount'),
  storageMessageCountLabel: document.getElementById('storageMessageCountLabel'),
  storageMediaUsage: document.getElementById('storageMediaUsage'),
  storageMediaUsageLabel: document.getElementById('storageMediaUsageLabel'),
  profileSettingsTitle: document.getElementById('profileSettingsTitle'),
  nicknameLabel: document.getElementById('nicknameLabel'),
  settingsTitle: document.getElementById('settingsTitle'),
  languageLabel: document.getElementById('languageLabel'),
  themeLabel: document.getElementById('themeLabel'),
  onlineTitle: document.getElementById('onlineTitle'),
  storageTitle: document.getElementById('storageTitle'),
  storagePathLabel: document.getElementById('storagePathLabel'),
  changeDataPath: document.getElementById('changeDataPath'),
  cleanupMediaLabel: document.getElementById('cleanupMediaLabel'),
  cleanupMedia: document.getElementById('cleanupMedia'),
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

function renderSelfProfile() {
  const name = String(state.self?.name || '').trim() || 'LanTalk';
  elements.brandName.textContent = name;
  elements.brandName.title = name;
  elements.brandMark.textContent = Array.from(name)[0] || 'L';
  elements.displayName.value = state.self?.name || '';

  if (state.self?.id) {
    const fullID = String(state.self.id);
    elements.profileId.textContent = fullID.slice(0, 8);
    elements.profileId.title = fullID;
    return;
  }

  elements.profileId.textContent = '';
  elements.profileId.title = '';
}

function t(key, vars = {}) {
  const language = state.settings.language || 'zh-CN';
  const dictionary = dictionaries[language] || dictionaries['zh-CN'];
  const keyStr = (typeof key === 'string') ? key : String(key ?? '');
  let value = dictionary[keyStr];
  if (value === undefined || value === null) {
    value = keyStr;
  }
  value = String(value);

  Object.entries(vars).forEach(([name, replacement]) => {
    if (String.prototype.replaceAll) {
      value = value.replaceAll(`{${name}}`, String(replacement));
    } else {
      value = value.split(`{${name}}`).join(String(replacement));
    }
  });

  return value;
}

function escapeHtml(value) {
  const s = value === undefined || value === null ? '' : String(value);
  if (String.prototype.replaceAll) {
    return s
      .replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll('"', '&quot;')
      .replaceAll("'", '&#39;');
  }
  return s
    .split('&').join('&amp;')
    .split('<').join('&lt;')
    .split('>').join('&gt;')
    .split('"').join('&quot;')
    .split("'").join('&#39;');
}

function escapeHtmlWithLineBreaks(value) {
  const escaped = escapeHtml(value);
  if (String.prototype.replaceAll) {
    return escaped.replaceAll('\n', '<br>');
  }
  return escaped.split('\n').join('<br>');
}

function isImageDataSource(value) {
  return String(value || '').trim().startsWith('data:image/');
}

function isFileDataSource(value) {
  return String(value || '').trim().startsWith('data:');
}

function imageSourceKey(message) {
  if (!message || message.kind !== 'image') {
    return '';
  }
  const value = String(message.text || '').trim();
  return isImageDataSource(value) ? '' : value;
}

function fileSourceKey(message) {
  if (!message || message.kind !== 'file') {
    return '';
  }
  const value = String(message.text || '').trim();
  return isFileDataSource(value) ? '' : value;
}

function cachedImageSource(message) {
  if (!message || message.kind !== 'image') {
    return '';
  }
  const value = String(message.text || '').trim();
  if (!value) {
    return '';
  }
  if (isImageDataSource(value)) {
    return value;
  }
  return state.imageCache[value] || '';
}

function cachedFileAvailability(message) {
  if (!message || message.kind !== 'file') {
    return null;
  }
  const fileKey = fileSourceKey(message);
  if (!fileKey) {
    return true;
  }
  if (Object.prototype.hasOwnProperty.call(state.fileAvailability, fileKey)) {
    return state.fileAvailability[fileKey];
  }
  return null;
}

function refreshImageMessageNodes(imageKey, source) {
  if (!imageKey || !source) {
    return;
  }

  document.querySelectorAll(`.image-message[data-image-ref="${CSS.escape(imageKey)}"]`).forEach((node) => {
    node.classList.remove('image-message-loading');
    const placeholder = node.querySelector('.message-image-placeholder');
    const image = node.querySelector('.message-image');
    placeholder?.classList.add('hidden');
    if (image) {
      image.src = source;
      image.dataset.previewSrc = source;
      image.classList.remove('hidden');
    }
  });
}

function applyFileMessageNodeAvailability(node, available) {
  if (!node) {
    return;
  }

  const isMissing = available === false;
  node.classList.toggle('is-missing', isMissing);
  node.dataset.fileStatus = isMissing ? 'missing' : 'ready';

  const metaNode = node.querySelector('.file-message-copy span');
  if (metaNode) {
    metaNode.textContent = isMissing ? t('fileUnavailableStatus') : (node.dataset.fileMeta || '');
  }
}

function refreshFileMessageNodes(fileKey, available) {
  if (!fileKey) {
    return;
  }

  document.querySelectorAll(`.file-message[data-file-ref="${CSS.escape(fileKey)}"]`).forEach((node) => {
    applyFileMessageNodeAvailability(node, available);
  });
}

async function resolveImageMessageSource(message) {
  if (!message || message.kind !== 'image') {
    return '';
  }

  const directSource = cachedImageSource(message);
  if (directSource) {
    return directSource;
  }

  const imageKey = imageSourceKey(message);
  if (!imageKey) {
    return '';
  }

  if (!state.imageLoads[imageKey]) {
    state.imageLoads[imageKey] = LoadImageSource(imageKey, message.mediaType || '')
      .then((source) => {
        state.imageCache[imageKey] = source;
        refreshImageMessageNodes(imageKey, source);
        return source;
      })
      .finally(() => {
        delete state.imageLoads[imageKey];
      });
  }

  return state.imageLoads[imageKey];
}

async function resolveFileMessageAvailability(message) {
  if (!message || message.kind !== 'file') {
    return false;
  }

  const fileKey = fileSourceKey(message);
  if (!fileKey) {
    return true;
  }

  if (Object.prototype.hasOwnProperty.call(state.fileAvailability, fileKey)) {
    return state.fileAvailability[fileKey];
  }

  if (!state.fileChecks[fileKey]) {
    state.fileChecks[fileKey] = Promise.resolve(CheckFileMessageAvailable(fileKey))
      .then((available) => {
        const resolved = Boolean(available);
        state.fileAvailability[fileKey] = resolved;
        refreshFileMessageNodes(fileKey, resolved);
        return resolved;
      })
      .finally(() => {
        delete state.fileChecks[fileKey];
      });
  }

  return state.fileChecks[fileKey];
}

function warmConversationImages(messages) {
  const uniqueMessages = new Map();
  messages.forEach((message) => {
    if (message.kind !== 'image') {
      return;
    }
    const imageKey = imageSourceKey(message);
    if (!imageKey || uniqueMessages.has(imageKey)) {
      return;
    }
    uniqueMessages.set(imageKey, message);
  });

  uniqueMessages.forEach((message) => {
    resolveImageMessageSource(message).catch((error) => {
      console.error('LoadImageSource failed', error);
    });
  });
}

function warmConversationFiles(messages) {
  const uniqueMessages = new Map();
  messages.forEach((message) => {
    if (message.kind !== 'file') {
      return;
    }
    const fileKey = fileSourceKey(message);
    if (!fileKey || uniqueMessages.has(fileKey)) {
      return;
    }
    uniqueMessages.set(fileKey, message);
  });

  uniqueMessages.forEach((message) => {
    resolveFileMessageAvailability(message).catch((error) => {
      console.error('CheckFileMessageAvailable failed', error);
    });
  });
}

function formatFileSize(bytes) {
  const size = Number(bytes) || 0;
  if (size <= 0) {
    return '0 B';
  }

  if (size < 1024) {
    return `${size} B`;
  }

  const units = ['KB', 'MB', 'GB'];
  let value = size / 1024;
  let unitIndex = 0;
  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024;
    unitIndex += 1;
  }
  const fractionDigits = value >= 100 ? 0 : value >= 10 ? 1 : 2;
  return `${value.toFixed(fractionDigits)} ${units[unitIndex]}`;
}

function renderStorageStats() {
  elements.storageMessageCountLabel.textContent = t('storageMessages');
  elements.storageMediaUsageLabel.textContent = t('storageMedia');

  if (!state.storageStats.loaded) {
    elements.storageMessageCount.textContent = '--';
    elements.storageMediaUsage.textContent = '--';
    return;
  }

  elements.storageMessageCount.textContent = String(Math.max(0, Number(state.storageStats.messageCount) || 0));
  elements.storageMediaUsage.textContent = formatFileSize(state.storageStats.mediaBytes);
}

async function refreshStorageStats() {
  try {
    const stats = await LoadStorageStats();
    state.storageStats.messageCount = Math.max(0, Number(stats?.messageCount) || 0);
    state.storageStats.mediaBytes = Math.max(0, Number(stats?.mediaBytes) || 0);
    state.storageStats.loaded = true;
    renderStorageStats();
  } catch (error) {
    console.error('LoadStorageStats failed', error);
  }
}

function scheduleStorageStatsRefresh() {
  if (state.storageStats.refreshTimer) {
    window.clearTimeout(state.storageStats.refreshTimer);
  }

  state.storageStats.refreshTimer = window.setTimeout(() => {
    state.storageStats.refreshTimer = 0;
    void refreshStorageStats();
  }, 180);
}

function fileMessageName(message) {
  return (message?.mediaName || '').trim() || t('fileAlt');
}

function fileMessageExtension(message) {
  const name = fileMessageName(message);
  const extIndex = name.lastIndexOf('.');
  if (extIndex < 0 || extIndex >= name.length - 1) {
    return '';
  }
  return name.slice(extIndex + 1).toLowerCase();
}

function fileMessageCategory(message) {
  const extension = fileMessageExtension(message);
  const mediaType = String(message?.mediaType || '').trim().toLowerCase();

  if (extension === 'pdf' || mediaType === 'application/pdf') {
    return 'pdf';
  }
  if (['doc', 'docx', 'odt', 'rtf', 'ppt', 'pptx', 'key'].includes(extension)) {
    return 'document';
  }
  if (['xls', 'xlsx', 'csv', 'ods'].includes(extension)) {
    return 'sheet';
  }
  if (['zip', 'rar', '7z', 'tar', 'gz', 'bz2', 'xz'].includes(extension)) {
    return 'archive';
  }
  if (mediaType.startsWith('image/')) {
    return 'image';
  }
  if (mediaType.startsWith('audio/') || mediaType.startsWith('video/')) {
    return 'media';
  }
  if (['js', 'ts', 'jsx', 'tsx', 'json', 'xml', 'yaml', 'yml', 'md', 'go', 'py', 'java', 'c', 'cpp', 'h', 'hpp', 'rs', 'sql', 'html', 'css', 'scss', 'vue'].includes(extension)) {
    return 'code';
  }
  if (['txt', 'log', 'ini', 'conf'].includes(extension) || mediaType.startsWith('text/')) {
    return 'text';
  }
  return 'generic';
}

function fileMessageBadge(message) {
  const extension = fileMessageExtension(message);
  if (extension) {
    return extension.slice(0, 4).toUpperCase();
  }

  const category = fileMessageCategory(message);
  if (category === 'pdf') {
    return 'PDF';
  }
  if (category === 'document') {
    return 'DOC';
  }
  if (category === 'sheet') {
    return 'XLS';
  }
  if (category === 'archive') {
    return 'ZIP';
  }
  if (category === 'image') {
    return 'IMG';
  }
  if (category === 'media') {
    return 'MED';
  }
  if (category === 'code') {
    return 'CODE';
  }
  if (category === 'text') {
    return 'TXT';
  }

  const subtype = String(message?.mediaType || '')
    .trim()
    .split(';')[0]
    .split('/')
    .pop();
  if (subtype) {
    return subtype.replace(/[^a-z0-9]/gi, '').slice(0, 4).toUpperCase() || 'FILE';
  }

  return 'FILE';
}

function fileMessageMeta(message) {
  const sizeLabel = formatFileSize(message?.mediaSize);
  if (sizeLabel) {
    return sizeLabel;
  }

  const mediaType = String(message?.mediaType || '').trim();
  if (mediaType && mediaType !== 'application/octet-stream') {
    return mediaType;
  }

  return t('fileMessageSummary');
}

function setFileAvailability(message, available) {
  const fileKey = fileSourceKey(message);
  if (!fileKey) {
    return;
  }
  state.fileAvailability[fileKey] = available;
  refreshFileMessageNodes(fileKey, available);
}

function transferProgressForMessage(messageId) {
  const value = Number(state.transferProgress[messageId]);
  if (!Number.isFinite(value)) {
    return null;
  }
  return Math.max(0, Math.min(100, Math.round(value)));
}

function clearTransferProgressForMessages(messages) {
  (messages || []).forEach((message) => {
    if (!message?.id || message.status === 'sending') {
      return;
    }
    delete state.transferProgress[message.id];
  });
}

function syncTransferProgressCache(previousMessages, nextMessages) {
  const nextIDs = new Set((nextMessages || []).map((message) => message?.id).filter(Boolean));
  (previousMessages || []).forEach((message) => {
    if (!message?.id || nextIDs.has(message.id)) {
      return;
    }
    delete state.transferProgress[message.id];
  });

  clearTransferProgressForMessages(nextMessages);
}

function messageCanRetry(message) {
  return Boolean(message)
    && message.direction === 'outbound'
    && message.status === 'failed';
}

function messageCanDelete(message) {
  return Boolean(message) && message.status !== 'sending';
}

function messageCanForward(message) {
  return Boolean(message) && message.status !== 'sending';
}

function messageStatusLabel(message) {
  if (!message) {
    return '';
  }
  if (message.status !== 'sending') {
    return t(message.status);
  }

  const progress = transferProgressForMessage(message.id);
  if (message.kind !== 'text' && progress !== null && progress > 0 && progress < 100) {
    return `${t('sending')} ${progress}%`;
  }

  return t('sending');
}

function renderMessageContent(message, searchClass = '') {
  if (message.kind === 'image') {
    const source = cachedImageSource(message);
    const imageKey = imageSourceKey(message);
    return `
      <div class="image-message ${source ? '' : 'image-message-loading'} ${searchClass}"${imageKey ? ` data-image-ref="${escapeHtml(imageKey)}"` : ''}>
        <div class="message-image-placeholder ${source ? 'hidden' : ''}">${escapeHtml(message.mediaName || t('imageAlt'))}</div>
        <img class="message-image ${source ? '' : 'hidden'}" src="${escapeHtml(source)}" alt="${escapeHtml(message.mediaName || t('imageAlt'))}" data-preview-src="${escapeHtml(source)}" data-preview-name="${escapeHtml(message.mediaName || t('imageAlt'))}" />
      </div>
    `;
  }

  if (message.kind === 'file') {
    const fileCategory = fileMessageCategory(message);
    const fileKey = fileSourceKey(message);
    const available = cachedFileAvailability(message);
    const fileMeta = fileMessageMeta(message);
    return `
      <button
        class="file-message ${available === false ? 'is-missing' : ''} ${searchClass}"
        type="button"
        title="${escapeHtml(fileMessageName(message))}"
        data-file-category="${escapeHtml(fileCategory)}"
        data-file-status="${available === false ? 'missing' : 'ready'}"
        data-file-meta="${escapeHtml(fileMeta)}"
        ${fileKey ? `data-file-ref="${escapeHtml(fileKey)}"` : ''}
      >
        <span class="file-message-badge">${escapeHtml(fileMessageBadge(message))}</span>
        <span class="file-message-copy">
          <strong>${escapeHtml(fileMessageName(message))}</strong>
          <span>${escapeHtml(available === false ? t('fileUnavailableStatus') : fileMeta)}</span>
        </span>
      </button>
    `;
  }

  return renderSearchHighlightedText(message.text);
}

function formatTime(timestamp) {
  return new Intl.DateTimeFormat(state.settings.language || 'zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    month: '2-digit',
    day: '2-digit',
  }).format(new Date(timestamp));
}

function clamp(value, min, max) {
  return Math.min(Math.max(value, min), max);
}

function imageViewerTitle() {
  return t('imagePreviewTitle');
}

function activeConversationMessages() {
  if (!state.activePeerId) {
    return [];
  }
  return state.conversations[state.activePeerId] || [];
}

function conversationImageMessages(peerId) {
  return (state.conversations[peerId] || []).filter((message) => message.kind === 'image');
}

function conversationMeta(peerId) {
  if (!peerId) {
    return {
      initialized: false,
      loading: false,
      hasMore: false,
    };
  }

  if (!state.conversationMeta[peerId]) {
    state.conversationMeta[peerId] = {
      initialized: false,
      loading: false,
      hasMore: false,
    };
  }

  return state.conversationMeta[peerId];
}

function findMessage(peerId, messageId) {
  const messages = state.conversations[peerId] || [];
  return messages.find((message) => message.id === messageId) || null;
}

function setConversationMessages(peerId, messages) {
  const previousMessages = state.conversations[peerId] || [];
  state.conversations[peerId] = messages;
  syncTransferProgressCache(previousMessages, messages);
  if (state.activePeerId === peerId) {
    rerenderConversation({ preserveViewport: true });
  }
}

function isConversationNearBottom(threshold = conversationStickToBottomThreshold) {
  const { scrollTop, scrollHeight, clientHeight } = elements.messageList;
  if (!scrollHeight || !clientHeight) {
    return true;
  }
  return (scrollHeight - clientHeight - scrollTop) <= threshold;
}

function captureConversationViewport() {
  if (!state.activePeerId || elements.messageList.classList.contains('empty')) {
    return null;
  }

  return {
    peerId: state.activePeerId,
    scrollTop: elements.messageList.scrollTop,
    scrollHeight: elements.messageList.scrollHeight,
    stickToBottom: isConversationNearBottom(),
  };
}

function restoreConversationViewport(snapshot) {
  if (!snapshot || snapshot.peerId !== state.activePeerId || elements.messageList.classList.contains('empty')) {
    return;
  }

  if (snapshot.stickToBottom) {
    elements.messageList.scrollTop = elements.messageList.scrollHeight;
    return;
  }

  const heightDelta = elements.messageList.scrollHeight - snapshot.scrollHeight;
  elements.messageList.scrollTop = Math.max(0, snapshot.scrollTop + heightDelta);
}

function removeMessageLocally(peerId, messageId) {
  const currentMessages = state.conversations[peerId] || [];
  const nextMessages = currentMessages.filter((message) => message.id !== messageId);
  if (nextMessages.length === currentMessages.length) {
    return null;
  }

  setConversationMessages(peerId, nextMessages);
  return currentMessages;
}

function truncateText(value, maxLength = 120) {
  const text = String(value || '').replace(/\s+/g, ' ').trim();
  if (text.length <= maxLength) {
    return text;
  }
  return `${text.slice(0, Math.max(0, maxLength - 1)).trimEnd()}…`;
}

async function loadFileMessageSource(message) {
  if (!message || message.kind !== 'file') {
    return '';
  }

  try {
    return await LoadFileSource(message.text || '', message.mediaType || '');
  } catch (error) {
    if (isMissingFileError(error)) {
      setFileAvailability(message, false);
    }
    throw error;
  }
}

function canForwardToPeer(peer, sourcePeerId = '') {
  if (!peer || !peer.id || peer.id === sourcePeerId) {
    return false;
  }

  return peerCanSend(peer);
}

function availableForwardPeers() {
  return state.peers.filter((peer) => canForwardToPeer(peer, state.forwardDialog.peerId));
}

function suggestedImageName(message) {
  const baseName = (message?.mediaName || '').trim();
  if (baseName) {
    return baseName;
  }

  const mediaType = (message?.mediaType || '').toLowerCase();
  if (mediaType.includes('jpeg') || mediaType.includes('jpg')) {
    return 'image.jpg';
  }
  if (mediaType.includes('webp')) {
    return 'image.webp';
  }
  if (mediaType.includes('gif')) {
    return 'image.gif';
  }
  if (mediaType.includes('bmp')) {
    return 'image.bmp';
  }
  return 'image.png';
}

function suggestedFileName(message) {
  const baseName = (message?.mediaName || '').trim();
  if (baseName) {
    return baseName;
  }
  return 'file';
}

async function copyImageToClipboard(source) {
  if (!source) {
    return false;
  }

  try {
    await CopyImageToClipboard(source);
    return true;
  } catch (error) {
    console.error('CopyImageToClipboard failed', error);
    return false;
  }
}

function isImageLikeFile(file) {
  const mediaType = String(file?.type || '').trim().toLowerCase();
  if (mediaType.startsWith('image/')) {
    return true;
  }

  const name = String(file?.name || '').trim().toLowerCase();
  return /\.(png|jpe?g|webp|gif|bmp|ico|svg)$/i.test(name);
}

function inferredFileMediaType(file, fallback = '') {
  const explicit = String(file?.type || '').trim().toLowerCase();
  if (explicit) {
    return explicit;
  }

  const name = String(file?.name || '').trim().toLowerCase();
  if (name.endsWith('.png')) return 'image/png';
  if (name.endsWith('.jpg') || name.endsWith('.jpeg')) return 'image/jpeg';
  if (name.endsWith('.webp')) return 'image/webp';
  if (name.endsWith('.gif')) return 'image/gif';
  if (name.endsWith('.bmp')) return 'image/bmp';
  if (name.endsWith('.ico')) return 'image/x-icon';
  if (name.endsWith('.svg')) return 'image/svg+xml';
  if (name.endsWith('.pdf')) return 'application/pdf';
  if (name.endsWith('.txt')) return 'text/plain';
  if (name.endsWith('.json')) return 'application/json';
  if (name.endsWith('.zip')) return 'application/zip';

  return fallback || 'application/octet-stream';
}

function dataURLMediaType(dataURL) {
  const match = /^data:([^;,]*)(?:;[^,]*)?,/i.exec(String(dataURL || '').trim());
  return String(match?.[1] || '').trim().toLowerCase();
}

function replaceDataURLMediaType(dataURL, mediaType) {
  const parts = String(dataURL || '').split(',', 2);
  if (parts.length !== 2 || !mediaType) {
    return dataURL;
  }
  return `data:${mediaType};base64,${parts[1]}`;
}

async function readAttachmentAsDataURL(file) {
  const dataURL = await new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result);
    reader.onerror = () => reject(reader.error || new Error('failed to read file'));
    reader.readAsDataURL(file);
  });

  const currentMediaType = dataURLMediaType(dataURL);
  const nextMediaType = inferredFileMediaType(file, currentMediaType);
  if (!nextMediaType || (currentMediaType && currentMediaType !== 'application/octet-stream')) {
    return dataURL;
  }

  return replaceDataURLMediaType(dataURL, nextMediaType);
}

async function sendAttachmentFile(file, peerId = state.activePeerId) {
  if (!file || !peerId) {
    return;
  }
  const peer = state.peers.find((item) => item.id === peerId);
  if (!peerCanSend(peer)) {
    showToast(t('peerUnavailable'));
    return;
  }

  const isImageFile = isImageLikeFile(file);
  if (isImageFile && file.size > maxImageBytes) {
    showToast(t('imageTooLarge'));
    return;
  }
  if (file.size > maxFileBytes) {
    showToast(t('fileTooLarge'));
    return;
  }

  try {
    const dataURL = await readAttachmentAsDataURL(file);
    await SendFileMessage(peerId, dataURL, file.name || '');
  } catch (error) {
    console.error('SendFileMessage failed', error);
    showToast(userFacingError(error, isImageFile ? 'sendImageFailed' : 'sendFileFailed'));
  }
}

async function sendAttachmentFiles(files, peerId = state.activePeerId) {
  if (!peerId) {
    showToast(t('composerIdle'));
    return;
  }
  const peer = state.peers.find((item) => item.id === peerId);
  if (!peerCanSend(peer)) {
    showToast(t('peerUnavailable'));
    return;
  }

  const items = Array.from(files || []).filter((file) => file && typeof file.size === 'number');
  for (const file of items) {
    // Keep the send order stable when multiple attachments are selected together.
    await sendAttachmentFile(file, peerId);
  }
}

async function sendLocalAttachmentPaths(paths, peerId = state.activePeerId) {
  if (!peerId) {
    showToast(t('composerIdle'));
    return;
  }
  const peer = state.peers.find((item) => item.id === peerId);
  if (!peerCanSend(peer)) {
    showToast(t('peerUnavailable'));
    return;
  }

  const items = Array.from(paths || []).map((path) => String(path || '').trim()).filter(Boolean);
  for (const path of items) {
    try {
      await SendLocalFileMessage(peerId, path);
    } catch (error) {
      console.error('SendLocalFileMessage failed', error);
      showToast(userFacingError(error, 'sendFileFailed'));
    }
  }
}

function peerStatusLabel(peer) {
  const address = (peer.address || '').trim();
  const port = Number(peer.listenPort) || 0;

  if (address && port > 0) {
    return `${address}:${port}`;
  }
  if (address && address !== ':0') {
    return address;
  }
  if (port > 0) {
    return String(port);
  }
  return '';
}

function peerIsOnline(peer) {
  if (!peer || !peer.id) {
    return false;
  }
  if (peer.source === 'debug') {
    return true;
  }

  const address = (peer.address || '').trim();
  const port = Number(peer.listenPort) || 0;
  if (!address || address === ':0' || port <= 0) {
    return false;
  }

  if (peer.source === 'lan') {
    const lastSeen = Number(peer.lastSeen) || 0;
    return lastSeen > 0 && (Date.now() - lastSeen) <= peerOfflineThresholdMs;
  }

  return true;
}

function peerCanSend(peer) {
  return peerIsOnline(peer);
}

function peerSubtitle(peer) {
  const label = peerStatusLabel(peer);
  if (peerIsOnline(peer)) {
    return label || t('onlineNow');
  }
  if (label) {
    return `${t('offline')} · ${label}`;
  }
  return t('offline');
}

function rawErrorMessage(error) {
  if (!error) {
    return '';
  }
  if (typeof error === 'string') {
    return error;
  }
  if (typeof error.message === 'string') {
    return error.message;
  }
  return String(error);
}

function isMissingFileError(error) {
  const raw = rawErrorMessage(error).trim();
  if (!raw) {
    return false;
  }

  return /no such file|cannot find the file|the system cannot find the file|file not found|does not exist|not exist/i.test(raw);
}

function userFacingError(error, fallbackKey) {
  const raw = rawErrorMessage(error).trim();
  if (!raw) {
    return t(fallbackKey);
  }

  if (/peer is offline|actively refused|i\/o timeout|connection refused|network is unreachable|cannot assign requested address/i.test(raw)) {
    return t('peerUnavailable');
  }
  if (/message is still sending|conversation has sending messages/i.test(raw)) {
    return t('messagesStillSending');
  }
  if (isMissingFileError(error)) {
    return t('fileUnavailable');
  }
  if (/image too large/i.test(raw)) {
    return t('imageTooLarge');
  }
  if (/file too large/i.test(raw)) {
    return t('fileTooLarge');
  }
  if (/media storage limit reached/i.test(raw)) {
    return t('mediaStorageLimitReached');
  }

  return t(fallbackKey);
}

function isExternalFileDropEvent(event) {
  const types = Array.from(event?.dataTransfer?.types || []);
  return types.includes('Files');
}

function removeToast(toast) {
  if (!toast || toast.dataset.closing === 'true') {
    return;
  }
  toast.dataset.closing = 'true';
  toast.classList.remove('visible');
  window.setTimeout(() => {
    toast.remove();
  }, 180);
}

function showToast(message, tone = 'error') {
  const text = String(message || '').trim();
  if (!text) {
    return;
  }

  const toast = document.createElement('div');
  toast.className = `toast toast-${tone}`;
  toast.innerHTML = `
    <div class="toast-copy">${escapeHtml(text)}</div>
    <button class="toast-close" type="button" aria-label="${escapeHtml(t('dismiss'))}">×</button>
  `;

  const closeButton = toast.querySelector('.toast-close');
  closeButton?.addEventListener('click', () => removeToast(toast));

  elements.toastRegion.appendChild(toast);
  requestAnimationFrame(() => toast.classList.add('visible'));
  window.setTimeout(() => removeToast(toast), tone === 'success' ? 2400 : 3600);
}

function unreadCountForPeer(peerId) {
  return Number(state.unreadCounts[peerId]) || 0;
}

function totalUnreadCount() {
  return Object.values(state.unreadCounts).reduce((sum, value) => sum + (Number(value) || 0), 0);
}

function isConversationVisible(peerId) {
  return Boolean(peerId) && state.activePeerId === peerId && state.windowVisible && state.windowFocused;
}

function syncWindowTitle() {
  const unreadTotal = totalUnreadCount();
  WindowSetTitle(unreadTotal > 0 ? `LanTalk (${unreadTotal})` : 'LanTalk');
}

function updateWindowFocusedState(focused) {
  state.windowFocused = Boolean(focused) && !document.hidden;
  clearActivePeerUnreadIfVisible();
}

function updateWindowVisibleState(visible) {
  state.windowVisible = Boolean(visible);
  clearActivePeerUnreadIfVisible();
}

function syncUnreadState({ render = true } = {}) {
  const unreadTotal = totalUnreadCount();
  if (render) {
    renderPeers();
  }
  syncWindowTitle();
  if (state.lastUnreadTotal === unreadTotal) {
    return;
  }
  state.lastUnreadTotal = unreadTotal;
  Promise.resolve(SetUnreadCount(unreadTotal)).catch((error) => {
    console.error('SetUnreadCount failed', error);
  });
}

function clearUnread(peerId, options = {}) {
  if (!peerId || !unreadCountForPeer(peerId)) {
    return;
  }
  delete state.unreadCounts[peerId];
  syncUnreadState(options);
}

function incrementUnread(peerId, count = 1) {
  if (!peerId || count <= 0) {
    return;
  }
  state.unreadCounts[peerId] = unreadCountForPeer(peerId) + count;
  syncUnreadState();
}

function clearActivePeerUnreadIfVisible() {
  if (!isConversationVisible(state.activePeerId)) {
    return;
  }
  clearUnread(state.activePeerId, { render: true });
}

function pruneUnreadCounts(validPeerIDs) {
  let changed = false;
  Object.keys(state.unreadCounts).forEach((peerId) => {
    if (validPeerIDs.has(peerId)) {
      return;
    }
    delete state.unreadCounts[peerId];
    changed = true;
  });
  return changed;
}

function collectNewInboundMessages(previousMessages, nextMessages) {
  const seen = new Set((previousMessages || []).map((message) => message.id));
  return (nextMessages || []).filter((message) => message.direction === 'inbound' && !seen.has(message.id));
}

function messageNotificationSummary(message) {
  if (!message) {
    return '';
  }
  if (message.kind === 'image') {
    return t('imageMessageSummary');
  }
  if (message.kind === 'file') {
    return fileMessageName(message);
  }
  return truncateText(message.text || '', 120);
}

async function initializeDesktopNotifications() {
  try {
    const available = await IsNotificationAvailable();
    if (!available) {
      return;
    }
    await InitializeNotifications();
    const authorised = await RequestNotificationAuthorization();
    state.notifications.available = Boolean(authorised);
  } catch (error) {
    console.error('InitializeNotifications failed', error);
  }
}

async function notifyIncomingMessages(peerId, messages) {
  if (!state.notifications.available || !messages.length) {
    return;
  }

  const latestMessage = messages[messages.length - 1];
  const peer = state.peers.find((item) => item.id === peerId);

  try {
    await SendNotification({
      id: `message-${latestMessage.id}`,
      title: peer?.name || 'LanTalk',
      body: messageNotificationSummary(latestMessage),
      data: {
        peerId,
        messageId: latestMessage.id,
      },
    });
  } catch (error) {
    console.error('SendNotification failed', error);
  }
}

function activeSearchQuery() {
  if (!state.search.open) {
    return '';
  }
  return String(state.search.query || '').trim();
}

function searchableMessageText(message) {
  if (!message) {
    return '';
  }
  if (message.kind === 'file' || message.kind === 'image') {
    return `${message.mediaName || ''} ${message.mediaType || ''}`.trim();
  }
  return String(message.text || '');
}

function messageMatchesSearch(message, query = activeSearchQuery()) {
  const normalizedQuery = String(query || '').trim().toLowerCase();
  if (!normalizedQuery || !message) {
    return false;
  }
  return searchableMessageText(message).toLowerCase().includes(normalizedQuery);
}

function currentSearchMessageID() {
  return state.search.results[state.search.activeIndex] || '';
}

function renderSearchHighlightedText(value) {
  const source = String(value ?? '');
  const query = activeSearchQuery().toLowerCase();
  if (!query) {
    return escapeHtmlWithLineBreaks(source);
  }

  const lowerSource = source.toLowerCase();
  let cursor = 0;
  let html = '';

  while (cursor < source.length) {
    const matchIndex = lowerSource.indexOf(query, cursor);
    if (matchIndex < 0) {
      html += escapeHtml(source.slice(cursor));
      break;
    }

    html += escapeHtml(source.slice(cursor, matchIndex));
    html += `<mark class="message-search-mark">${escapeHtml(source.slice(matchIndex, matchIndex + query.length))}</mark>`;
    cursor = matchIndex + query.length;
  }

  return html.split('\n').join('<br>');
}

function syncSearchUI() {
  const isOpen = state.search.open;
  const hasPeer = Boolean(state.activePeerId);
  const hasResults = state.search.results.length > 0;

  elements.searchToggle.textContent = t('findConversationContent');
  elements.searchToggle.disabled = !hasPeer;
  elements.chatSearchBar.classList.toggle('hidden', !isOpen);
  elements.chatSearchBar.classList.toggle('is-busy', state.search.loading);
  elements.chatSearchInput.placeholder = t('searchPlaceholder');
  elements.chatSearchInput.value = state.search.query;
  elements.chatSearchPrev.textContent = t('searchPrev');
  elements.chatSearchNext.textContent = t('searchNext');
  elements.chatSearchClose.textContent = t('searchClose');
  elements.chatSearchPrev.disabled = !isOpen || state.search.loading || state.search.results.length <= 1;
  elements.chatSearchNext.disabled = !isOpen || state.search.loading || state.search.results.length <= 1;

  if (!isOpen || !hasPeer) {
    elements.chatSearchMeta.textContent = '';
    return;
  }
  if (state.search.loading) {
    elements.chatSearchMeta.textContent = t('searchLoading');
    return;
  }
  if (!activeSearchQuery()) {
    elements.chatSearchMeta.textContent = '';
    return;
  }
  if (!hasResults) {
    elements.chatSearchMeta.textContent = t('searchNoResult');
    return;
  }

  elements.chatSearchMeta.textContent = t('searchCount', {
    current: String(state.search.activeIndex + 1),
    total: String(state.search.results.length),
  });
}

function scrollToActiveSearchResult() {
  const messageId = currentSearchMessageID();
  if (!messageId) {
    return;
  }

  requestAnimationFrame(() => {
    const selector = `.message[data-message-id="${CSS.escape(messageId)}"]`;
    const target = elements.messageList.querySelector(selector);
    target?.scrollIntoView({
      block: 'center',
      behavior: 'smooth',
    });
  });
}

function setSearchResultIndex(index) {
  const total = state.search.results.length;
  if (!total) {
    return;
  }

  const normalizedIndex = ((index % total) + total) % total;
  if (state.search.activeIndex === normalizedIndex) {
    scrollToActiveSearchResult();
    return;
  }

  state.search.activeIndex = normalizedIndex;
  renderConversation();
  syncSearchUI();
  scrollToActiveSearchResult();
}

function closeConversationSearch() {
  state.search.open = false;
  state.search.query = '';
  state.search.results = [];
  state.search.activeIndex = -1;
  state.search.loading = false;
  state.search.peerId = '';
  state.search.token += 1;
  rerenderConversation({ preserveViewport: true });
  syncSearchUI();
}

function openConversationSearch() {
  if (!state.activePeerId) {
    return;
  }

  state.search.open = true;
  state.search.peerId = state.activePeerId;
  syncSearchUI();
  requestAnimationFrame(() => {
    elements.chatSearchInput.focus();
    elements.chatSearchInput.select();
  });
  if (activeSearchQuery()) {
    void refreshConversationSearch();
  }
}

async function ensureConversationSearchReady(peerId, token) {
  while (token === state.search.token) {
    const meta = conversationMeta(peerId);
    if (!meta.hasMore) {
      break;
    }
    if (meta.loading) {
      await new Promise((resolve) => window.setTimeout(resolve, 48));
      continue;
    }
    await loadConversationHistory(peerId);
  }
}

async function refreshConversationSearch({ ensureLoaded = true, preserveSelection = true } = {}) {
  if (!state.search.open) {
    syncSearchUI();
    return;
  }

  const peerId = state.activePeerId;
  const query = activeSearchQuery();
  const token = state.search.token + 1;
  state.search.token = token;
  state.search.peerId = peerId || '';

  if (!peerId || !query) {
    state.search.loading = false;
    state.search.results = [];
    state.search.activeIndex = -1;
    rerenderConversation({ preserveViewport: true });
    syncSearchUI();
    return;
  }

  const preservedMessageID = preserveSelection ? currentSearchMessageID() : '';
  state.search.loading = ensureLoaded;
  syncSearchUI();

  if (ensureLoaded) {
    await ensureConversationSearchReady(peerId, token);
  }

  if (token !== state.search.token || !state.search.open || state.activePeerId !== peerId) {
    return;
  }

  const results = (state.conversations[peerId] || [])
    .filter((message) => messageMatchesSearch(message, query))
    .map((message) => message.id);

  state.search.loading = false;
  state.search.results = results;
  if (!results.length) {
    state.search.activeIndex = -1;
  } else if (preservedMessageID && results.includes(preservedMessageID)) {
    state.search.activeIndex = results.indexOf(preservedMessageID);
  } else {
    state.search.activeIndex = results.length - 1;
  }

  rerenderConversation({ preserveViewport: true });
  syncSearchUI();
  scrollToActiveSearchResult();
}

function setActivePeer(peerId) {
  state.activePeerId = peerId;
  closeChatMenu();
  clearUnread(peerId, { render: false });
  renderPeers();
  renderConversation();
  ensureConversationLoaded(peerId);
  if (state.search.open) {
    void refreshConversationSearch({ preserveSelection: false });
  } else {
    syncSearchUI();
  }
}

function activePeer() {
  return state.peers.find((peer) => peer.id === state.activePeerId) || null;
}

async function loadConversationHistory(peerId, options = {}) {
  const { reset = false } = options;
  if (!peerId) {
    return;
  }

  const meta = conversationMeta(peerId);
  if (meta.loading) {
    return;
  }

  const currentMessages = state.conversations[peerId] || [];
  const previousMessages = currentMessages;
  const oldestMessage = currentMessages[0] || null;
  const beforeTimestamp = !reset && oldestMessage ? Number(oldestMessage.timestamp) || 0 : 0;
  const beforeID = !reset && oldestMessage ? oldestMessage.id || '' : '';
  const shouldPreserveScroll = !reset && state.activePeerId === peerId;
  const previousScrollHeight = shouldPreserveScroll ? elements.messageList.scrollHeight : 0;
  const previousScrollTop = shouldPreserveScroll ? elements.messageList.scrollTop : 0;

  meta.loading = true;
  state.preserveConversationScroll = shouldPreserveScroll;
  renderConversation();
  if (shouldPreserveScroll) {
    elements.messageList.scrollTop = previousScrollTop;
  }

  try {
    const page = await LoadConversationPage(peerId, beforeTimestamp, beforeID, conversationPageSize);
    state.conversations[peerId] = page.messages || [];
    syncTransferProgressCache(previousMessages, state.conversations[peerId]);
    meta.initialized = true;
    meta.hasMore = Boolean(page.hasMore);
  } catch (error) {
    console.error('LoadConversationPage failed', error);
    showToast(t('loadHistoryFailed'));
  } finally {
    meta.loading = false;
    renderConversation();
    if (shouldPreserveScroll) {
      const nextScrollHeight = elements.messageList.scrollHeight;
      elements.messageList.scrollTop = Math.max(0, nextScrollHeight - previousScrollHeight + previousScrollTop);
    }
    state.preserveConversationScroll = false;
  }
}

function ensureConversationLoaded(peerId) {
  const meta = conversationMeta(peerId);
  if (meta.initialized || meta.loading) {
    return;
  }
  void loadConversationHistory(peerId, { reset: true });
}

function maybeAutoloadEarlierMessages(peerId) {
  if (!peerId || state.activePeerId !== peerId) {
    return;
  }

  const meta = conversationMeta(peerId);
  if (!meta.hasMore || meta.loading) {
    return;
  }

  if (elements.messageList.scrollHeight <= elements.messageList.clientHeight + 12) {
    void loadConversationHistory(peerId);
  }
}

function renderStaticText() {
  document.documentElement.dataset.theme = state.settings.theme || 'midnight';
  elements.openSettings.textContent = t('settings');
  elements.chatMenuTrigger.textContent = '...';
  elements.chatMenuTrigger.title = t('moreActions');
  elements.chatMenuTrigger.setAttribute('aria-label', t('moreActions'));
  elements.profileSettingsTitle.textContent = t('profileSettings');
  elements.nicknameLabel.textContent = t('nicknameLabel');
  elements.profileIdLabel.textContent = t('profileIdLabel');
  elements.saveName.textContent = t('save');
  elements.settingsDialogTitle.textContent = t('settings');
  elements.settingsDialogClose.textContent = t('dismiss');
  elements.settingsTitle.textContent = t('appearanceSettings');
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
  elements.onlineTitle.textContent = t('contacts');
  elements.storageTitle.textContent = t('storage');
  elements.storagePathLabel.textContent = t('storagePathLabel');
  elements.changeDataPath.textContent = t('change');
  elements.cleanupMediaLabel.textContent = t('cleanupMediaLabel');
  elements.cleanupMedia.textContent = t('cleanupMedia');
  elements.storageHint.textContent = t('storageHint');
  elements.clearConversation.textContent = t('clearConversationRecords');
  renderStorageStats();
  elements.searchToggle.textContent = t('findConversationContent');
  elements.sendButton.textContent = t('send');
  elements.fileTrigger.textContent = t('file');
  elements.emojiTrigger.textContent = t('emoji');
  elements.imageViewerTitle.textContent = imageViewerTitle();
  elements.imageViewerCounter.textContent = '';
  elements.imageViewerPrev.setAttribute('aria-label', t('previousImage'));
  elements.imageViewerNext.setAttribute('aria-label', t('nextImage'));
  elements.imageViewerClose.textContent = t('closePreview');
  elements.imageViewerHint.textContent = t('imageViewerHint');
  elements.imageViewerZoomOut.textContent = t('zoomOut');
  elements.imageViewerZoomIn.textContent = t('zoomIn');
  elements.imageViewerReset.textContent = t('zoomReset');
  elements.messageContextSave.textContent = t('saveAs');
  elements.messageContextReveal.textContent = t('openFolder');
  elements.messageContextRetry.textContent = t('retryMessage');
  elements.messageContextForward.textContent = t('forward');
  elements.deleteConfirmTitle.textContent = t('deleteConfirmTitle');
  elements.deleteConfirmBody.textContent = t('deleteConfirmBody');
  elements.deleteConfirmCancel.textContent = t('cancel');
  elements.deleteConfirmSubmit.textContent = t('confirmDelete');
  elements.actionConfirmCancel.textContent = t('cancel');
  syncActionDialogText();
  elements.forwardDialogTitle.textContent = t('forwardDialogTitle');
  elements.forwardDialogBody.textContent = t('forwardDialogBody');
  elements.forwardPreviewLabel.textContent = t('forwardPreviewLabel');
  elements.forwardTargetLabel.textContent = t('forwardTargetLabel');
  elements.forwardDialogEmpty.textContent = t('forwardEmpty');
  elements.forwardDialogCancel.textContent = t('cancel');
  syncForwardSubmitButtonText();
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

  renderSelfProfile();

  if (!elements.dataPath.dataset.loaded) {
    elements.dataPath.textContent = t('loadingPath');
  }

  syncSearchUI();
}

function updateImageViewerScaleLabel() {
  elements.imageViewerScale.textContent = `${Math.round(state.viewer.scale * 100)}%`;
}

function syncImageViewerNavigation() {
  const images = conversationImageMessages(state.viewer.peerId);
  const currentIndex = images.findIndex((message) => message.id === state.viewer.messageId);
  const hasSequence = images.length > 1 && currentIndex >= 0;

  elements.imageViewerCounter.classList.toggle('hidden', !hasSequence);
  elements.imageViewerPrev.classList.toggle('hidden', !hasSequence);
  elements.imageViewerNext.classList.toggle('hidden', !hasSequence);

  if (!hasSequence) {
    elements.imageViewerCounter.textContent = '';
    elements.imageViewerPrev.disabled = true;
    elements.imageViewerNext.disabled = true;
    return;
  }

  elements.imageViewerCounter.textContent = `${currentIndex + 1} / ${images.length}`;
  elements.imageViewerPrev.disabled = currentIndex <= 0 || state.viewer.loading;
  elements.imageViewerNext.disabled = currentIndex >= images.length - 1 || state.viewer.loading;
}

function syncForwardSubmitButtonText() {
  elements.forwardDialogSubmit.textContent = t('confirmForward');
}

function closeDeleteDialog() {
  state.deleteDialog.open = false;
  state.deleteDialog.peerId = '';
  state.deleteDialog.messageId = '';
  elements.deleteConfirmDialog.classList.add('hidden');
  elements.deleteConfirmDialog.setAttribute('aria-hidden', 'true');
}

function closeSettingsDialog() {
  state.settingsDialog.open = false;
  closeChatMenu();
  closeLanguageMenu();
  closeThemeMenu();
  elements.settingsDialog.classList.add('hidden');
  elements.settingsDialog.setAttribute('aria-hidden', 'true');
}

function openSettingsDialog() {
  state.settingsDialog.open = true;
  closeChatMenu();
  closeEmojiPanel();
  elements.settingsDialog.classList.remove('hidden');
  elements.settingsDialog.setAttribute('aria-hidden', 'false');
}

function closeChatMenu() {
  state.chatMenu.open = false;
  elements.chatMenuPanel.classList.add('hidden');
  elements.chatMenuTrigger.setAttribute('aria-expanded', 'false');
}

function openChatMenu() {
  if (elements.chatMenuTrigger.disabled) {
    return;
  }
  state.chatMenu.open = true;
  closeEmojiPanel();
  elements.chatMenuPanel.classList.remove('hidden');
  elements.chatMenuTrigger.setAttribute('aria-expanded', 'true');
}

function openDeleteDialog(peerId, messageId) {
  state.deleteDialog.open = true;
  state.deleteDialog.peerId = peerId;
  state.deleteDialog.messageId = messageId;
  elements.deleteConfirmDialog.classList.remove('hidden');
  elements.deleteConfirmDialog.setAttribute('aria-hidden', 'false');
}

function actionDialogConfig() {
  switch (state.actionDialog.mode) {
    case 'clearConversation':
      return {
        titleKey: 'clearConversationTitle',
        bodyKey: 'clearConversationBody',
        submitKey: 'confirmClearConversation',
        danger: true,
      };
    case 'cleanupMedia':
      return {
        titleKey: 'cleanupMediaTitle',
        bodyKey: 'cleanupMediaBody',
        submitKey: 'confirmCleanupMedia',
        danger: false,
      };
    default:
      return null;
  }
}

function syncActionDialogText() {
  const config = actionDialogConfig();
  if (!config) {
    elements.actionConfirmTitle.textContent = '';
    elements.actionConfirmBody.textContent = '';
    elements.actionConfirmSubmit.textContent = '';
    elements.actionConfirmSubmit.disabled = true;
    elements.actionConfirmCancel.disabled = false;
    elements.actionConfirmSubmit.classList.remove('danger-btn');
    return;
  }

  const peer = state.peers.find((item) => item.id === state.actionDialog.peerId);
  elements.actionConfirmTitle.textContent = t(config.titleKey);
  elements.actionConfirmBody.textContent = t(config.bodyKey, {
    name: peer?.name || t('noConversation'),
  });
  elements.actionConfirmSubmit.textContent = t(config.submitKey);
  elements.actionConfirmSubmit.disabled = state.actionDialog.busy;
  elements.actionConfirmCancel.disabled = state.actionDialog.busy;
  elements.actionConfirmSubmit.classList.toggle('danger-btn', config.danger);
}

function closeActionDialog() {
  if (state.actionDialog.busy) {
    return;
  }

  state.actionDialog.open = false;
  state.actionDialog.mode = '';
  state.actionDialog.peerId = '';
  state.actionDialog.busy = false;
  syncActionDialogText();
  elements.actionConfirmDialog.classList.add('hidden');
  elements.actionConfirmDialog.setAttribute('aria-hidden', 'true');
}

function openActionDialog(mode, peerId = '') {
  state.actionDialog.open = true;
  state.actionDialog.mode = mode;
  state.actionDialog.peerId = peerId;
  state.actionDialog.busy = false;
  syncActionDialogText();
  elements.actionConfirmDialog.classList.remove('hidden');
  elements.actionConfirmDialog.setAttribute('aria-hidden', 'false');
}

function closeForwardDialog() {
  state.forwardDialog.open = false;
  state.forwardDialog.peerId = '';
  state.forwardDialog.messageId = '';
  state.forwardDialog.targetPeerId = '';
  elements.forwardPreview.innerHTML = '';
  elements.forwardDialogSubmit.dataset.busy = 'false';
  syncForwardSubmitButtonText();
  elements.forwardDialog.classList.add('hidden');
  elements.forwardDialog.setAttribute('aria-hidden', 'true');
}

function renderForwardPreview() {
  const message = findMessage(state.forwardDialog.peerId, state.forwardDialog.messageId);
  if (!message) {
    elements.forwardPreview.innerHTML = '';
    return;
  }

  if (message.kind === 'file') {
    elements.forwardPreview.innerHTML = `
      <div class="forward-preview-thumb forward-preview-thumb-file">${escapeHtml(fileMessageBadge(message))}</div>
      <div class="forward-preview-copy">
        <strong>${escapeHtml(fileMessageName(message))}</strong>
        <span>${escapeHtml(fileMessageMeta(message))}</span>
      </div>
    `;
    return;
  }

  if (message.kind === 'image') {
    const source = cachedImageSource(message);
    if (!source) {
      resolveImageMessageSource(message).then(() => {
        if (
          state.forwardDialog.open &&
          state.forwardDialog.peerId === message.peerId &&
          state.forwardDialog.messageId === message.id
        ) {
          renderForwardPreview();
        }
      }).catch((error) => {
        console.error('Load forward preview image failed', error);
      });
    }
    elements.forwardPreview.innerHTML = `
      <div class="forward-preview-thumb">
        ${source
          ? `<img src="${escapeHtml(source)}" alt="${escapeHtml(message.mediaName || t('imageAlt'))}" />`
          : `<div class="forward-preview-placeholder">${escapeHtml(message.mediaName || t('imageAlt'))}</div>`}
      </div>
      <div class="forward-preview-copy">
        <strong>${escapeHtml(message.mediaName || t('imageAlt'))}</strong>
        <span>${escapeHtml(message.mediaType || t('imageMessageSummary'))}</span>
      </div>
    `;
    return;
  }

  const previewText = truncateText(message.text || '');
  elements.forwardPreview.innerHTML = `
    <div class="forward-preview-copy forward-preview-copy-text">
      <strong>${escapeHtml(message.senderName || '')}</strong>
      <p>${escapeHtml(previewText || '...')}</p>
    </div>
  `;
}

function renderForwardPeerList() {
  const peers = availableForwardPeers();
  if (!peers.length) {
    elements.forwardPeerList.innerHTML = '';
    elements.forwardDialogEmpty.classList.remove('hidden');
    elements.forwardDialogSubmit.disabled = true;
    return;
  }

  if (!state.forwardDialog.targetPeerId || !peers.some((peer) => peer.id === state.forwardDialog.targetPeerId)) {
    const preferredPeer = peers.find((peer) => peer.id !== state.forwardDialog.peerId) || peers[0];
    state.forwardDialog.targetPeerId = preferredPeer?.id || '';
  }

  elements.forwardDialogEmpty.classList.add('hidden');
  elements.forwardDialogSubmit.disabled = !state.forwardDialog.targetPeerId;
  elements.forwardPeerList.innerHTML = peers.map((peer) => `
    <button
      class="forward-peer-option ${peer.id === state.forwardDialog.targetPeerId ? 'active' : ''}"
      type="button"
      data-forward-peer-id="${escapeHtml(peer.id)}"
    >
      <div class="forward-peer-avatar">${escapeHtml(peer.name.slice(0, 1).toUpperCase())}</div>
      <div class="forward-peer-copy">
        <strong>${escapeHtml(peer.name)}</strong>
        <span>${escapeHtml(peerStatusLabel(peer) || t('onlineNow'))}</span>
      </div>
    </button>
  `).join('');

  elements.forwardPeerList.querySelectorAll('[data-forward-peer-id]').forEach((node) => {
    node.addEventListener('click', () => {
      state.forwardDialog.targetPeerId = node.dataset.forwardPeerId;
      renderForwardPeerList();
    });
  });
}

function openForwardDialog(peerId, messageId) {
  state.forwardDialog.open = true;
  state.forwardDialog.peerId = peerId;
  state.forwardDialog.messageId = messageId;
  state.forwardDialog.targetPeerId = '';
  elements.forwardDialogSubmit.dataset.busy = 'false';
  syncForwardSubmitButtonText();
  renderForwardPreview();
  renderForwardPeerList();
  elements.forwardDialog.classList.remove('hidden');
  elements.forwardDialog.setAttribute('aria-hidden', 'false');
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
    const isOnline = peerIsOnline(peer);
    const unread = unreadCountForPeer(peer.id);
    const tags = [];
    if (peer.source === 'debug') {
      tags.push(`<span class="peer-tag">${escapeHtml(t('debugTag'))}</span>`);
    } else if (peer.source === 'manual') {
      tags.push(`<span class="peer-tag">${escapeHtml(t('manualTag'))}</span>`);
    }
    if (!isOnline) {
      tags.push(`<span class="peer-tag peer-tag-offline">${escapeHtml(t('offline'))}</span>`);
    }
    const unreadBadge = unread > 0
      ? `<span class="peer-unread">${escapeHtml(unread > 99 ? '99+' : String(unread))}</span>`
      : '';
    return `
      <button class="peer-card ${isActive ? 'active' : ''} ${isOnline ? '' : 'offline'}" data-peer-id="${peer.id}">
        <div class="peer-avatar">${escapeHtml(peer.name.slice(0, 1).toUpperCase())}</div>
        <div class="peer-copy">
          <strong><span class="peer-name">${escapeHtml(peer.name)}</span>${tags.join('')}${unreadBadge}</strong>
          <span>${escapeHtml(peerSubtitle(peer))}</span>
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
  const meta = peer ? conversationMeta(peer.id) : null;
  const activeSearchMessage = currentSearchMessageID();

  if (!peer) {
    elements.activePeer.textContent = t('noConversation');
    elements.activeStatus.textContent = t('choosePeer');
    elements.chatMenuTrigger.disabled = true;
    closeChatMenu();
    elements.clearConversation.disabled = true;
    elements.messageList.className = 'message-list empty';
    elements.messageList.innerHTML = `
      <div class="empty-state">
        <h3>${escapeHtml(t('waitingTitle'))}</h3>
        <p>${escapeHtml(t('waitingBody'))}</p>
      </div>
    `;
    elements.messageInput.disabled = true;
    elements.sendButton.disabled = true;
    elements.fileInput.disabled = true;
    elements.fileTrigger.classList.add('disabled');
    elements.emojiTrigger.disabled = true;
    elements.composerHint.textContent = t('composerIdle');
    closeEmojiPanel();
    return;
  }

  elements.activePeer.textContent = peer.name;
  elements.activeStatus.textContent = peerSubtitle(peer);
  elements.chatMenuTrigger.disabled = false;
  elements.clearConversation.disabled = !messages.length;
  const canSendToPeer = peerCanSend(peer);
  elements.messageInput.disabled = !canSendToPeer;
  elements.sendButton.disabled = !canSendToPeer;
  elements.fileInput.disabled = !canSendToPeer;
  elements.fileTrigger.classList.toggle('disabled', !canSendToPeer);
  elements.emojiTrigger.disabled = !canSendToPeer;
  elements.composerHint.textContent = canSendToPeer ? t('composerReady') : t('composerPeerOffline');
  if (!canSendToPeer) {
    closeEmojiPanel();
  }

  if (meta && meta.loading && !meta.initialized && !messages.length) {
    elements.messageList.className = 'message-list empty';
    elements.messageList.innerHTML = `
      <div class="empty-state">
        <h3>${escapeHtml(t('loadingConversation'))}</h3>
      </div>
    `;
    return;
  }

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
  const historyControl = meta && (meta.hasMore || meta.loading)
    ? `
      <div class="message-history-strip">
        <span class="message-history-label">${escapeHtml(meta.loading ? t('loadingOlderMessages') : t('loadOlderMessages'))}</span>
      </div>
    `
    : '';
  elements.messageList.innerHTML = historyControl + messages.map((message) => `
    <article class="message ${message.direction} ${messageMatchesSearch(message) ? 'message-search-hit' : ''}" data-message-id="${escapeHtml(message.id)}" data-message-kind="${escapeHtml(message.kind)}">
      <div class="message-meta">
        <span>${escapeHtml(message.senderName)}</span>
        <span>${formatTime(message.timestamp)}</span>
        <span class="status status-${message.status}">${escapeHtml(messageStatusLabel(message))}</span>
      </div>
      <div class="message-bubble ${message.kind === 'image' ? 'message-bubble-image' : message.kind === 'file' ? 'message-bubble-file' : ''} ${message.kind === 'text' ? (message.id === activeSearchMessage ? 'message-bubble-search-active' : messageMatchesSearch(message) ? 'message-bubble-search-hit' : '') : ''}">${renderMessageContent(message, message.kind === 'text' ? '' : message.id === activeSearchMessage ? 'message-attachment-search-active' : messageMatchesSearch(message) ? 'message-attachment-search-hit' : '')}</div>
    </article>
  `).join('');

  if (!state.preserveConversationScroll && !(meta && meta.loading)) {
    elements.messageList.scrollTop = elements.messageList.scrollHeight;
  }
  warmConversationImages(messages);
  warmConversationFiles(messages);
  maybeAutoloadEarlierMessages(peer.id);
}

function rerenderConversation(options = {}) {
  const snapshot = options.preserveViewport ? captureConversationViewport() : null;
  renderConversation();
  if (snapshot) {
    restoreConversationViewport(snapshot);
  }
}

function closeImageViewer() {
  state.viewer.open = false;
  state.viewer.peerId = '';
  state.viewer.messageId = '';
  state.viewer.name = '';
  state.viewer.loading = false;
  state.viewer.scale = viewerMinScale;
  state.viewer.naturalWidth = 0;
  state.viewer.naturalHeight = 0;
  state.viewer.baseWidth = 0;
  state.viewer.baseHeight = 0;
  elements.imageViewer.classList.add('hidden');
  elements.imageViewer.setAttribute('aria-hidden', 'true');
  elements.imageViewerImg.classList.remove('zoomed');
  elements.imageViewerImg.removeAttribute('src');
  elements.imageViewerImg.style.removeProperty('width');
  elements.imageViewerImg.style.removeProperty('height');
  elements.imageViewerStage.style.removeProperty('width');
  elements.imageViewerStage.style.removeProperty('height');
  elements.imageViewerFrame.classList.remove('is-pannable');
  elements.imageViewerFrame.scrollLeft = 0;
  elements.imageViewerFrame.scrollTop = 0;
  updateImageViewerScaleLabel();
  elements.imageViewerTitle.textContent = imageViewerTitle();
  elements.imageViewerCounter.textContent = '';
  elements.imageViewerCounter.classList.add('hidden');
  elements.imageViewerPrev.classList.add('hidden');
  elements.imageViewerNext.classList.add('hidden');
}

function presentImageViewerMessage(peerId, messageId, src, name) {
  state.viewer.open = true;
  state.viewer.peerId = peerId;
  state.viewer.messageId = messageId;
  state.viewer.src = src;
  state.viewer.name = name || '';
  state.viewer.scale = viewerMinScale;
  state.viewer.naturalWidth = 0;
  state.viewer.naturalHeight = 0;
  state.viewer.baseWidth = 0;
  state.viewer.baseHeight = 0;
  elements.imageViewer.classList.remove('hidden');
  elements.imageViewer.setAttribute('aria-hidden', 'false');
  elements.imageViewerTitle.textContent = imageViewerTitle();
  elements.imageViewerImg.alt = name || t('imageAlt');
  elements.imageViewerImg.classList.remove('zoomed');
  elements.imageViewerImg.style.removeProperty('width');
  elements.imageViewerImg.style.removeProperty('height');
  elements.imageViewerStage.style.removeProperty('width');
  elements.imageViewerStage.style.removeProperty('height');
  elements.imageViewerFrame.classList.remove('is-pannable');
  elements.imageViewerFrame.scrollLeft = 0;
  elements.imageViewerFrame.scrollTop = 0;
  updateImageViewerScaleLabel();
  syncImageViewerNavigation();
  elements.imageViewerImg.src = src;

  if (elements.imageViewerImg.complete && elements.imageViewerImg.naturalWidth > 0) {
    syncImageViewerBaseSize();
  }
}

async function openStoredFileMessage(message) {
  if (!message || message.kind !== 'file' || !message.text) {
    return;
  }
  if (cachedFileAvailability(message) === false) {
    showToast(t('fileUnavailable'));
    return;
  }

  try {
    await OpenFileMessage(message.text);
    setFileAvailability(message, true);
  } catch (error) {
    console.error('OpenFileMessage failed', error);
    if (isMissingFileError(error)) {
      setFileAvailability(message, false);
    }
    showToast(userFacingError(error, 'openFileFailed'));
  }
}

async function openImageViewerForMessage(peerId, messageId, preferredSource = '') {
  const message = findMessage(peerId, messageId);
  if (!message || message.kind !== 'image') {
    return;
  }

  state.viewer.loading = true;
  syncImageViewerNavigation();

  try {
    const source = preferredSource || cachedImageSource(message) || await resolveImageMessageSource(message);
    if (!source) {
      showToast(t('loadImageFailed'));
      return;
    }
    presentImageViewerMessage(peerId, messageId, source, message.mediaName || t('imageAlt'));
  } catch (error) {
    console.error('Open image viewer failed', error);
    showToast(t('loadImageFailed'));
  } finally {
    state.viewer.loading = false;
    syncImageViewerNavigation();
  }
}

async function stepImageViewer(direction) {
  if (!state.viewer.open || !state.viewer.peerId || !state.viewer.messageId || state.viewer.loading) {
    return;
  }

  const images = conversationImageMessages(state.viewer.peerId);
  const currentIndex = images.findIndex((message) => message.id === state.viewer.messageId);
  if (currentIndex < 0) {
    return;
  }

  const targetMessage = images[currentIndex + direction];
  if (!targetMessage) {
    return;
  }

  await openImageViewerForMessage(state.viewer.peerId, targetMessage.id);
}

function renderImageViewerLayout() {
  if (!state.viewer.baseWidth || !state.viewer.baseHeight) {
    return;
  }

  const frame = elements.imageViewerFrame;
  const isBaseScale = Math.abs(state.viewer.scale - viewerMinScale) < 0.001;
  const renderedWidth = state.viewer.baseWidth * state.viewer.scale;
  const renderedHeight = state.viewer.baseHeight * state.viewer.scale;
  if (isBaseScale) {
    elements.imageViewerFrame.classList.remove('is-pannable');
    elements.imageViewerStage.style.width = '100%';
    elements.imageViewerStage.style.height = '100%';
    elements.imageViewerImg.style.width = `${Math.round(state.viewer.baseWidth)}px`;
    elements.imageViewerImg.style.height = `${Math.round(state.viewer.baseHeight)}px`;
    elements.imageViewerFrame.scrollLeft = 0;
    elements.imageViewerFrame.scrollTop = 0;
  } else {
    const stageWidth = Math.max(frame.clientWidth, renderedWidth);
    const stageHeight = Math.max(frame.clientHeight, renderedHeight);
    elements.imageViewerFrame.classList.add('is-pannable');
    elements.imageViewerStage.style.width = `${Math.ceil(stageWidth)}px`;
    elements.imageViewerStage.style.height = `${Math.ceil(stageHeight)}px`;
    elements.imageViewerImg.style.width = `${Math.ceil(renderedWidth)}px`;
    elements.imageViewerImg.style.height = `${Math.ceil(renderedHeight)}px`;
  }
  elements.imageViewerImg.classList.toggle('zoomed', state.viewer.scale > viewerMinScale);
  updateImageViewerScaleLabel();
}

function captureImageViewerAnchor(anchor) {
  if (!state.viewer.baseWidth || !state.viewer.baseHeight) {
    return null;
  }

  const frame = elements.imageViewerFrame;
  const viewportX = anchor ? anchor.clientX - frame.getBoundingClientRect().left : frame.clientWidth / 2;
  const viewportY = anchor ? anchor.clientY - frame.getBoundingClientRect().top : frame.clientHeight / 2;
  const renderedWidth = state.viewer.baseWidth * state.viewer.scale;
  const renderedHeight = state.viewer.baseHeight * state.viewer.scale;
  const stageWidth = Math.max(frame.clientWidth, renderedWidth);
  const stageHeight = Math.max(frame.clientHeight, renderedHeight);
  const imageLeft = Math.max((stageWidth - renderedWidth) / 2, 0);
  const imageTop = Math.max((stageHeight - renderedHeight) / 2, 0);
  const stageX = frame.scrollLeft + viewportX;
  const stageY = frame.scrollTop + viewportY;

  return {
    viewportX,
    viewportY,
    imageX: clamp((stageX - imageLeft) / state.viewer.scale, 0, state.viewer.baseWidth),
    imageY: clamp((stageY - imageTop) / state.viewer.scale, 0, state.viewer.baseHeight),
  };
}

function restoreImageViewerAnchor(anchor) {
  if (!anchor || !state.viewer.baseWidth || !state.viewer.baseHeight) {
    return;
  }

  const frame = elements.imageViewerFrame;
  const renderedWidth = state.viewer.baseWidth * state.viewer.scale;
  const renderedHeight = state.viewer.baseHeight * state.viewer.scale;
  const stageWidth = Math.max(frame.clientWidth, renderedWidth);
  const stageHeight = Math.max(frame.clientHeight, renderedHeight);
  const imageLeft = Math.max((stageWidth - renderedWidth) / 2, 0);
  const imageTop = Math.max((stageHeight - renderedHeight) / 2, 0);
  const maxScrollLeft = Math.max(stageWidth - frame.clientWidth, 0);
  const maxScrollTop = Math.max(stageHeight - frame.clientHeight, 0);

  frame.scrollLeft = clamp(imageLeft + anchor.imageX * state.viewer.scale - anchor.viewportX, 0, maxScrollLeft);
  frame.scrollTop = clamp(imageTop + anchor.imageY * state.viewer.scale - anchor.viewportY, 0, maxScrollTop);
}

function syncImageViewerBaseSize(anchor = null) {
  if (!state.viewer.open) {
    return;
  }

  const image = elements.imageViewerImg;
  const naturalWidth = image.naturalWidth;
  const naturalHeight = image.naturalHeight;
  if (!naturalWidth || !naturalHeight) {
    return;
  }

  const preservedAnchor = anchor || captureImageViewerAnchor();
  const frame = elements.imageViewerFrame;
  const fitScale = Math.min(
    1,
    frame.clientWidth / naturalWidth,
    frame.clientHeight / naturalHeight,
  );

  state.viewer.naturalWidth = naturalWidth;
  state.viewer.naturalHeight = naturalHeight;
  state.viewer.baseWidth = Math.max(1, naturalWidth * fitScale);
  state.viewer.baseHeight = Math.max(1, naturalHeight * fitScale);

  renderImageViewerLayout();
  restoreImageViewerAnchor(preservedAnchor);
}

function zoomImageViewer(scale, anchor) {
  if (!state.viewer.open || !elements.imageViewerImg.complete || !state.viewer.baseWidth || !state.viewer.baseHeight) {
    return;
  }

  const nextScale = clamp(scale, viewerMinScale, viewerMaxScale);
  const previousScale = state.viewer.scale;
  if (Math.abs(nextScale - previousScale) < 0.001) {
    return;
  }

  const preservedAnchor = captureImageViewerAnchor(anchor);
  state.viewer.scale = nextScale;
  renderImageViewerLayout();
  restoreImageViewerAnchor(preservedAnchor);

  if (nextScale === viewerMinScale) {
    elements.imageViewerFrame.scrollLeft = 0;
    elements.imageViewerFrame.scrollTop = 0;
  }
}

function toggleImageViewerZoom(anchor) {
  if (!state.viewer.open) {
    return;
  }

  const nextScale = state.viewer.scale > viewerMinScale ? viewerMinScale : viewerDoubleClickScale;
  zoomImageViewer(nextScale, anchor);
}

async function loadBootstrap() {
  const snapshot = await Bootstrap();
  state.self = snapshot.self;
  state.settings = snapshot.settings || state.settings;
  state.peers = snapshot.peers || [];
  state.conversations = snapshot.conversations || {};
  state.conversationMeta = {};

  if (!state.activePeerId && state.peers.length) {
    state.activePeerId = state.peers[0].id;
  }

  renderStaticText();
  renderPeers();
  renderConversation();
  ensureConversationLoaded(state.activePeerId);
  syncUnreadState({ render: false });
  clearActivePeerUnreadIfVisible();
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
  renderSelfProfile();
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

elements.chatMenuTrigger.addEventListener('click', () => {
  if (state.chatMenu.open) {
    closeChatMenu();
    return;
  }
  openChatMenu();
});

elements.openSettings.addEventListener('click', () => {
  openSettingsDialog();
});

elements.settingsDialogClose.addEventListener('click', () => {
  closeSettingsDialog();
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
    rerenderConversation({ preserveViewport: true });
  });
});

  elements.themeOptions.forEach((option) => {
  option.addEventListener('click', async () => {
    closeThemeMenu();
    const settings = await UpdateTheme(option.dataset.theme);
    state.settings = settings;
    renderStaticText();
    renderPeers();
    rerenderConversation({ preserveViewport: true });
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
  if (!event.target.closest('.chat-more')) {
    closeChatMenu();
  }
});

elements.settingsDialog.addEventListener('click', (event) => {
  if (event.target === elements.settingsDialog || event.target.classList.contains('dialog-backdrop')) {
    closeSettingsDialog();
  }
});

elements.changeDataPath.addEventListener('click', async () => {
  try {
    const path = await ChooseDataDirectory();
    elements.dataPath.textContent = path;
    elements.dataPath.title = path;
    elements.dataPath.dataset.loaded = 'true';
    state.storageStats.loaded = false;
    renderStorageStats();
    scheduleStorageStatsRefresh();
    showToast(t('storageUpdated'), 'success');
  } catch (error) {
    console.error('ChooseDataDirectory failed', error);
  }
});

elements.cleanupMedia.addEventListener('click', () => {
  openActionDialog('cleanupMedia');
});

elements.addDebugBot.addEventListener('click', async () => {
  const peer = await EnsureDebugPeer();
  setActivePeer(peer.id);
});

elements.fileInput.addEventListener('change', async (event) => {
  try {
    await sendAttachmentFiles(event.target.files || []);
  } finally {
    elements.fileInput.value = '';
  }
});

elements.messageList.addEventListener('dblclick', (event) => {
  const image = event.target.closest('.message-image');
  if (!image) {
    return;
  }
  const article = image.closest('.message');
  const messageId = article?.dataset.messageId;
  if (!state.activePeerId || !messageId) {
    return;
  }
  void openImageViewerForMessage(state.activePeerId, messageId, image.dataset.previewSrc || '');
});

elements.messageList.addEventListener('click', (event) => {
  const fileCard = event.target.closest('.file-message');
  if (!fileCard) {
    return;
  }

  const article = fileCard.closest('.message');
  const messageId = article?.dataset.messageId;
  if (!state.activePeerId || !messageId) {
    return;
  }

  const message = findMessage(state.activePeerId, messageId);
  if (!message || message.kind !== 'file') {
    return;
  }

  void openStoredFileMessage(message);
});

elements.messageList.addEventListener('scroll', () => {
  if (!state.activePeerId) {
    return;
  }

  const meta = conversationMeta(state.activePeerId);
  if (!meta.hasMore || meta.loading) {
    return;
  }

  if (elements.messageList.scrollTop <= historyAutoLoadThreshold) {
    void loadConversationHistory(state.activePeerId);
  }
});

elements.imageViewerImg.addEventListener('dblclick', (event) => {
  event.stopPropagation();
  toggleImageViewerZoom(event);
});

elements.imageViewerImg.addEventListener('load', () => {
  syncImageViewerBaseSize();
});

elements.imageViewerPrev.addEventListener('click', () => {
  void stepImageViewer(-1);
});

elements.imageViewerNext.addEventListener('click', () => {
  void stepImageViewer(1);
});

function hideMessageMenu() {
  state.messageMenu.open = false;
  elements.messageContextMenu.classList.add('hidden');
  elements.messageContextMenu.setAttribute('aria-hidden', 'true');
}

function showMessageMenu(x, y, peerId, messageId, kind) {
  const message = findMessage(peerId, messageId);
  const canRetry = messageCanRetry(message);
  const canDelete = messageCanDelete(message);
  const canForward = messageCanForward(message);

  state.messageMenu.open = true;
  state.messageMenu.peerId = peerId;
  state.messageMenu.messageId = messageId;
  state.messageMenu.kind = kind || 'text';

  elements.messageContextCopy.textContent = kind === 'image'
    ? t('copyImage')
    : kind === 'file'
      ? t('copyFileName')
      : t('copy');
  elements.messageContextSave.classList.toggle('hidden', kind !== 'image' && kind !== 'file');
  elements.messageContextReveal.classList.toggle('hidden', kind !== 'image' && kind !== 'file');
  elements.messageContextRetry.classList.toggle('hidden', !canRetry);
  elements.messageContextForward.classList.toggle('hidden', !canForward);
  elements.messageContextDelete.classList.toggle('hidden', !canDelete);
  elements.messageContextDelete.textContent = t('deleteMessage');

  elements.messageContextMenu.style.left = `${x}px`;
  elements.messageContextMenu.style.top = `${y}px`;
  elements.messageContextMenu.classList.remove('hidden');
  elements.messageContextMenu.setAttribute('aria-hidden', 'false');
}

document.addEventListener('click', (event) => {
  if (state.messageMenu.open && !event.target.closest('.message-context-menu')) {
    hideMessageMenu();
  }
});

elements.messageList.addEventListener('contextmenu', (event) => {
  event.preventDefault();
  const article = event.target.closest('.message');
  if (!article || !state.activePeerId) {
    return;
  }

  const messageId = article.dataset.messageId;
  const kind = article.dataset.messageKind || 'text';
  const rect = article.getBoundingClientRect();
  const x = Math.max(8, Math.min(window.innerWidth - 200, event.clientX));
  const y = Math.max(8, Math.min(window.innerHeight - 100, event.clientY));
  showMessageMenu(x, y, state.activePeerId, messageId, kind);
});

elements.messageContextCopy.addEventListener('click', async () => {
  const peerId = state.messageMenu.peerId;
  const messageId = state.messageMenu.messageId;
  const kind = state.messageMenu.kind;
  if (!peerId || !messageId) return hideMessageMenu();

  const message = findMessage(peerId, messageId);
  if (!message) return hideMessageMenu();

  try {
    if (kind === 'image') {
      const imageSource = await resolveImageMessageSource(message);
      const copied = imageSource
        ? await copyImageToClipboard(imageSource)
        : false;

      if (!copied) {
        showToast(t('copyImageFailed'));
      }
    } else if (kind === 'file') {
      await ClipboardSetText(fileMessageName(message));
    } else {
      await ClipboardSetText(message.text || '');
    }
  } catch (err) {
    // best-effort
  }
  hideMessageMenu();
});

elements.messageContextSave.addEventListener('click', async () => {
  const peerId = state.messageMenu.peerId;
  const messageId = state.messageMenu.messageId;
  if (!peerId || !messageId) {
    hideMessageMenu();
    return;
  }

  const message = findMessage(peerId, messageId);
  if (!message || (message.kind !== 'image' && message.kind !== 'file')) {
    hideMessageMenu();
    return;
  }

  hideMessageMenu();
  try {
    if (message.kind === 'image') {
      await SaveImageMessage(message.text || '', suggestedImageName(message));
    } else {
      await SaveFileMessage(message.text || '', message.mediaType || '', suggestedFileName(message));
    }
  } catch (error) {
    if (message.kind === 'file' && isMissingFileError(error)) {
      setFileAvailability(message, false);
    }
    showToast(userFacingError(error, message.kind === 'image' ? 'saveImageFailed' : 'saveFileFailed'));
  }
});

elements.messageContextReveal.addEventListener('click', async () => {
  const peerId = state.messageMenu.peerId;
  const messageId = state.messageMenu.messageId;
  if (!peerId || !messageId) {
    hideMessageMenu();
    return;
  }

  const message = findMessage(peerId, messageId);
  if (!message || (message.kind !== 'image' && message.kind !== 'file')) {
    hideMessageMenu();
    return;
  }

  hideMessageMenu();
  try {
    if (message.kind === 'image') {
      await RevealImageMessage(message.text || '');
    } else {
      await RevealFileMessage(message.text || '');
    }
  } catch (error) {
    showToast(userFacingError(error, 'openFolderFailed'));
  }
});

elements.messageContextRetry.addEventListener('click', async () => {
  const peerId = state.messageMenu.peerId;
  const messageId = state.messageMenu.messageId;
  if (!peerId || !messageId) {
    hideMessageMenu();
    return;
  }

  const message = findMessage(peerId, messageId);
  hideMessageMenu();
  if (!message || !messageCanRetry(message)) {
    return;
  }

  try {
    await RetryMessage(peerId, messageId);
  } catch (error) {
    console.error('RetryMessage failed', error);
    if ((message.kind === 'file' || message.kind === 'image') && isMissingFileError(error)) {
      setFileAvailability(message, false);
    }
    showToast(userFacingError(error, message.kind === 'text' ? 'sendMessageFailed' : message.kind === 'image' ? 'sendImageFailed' : 'sendFileFailed'));
  }
});

elements.messageContextForward.addEventListener('click', () => {
  const peerId = state.messageMenu.peerId;
  const messageId = state.messageMenu.messageId;
  if (!peerId || !messageId) {
    hideMessageMenu();
    return;
  }

  hideMessageMenu();
  openForwardDialog(peerId, messageId);
});

elements.messageContextDelete.addEventListener('click', async () => {
  const peerId = state.messageMenu.peerId;
  const messageId = state.messageMenu.messageId;
  if (!peerId || !messageId) {
    hideMessageMenu();
    return;
  }
  hideMessageMenu();
  openDeleteDialog(peerId, messageId);
});

elements.deleteConfirmCancel.addEventListener('click', () => {
  closeDeleteDialog();
});

elements.deleteConfirmSubmit.addEventListener('click', async () => {
  const peerId = state.deleteDialog.peerId;
  const messageId = state.deleteDialog.messageId;
  if (!peerId || !messageId) {
    closeDeleteDialog();
    return;
  }

  const previousMessages = removeMessageLocally(peerId, messageId);
  try {
    closeDeleteDialog();
    await DeleteMessage(peerId, messageId);
  } catch (err) {
    console.error('DeleteMessage failed', err);
    if (previousMessages) {
      setConversationMessages(peerId, previousMessages);
    }
    closeDeleteDialog();
  }
});

elements.deleteConfirmDialog.addEventListener('click', (event) => {
  if (event.target === elements.deleteConfirmDialog || event.target.classList.contains('dialog-backdrop')) {
    closeDeleteDialog();
  }
});

elements.clearConversation.addEventListener('click', () => {
  if (!state.activePeerId || elements.clearConversation.disabled) {
    return;
  }
  closeChatMenu();
  openActionDialog('clearConversation', state.activePeerId);
});

elements.actionConfirmCancel.addEventListener('click', () => {
  closeActionDialog();
});

elements.actionConfirmSubmit.addEventListener('click', async () => {
  if (state.actionDialog.busy) {
    return;
  }

  const { mode, peerId } = state.actionDialog;
  if (!mode) {
    closeActionDialog();
    return;
  }

  state.actionDialog.busy = true;
  syncActionDialogText();

  try {
    if (mode === 'clearConversation') {
      const removedCount = await ClearConversation(peerId);
      const meta = conversationMeta(peerId);
      meta.initialized = true;
      meta.hasMore = false;
      meta.loading = false;
      if (state.search.open && state.search.peerId === peerId) {
        closeConversationSearch();
      }
      state.actionDialog.busy = false;
      closeActionDialog();
      showToast(t('conversationCleared', {
        count: String(Math.max(0, Number(removedCount) || 0)),
      }), 'success');
      scheduleStorageStatsRefresh();
      return;
    }

    if (mode === 'cleanupMedia') {
      const result = await CleanupUnusedMedia();
      state.storageStats.mediaBytes = Math.max(0, Number(result?.mediaBytes) || 0);
      state.storageStats.loaded = true;
      renderStorageStats();
      state.actionDialog.busy = false;
      closeActionDialog();
      if ((Number(result?.removedFiles) || 0) > 0) {
        showToast(t('cleanupMediaDone', {
          count: String(Math.max(0, Number(result.removedFiles) || 0)),
          size: formatFileSize(result.reclaimedBytes),
        }), 'success');
      } else {
        showToast(t('cleanupMediaEmpty'), 'success');
      }
      scheduleStorageStatsRefresh();
      return;
    }

    state.actionDialog.busy = false;
    closeActionDialog();
  } catch (error) {
    console.error('Action dialog submit failed', error);
    state.actionDialog.busy = false;
    syncActionDialogText();
    showToast(userFacingError(error, mode === 'clearConversation' ? 'clearConversationFailed' : 'cleanupMediaFailed'));
  }
});

elements.actionConfirmDialog.addEventListener('click', (event) => {
  if (event.target === elements.actionConfirmDialog || event.target.classList.contains('dialog-backdrop')) {
    closeActionDialog();
  }
});

elements.forwardDialogCancel.addEventListener('click', () => {
  closeForwardDialog();
});

elements.forwardDialogSubmit.addEventListener('click', () => {
  const sourcePeerId = state.forwardDialog.peerId;
  const messageId = state.forwardDialog.messageId;
  const targetPeerId = state.forwardDialog.targetPeerId;
  if (!sourcePeerId || !messageId || !targetPeerId) {
    closeForwardDialog();
    return;
  }

  const message = findMessage(sourcePeerId, messageId);
  if (!message) {
    closeForwardDialog();
    return;
  }

  elements.forwardDialogSubmit.disabled = true;

  closeForwardDialog();

  const sendPromise = message.kind === 'image'
    ? resolveImageMessageSource(message).then((imageSource) => {
      if (!imageSource) {
        throw new Error('image source is unavailable');
      }
      return SendFileMessage(targetPeerId, imageSource, suggestedImageName(message));
    })
    : message.kind === 'file'
      ? loadFileMessageSource(message).then((fileSource) => {
        if (!fileSource) {
          throw new Error('file source is unavailable');
        }
        return SendFileMessage(targetPeerId, fileSource, suggestedFileName(message));
      })
      : SendChatMessage(targetPeerId, message.text || '');

  sendPromise.catch((error) => {
    console.error('Forward message failed', error);
    if (message.kind === 'file' && isMissingFileError(error)) {
      setFileAvailability(message, false);
    }
    showToast(userFacingError(error, 'forwardFailed'));
  });
});

elements.forwardDialog.addEventListener('click', (event) => {
  if (event.target === elements.forwardDialog || event.target.classList.contains('dialog-backdrop')) {
    closeForwardDialog();
  }
});

elements.imageViewerZoomOut.addEventListener('click', () => {
  zoomImageViewer(state.viewer.scale / 1.2);
});

elements.imageViewerZoomIn.addEventListener('click', () => {
  zoomImageViewer(state.viewer.scale * 1.2);
});

elements.imageViewerReset.addEventListener('click', () => {
  zoomImageViewer(viewerMinScale);
});

elements.imageViewerFrame.addEventListener('wheel', (event) => {
  if (!state.viewer.open || !elements.imageViewerImg.src) {
    return;
  }

  event.preventDefault();
  const multiplier = Math.exp(-event.deltaY * viewerWheelIntensity);
  zoomImageViewer(state.viewer.scale * multiplier, event);
}, { passive: false });

elements.imageViewerClose.addEventListener('click', () => {
  closeImageViewer();
});

elements.imageViewer.addEventListener('click', (event) => {
  if (event.target === elements.imageViewer || event.target.classList.contains('image-viewer-backdrop')) {
    closeImageViewer();
  }
});

document.addEventListener('keydown', (event) => {
  if (event.key === 'Escape' && state.viewer.open) {
    closeImageViewer();
    return;
  }
  if (state.viewer.open && event.key === 'ArrowLeft') {
    event.preventDefault();
    void stepImageViewer(-1);
    return;
  }
  if (state.viewer.open && event.key === 'ArrowRight') {
    event.preventDefault();
    void stepImageViewer(1);
    return;
  }
  if (event.key === 'Escape' && state.deleteDialog.open) {
    closeDeleteDialog();
    return;
  }
  if (event.key === 'Escape' && state.actionDialog.open) {
    closeActionDialog();
    return;
  }
  if (event.key === 'Escape' && state.forwardDialog.open) {
    closeForwardDialog();
    return;
  }
  if (event.key === 'Escape' && state.settingsDialog.open) {
    closeSettingsDialog();
    return;
  }
  if (event.key === 'Escape' && state.chatMenu.open) {
    closeChatMenu();
    return;
  }
  if (event.key === 'Escape' && state.search.open && document.activeElement !== elements.chatSearchInput) {
    closeConversationSearch();
  }
});

window.addEventListener('resize', () => {
  if (!state.viewer.open) {
    return;
  }
  syncImageViewerBaseSize();
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
  setActivePeer(peer.id);
});

EventsOn('data:pathUpdated', (path) => {
  elements.dataPath.textContent = path;
  elements.dataPath.title = path;
  elements.dataPath.dataset.loaded = 'true';
  state.storageStats.loaded = false;
  renderStorageStats();
  scheduleStorageStatsRefresh();
});

EventsOn('settings:updated', (settings) => {
  state.settings = settings;
  renderStaticText();
  if (state.forwardDialog.open) {
    renderForwardPreview();
    renderForwardPeerList();
  }
  renderPeers();
  rerenderConversation({ preserveViewport: true });
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

  try {
    await SendChatMessage(state.activePeerId, text);
    elements.messageInput.value = '';
    closeEmojiPanel();
  } catch (error) {
    console.error('SendChatMessage failed', error);
    showToast(userFacingError(error, 'sendMessageFailed'));
  }
});

elements.messageInput.addEventListener('keydown', async (event) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault();
    elements.composer.requestSubmit();
  }
});

elements.messageInput.addEventListener('paste', (event) => {
  const files = Array.from(event.clipboardData?.files || []);
  if (!files.length) {
    return;
  }

  event.preventDefault();
  void sendAttachmentFiles(files);
});

elements.searchToggle.addEventListener('click', () => {
  closeChatMenu();
  if (state.search.open) {
    closeConversationSearch();
    return;
  }
  openConversationSearch();
});

elements.chatSearchInput.addEventListener('input', (event) => {
  state.search.query = event.target.value || '';
  void refreshConversationSearch({ preserveSelection: false });
});

elements.chatSearchInput.addEventListener('keydown', (event) => {
  if (event.key === 'Enter') {
    event.preventDefault();
    setSearchResultIndex(state.search.activeIndex + (event.shiftKey ? -1 : 1));
  }
  if (event.key === 'Escape') {
    event.preventDefault();
    closeConversationSearch();
  }
});

elements.chatSearchPrev.addEventListener('click', () => {
  setSearchResultIndex(state.search.activeIndex - 1);
});

elements.chatSearchNext.addEventListener('click', () => {
  setSearchResultIndex(state.search.activeIndex + 1);
});

elements.chatSearchClose.addEventListener('click', () => {
  closeConversationSearch();
});

EventsOn('profile:updated', (profile) => {
  state.self = profile;
  renderSelfProfile();
});

window.addEventListener('focus', () => {
  updateWindowFocusedState(true);
});

window.addEventListener('blur', () => {
  updateWindowFocusedState(false);
});

document.addEventListener('visibilitychange', () => {
  updateWindowFocusedState(document.hasFocus());
});

EventsOn('transfer:progress', (payload) => {
  const messageId = String(payload?.messageId || '').trim();
  if (!messageId) {
    return;
  }

  state.transferProgress[messageId] = Math.max(0, Math.min(100, Number(payload?.progress) || 0));
  if (payload?.peerId && payload.peerId === state.activePeerId) {
    rerenderConversation({ preserveViewport: true });
  }
});

OnFileDrop((_x, _y, paths) => {
  if (!paths?.length) {
    return;
  }
  void sendLocalAttachmentPaths(paths);
}, false);

window.addEventListener('drop', (event) => {
  if (CanResolveFilePaths() || !isExternalFileDropEvent(event)) {
    return;
  }
  event.preventDefault();
  event.stopPropagation();
  void sendAttachmentFiles(event.dataTransfer?.files || []);
});

EventsOn('window:visibilityChanged', (payload) => {
  updateWindowVisibleState(payload?.visible !== false);
  if (payload?.visible !== false) {
    updateWindowFocusedState(document.hasFocus());
  }
});

EventsOn('peers:updated', (peers) => {
  state.peers = peers || [];
  const knownPeerIDs = new Set(state.peers.map((peer) => peer.id));
  const unreadChanged = pruneUnreadCounts(knownPeerIDs);

  if (state.activePeerId && !state.peers.some((peer) => peer.id === state.activePeerId)) {
    state.activePeerId = state.peers[0]?.id || null;
  }
  if (!state.activePeerId && state.peers.length) {
    state.activePeerId = state.peers[0].id;
  }

  if (state.forwardDialog.open) {
    renderForwardPeerList();
  }
  if (unreadChanged) {
    syncUnreadState({ render: false });
  }
  renderPeers();
  renderConversation();
  clearActivePeerUnreadIfVisible();
  ensureConversationLoaded(state.activePeerId);
  if (!state.activePeerId && state.search.open) {
    closeConversationSearch();
    return;
  }
  if (state.search.open) {
    void refreshConversationSearch({ preserveSelection: false });
  } else {
    syncSearchUI();
  }
});

EventsOn('conversation:updated', (payload) => {
  const peerId = String(payload?.peerId || '').trim();
  if (!peerId) {
    return;
  }

  const previousMessages = state.conversations[peerId] || [];
  const nextMessages = payload.messages || [];
  const newInboundMessages = collectNewInboundMessages(previousMessages, nextMessages);
  state.conversations[peerId] = nextMessages;
  syncTransferProgressCache(previousMessages, nextMessages);
  scheduleStorageStatsRefresh();

  if (!state.activePeerId) {
    state.activePeerId = peerId;
  }

  if (newInboundMessages.length > 0) {
    if (isConversationVisible(peerId)) {
      clearUnread(peerId, { render: false });
    } else {
      incrementUnread(peerId, newInboundMessages.length);
      void notifyIncomingMessages(peerId, newInboundMessages);
    }
  }

  rerenderConversation({ preserveViewport: true });
  clearActivePeerUnreadIfVisible();
  if (state.activePeerId === peerId) {
    ensureConversationLoaded(peerId);
  }
  if (state.viewer.open && state.viewer.peerId === peerId) {
    const currentImageExists = conversationImageMessages(peerId)
      .some((message) => message.id === state.viewer.messageId);
    if (!currentImageExists) {
      closeImageViewer();
      return;
    }
    syncImageViewerNavigation();
  }
  if (state.search.open && state.search.peerId === peerId) {
    void refreshConversationSearch({ ensureLoaded: false });
  } else {
    syncSearchUI();
  }
});

void initializeDesktopNotifications();

Promise.all([loadBootstrap(), loadDataPath(), refreshStorageStats()]).catch((error) => {
  console.error(error);
  elements.messageList.className = 'message-list empty';
  elements.messageList.innerHTML = `
    <div class="empty-state">
      <h3>${escapeHtml(t('bootFailed'))}</h3>
      <p>${escapeHtml(String(error))}</p>
    </div>
  `;
});
