const process = require('process');
const count = process.argv.length > 2 ? process.argv[2] : 1;
const monthsAbbr = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

for (let i = 0; i < count; i ++) {
  let d = new Date(Date.now() - 1000 * count + i * 1000);
  const time = d.getDate() + '/' + monthsAbbr[d.getMonth()] + '/' + d.getFullYear() +':' + d.getHours() + ':' + d.getMinutes() + ':' + d.getSeconds() + ' +0200';
  const log = getOneOfThree(i, 'foo bar', 'error', 'exception');
  const cluster = getOneOfThree(Math.floor(i / 3), 'A', 'B', 'C');
  console.log(`{"time": "${time}", "log": "${log}", "ecs_cluster": "${cluster}", "container_name": "/ecs-container-a1b2c3"}`);
}

function getOneOfThree(i, one, two, three) {
  let n = Math.floor(100 * Math.random());
  return n < 20 ? one : n < 60 ? two : three;
}
