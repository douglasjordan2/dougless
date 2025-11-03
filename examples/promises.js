// Test Promise functionality
console.log('Testing promises...');

// Basic promise
const basicPromise = new Promise((resolve, reject) => {
  setTimeout(() => {
    resolve('Basic promise resolved!');
  }, 50);
});

basicPromise.then(value => {
  console.log('then:', value);
}).catch(err => {
  console.log('catch:', err);
});

// Promise rejection
const rejectedPromise = new Promise((resolve, reject) => {
  setTimeout(() => {
    reject('Promise rejected!');
  }, 100);
});

rejectedPromise
  .then(value => {
    console.log('Should not print');
  })
  .catch(err => {
    console.log('Caught rejection:', err);
  });

// Promise chaining
Promise.resolve(5)
  .then(x => {
    console.log('Chain step 1:', x);
    return x * 2;
  })
  .then(x => {
    console.log('Chain step 2:', x);
    return x + 10;
  })
  .then(x => {
    console.log('Chain step 3:', x);
  });

// Promise.all
Promise.all([
  Promise.resolve(1),
  Promise.resolve(2),
  Promise.resolve(3)
]).then(values => {
  console.log('Promise.all result:', values);
});

// Promise.race
Promise.race([
  new Promise(resolve => setTimeout(() => resolve('slow'), 100)),
  new Promise(resolve => setTimeout(() => resolve('fast'), 10))
]).then(winner => {
  console.log('Promise.race winner:', winner);
});

// Promise.any
Promise.any([
  Promise.reject('error 1'),
  Promise.resolve('success!'),
  Promise.reject('error 2')
]).then(value => {
  console.log('Promise.any result:', value);
});

// Promise.allSettled
Promise.allSettled([
  Promise.resolve('success'),
  Promise.reject('failure'),
  Promise.resolve('another success')
]).then(results => {
  console.log('Promise.allSettled results:', results);
});

console.log('All promises created');
