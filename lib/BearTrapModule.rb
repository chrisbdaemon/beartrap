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

class BearTrapModuleException < StandardError ; end
class BearTrapInvalidModuleException < BearTrapModuleException ; end
class BearTrapMissingParameterException < BearTrapModuleException ; end
class BearTrapInvalidParameterException < BearTrapModuleException ; end


class BearTrapModule
	
	@parameters = []
	
	@module_type = ''
	
	def check_parameters( config )
		
		puts_d "Checking parameters for #{self}"
		
		@parameters.each do |p|
			if p[ 'required' ] == true
				if ! config.keys.include? p[ 'name' ]
					raise BearTrapMissingParameterException, "#{self.class}: Required parameter '#{p[ 'name' ]}' not found."
				end
			end
				
		end
		
		config.each do |cfg_key, cfg_val|
			
			if cfg_key == 'type'
				next
			end
			
			invalid_param = true
			@parameters.each do |p|
				if cfg_key == p[ 'name' ]
					invalid_param = false
					break
				end
			end
			
			if invalid_param == true
				raise BearTrapInvalidParameterException, "#{self.class}: Encountered unknown parameter '#{cfg_key}'."
			end
		end
		
		puts_d "Parameters verified for #{self}"
		
		return true
	end
	
	def self.load( params, arguments=nil )
		class_name = params[ 'type' ] + @module_type
		filename = File.dirname( __FILE__ ) + File::SEPARATOR + @module_type + File::SEPARATOR + class_name + ".rb"
		
		puts_d "Loading #{filename}"
		
		if ! $".include? filename
			begin
				require filename
			rescue LoadError
				raise BearTrapInvalidModuleException, "#{self}: Unable to load module #{filename}."
			end
			puts_d "#{filename} loaded successfully"
		else 
			puts_d "#{filename} already loaded"
		end
		
		if arguments != nil
			return Kernel.const_get( class_name ).new( params, arguments )
		else
			return Kernel.const_get( class_name ).new( params )
		end
		
	end
end