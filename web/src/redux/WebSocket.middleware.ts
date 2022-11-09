import {Middleware} from 'redux';
import {createAction} from '@reduxjs/toolkit';

import webSocketGateway, {IListenerFunction, IMessage} from 'gateways/WebSocket.gateway';
import {RootState} from 'redux/store';

export const webSocketSendMessage = createAction<string>('websocket/MESSAGE_SEND');
export const webSocketMessageReceived = createAction<string>('websocket/MESSAGE_RECEIVED');

function createWebSocketMiddleware(): Middleware<{}, RootState> {
  return storeApi => {
    // Connect the websocket
    webSocketGateway.connect();

    const listener: IListenerFunction = data => {
      // Here we can dispatch an action with the websocket payload
      // or apply a manual cache update for the query endpoint
      storeApi.dispatch(webSocketMessageReceived(data.event));
    };

    return next => action => {
      if (action.type === 'tests/executeQuery/fulfilled' && action?.meta?.arg?.endpointName === 'getRunById') {
        const args = action?.meta?.arg?.originalArgs ?? {};
        webSocketGateway.subscribe(`test/${args.testId}/run/${args.runId}`, listener);
      }

      if (webSocketSendMessage.match(action)) {
        console.log('send socket middleware');
        const message: IMessage = {type: 'subscribe', resource: action.payload};
        webSocketGateway.send(message);
        return;
      }

      return next(action);
    };
  };
}

export default createWebSocketMiddleware;
