import {PromiseWithKnownReason} from '@reduxjs/toolkit/dist/query/core/buildMiddleware/types';
import webSocketGateway, {IListenerFunction} from 'gateways/WebSocket.gateway';

export type {IListenerFunction} from 'gateways/WebSocket.gateway';

interface IInitWebSocketSubscription {
  listener: IListenerFunction;
  resource: string;
  waitToCleanSubscription: Promise<void>;
  waitToInitSubscription: PromiseWithKnownReason<any, any>;
}

/**
 * Implements RTK Query logic to receive streaming updates for persistent queries.
 * It allows a query to establish a WebSocket subscription, and apply updates to
 * the cached data when the information is received from the server.
 * For more information please take a look: https://redux-toolkit.js.org/rtk-query/usage/streaming-updates
 */
export async function initWebSocketSubscription({
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
  webSocketGateway.unsubscribe(resource, '1234');
}
