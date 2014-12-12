class go (
    $version = $go::params::version,
    $package_url = $go::params::package_url,
    $package_sha1 = $go::params::package_sha1,
    $tmp_path = $go::params::tmp_path,
    $install_root = $go::params::install_root,
    $gouser = $go::params::gouser,
    $gouser_group = $go::params::gouser_group,
    $gouser_home = $go::params::gouser_home,
    $gopath = $go::params::gopath
) inherits go::params {
    require git
    
    if $osfamily != 'RedHat' { fail('Unsupported OS') }
    
    #
    # See: https://golang.org/doc/install
    #
    $prereqs = [ 'mercurial' ]
    
    package { $prereqs :
        ensure        => 'present',
        allow_virtual => false
    }
    
    if $::goversion != "go${version}" {
        exec { 'DownloadGoPackage' :
            cwd     => $tmp_path,
            command => "/usr/bin/curl -sLf --retry 5 -O ${package_url}",
            unless  => "/usr/bin/sha1sum -b ${tmp_path}/go${version}.linux-amd64.tar.gz | /bin/grep ${package_sha1}"
        } ~>
        exec { 'ExtractGoSources':
            cwd     => $install_root,
            command => "/bin/tar -xzf ${tmp_path}/go${version}.linux-amd64.tar.gz",
            creates => "${install_root}/go/bin/go"
        }
    }
    
    file { '/etc/profile.d/go.sh' :
        ensure  => 'present',
        content => "export GOROOT=${install_root}/go\nexport PATH=\$PATH:\$GOROOT/bin\n"
    }
    
    #
    # See: https://golang.org/doc/code.html#GOPATH
    #
    if $gouser != undef {
      validate_string($gouser_home)
      validate_string($gopath)
      
      if $gouser_group == undef { $gouser_group = $gouser }
      
      file { $gopath :
        ensure => 'directory',
        owner  => $gouser,
        group  => $gouser_group
      }
      
      file { "${gouser_home}/.bashrc":
        ensure => 'file',
        owner  => $gouser,
        group  => $gouser_group
      }
      
      file_line { 'GoPathEnvVar' :
        line    => "export GOPATH=${gopath}",
        path    => "${gouser_home}/.bashrc",
        require => File["${gouser_home}/.bashrc"]
      }
      
      file_line { 'GoPathBinEnvVar' :
        line    => "export PATH=\$PATH:${gopath}/bin",
        path    => "${gouser_home}/.bashrc",
        require => File["${gouser_home}/.bashrc"]
      }
    }
}