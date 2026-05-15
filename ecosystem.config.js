module.exports = {
  apps: [{
    name: "spamgo",
    script: "./main",
    autorestart: true,
    watch: false,
    max_restarts: 999,
    restart_delay: 2000,
    stop_exit_codes: [0],        // sirf exit(0) te BAND raho
    // exit(2) te pm2 dobara start karega
  }]
}
