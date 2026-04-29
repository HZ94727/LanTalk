export namespace chat {
	
	export class ChatMessage {
	    id: string;
	    peerId: string;
	    senderId: string;
	    senderName: string;
	    kind: string;
	    text: string;
	    mediaName: string;
	    mediaType: string;
	    mediaSize: number;
	    timestamp: number;
	    direction: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.peerId = source["peerId"];
	        this.senderId = source["senderId"];
	        this.senderName = source["senderName"];
	        this.kind = source["kind"];
	        this.text = source["text"];
	        this.mediaName = source["mediaName"];
	        this.mediaType = source["mediaType"];
	        this.mediaSize = source["mediaSize"];
	        this.timestamp = source["timestamp"];
	        this.direction = source["direction"];
	        this.status = source["status"];
	    }
	}
	export class ConversationPage {
	    peerId: string;
	    messages: ChatMessage[];
	    hasMore: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ConversationPage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.peerId = source["peerId"];
	        this.messages = this.convertValues(source["messages"], ChatMessage);
	        this.hasMore = source["hasMore"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MediaCleanupResult {
	    removedFiles: number;
	    reclaimedBytes: number;
	    mediaBytes: number;
	
	    static createFrom(source: any = {}) {
	        return new MediaCleanupResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.removedFiles = source["removedFiles"];
	        this.reclaimedBytes = source["reclaimedBytes"];
	        this.mediaBytes = source["mediaBytes"];
	    }
	}
	export class Peer {
	    id: string;
	    name: string;
	    address: string;
	    listenPort: number;
	    lastSeen: number;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new Peer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.listenPort = source["listenPort"];
	        this.lastSeen = source["lastSeen"];
	        this.source = source["source"];
	    }
	}
	export class Profile {
	    id: string;
	    name: string;
	    listenPort: number;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.listenPort = source["listenPort"];
	    }
	}
	export class Settings {
	    language: string;
	    theme: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.theme = source["theme"];
	    }
	}
	export class Snapshot {
	    self: Profile;
	    settings: Settings;
	    peers: Peer[];
	    conversations: Record<string, Array<ChatMessage>>;
	
	    static createFrom(source: any = {}) {
	        return new Snapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.self = this.convertValues(source["self"], Profile);
	        this.settings = this.convertValues(source["settings"], Settings);
	        this.peers = this.convertValues(source["peers"], Peer);
	        this.conversations = this.convertValues(source["conversations"], Array<ChatMessage>, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StorageStats {
	    messageCount: number;
	    mediaBytes: number;
	
	    static createFrom(source: any = {}) {
	        return new StorageStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.messageCount = source["messageCount"];
	        this.mediaBytes = source["mediaBytes"];
	    }
	}

}

