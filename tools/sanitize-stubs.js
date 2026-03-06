#!/usr/bin/env node
// sanitizes recorded WireMock stubs by removing non-deterministic fields
// from request body patterns so they match on replay.
const fs = require('fs');
const path = require('path');

const MAPPINGS_DIR = path.join(__dirname, '..', 'wiremock', 'mappings');
// fields that change per-request and break replay matching
const STRIP_FIELDS = ['ClientToken', 'clientToken', 'RequestToken', 'requestToken', 'IdempotencyToken'];

let modified = 0;
for (const file of fs.readdirSync(MAPPINGS_DIR)) {
  if (!file.endsWith('.json') || file.startsWith('proxy-')) {
    continue;
  }
  const filePath = path.join(MAPPINGS_DIR, file);
  const stub = JSON.parse(fs.readFileSync(filePath, 'utf-8'));
  const patterns = stub.request?.bodyPatterns;
  if (!Array.isArray(patterns)) {
    continue;
  }

  let changed = false;
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
  if (changed) {
    fs.writeFileSync(filePath, JSON.stringify(stub, null, 2) + '\n');
    modified++;
  }
}
if (modified) {
  console.log('sanitized ' + modified + ' stub(s)');
}
