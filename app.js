import WebSocket from 'ws';

const ws = new WebSocket("ws://localhost:18115");

ws.onopen = function (evt) {
  console.log("CONNECTED");
  ws.send(
    '{"id": 2, "jsonrpc": "2.0", "method": "subscribe", "params": ["new_transaction"]}'
  );
};
ws.onmessage = function (evt) {
  console.log(evt.data);
};
ws.onerror = function (evt) {
  console.log("ERROR:" + evt);
};
ws.onclose = function (evt) {
  console.log("DISCONNECTED");
};