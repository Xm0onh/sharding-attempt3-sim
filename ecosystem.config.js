module.exports = {
    apps: [
      {
        name: 'frontend',
        cwd: './sharding-simulator',
        script: 'npm',
        args: 'run dev',
        env: {
          PORT: 3000,
        },
      },
      {
        name: 'backend',
        cwd: './',
        script: 'go',
        args: 'run .',
        env: {
          PORT: 8080,
        },
      }
	],
}
