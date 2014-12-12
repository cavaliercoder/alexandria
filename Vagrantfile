# -*- mode: ruby -*-
# vi: set ft=ruby :

# Alexandria CMDB - Open source configuration management database
# Copyright (C) 2014  Ryan Armstrong <ryan@cavaliercoder.com>
# 
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
# 
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
# package controllers

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

$script = <<SCRIPT
puppet module install puppetlabs/stdlib
puppet module install puppetlabs/mongodb
SCRIPT

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "puppetlabs/centos-6.5-64-puppet"
 
  config.vm.network "forwarded_port", guest: 3000, host: 3000 

  config.vm.provision "shell", inline: $script
  config.vm.provision "puppet" do |puppet|
    puppet.module_path = "./puppet/modules"
    puppet.manifests_path = "./puppet/manifests"
    puppet.hiera_config_path = "./puppet/hiera.yaml"
    puppet.facter = {
      "vagrant" => "1"
    }
    # puppet.options = "--verbose --debug"
  end
end
