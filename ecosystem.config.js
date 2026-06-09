module.exports = {
  apps: [
    {
      name: 'umbra-cms',
      script: './build-and-run.js',

      // Only watch source files — NOT tmp/ (avoid rebuild loop)
      watch: [
        '*.go',
        'templates/',
        'static/',
        'data/',
      ],
      ignore_watch: ['tmp/', '.git', 'node_modules', '.env'],
      watch_options: { followSymlinks: false },
      restart_delay: 1000,
      max_restarts: 10,

      env: {
        PORT: '8080',
      },
    },
  ],
};
