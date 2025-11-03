// Test file system operations
console.log('Testing file operations...');

const testFile = '/tmp/dougless-test.txt';
const testDir = '/tmp/dougless-test-dir/';

// Test write file
console.log('Writing file...');
files.write(testFile, 'Hello from Dougless!').then(() => {
  console.log('File written successfully');
  
  // Test read file
  return files.read(testFile);
}).then((err, data) => {
  if (err) {
    console.log('Error reading file:', err);
  } else {
    console.log('File contents:', data);
  }
  
  // Test create directory
  console.log('Creating directory...');
  return files.write(testDir);
}).then((err) => {
  if (err) {
    console.log('Error creating directory:', err);
  } else {
    console.log('Directory created');
  }
  
  // Test read directory
  return files.read(testDir);
}).then((err, entries) => {
  if (err) {
    console.log('Error reading directory:', err);
  } else {
    console.log('Directory contents:', entries);
  }
  
  // Cleanup - delete file
  return files.rm(testFile);
}).then((err) => {
  if (err) {
    console.log('Error deleting file:', err);
  } else {
    console.log('File deleted');
  }
  
  // Cleanup - delete directory
  return files.rm(testDir);
}).then((err) => {
  if (err) {
    console.log('Error deleting directory:', err);
  } else {
    console.log('Directory deleted');
  }
  
  console.log('File operations complete!');
}).catch(err => {
  console.log('Unexpected error:', err);
});

console.log('File operations started...');
