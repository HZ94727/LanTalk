package chat

type Profile struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ListenPort int    `json:"listenPort"`
}

type Peer struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	ListenPort int    `json:"listenPort"`
	LastSeen   int64  `json:"lastSeen"`
	Source     string `json:"source"`
}

type ChatMessage struct {
	ID         string `json:"id"`
	PeerID     string `json:"peerId"`
	SenderID   string `json:"senderId"`
	SenderName string `json:"senderName"`
	Kind       string `json:"kind"`
	Text       string `json:"text"`
	MediaName  string `json:"mediaName"`
	MediaType  string `json:"mediaType"`
	MediaSize  int64  `json:"mediaSize"`
	MediaHash  string `json:"-"`
	Timestamp  int64  `json:"timestamp"`
	Direction  string `json:"direction"`
	Status     string `json:"status"`
}

type TransferProgress struct {
	PeerID    string `json:"peerId"`
	MessageID string `json:"messageId"`
	Progress  int    `json:"progress"`
}

type StorageStats struct {
	MessageCount int   `json:"messageCount"`
	MediaBytes   int64 `json:"mediaBytes"`
}

type MediaCleanupResult struct {
	RemovedFiles   int   `json:"removedFiles"`
	ReclaimedBytes int64 `json:"reclaimedBytes"`
	MediaBytes     int64 `json:"mediaBytes"`
}

type Snapshot struct {
	Self          Profile                  `json:"self"`
	Settings      Settings                 `json:"settings"`
	Peers         []Peer                   `json:"peers"`
	Conversations map[string][]ChatMessage `json:"conversations"`
}

type ConversationPage struct {
	PeerID   string        `json:"peerId"`
	Messages []ChatMessage `json:"messages"`
	HasMore  bool          `json:"hasMore"`
}

type Settings struct {
	Language string `json:"language"`
	Theme    string `json:"theme"`
}

type announcement struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ListenPort int    `json:"listenPort"`
}

type directMessage struct {
	ID         string `json:"id"`
	SenderID   string `json:"senderId"`
	SenderName string `json:"senderName"`
	Kind       string `json:"kind"`
	Text       string `json:"text"`
	MediaName  string `json:"mediaName"`
	MediaType  string `json:"mediaType"`
	MediaSize  int64  `json:"mediaSize"`
	Timestamp  int64  `json:"timestamp"`
}
