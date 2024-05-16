/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Method to check whether the secret value is a nested object (returns true)
 * All other values return false
 * @param value string or stringified JSON
 * @returns boolean
 */
export function isAdvancedSecret(value) {
  try {
    const obj = typeof value === 'string' ? JSON.parse(value) : value;
    if (Array.isArray(obj)) return false;
    return Object.values(obj).any((value) => typeof value !== 'string');
  } catch (e) {
    return false;
  }
}

/**
 * Method to obfuscate all values in a map, including nested values and arrays
 * @param obj object
 * @returns object
 */
export function obfuscateData(obj) {
  if (typeof obj !== 'object' || Array.isArray(obj)) return obj;
  const newObj = {};
  for (const key of Object.keys(obj)) {
    if (Array.isArray(obj[key])) {
      newObj[key] = obj[key].map(() => '********');
    } else if (typeof obj[key] === 'object') {
      // unfortunately in javascript if the value of a key is null
      // calling typeof on this value will return object even if it is a string ex: { "test" : null }
      // this is due to a "historical js bug that will never be fixed"
      // we handle this situation here
      if (obj[key] === null) {
        newObj[key] = '********';
      } else {
        newObj[key] = obfuscateData(obj[key]);
      }
    } else {
      newObj[key] = '********';
    }
  }
  return newObj;
}
