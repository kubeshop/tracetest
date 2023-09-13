import Env from 'utils/Env';

interface IData<T> {
  event: T;
  resource: string;
  type: string;
}

export interface IListenerFunction<T = any> {
  (data: IData<T>): void;
}

interface IMessage {
  type: 'subscribe' | 'unsubscribe';
  resource: string;
  subscriptionId?: string;
}

interface IWebSocketGateway {
  /** Opens the websocket connection and sets default event listeners */
  connect(): void;
  /** Closes the websocket connection */
  disconnect(): void;
  /** Listens on the given `event` with `listener` */
  on(event: string, listener: IListenerFunction): void;
  /** Removes all registered listeners for `event` */
  off(event: string): void;
  /** Sends `message` to the server */
  send(message: IMessage): void;
  /** Subscribes to updates on the given `resource` */
  subscribe(resource: string, listener: IListenerFunction): void;
  /** Cancels a subscription on the given `resource` */
  unsubscribe(resource: string): void;
}

interface IParams {
  url: string;
}

const MAX_RECONNECTION_ATTEMPTS = 4;
const DELAY_RECONNECTION_ATTEMPTS = 1000;
enum MESSAGE_REQUEST_TYPE {
  SUBSCRIBE = 'subscribe',
  UNSUBSCRIBE = 'unsubscribe',
}
enum MESSAGE_RESPONSE_TYPE {
  SUCCESS = 'success',
  UPDATE = 'update',
}

const WebSocketGateway = ({url}: IParams): IWebSocketGateway => {
  let socket: WebSocket | null = null;
  let isConnected: boolean = false;
  let pendingToSend: string[] = [];
  let listeners: Record<string, IListenerFunction[]> = {};
  let subscriptionIds: Record<string, string> = {};
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

  const closeListener = () => {
    isConnected = false;
    cleanUpAttempts();
    reconnect();
  };

  const messageListener = (event: MessageEvent) => {
    const data = JSON.parse(event.data);

    if (data.type === MESSAGE_RESPONSE_TYPE.UPDATE) {
      const eventListeners = listeners[data.resource] || [];
      eventListeners.forEach(listener => listener(data));
    }

    if (data.type === MESSAGE_RESPONSE_TYPE.SUCCESS && data.resource) {
      subscriptionIds[data.resource] = data?.message?.subscriptionId;
    }
  };

  const connect = () => {
    if (socket !== null) socket.close();
    socket = new WebSocket(url);
    socket.addEventListener('open', openListener);
    socket.addEventListener('close', closeListener);
    socket.addEventListener('message', messageListener);
  };

  const disconnect = () => {
    if (socket !== null) socket.close();
    socket = null;
    isConnected = false;
    pendingToSend = [];
    listeners = {};
    subscriptionIds = {};
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
    connect() {
      connect();
    },
    disconnect() {
      disconnect();
    },
    on(event, listener) {
      const eventListeners = listeners[event] || [];
      eventListeners.push(listener);
      listeners = {...listeners, [event]: eventListeners};
    },
    off(event) {
      delete listeners[event];
    },
    send(message) {
      if (!socket) {
        return;
      }
      if (!isConnected) {
        pendingToSend.push(JSON.stringify(message));
        return;
      }
      socket.send(JSON.stringify(message));
    },
    subscribe(resource, listener) {
      const message: IMessage = {type: MESSAGE_REQUEST_TYPE.SUBSCRIBE, resource};
      this.send(message);
      this.on(resource, listener);
    },
    unsubscribe(resource) {
      const subscriptionId = subscriptionIds?.[resource] ?? '';
      delete subscriptionIds[resource];
      const message: IMessage = {type: MESSAGE_REQUEST_TYPE.UNSUBSCRIBE, resource, subscriptionId};
      this.send(message);
      this.off(resource);
    },
  };
};

function getWebSocketURL() {
  const serverPathPrefix = Env.get('serverPathPrefix');
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
  const hostname = window.location.hostname;
  const port = process.env.NODE_ENV === 'development' ? '11633' : window.location.port;
  const pathname = serverPathPrefix === '/' ? '/ws' : `${serverPathPrefix}/ws`;
  return `${protocol}://${hostname}:${port}${pathname}`;
}

const webSocketGateway = WebSocketGateway({url: getWebSocketURL()});
// Disable websocket connection for now
// webSocketGateway.connect();

export default webSocketGateway;
