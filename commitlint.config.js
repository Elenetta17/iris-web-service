// commitlint.config.js
module.exports = {
  parserPreset: {
    parserOpts: {
      // Standard conventional commit header pattern: type(scope?): subject
      headerPattern: /^(\w*)(?:\((.*)\))?: (.*)$/,
      headerCorrespondence: ['type', 'scope', 'subject'],
    },
  },
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      2,
      'always',
      ['feat', 'fix', 'docs', 'test', 'chore', 'refactor', 'perf'],
    ],
    'header-max-length': [2, 'always', 72],
    'body-empty': [2, 'never'],
    'body-min-length': [2, 'always', 20],
  },
  plugins: [
    {
      rules: {
        'body-contains-author-and-testing': (parsed, when, value) => {
          const body = parsed.body || '';
          const hasAuthor = /Author: .+ <.+@.+>/.test(body);
          const hasTesting = /Testing:/.test(body);
          if (!hasAuthor || !hasTesting) {
            return [false, "Body must include 'Author: name <email>' and 'Testing:' sections"];
          }
          return [true, ''];
        },
      },
    },
  ],
  rules: {
    'body-contains-author-and-testing': [2, 'always'],
  },
};
