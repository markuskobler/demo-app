provider "google" {
  credentials = "./config/demo.json"
  project     = "desource-demo"
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
ExecStart=/usr/bin/docker run --name %p -p 8888:80 quay.io/markus/demo-app

ExecStop=/usr/bin/docker stop %p

[Install]
WantedBy=multi-user.target
EOF
}

resource "google_dns_record_set" "demo" {
  managed_zone = "distinctiveco"
  project = "desource-net"
  name = "demo.distinctive.co."
  type = "A"
  ttl  = 10
  rrdatas = ["${google_compute_instance.demo.network_interface.0.access_config.0.assigned_nat_ip}"]
}
