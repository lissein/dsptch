# DSPTCH

Small services to dispatch realtime messages to clients. This is useful when your backend is
using languages or frameworks that makes realtime messages hard or not performant/scalable.

API:

- Publish messages to broker/queue => Redis, sqs ...

DSPTCH:

- listen for messages in sources
- handle websocket connections / push notif registrations
- when message read from queue => dispatch to destinations (ws, push, ...)
