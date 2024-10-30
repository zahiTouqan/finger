terraform {
    required_providers {
        docker = {
            source = "kreuzwerker/docker"
            version = "3.0.2"
        }
    }
}

provider "docker" {
    host = "unix:///var/run/docker.sock"
}

resource "docker_network" "private_network" {
    name = "finger-network"
}

resource "docker_volume" "finger_data" {
    name = "finger-data"
}

resource "docker_image" "server_image" {
    name = "finger-server"
    keep_locally = true
}

resource "docker_image" "client_image" {
    name = "finger-client"
    keep_locally = true
}

resource "docker_container" "server_container" {
    name = "server"
    image = docker_image.server_image.name
    restart = "unless-stopped"

    network_advanced {
        name = docker_network.private_network.name
    }
    
    ports {
        internal = 79
        external = 79
    }

    volumes {
        volume_name = docker_volume.finger_data.name
        container_path = "/app/data/server"
    }
}

resource "docker_container" "client_container" {
    name = "client"
    image = docker_image.client_image.name
    restart = "unless-stopped"

    network_advanced {
        name = docker_network.private_network.name
    }

    volumes {
        volume_name = docker_volume.finger_data.name
        container_path = "/app/data/client"
    }

    environment = [
        "SERVER_HOST=sever"
    ]
}