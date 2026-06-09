module.exports = {
  apps: [
    {
      name: 'umbra-cms',
      script: './build-and-run.js',

      // Watch all source — rebuild & restart on any change
      watch: ['.'],
      ignore_watch: ['tmp/', '.git', 'node_modules', '.env'],
      watch_options: { followSymlinks: false },
      restart_delay: 1000,
      max_restarts: 10,

      // Env — .env file is loaded by loadEnv() at startup
      env: {
        PORT: '8080',
      },
    },
  ],
};
