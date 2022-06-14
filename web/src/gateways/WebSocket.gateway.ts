interface ListenerFunction {
  (data: any): void;
}

interface IWebSocketGateway {
  connect(): void;
  disconnect(): void;
  on(event: string, listener: ListenerFunction): void;
  off(event: string): void;
  send(data: any): void;
  subscribe(resource: string, listener: ListenerFunction): void;
  unsubscribe(resource: string, subscriptionId: string): void;
}

interface IParams {
  url: string;
}

const MAX_RECONNECTION_ATTEMPTS = 4;
const DELAY_RECONNECTION_ATTEMPTS = 1000;

const WebSocketGateway = ({url}: IParams): IWebSocketGateway => {
  let socket: WebSocket | null = null;
  let isConnected: boolean = false;
  let pendingToSend: string[] = [];
  let listeners: Record<string, ListenerFunction[]> = {};
  let reconnectionAttempts: number = 0;
  let attempts: VoidFunction[] = [];

  const cleanUpAttempts = () => {
    attempts.forEach(attempt => attempt());
    attempts = [];
  };

  const openListener = () => {
    isConnected = true;
    cleanUpAttempts();
    pendingToSend.forEach(pending => socket?.send(pending));
    pendingToSend = [];
  };

  const errorListener = () => {
    isConnected = false;
    cleanUpAttempts();
    reconnect();
  };

  const messageListener = (event: MessageEvent) => {
    const data = JSON.parse(event.data);
    const eventListeners = listeners[data.resource] || [];
    eventListeners.forEach(listener => listener(data));
  };

  const connect = () => {
    if (socket !== null) socket.close();
    socket = new WebSocket(url);
    socket.addEventListener('open', openListener);
    socket.addEventListener('error', errorListener);
    socket.addEventListener('message', messageListener);
  };

  const disconnect = () => {
    if (socket !== null) socket.close();
    socket = null;
    isConnected = false;
    pendingToSend = [];
    listeners = {};
    reconnectionAttempts = 0;
  };

  const reconnect = () => {
    if (reconnectionAttempts >= MAX_RECONNECTION_ATTEMPTS) {
      disconnect();
      return;
    }

    const timer = setTimeout(() => {
      reconnectionAttempts += 1;
      connect();
    }, DELAY_RECONNECTION_ATTEMPTS);

    attempts.push(() => {
      clearTimeout(timer);
    });
  };

  return {
    /** Opens the websocket connection and sets default event listeners */
    connect() {
      connect();
    },
    /** Closes the websocket connection */
    disconnect() {
      disconnect();
    },
    /** Listens on the given `event` with `listener` */
    on(event, listener) {
      const eventListeners = listeners[event] || [];
      eventListeners.push(listener);
      listeners = {...listeners, [event]: eventListeners};
    },
    /** Removes all registered callbacks for `event` */
    off(event) {
      delete listeners[event];
    },
    /** Sends `data` to the server */
    send(data) {
      if (!socket) {
        return;
      }
      if (!isConnected) {
        pendingToSend.push(JSON.stringify(data));
        return;
      }
      socket.send(JSON.stringify(data));
    },
    /** Subscribes to updates on the given `resource` */
    subscribe(resource, listener) {
      const data = {type: 'subscribe', resource};
      this.send(data);
      this.on(resource, listener);
    },
    /** Cancels a subscription on the given `resource` */
    unsubscribe(resource, subscriptionId) {
      const data = {type: 'unsubscribe', resource, subscriptionId};
      this.send(data);
      this.off(resource);
    },
  };
};

export default WebSocketGateway;
