#!/usr/bin/env node

/*  
	NetIO by @sha0coder

	./netio.js [plain/tls] [host] [port] [nodes]

	press h for help

*/

(function(args) {

	var netio = {
		port: 80,
		host: '',
		isTls: false,
		conns: 1,
		buffer: [],
		wordlist: [],
		readCommand: null,
		certfile: 'cert.pem',
		keyfile: 'priv.pem',
		lineReader: require('line-reader'),
		clients: [],
		tls: require('tls'),
		net: require('net'),
		os: require('os'),
		fs: require('fs'),
		tlsCfg: {
			key: null,
			cert: null,
			//pfx: null,
			passphrase: '1234',
			requestCert: true,
			rejectUnauthorized: false
		}
	};

	netio.parseArgs = function(args) {
		console.log('NetIO by @sha0coder\npress h for help');

		if (args.length != 6) {
			netio.usage();
		} else {
			netio.isTls = false;
			if (args[2] == 'tls')
				netio.isTls = true;
			netio.main(args[3], args[4], args[5]);
		}
	}

	netio.usage = function() {
		console.log('./netio.js [plain/tls] [host] [port] [num of connections]');
	}

	netio.main = function(host,port,conns)  {
		netio.host = host;
		netio.port = port;
		netio.conns = conns;

		if (netio.isTls) 
			netio.initTLS();

		netio.initStdin();
		
		for (var i=0; i<conns; i++)
			netio.connect(i+1);
	}

	netio.initStdin = function() {
		process.stdin.resume();
		process.stdin.setEncoding('utf-8');
		process.stdin.on('data', netio.onstdin);
	}

	netio.initTLS = function() {
		print.start('Loading certificate');
		netio.tlsCfg.cert = netio.fs.loadFileSync(netio.certfile, encoding='ascii');
		netio.tlsCfg.key = netio.fs.loadFileSync(netio.keyfile, encoding='ascii');
		print.end(true);
	}

	netio.onstdin = function(chunk) {
		var data = chunk.toString();

		switch(data[0]) {
			case '>': // raw send
				netio.sendAll(data.substring(1).replace(/\\n/,'\n'));
				break;

			case 'r': // repeated send
				var spl = data.split(' ');
				times = parseInt(spl[1]); // quality check
				spl.splice(0,2);
				data = spl.join(' ');
				data = data.replace(/\\n/,'\n');
				for (var i=0; i<times; i++)
					netio.sendAll(data);
				break;

			case 'h':
				netio.showCommands();
				break;

			case 'l':
				var spl = data.split(' ');
				netio.load(spl[1].replace(/\n/,''));
				break;

			case 'q':
				netio.quit();
				break;

			case 'f':
				var spl = data.split(' ');
				spl.splice(0,1);
				data = spl.join(' ');
				data = data.replace(/\\n/,'\n');
				netio.doFuzz(data);
				break;

			case 'b':
				var spl = data.split(' ');
				spl.splice(0,1);
				data = spl.join(' ');
				data = data.replace(/\\n/,'\n');
				netio.doBrute(data);
				break;
		}
	}

	netio.sendAll = function(data) {
		netio.clients.forEach(function(client) {
			if (client)
				client.write(data);
			else 
				console.log('trying to send to a closed connection');
		});
	}

	netio.tlsConnect = function(node) {
		var client = netio.tls.connect(netio.port, netio.host, netio.tlsCfg, function() {
			if (!client.authorized) {
					console.log('Cant stablish tls session :(');
					process.exit(1); // should close the sock?
			}
			console.log('node %d connected.',node);
		});
		return client;
	}

	netio.plainConnect = function(node) {
		var client = netio.net.connect({host:netio.host, port:netio.port}, function() {
			console.log('node %d connected.',node);
		});
		return client;
	}

	netio.connect = function(node) {
		var client;

		if (netio.isTls)
			client = netio.tlsConnect(node);
		else 
			client = netio.plainConnect(node);

		client.setEncoding('utf-8');
		netio.clients.push(client);

		if (!client.on)
			client.on = client.addListener;

		client.on('data', function(data) {
			console.log('<<<<%d<<<<<',node);
			console.log(data.toString());
			console.log('-----------');
		});

		client.on('end', function() {
			console.log('node %d disconnected!',node);
		});
	}

	netio.showCommands = function() {
		console.log('commands:');
		console.log('>SEND\\nRAW\\nDATA\\n');
		console.log('r 10 SEND\\nTHIS\\n10\\nTIMES');
		console.log('h  (this help)');
		//console.log('f FUZZ\\nHERE: ##   (bofs,fmts,sql,xss,...)');
		console.log('l wordlist.txt    (load wordlist)');
		console.log('b BRUTEFORCE\\nHERE: ## ');
		console.log('q (disconnect all nodes and quit)');
	}

	netio.quit = function() {
		netio.clients.forEach(function(client) {
			client.end();
		});
		console.log('closing connections');
	}

	netio.load = function(filepath) {
		netio.lineReader.eachLine(filepath, function(line, last) {
			if (last) {
				console.log('%d words loaded!',netio.wordlist.length);
			} else {
				netio.wordlist.push(line);
			}
		});	
	}

	netio.getReadyClient = function() {
		//todo: get a ready client instead of random
		var node = parseInt(Math.random()*netio.clients.length);
		return netio.clients[node];
	}

	netio.doFuzz = function(data) {
		netio.todo();
	}

	netio.doBrute = function(data) {
		netio.wordlist.forEach(function(word) {
			//process.stdout.write('>%s',data.replace(/##/,word));
			netio.getReadyClient().write(data.replace(/##/,word));
		});
	}

	netio.todo = function() {
		console.log('not implemented by now');
	}

	netio.parseArgs(args);

})(process.argv);

