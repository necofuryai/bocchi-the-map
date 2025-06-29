#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

function collectSpecs(suites) {
  if (!suites) return [];
  return suites.flatMap(suite => [
    ...(suite.specs || []),
    ...collectSpecs(suite.suites)
  ]);
}

function parseE2EResults(resultsFile) {
  if (!fs.existsSync(resultsFile)) {
    console.log('No E2E results found, creating placeholder');
    const placeholder = {
      stats: { expected: 0, passed: 0, failed: 0, flaky: 0, skipped: 0 }
    };
    fs.writeFileSync(resultsFile, JSON.stringify(placeholder, null, 2));
  }

  const data = JSON.parse(fs.readFileSync(resultsFile, 'utf8'));
  const specs = collectSpecs(data.suites);

  const totalTests = specs.length || 0;
  const passedTests = specs.filter(spec => 
    Array.isArray(spec.tests) && spec.tests.every(test => test.status === 'passed')
  ).length || 0;
  const failedTests = specs.filter(spec => 
    Array.isArray(spec.tests) && spec.tests.some(test => test.status === 'failed')
  ).length || 0;

  const successRate = totalTests > 0 ? Math.round((passedTests / totalTests) * 100) : 0;

  return {
    totalTests,
    passedTests,
    failedTests,
    successRate
  };
}

function main() {
  const resultsFile = process.argv[2] || 'e2e-results.json';
  const results = parseE2EResults(resultsFile);

  // Output shell-friendly environment variable assignments
  console.log(`TOTAL_TESTS=${results.totalTests}`);
  console.log(`PASSED_TESTS=${results.passedTests}`);
  console.log(`FAILED_TESTS=${results.failedTests}`);
  console.log(`SUCCESS_RATE=${results.successRate}`);
}

if (require.main === module) {
  main();
}

module.exports = { parseE2EResults, collectSpecs };