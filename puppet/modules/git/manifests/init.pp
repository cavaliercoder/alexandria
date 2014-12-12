class git (
    $version = $git::params::version,
    $package_url = $git::params::package_url,
    $package_sum1 = $git::params::package_sha1,
    $tmp_path = $git::params::tmp_path
) inherits git::params {
    if $osfamily != 'RedHat' { fail('Unsupported OS') }
    
    #
    # See: http://git-scm.com/book/en/v2/Getting-Started-Installing-Git
    #    
    if $::gitversion != $version {
        if $::gitinstalled {
          notify { "Upgrading Git from version ${::gitversion} to ${version}" : 
            withpath => true,
            before    => Exec['Install Git']
          }
        } else {
          notify { "Installing Git version ${version}" :
            withpath  => true,
            before    => Exec['Install Git']
          }
        }
        
        package { 'git' :
          ensure        => 'purged',
          allow_virtual => false
        }
        
        $prereqs = [
          'tar',
          'grep',
          'curl',
          'make',
          'gcc',
          'libcurl-devel',
          'expat-devel',
          'gettext-devel',
          'openssl-devel',
          'zlib-devel',
          'perl-ExtUtils-MakeMaker'
        ]
        
        package { $prereqs :
            ensure        => 'present',
            allow_virtual => false
        }
        
        exec { 'Download Git package' :
            cwd     => $tmp_path,
            command => "/usr/bin/curl -sLf --retry 5 -O ${package_url}",
            unless  => "/usr/bin/sha1sum -b ${tmp_path}/git-${version}.tar.gz | /bin/grep ${package_sha1}",
            require => Package[$prereqs]
        } ->
        exec { 'Extract Git sources':
            cwd     => $tmp_path,
            command => "/bin/tar -xzf ${tmp_path}/git-${version}.tar.gz",
            creates => "/tmp/git-${version}/configure",
            require => Package[$prereqs]
        } ->
        exec { 'Configure Git sources':
            cwd     => "/tmp/git-${version}",
            command => '/usr/bin/make configure && ./configure',
            require => Package[$prereqs]
        } ->
        exec { 'Build Git sources' :
            cwd     => "/tmp/git-${version}",
            command => '/usr/bin/make all',
            require => Package[$prereqs]
        } ->
        exec { 'Install Git' :
            cwd     => "/tmp/git-${version}",
            command => '/usr/bin/make install',
            require => Package[$prereqs]
        }
    }
}