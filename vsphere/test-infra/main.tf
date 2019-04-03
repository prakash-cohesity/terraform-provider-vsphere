variable "plan" {
  default = "c1.xlarge.x86"
}

variable "esxi_version" {
  default = "6.5"
}

variable "facility" {
  default = "ams1"
}

variable "ovftool_url" {
  description = "URL from which to download ovftool"
}
variable "vcsa_iso_url" {
  description = "URL from which to download VCSA ISO"
}

locals {
  esxi_ssl_cert_thumbprint_path = "ssl_cert_thumbprint.txt"
}

provider "packet" {
}

resource "packet_project" "test" {
  name = "Terraform Acc Test vSphere"
}

data "packet_operating_system" "helper" {
  name             = "CentOS"
  distro           = "centos"
  version          = "7"
  provisionable_on = "t1.small.x86"
}

data "local_file" "esxi_thumbprint" {
  filename   = "${path.module}/${local.esxi_ssl_cert_thumbprint_path}"
  depends_on = ["packet_device.esxi"]
}

resource "random_string" "password" {
  length           = 16
  special          = true
  min_lower        = 1
  min_numeric      = 1
  min_upper        = 1
  min_special      = 1
  override_special = "@_"
}

data "template_file" "vcsa" {
  template = "${file("vcsa-template.json")}"
  vars = {
    esxi_host                = "${packet_device.esxi.access_public_ipv4}"
    esxi_username            = "root"
    esxi_password            = "${packet_device.esxi.root_password}"
    esxi_ssl_cert_thumbprint = "${chomp(data.local_file.esxi_thumbprint.content)}"
    ipv4_address             = "${cidrhost(format("%s/%s", packet_device.esxi.network[0].gateway, packet_device.esxi.public_ipv4_subnet_size), 3)}"
    ipv4_prefix              = "${packet_device.esxi.public_ipv4_subnet_size}"
    ipv4_gateway             = "${packet_device.esxi.network[0].gateway}"
    network_name             = "${cidrhost(format("%s/%s", packet_device.esxi.network[0].gateway, packet_device.esxi.public_ipv4_subnet_size), 3)}"
    os_password              = "${random_string.password.result}"
    sso_password             = "${random_string.password.result}"
  }
}

resource "local_file" "vcsa" {
  content  = "${data.template_file.vcsa.rendered}"
  filename = "${path.module}/template.json"
}

resource "tls_private_key" "test" {
  algorithm = "RSA"
}

resource "packet_project_ssh_key" "test" {
  name       = "tf-acc-test"
  public_key = "${tls_private_key.test.public_key_openssh}"
  project_id = "${packet_project.test.id}"
}


resource "packet_device" "helper" {
  hostname            = "tf-acc-vmware-helper"
  plan                = "t1.small.x86"
  facilities          = ["${var.facility}"]
  operating_system    = "${data.packet_operating_system.helper.id}"
  billing_cycle       = "hourly"
  project_id          = "${packet_project.test.id}"
  project_ssh_key_ids = ["${packet_project_ssh_key.test.id}"]

  provisioner "file" {
    connection {
      type        = "ssh"
      host        = "${self.access_public_ipv4}"
      user        = "root"
      private_key = "${tls_private_key.test.private_key_pem}"
      agent       = false
    }

    source      = "./install-vcsa.sh"
    destination = "/tmp/install-vcsa.sh"
  }

  provisioner "file" {
    connection {
      type        = "ssh"
      host        = "${self.access_public_ipv4}"
      user        = "root"
      private_key = "${tls_private_key.test.private_key_pem}"
      agent       = false
    }

    source      = "${local_file.vcsa.filename}"
    destination = "/tmp/vcsa-template.json"
  }

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      host        = "${self.access_public_ipv4}"
      user        = "root"
      private_key = "${tls_private_key.test.private_key_pem}"
      agent       = false
    }

    inline = [
      <<SCRIPT
export OVFTOOL_URL="${var.ovftool_url}"
export VCSA_ISO_URL="${var.vcsa_iso_url}"
export VCSA_TPL_PATH=/tmp/vcsa-template.json
echo "Installing vCenter Server Appliance..."
chmod a+x /tmp/install-vcsa.sh
/tmp/install-vcsa.sh
SCRIPT
    ]
  }
}


data "packet_operating_system" "esxi" {
  name = "VMware ESXi"
  distro = "vmware"
  version = "${var.esxi_version}"
  provisionable_on = "${var.plan}"
}

resource "packet_device" "esxi" {
  hostname = "tf-acc-vmware-esxi"
  plan = "${var.plan}"
  facilities = ["${var.facility}"]
  operating_system = "${data.packet_operating_system.esxi.id}"
  billing_cycle = "hourly"
  project_id = "${packet_project.test.id}"
  project_ssh_key_ids = ["${packet_project_ssh_key.test.id}"]

  provisioner "remote-exec" {
    connection {
      type = "ssh"
      host = "${self.access_public_ipv4}"
      user = "root"
      private_key = "${tls_private_key.test.private_key_pem}"
      agent = false
      timeout = "7m"
    }

    inline = [
      "openssl x509 -in /etc/vmware/ssl/rui.crt -fingerprint -sha1 -noout | awk -F= '{print $2}' > /tmp/ssl-rui-thumbprint.txt"
    ]
  }

  provisioner "local-exec" {
    environment = {
      SSH_PRIV_KEY = "${tls_private_key.test.private_key_pem}"
      FROM = "root@${self.access_public_ipv4}:/tmp/ssl-rui-thumbprint.txt"
      TO = "${local.esxi_ssl_cert_thumbprint_path}"
    }
    command = "./scp.sh"
  }
}
