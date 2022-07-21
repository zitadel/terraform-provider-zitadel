module.exports = {
    branches: [
        {name: 'main'},
        {name: 'alpha', prerelease: true},
    ],
    plugins: [
        "@semantic-release/commit-analyzer"
    ]
};
