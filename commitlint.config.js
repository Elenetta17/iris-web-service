// commitlint.config.js
module.exports = {
  extends: ['@commitlint/config-conventional'], // Use conventional config as base
  rules: {
    // Enforce commit message type (e.g. feat, fix, etc.)
    'type-enum': [
      2,
      'always',
      ['feat', 'fix', 'docs', 'test', 'chore', 'refactor', 'perf']
    ],

    // Enforce commit message header length
    'header-max-length': [2, 'always', 72],

    // Custom rule: ensure the commit message includes an 'Author:' line
    'body-pattern': [
      2,
      'always',
      /^(?=.*\bAuthor: .+ <.+@.+>\b)(?=.*\bTesting:\b).+/m
    ],
  },
}
