import WebSocket from 'ws';

const ws = new WebSocket("ws://localhost:18115");
let data = {txs: []};

ws.on('open', function () {
  console.log("CONNECTED");
  ws.send(
    '{"id": 2, "jsonrpc": "2.0", "method": "subscribe", "params": ["new_transaction"]}'
  );
});
ws.on('message', function (e) {
  // console.log(e);
  data.txs.push(e);
});
ws.on('error', function (e) {
  console.log("ERROR:" + e);
});
ws.on('close', function () {
  console.log("DISCONNECTED");
});

const sendData = async() => {
  try {
    await fetch(`${server_path}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: data
    });
    data.txs = [];
  } catch (err) {
    console.error(err.message);
  }
}

setInterval(async function() {
  await sendData();
}, 60000);
