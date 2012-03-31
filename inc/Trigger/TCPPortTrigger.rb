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

require 'socket'

class TCPPortTrigger < Trigger
	
	attr_accessor :port, :address, :server
	
	def initialize( params, arguments )
		
		@parameters = [ { 'name' => 'address', 'required' => false },
		                { 'name' => 'port', 'required' => true } ]
		
		self.check_parameters( params )
		
		@callback = arguments[ :callback ]
		@port     = params[ 'port' ]
		@address  = ( params[ 'address' ] != nil ) ? params[ 'address' ] : '0.0.0.0'
	end
	
	def set_trigger
		
		puts_d "Binding TCP socket to: #{@address}:#{@port}"
		
		@server = TCPServer.open( @address, @port )
		thread = Thread.new {
			puts_d "Created thread for #{@server}"
			loop {
				begin
					client = @server.accept
				rescue Exception => e
					puts e
					puts e.message
				end
				
				# When scanned at high speeds, connection doesn't stay open long
				# enough to get ip address, exception needs to be caught.
				begin
					ip = client.peeraddr.last
					@callback.got_alert ip
				rescue Errno::EINVAL
					puts_d "Unable to retrieve peer's IP address, socket closed too soon."
				end
					
				client.close
			}
		}
		
		return thread
	end
end

#$trigger_types['tcp_port'] = TCPPortTrigger