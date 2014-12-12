# -*- mode: ruby -*-
# vi: set ft=ruby :

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
