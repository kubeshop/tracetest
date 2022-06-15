import debugModule from 'debug';
import {LOCALHOST_URL_WEB_SOCKET} from 'constants/Common.constants';

const debug = debugModule('WebSocketGateway');

interface IData<T> {
  event: T;
  resource: string;
  type: string;
}

export interface IListenerFunction<T = any> {
  (data: IData<T>): void;
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
  /** Sends `data` to the server */
  send(data: any): void;
  /** Subscribes to updates on the given `resource` */
  subscribe(resource: string, listener: IListenerFunction): void;
  /** Cancels a subscription on the given `resource` */
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
  let listeners: Record<string, IListenerFunction[]> = {};
  let reconnectionAttempts: number = 0;
  let attempts: VoidFunction[] = [];

  const cleanUpAttempts = () => {
    attempts.forEach(attempt => attempt());
    attempts = [];
  };

  const openListener = () => {
    debug('openListener');
    isConnected = true;
    cleanUpAttempts();
    pendingToSend.forEach(pending => socket?.send(pending));
    pendingToSend = [];
  };

  const closeListener = () => {
    debug('closeListener');
    connect();
  };

  const errorListener = () => {
    debug('errorListener');
    isConnected = false;
    cleanUpAttempts();
    reconnect();
  };

  const messageListener = (event: MessageEvent) => {
    const data = JSON.parse(event.data);
    debug('messageListener: %O', data);
    const eventListeners = listeners[data.resource] || [];
    eventListeners.forEach(listener => listener(data));
  };

  const connect = () => {
    if (socket !== null) socket.close();
    socket = new WebSocket(url);
    socket.addEventListener('open', openListener);
    socket.addEventListener('close', closeListener);
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
    debug('reconnect');
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
      debug('connect');
      connect();
    },
    disconnect() {
      debug('disconnect');
      disconnect();
    },
    on(event, listener) {
      debug('on');
      const eventListeners = listeners[event] || [];
      eventListeners.push(listener);
      listeners = {...listeners, [event]: eventListeners};
    },
    off(event) {
      debug('off');
      delete listeners[event];
    },
    send(data) {
      debug('send %O', data);
      if (!socket) {
        return;
      }
      if (!isConnected) {
        debug('send pending');
        pendingToSend.push(JSON.stringify(data));
        return;
      }
      socket.send(JSON.stringify(data));
    },
    subscribe(resource, listener) {
      debug('subscribe %s', resource);
      const data = {type: 'subscribe', resource};
      this.send(data);
      this.on(resource, listener);
    },
    unsubscribe(resource, subscriptionId) {
      debug('unsubscribe %s', resource);
      const data = {type: 'unsubscribe', resource, subscriptionId};
      this.send(data);
      this.off(resource);
    },
  };
};

const URL = document.baseURI.includes('localhost') ? LOCALHOST_URL_WEB_SOCKET : `${document.baseURI}ws/`;
const webSocketGateway = WebSocketGateway({url: URL});
webSocketGateway.connect();

export default webSocketGateway;
