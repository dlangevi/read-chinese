const fetch = require('node-fetch-commonjs');

async function run(ipcCall, params) {
  const jsonArgs = JSON.stringify(params);
  const response = await fetch('http://localhost:3451/ipc', {
    method: 'Post',
    headers: {
      'Content-Type': 'application/json;charset=utf-8',
    },
    body: JSON.stringify({
      function: ipcCall,
      args: jsonArgs,
    }),
  });
  return response.json();
}

async function main() {
  const res = await run('hskWords', ['2.0', 1]);
  console.log(res);
}

main();
