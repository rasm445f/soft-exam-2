# Create manager node
resource "digitalocean_droplet" "manager" {
  count  = 1
  image  = "docker-20-04"
  name   = "swarm-manager-${count.index + 1}"
  region = "fra1"
  size   = "s-1vcpu-1gb"
  tags   = ["manager"]
  ssh_keys = [
    data.digitalocean_ssh_key.RasmusMac.id
  ]
   provisioner "file" {
    source      = "docker-compose.yaml"
    destination = "/srv/docker-compose.yml"

     connection {
      type        = "ssh"
      user        = "root"
      host        = self.ipv4_address
      private_key = file(var.pvt_key)
    }
  }

  provisioner "remote-exec" {
    connection {
      host        = self.ipv4_address
      user        = "root"
      type        = "ssh"
      private_key = file(var.pvt_key)
      timeout     = "2m"
    }

    inline = ["echo 'SSH is ready'"]
  }
}

# Create worker nodes
resource "digitalocean_droplet" "worker" {
  count  = 2
  image  = "docker-20-04"
  name   = "swarm-worker-${count.index + 1}"
  region = "fra1"
  size   = "s-1vcpu-1gb"
  tags   = ["worker"]
  ssh_keys = [
    data.digitalocean_ssh_key.RasmusMac.id
  ]

  provisioner "remote-exec" {
    connection {
      host        = self.ipv4_address
      user        = "root"
      type        = "ssh"
      private_key = file(var.pvt_key)
      timeout     = "2m"
    }

    inline = ["echo 'SSH is ready'"]
  }
}


# Configure Swarm cluster
resource "swarm_cluster" "cluster" {
  skip_manager_validation = true
  dynamic "nodes" {
    for_each = concat(digitalocean_droplet.manager, digitalocean_droplet.worker)
    content {
        hostname = nodes.value.name
        tags = tomap({
          "role"   = sort(nodes.value.tags)[0]
        })
        public_address  = nodes.value.ipv4_address
        private_address = nodes.value.ipv4_address_private
    }
  }
  lifecycle {
    prevent_destroy = false
  }
}

# resource "digitalocean_droplet" "www-1" {
#   image = "ubuntu-20-04-x64"
#   name = "www-1"
#   region = "fra1"
#   size = "s-1vcpu-512mb-10gb"
#   ssh_keys = [
#     data.digitalocean_ssh_key.RasmusMac.id
#   ]
  
#   connection {
#     host = self.ipv4_address
#     user = "root"
#     type = "ssh"
#     private_key = file(var.pvt_key)
#     timeout = "2m"
#   }
  
#   provisioner "remote-exec" {
#     inline = [
#       "export PATH=$PATH:/usr/bin",
#       # install nginx
#       "sudo apt update",
#       "sudo apt install -y nginx"
#     ]
#  }
# }

