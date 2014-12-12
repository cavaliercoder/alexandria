include ::go
include ::git

include ::alexandria::vagrant

$packages = [
  'github.com/go-martini/martini',      # Martini application stack
  'github.com/martini-contrib/binding', # Martini request binding
  'github.com/revel/cmd/revel',         # Revel MVC application stack
  'gopkg.in/mgo.v2',                    # MongoDB driver
  'gopkg.in/redis.v2',                  # Redis driver
  'github.com/codegangsta/gin',         # Gin live reload
  'github.com/codegangsta/cli'		# CLI framework
]

go::package { $packages :
  ensure  => 'present'
}

class {'::mongodb::globals':
  manage_package_repo => true,
}->
class {'::mongodb::server': }->
class {'::mongodb::client': }
