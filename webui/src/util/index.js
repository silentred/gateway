export default {
    camelToSnake: camelToSnake,
    snakeToCamel: snakeToCamel
}

function isObject(obj) {
  return obj !== null && typeof obj === 'object';
}

function isNumber(n) {
  return !isNaN(parseFloat(n)) && isFinite(n);
}

function camelToSnake(data, depth) {
  if (isObject(data)) {
    if (typeof depth === 'undefined') {
      depth = 1;
    }
    return _processKeys(data, _snakelize, depth);
  } else {
    return _snakelize(data);
  }
};

function snakeToCamel(data, depth) {
  if (isObject(data)) {
    if (typeof depth === 'undefined') {
      depth = 1;
    }
    return _processKeys(data, _camelize, depth);
  } else {
    return _camelize(data);
  }
};

// snakelize a string formed in underscore
function _snakelize(key) {
  let separator = '_';
  let split = /(?=[A-Z])/;

  return key.split(split).join(separator).toLowerCase();
}

// camelize a string formed in underscore
function _camelize(key) {
  if (isNumber(key)) {
    return key;
  }
  key = key.replace(/[\-_\s]+(.)?/g, function(match, ch) {
    return ch ? ch.toUpperCase() : '';
  });
  // Ensure 1st char is always lowercase
  return key.substr(0, 1).toLowerCase() + key.substr(1);
}

// camelize/snakelize keys of an object
// @param {number} depth to which level of keys should it process
function _processKeys(obj, processer, depth) {
  if (depth === 0 || !isObject(obj)) {
    return obj;
  }

  let result = {};
  let keys = Object.keys(obj);

  for (let i = 0; i < keys.length; i++) {
    result[processer(keys[i])] = _processKeys(obj[keys[i]], processer, depth - 1);
  }

  return result;
}