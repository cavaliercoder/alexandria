class alexandria::vagrant (
  
) {  
    file { '/home/vagrant/gocode/src/alexandria' :
        ensure => 'link',
        target => '/vagrant/src',
    }
    
    file { '/etc/motd' :
        ensure  => file,
        owner   => 'root',
        group   => 'root',
        mode    => '0644',
        content => template('alexandria/vagrant/motd.erb')
    }
}