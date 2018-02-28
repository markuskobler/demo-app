variable "project" { type = "string" }
variable "dns_project" { type = "string" }
variable "dns_zone" { type = "string" }
variable "dns_name" { type = "string" }

provider "google" {
  credentials = "./config/demo.json"
  project     = "${var.project}"
  region      = "us-east1"
}

resource "google_compute_instance" "demo" {
  name         = "demo"
  machine_type = "f1-micro"
  zone         = "us-east1-d"

  tags = ["demo"]

  boot_disk {
    initialize_params {
      image = "coreos-cloud/coreos-stable"
    }
  }

  network_interface {
    network = "default"

    access_config {}
  }

  metadata {
    "user-data" = "${data.ignition_config.demo.rendered}"
  }
}

data "ignition_config" "demo" {
  systemd = [
    "${data.ignition_systemd_unit.demo.id}",
  ]
}

data "ignition_systemd_unit" "demo" {
  name   = "demo.service"
  enabled = true
  content = <<EOF
[Unit]
Description=demo
Requires=docker.service
After=network-online.target

[Service]
Restart=always

ExecStartPre=-/usr/bin/docker pull quay.io/markus/demo-app
ExecStartPre=-/usr/bin/docker rm %p
ExecStart=/usr/bin/docker run --name %p -p 80:8888 quay.io/markus/demo-app

ExecStop=/usr/bin/docker stop %p

[Install]
WantedBy=multi-user.target
EOF
}

resource "google_dns_record_set" "demo" {
  managed_zone = "${var.dns_zone}"
  project = "${var.dns_project}"
  name = "${var.dns_name}"
  type = "A"
  ttl  = 10
  rrdatas = ["${google_compute_instance.demo.network_interface.0.access_config.0.assigned_nat_ip}"]
}
