define go::package (
    $path = undef,
    $ensure = 'present',
    $build_flags = undef,
    $timeout = 300,
    $gouser = $go::params::gouser,
    $gouser_group = $go::params::gouser_group,
    $gouser_home = $go::params::gouser_home,
    $gopath = $go::params::gopath
) {
  require ::go
  require ::git
  
  $_path = $path ? { undef => $name, default => $path }
  $_goroot = $::goroot ? { undef => '/usr/local/go', default => $::goroot }
  
  validate_string($gouser)
  validate_string($gopath)
  validate_re($timeout, '^\d+$')
  
  case $ensure {
    present: {
      exec { "Install package ${name}" :
        user        => $gouser,
        cwd         => $gopath,
        environment => [
          "GOPATH=${gopath}",
          'PATH=/usr/local/bin:/usr/bin'  # Ensure source built git is used
        ],
        command     => "${_goroot}/bin/go get -x -v ${build_flags} ${_path}",
        creates     => "${gopath}/src/${_path}",
        timeout     => $timeout
      }
    }
    absent: {
      file { "${gopath}/src/${_path}" :
        ensure => 'absent',
        force  => true
      }
      
      file { "${gopath}/pkg/${_path}" :
        ensure => 'absent',
        force  => true
      }
    }
    default: {
      fail("Invalid ensure parameter: ${ensure}")
    }
  }
}