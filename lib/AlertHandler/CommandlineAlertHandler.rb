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

class CommandlineAlertHandler < AlertHandler

	def initialize( params )

		@parameters = [ { 'name' => 'block_command', 'required' => true },
		                { 'name' => 'unblock_command', 'required' => true },
		                { 'name' => 'action_command', 'required' => false },
		                { 'name' => 'regexp', 'required' => false } ]

		self.check_parameters( params )

		@block_command    = params[ 'block_command' ]
		@unblock_command  = params[ 'unblock_command' ]

		if params[ 'action_command' ] != nil
			@action_command = params[ 'action_command' ]
		end		

		if params[ 'regexp' ] != nil
			regexp_str = params[ 'regexp' ]
		else
			regexp_str = '[^\.0-9]'
		end

		@regexp = Regexp.new( regexp_str )

	end

	def handle_alert( address )

		block_command = self.build_command( @block_command, address )
		action_command = self.build_command( @action_command, address )

		if defined? @action_command
	
			puts "Command: #{action_command}"
			`#{action_command}`

			unless $?.success?
				puts "Command failed with status #{$?.exitstatus}"
			end		

		end

		puts "Command: #{block_command}"
		`#{block_command}`

		unless $?.success?
			puts "Command failed with status #{$?.exitstatus}"
		end		
	end

	def unblock_address( address )

		command = self.build_command( @unblock_command, address )

		puts "Command: #{command}"
		`#{command}`

		unless $?.success?
			puts "Command failed with status #{$?.exitstatus}"
		end

	end

	def build_command( command, address )

		if command.match( /\$RAW_IP/ ) != nil
			return command.gsub( /\$RAW_IP/, address )
		elsif command.match( /\$IP/ ) != nil
			return command.gsub( /\$IP/, address.gsub( @regexp, '' ) )
		end

	end

end
