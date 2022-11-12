const http = require("http")

const server = http.createServer((req, res) => {
    res.setHeader("Access-Control-Allow-Origin", '*');
});

server.on('request', (req, res) => {
    res.end("SERVER:I received");
})

server.listen(3000);
console.log("服务器已启动，监听3000端口，请使用localhost:3000访问");
