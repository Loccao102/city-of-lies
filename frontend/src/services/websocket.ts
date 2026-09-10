export interface WSEventMessage {
  sequence: number;
  type: string;
  game_second: number;
  payload: any;
}

export type WSEventHandler = (event: WSEventMessage) => void;

export class GameWebSocket {
  private ws: WebSocket | null = null;
  private url: string;
  private sessionId: string;
  private handlers: Set<WSEventHandler> = new Set();
  private reconnectTimeout: any = null;
  private shouldReconnect = true;

  constructor(sessionId: string) {
    this.sessionId = sessionId;
    const wsBase = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/ws";
    this.url = `${wsBase}?session_id=${sessionId}`;
  }

  connect() {
    if (this.ws && (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING)) {
      return;
    }

    try {
      this.ws = new WebSocket(this.url);

      this.ws.onopen = () => {
        console.log("[WS] Connected to City of Lies session:", this.sessionId);
        if (this.reconnectTimeout) {
          clearTimeout(this.reconnectTimeout);
          this.reconnectTimeout = null;
        }
      };

      this.ws.onmessage = (event) => {
        try {
          const msg: WSEventMessage = JSON.parse(event.data);
          this.handlers.forEach((h) => h(msg));
        } catch (e) {
          console.error("[WS] Parse error:", e);
        }
      };

      this.ws.onclose = () => {
        console.log("[WS] Disconnected");
        if (this.shouldReconnect) {
          this.reconnectTimeout = setTimeout(() => this.connect(), 2000);
        }
      };

      this.ws.onerror = (err) => {
        console.warn("[WS] Error:", err);
      };
    } catch (err) {
      console.error("[WS] Connection attempt failed:", err);
    }
  }

  subscribe(handler: WSEventHandler) {
    this.handlers.add(handler);
    return () => {
      this.handlers.delete(handler);
    };
  }

  disconnect() {
    this.shouldReconnect = false;
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
    }
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}
