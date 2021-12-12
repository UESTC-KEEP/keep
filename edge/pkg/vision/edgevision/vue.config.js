module.exports = {
  devServer: {
    host:"localhost",
    port:8080,
    https:false,
    open:true,
    proxy: {
      // 配置跨域
      '/api': {
        target: 'http://localhost:8081/',
        ws:true,
        changOrigin: true,
        pathRewrite: {
          '^/api': ''
        }
      }
    }
  },
}
