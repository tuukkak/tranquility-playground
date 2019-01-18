var net = require('net');

var server = net.createServer(function(client) {
    console.log('Sever started')
    client.on('data', function (data) {
        console.log('Packet received!')
        client.write('OK\n');
    });

    client.on('timeout', function () {
        console.log('Client request time out.');
    })
});

server.listen(3002);