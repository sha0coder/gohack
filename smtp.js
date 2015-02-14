#!/usr/bin/env node

// by @sha0coder

(function(args) {

	var smtp = {
		host: '',
		port: 25,
		nodes: 2,
		queue: [],
		net: require('net')
	};

	smtp.start = function(n) {
		this.nodes = n;
		for (i=0; i<n; i++) {
			smtp.node(i);
		}
	}

	smtp.node = function(n) {
		var client = smtp.net.connect({host:smtp.host, port:smtp.port}, function() {
			console.log('%d connected.',n);
		});
		client.on('data', function(data) {
			console.log('data: %s',data.toString());
			if (smtp.queue.length>0) {
				var w = smtp.queue.pop();
				console.log('>>%s',w);
				client.write('VRFY '+w+'\n');
			} else 
				client.write('quit\n');
				return;
		});
	}

	smtp.check = function(w) {
		console.log('checking %s',w);
	}

	function main(wordlist,host) {
		smtp.host = host;

		var lineReader = require('line-reader');
		lineReader.eachLine(wordlist, function(line, last) {
			if (last) {
				smtp.start(2);

			} else {
				smtp.queue.push(line);
			}
		});	
	}

	main(args[2],args[3]);

})(process.argv);

