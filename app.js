import WebSocket from 'ws';
import axios from "axios";

const node_id = 1; // id of this node stored in psql
const ws = new WebSocket("ws://localhost:18115");
let data = {txs: []};

ws.on('open', function () {
  console.log("CONNECTED");
  ws.send(
    '{"id": 2, "jsonrpc": "2.0", "method": "subscribe", "params": ["new_transaction"]}'
  );
});
ws.on('message', function (e) {
  const res = JSON.parse(e);
  if (res.result) return;

  const tx = JSON.parse(res.params.result).transaction;
  const tx_hash = tx.hash;
  const unix_timestamp = Date.now();
  data.txs.push({node_id, tx_hash, unix_timestamp});
});
ws.on('error', function (e) {
  console.log("ERROR:" + e);
});
ws.on('close', function () {
  console.log("DISCONNECTED");
});

const sendData = async() => {
  try {
    await axios.post("http://localhost:3006/txs", data).then(data.txs = []);
  } catch (err) {
    console.error(err.message);
  }
}

setInterval(async function() {
  await sendData();
}, 60000);
