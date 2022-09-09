console.log('hi');

const obj = {};
function meta(funs) {
  funs.forEach(([name]) => {
    obj[name] = (...args) => {
      console.log(name, ...args);
    };
  });
}

meta([
  ['hello', 2],
  ['hi', 3],
  ['cool', 0],
]);

console.log(obj);
obj.hello(1, 2, 3);
