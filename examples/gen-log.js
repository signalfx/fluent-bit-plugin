const process = require('process');
const count = process.argv.length > 2 ? process.argv[2] : 1;
const monthsAbbr = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

const startTime = Date.now() - count * 1000;
for (let i = 0; i < count; i ++) {
  let d = new Date(startTime + i * 1000);
  const time = d.getDate() + '/' + monthsAbbr[d.getMonth()] + '/' + d.getFullYear() +':' + d.getHours() + ':' + d.getMinutes() + ':' + d.getSeconds() + ' +0200';
  const log = getOneOfThree('foo bar', 'error', 'exception');
  const cluster = getOneOfThree('A', 'B', 'C');
  console.log(`{"time": "${time}", "log": "${log}", "ecs_cluster": "${cluster}", "container_name": "/ecs-container-a1b2c3"}`);
}

function getOneOfThree(one, two, three) {
  let n = Math.floor(100 * Math.random());
  return n < 20 ? one : n < 60 ? two : three;
}
