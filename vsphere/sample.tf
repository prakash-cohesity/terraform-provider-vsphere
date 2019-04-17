#===============================================================================
# vSphere Provider
#===============================================================================
provider "vsphere" {
# provider config here
}
#===============================================================================

#===============================================================================
# vSphere Resources
#===============================================================================

variable "hosts" {
  default = [
# ESXi hosts list here
  ]
}

# datacenter/resourcepool/datastore etc

resource "vsphere_nas_datastore" "datastore" {
# NAS datastore config details here
  name            = "prakash-test"
# ...

  type         = "NFS"
# ...
}

resource "null_resource" "replace_ds" {
  provisioner "local-exec" {
    command = "python3 replace_ds.py"
  }
  depends_on = ["vsphere_nas_datastore.datastore"]
}

resource "vsphere_virtual_machine" "vm" {
# vSphere VM config details here

  name             = "tf-db-instance"
# ...

  # VM storage #
  disk {
# ...
    attach           = true
# ...
  }

  depends_on = ["vsphere_nas_datastore.datastore"]
}

resource "null_resource" "replace_vm" {
  provisioner "local-exec" {
    command = "python3 replace_vm.py"
  }
  depends_on = ["vsphere_virtual_machine.vm"]
}

resource "vsphere_cohesity_migrate_virtual_machine" "vm" {
  vm_uuid             = "${vsphere_virtual_machine.vm.id}"

# vSphere VM config details here
  name             = "tf-db-instance"
# ...

  # VM storage #
  disk {
# ...
    attach           = true
# ...
  }

  depends_on = ["null_resource.replace_vm"]
}

# NAS datastore unmount
resource "vsphere_cohesity_datastore_destroy" "datastore" {
  ds_id           = "${vsphere_nas_datastore.datastore.id}"
# NAS datastore config details here
  name            = "prakash-test"
# ...
  depends_on = ["vsphere_cohesity_migrate_virtual_machine.vm","null_resource.replace_ds"]
}

