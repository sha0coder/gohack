#!/usr/bin/expect -f


set host [lindex $argv 0]
set wordlist [lindex $argv 1]
set port 25
log_user 0

set fd [open $wordlist r]


while {[gets $fd user] >= 0} {


	spawn nc $host $port -vv

	expect "220 "
	send "helo $host\n"
	expect "250 "

		
	send "mail from: <test@$host>\n"

	expect "250 "
	send "rcpt to: <$user@$host>\n"

	expect {
		"250 " {
			send_user "$user exists\n"
		}

		"554 " {
			send_user "$user may exist\n"
		}

		"550 " {
			#send_user "$user don't exist\n"
		}

		"421 " {
			spawn nc $host $port -vv
			expect "220 "
			send "helo $host\n"
			expect "250 "
		}
	}


	send "quit\n"

	#sleep 1
	#send \003

	#send "rset\n"
}

expect eof