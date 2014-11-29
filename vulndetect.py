#!/usr/bin/python

class VulnDetect:

	def n(self,value):
		return True

v = VulnDetect()
#v.setUrl('')

sqli = v.setTests({
	'sqli': lambda x:  v.n("'") == v.n("'''") and v.n("''") == v.n("''''") and v.n("''") == v.n("'+'") and v.n("''") != v.n("'"),
})



if sqli(1):
	print "yes"