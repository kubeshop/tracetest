import {PromiseWithKnownReason} from '@reduxjs/toolkit/dist/query/core/buildMiddleware/types';
import webSocketGateway, {IListenerFunction} from 'gateways/WebSocket.gateway';

export type {IListenerFunction} from 'gateways/WebSocket.gateway';

interface IInitWebSocketSubscription {
  listener: IListenerFunction;
  resource: string;
  waitToCleanSubscription: Promise<void>;
  waitToInitSubscription: PromiseWithKnownReason<any, any>;
}

const WebSocketService = () => ({
  async initWebSocketSubscription({
    listener,
    resource,
    waitToCleanSubscription,
    waitToInitSubscription,
  }: IInitWebSocketSubscription) {
    try {
      await waitToInitSubscription;
      webSocketGateway.subscribe(resource, listener);
    } catch {
      // no-op in case `waitToCleanSubscription` resolves before `waitToInitSubscription`,
      // in which case `waitToInitSubscription` will throw
    }
    await waitToCleanSubscription;
    webSocketGateway.unsubscribe(resource);
  },
});

export default WebSocketService();
