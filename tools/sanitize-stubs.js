#!/usr/bin/env node
// sanitizes recorded WireMock stubs by removing non-deterministic fields
// from request body patterns and stripping sensitive data from responses.
const fs = require('fs');
const path = require('path');

const WIREMOCK_DIR = path.join(__dirname, '..', 'wiremock');
const MAPPINGS_DIR = path.join(WIREMOCK_DIR, 'mappings');
const FILES_DIR = path.join(WIREMOCK_DIR, '__files');

// fields that change per-request and break replay matching
const STRIP_FIELDS = ['ClientToken', 'clientToken', 'RequestToken', 'requestToken', 'IdempotencyToken'];

// response headers that leak token metadata
const STRIP_HEADERS = [
  'github-authentication-token-expiration',
  'X-OAuth-Scopes',
  'X-GitHub-Request-Id',
];

// sanitize mapping files
let modifiedMappings = 0;
for (const file of fs.readdirSync(MAPPINGS_DIR)) {
  if (!file.endsWith('.json') || file.startsWith('proxy-')) {
    continue;
  }
  const filePath = path.join(MAPPINGS_DIR, file);
  const stub = JSON.parse(fs.readFileSync(filePath, 'utf-8'));
  let changed = false;

  // strip non-deterministic request body fields
  const patterns = stub.request?.bodyPatterns;
  if (Array.isArray(patterns)) {
    for (const pattern of patterns) {
      if (!pattern.equalToJson) {
        continue;
      }
      try {
        const body = JSON.parse(pattern.equalToJson);
        for (const field of STRIP_FIELDS) {
          if (field in body) {
            delete body[field];
            changed = true;
          }
        }
        if (changed) {
          pattern.equalToJson = JSON.stringify(body);
        }
      } catch {}
    }
  }

  // strip sensitive response headers
  const headers = stub.response?.headers;
  if (headers) {
    for (const h of STRIP_HEADERS) {
      if (h in headers) {
        delete headers[h];
        changed = true;
      }
    }
  }

  if (changed) {
    fs.writeFileSync(filePath, JSON.stringify(stub, null, 2) + '\n');
    modifiedMappings++;
  }
}

// sanitize response body files (strip PII from user objects)
let modifiedBodies = 0;
for (const file of fs.readdirSync(FILES_DIR)) {
  if (!file.endsWith('.json')) {
    continue;
  }
  const filePath = path.join(FILES_DIR, file);
  const raw = fs.readFileSync(filePath, 'utf-8');
  let data;
  try {
    data = JSON.parse(raw);
  } catch {
    continue;
  }

  let changed = false;

  // recursively walk the object and sanitize user-like nodes
  function walk(obj) {
    if (obj === null || typeof obj !== 'object') {
      return;
    }
    if (Array.isArray(obj)) {
      obj.forEach(walk);
      return;
    }
    // detect user-like objects (have login field)
    if ('login' in obj) {
      if ('email' in obj && obj.email !== '') {
        obj.email = '';
        changed = true;
      }
      if ('name' in obj && obj.name !== null && obj.name !== '') {
        obj.name = obj.login;
        changed = true;
      }
      if ('company' in obj && obj.company !== null && obj.company !== '') {
        obj.company = '';
        changed = true;
      }
    }
    for (const val of Object.values(obj)) {
      walk(val);
    }
  }

  walk(data);
  if (changed) {
    fs.writeFileSync(filePath, JSON.stringify(data) + '\n');
    modifiedBodies++;
  }
}

if (modifiedMappings || modifiedBodies) {
  console.log(`sanitized ${modifiedMappings} mapping(s) and ${modifiedBodies} body file(s)`);
}
