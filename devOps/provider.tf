terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    swarm = {
      source = "aucloud/swarm"
      version = "~> 1.2"
    }
  }
}

variable "do_token" {}
variable "pvt_key" {}
variable "ssh_user" {
  description = "SSH user for connecting to swarm nodes"
  type        = string
  default     = "root"  # Default SSH user for DigitalOcean droplets
}

provider "digitalocean" {
  token = var.do_token
}
provider "swarm" {
  ssh_user = var.ssh_user
  ssh_key  = var.pvt_key
}

data "digitalocean_ssh_key" "RasmusMac" {
  name = "RasmusMac"
}