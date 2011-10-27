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
require 'openssl'

class SignedBroadcastDistributor < Distributor
	
	def initialize( params, arguments )
		@parameters = [ { 'name' => 'protocol', 'required' => true },
		                { 'name' => 'address', 'required' => true },
		                { 'name' => 'bind', 'required' => false },
		                { 'name' => 'port', 'required' => true },
		                { 'name' => 'ca_cert', 'required' => true },
		                { 'name' => 'cert', 'required' => true },
		                { 'name' => 'key', 'required' => true } ]
		
		self.check_parameters( params )
		
		@callback = arguments[ :callback ]
		@protocol = params[ 'protocol' ]
		@address  = params[ 'address' ]
		@bind     = params[ 'bind' ]
		@port     = params[ 'port' ]
		
		setup_signing( params )
	end
	
	def setup_signing( params )
		@ca_cert = OpenSSL::X509::Certificate.new( File.read( params[ 'ca_cert' ] ) )
		@cert    = OpenSSL::X509::Certificate.new( File.read( params[ 'cert' ] ) )
		@key     = OpenSSL::PKey::RSA.new( File.read( params[ 'key' ] ) )
		
		@cert_store = OpenSSL::X509::Store.new
		@cert_store.add_cert( @ca_cert )
	end
	
	def send_alert( message )
		
		data = [ message.length ].pack('L') + message
		
		puts_d "Message length before signing: #{data.size}"
		
		signature = OpenSSL::PKCS7::sign(@cert, @key, data, [ ], OpenSSL::PKCS7::BINARY|OpenSSL::PKCS7::DETACHED ).to_der
		
		signed_data = data + signature
		
		if @protocol == 'udp'
			self.send_alert_udp( signed_data )
		elsif @protocol == 'icmp'
			self.send_alert_icmp( signed_data )
		end
	end
	
	def start_listening
		
		puts_d "Binding UDP socket to #{@bind}:#{@port}"
		
		sock = UDPSocket.new
		sock.setsockopt( Socket::SOL_SOCKET, Socket::SO_BROADCAST, true )
		sock.bind( @bind, @port )
		
		thread = Thread.new {
			puts_d "Created thread for #{sock}"
			loop {
				signed_data, sender = sock.recvfrom( 65507 )
				
				if signed_data.length < 10
					continue
				end
								
				message_length = signed_data[ 0, 4 ].unpack( 'L' ).first + 4
				
				puts_d "Recieved message length: #{message_length}"
				
				data = signed_data[ 0, message_length ]
				
				signature = signed_data[ message_length, signed_data.size - message_length ]
				
				puts_d("Verifying #{data.length} byte message from #{sender}...")
				
				message = OpenSSL::PKCS7.new( signature )
				if message.verify( [ @ca_cert ], @cert_store, data, OpenSSL::PKCS7::BINARY|OpenSSL::PKCS7::DETACHED )
					puts_d("Message verification successful")
					
					ip = data[ 4, data.size ]
					
					@callback.got_alert( ip )
				else
					puts_v("Message verification from #{sender} failed!")
					continue
				end
			}
		}
		
		return thread
	end
	
	def send_alert_udp( message )
		sock = UDPSocket.new
		sock.setsockopt( Socket::SOL_SOCKET, Socket::SO_BROADCAST, true )
		
		puts_v "Sending alert to #{@address}:#{@port}"
		
		begin
			bytes_sent = sock.send( message, 0, @address, @port )
		rescue Errno::EMSGSIZE
			puts "Unable to send alert, the #{message.size} byte message is too large"
		end
		
		if bytes_sent != message.length
			stderr.puts "Only sent #{bytes_sent} out of #{message.length} to #{@address}:#{@port}"
		else
			puts_d "Successfully sent #{bytes_sent} bytes to #{@address}:#{@port}"
		end
		
		sock.close
	end
	
	def send_alert_icmp( message )
		
	end
	
end