[![Build Status](https://travis-ci.org/RutledgePaulV/kube-vulcano.svg)](https://travis-ci.org/RutledgePaulV/kube-vulcano)

### Kube-Vulcano

A simple GO process as a docker container for providing realtime service registration/deregistration to an external 
vulcand load balancer for a kubernetes cluster. Based off of 
[Kelsey Hightower's Motorboat](https://github.com/kelseyhightower/motorboat) which does the same thing except for 
nginx plus.


### Why use it
Kubernetes has built-in support for external load balancers at some cloud providers (GCE & AWS), but I like to run
my clusters at digitalocean.


### License
This project is licensed under [MIT license](http://opensource.org/licenses/MIT).
