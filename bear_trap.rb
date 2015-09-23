# Copyright (c) 2011, chrisbdaemon <chrisbdaemon@gmail.com>
# All rights reserved.
# 
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions are met:
# 
# Redistributions of source code must retain the above copyright notice, this
# list of conditions and the following disclaimer.
# 
# Redistributions in binary form must reproduce the above copyright notice,
# this list of conditions and the following disclaimer in the documentation
# and/or other materials provided with the distribution.
# 
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
# AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
# IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
# ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
# LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
# CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
# SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
# INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
# CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
# ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
# POSSIBILITY OF SUCH DAMAGE.

$:.unshift(File.dirname(__FILE__) + File::SEPARATOR + 'lib')

require 'yaml'
require 'getopt/long'

require 'Common'
require 'Trigger'
require 'AlertHandler'

class BearTrap
	
	@config = {}
	@threads = []
	
	def initialize( config )
		@config = config

		@threads = [ ]
	end
	
	def load_triggers
		@triggers = [ ]
		
		@config['triggers'].each do |t|
			puts_d "Trying to load trigger: #{t}"
			@triggers << Trigger.load( t, { :callback => self } )
		end
	end
	
	def load_alert_handlers
		@alert_handlers = [ ]
		
		@config[ 'alert_handlers' ].each do |ah|
			@alert_handlers << AlertHandler.load( ah )
		end
	end

	def got_alert( ip )
				
		# Hack needed when scan is being executed at high speeds to prevent
		# duplication
		if ! defined? @blocked_addresses
			@blocked_addresses = [ ]
		end
		
		# If address is already being blocked, ignore it
		if @blocked_addresses.include? ip
			return
		end
		
		# Add the ip recently blocked
		@blocked_addresses << ip
		
		# After 3 seconds, remove the address
		Thread.new {
			sleep($opt['t'])
			@blocked_addresses.delete ip
			@alert_handlers.each do |f|
				f.unblock_address( ip )
			end
		}
		
		@alert_handlers.each do |f|
			f.handle_alert( ip )
		end
		
	end
	
	def run
		@triggers.each do |t|
			@threads << t.set_trigger
		end
		
		@threads.each do |t|
			t.join
		end
	end
	
end

def usage
	$stderr.puts 'BearTrap v0.2-beta'
	$stderr.puts 'Usage: bear_trap.rb [-vd] -c <config file>'
	$stderr.puts 'OPTIONS:'
	$stderr.puts '  --config  -c <config file>    Filename to load configuration from [REQUIRED]'
	$stderr.puts '  --verbose -v                  Verbose output'
	$stderr.puts '  --debug   -d                  Debug output (very noisy)'
 	$stderr.puts '  --timeout -t                  Ban timeout in seconds'
	exit
end

$opt = Getopt::Long.getopts(
	['--config',  '-c', Getopt::REQUIRED],
	['--verbose', '-v', Getopt::BOOLEAN],
        ['--debug',   '-d'],
        ['--timeout', '-t', Getopt::OPTIONAL]
)

if $opt['c'] == nil
	usage()
end

begin
	if $opt['t'] != nil
		$opt['t'] = Integer($opt['t'])
	else
		$opt['t'] = 3
	end
	rescue ArgumentError
		$opt['t'] = 3
	end


config = YAML::load( File.open( $opt[ 'c' ] ) )

s = BearTrap.new( config )

s.load_triggers
s.load_alert_handlers

s.run
