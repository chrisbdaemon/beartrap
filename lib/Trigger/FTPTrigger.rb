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

require 'Trigger'
require 'socket'

class FTPTrigger < Trigger

	def initialize( params, arguments )

		@parameters = [ { 'name' => 'banner', 'required' => true },
		                { 'name' => 'address', 'required' => false },
		                { 'name' => 'port', 'required' => false } ]

		self.check_parameters( params )

		@callback = arguments[ :callback ]
		@banner   = params[ 'banner' ]
		@port     = ( params[ 'port' ] != nil ) ? params[ 'port' ].to_i : 21
		@address  = ( params[ 'address' ] != nil ) ? params[ 'address' ] : '0.0.0.0'
	end

	def set_trigger

		puts_d "Binding FTP server to #{@address}:#{@port}"

		@server = TCPServer.open( @address, @port )
		thread = Thread.new {
			puts_d "Created thread for #{@server}"
			loop {

				Thread.start(@server.accept) do |client|

					client.send( "220 #{@banner}\x0d\x0a", 0 )

					loop do
						data = client.recv(100)

						# Username being sent for the first time
						if data.match(/^USER /)

							# Trigger an alert
							begin
								@callback.got_alert client.peeraddr.last
							ensure
								client.close
							end

							break
						else

							# Error on any other command
							client.send( "530\x0d\x0a", 0 ) # Not logged in
						end
					end
				end
			}
		}

		return thread
	end
end
