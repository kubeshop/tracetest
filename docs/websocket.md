# Websocket endpoint

Tracetest allow you to subscribe to updates of resources using websockets. There are two endpoints that you can use to manage subscriptions:

## Endpoint
 
You can open a websocket connection by sending a request the path `/ws`. Example: `ws://localhost:11633/ws`.

## Messages

### Subscribing to updates

Once the connection is open, you can send a message with the format:

```json
{
    "type": "subscribe",
    "resource": "test/{testID}/run/{runID}"
}
```

If a problem happens, you will see an error like:

```json
{
    "type": "error",
    "message": "details of the error"
}
```

If the operation executes successfully, you will see a response like:

```json
{
    "type": "success",
    "resource": "test/{testID}/run/{runID}",
    "message": {
        "subscriptionId": "bdbc6cc8-bba6-4208-a8d3-d3c2c5b3e38b"
    }
}
```

The `subscriptionId` is an important field because it is required to cancel the subscription. You should store it, otherwise you will keep receiving updates of a resource that you might not want.

### Cancel a susbcription

Once you have a subscription to a resource, you might want to stop receiving events from that resource. So, there is a `unsubscribe` message that you can send to achieve that.

But send this message to the websocket connection:

```json
{
    "type": "unsubscribe",
    "resource": "test/{testID}/run/{runID}",
    "subscriptionId": "id returned in the subscription command"
}
```

You will receive an error message if any required field is not field. But in case of a successful operation, you will receive a message like:

```json
{
    "type":"success",
    "message":"ok"
}
```

This message will be sent regardless if the subscription exists or not.
