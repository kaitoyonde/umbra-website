module.exports = {
    apps: [
      {
        name: "go-air-app",
        script: "~/go/bin/air",            // PM2 executes the 'air' tool
        interpreter: "none",      // Tells PM2 it's a binary executible, not a JS file
        merge_logs: true,         // Aggregates standard logs
        env: {
          NODE_ENV: "production",
          // Ensure your system Go paths are accessible by PM2
          PATH: process.env.PATH + ":" + process.env.GOPATH + "/bin"
        }
      }
    ]
  };