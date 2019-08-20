#!/usr/bin/env python3
import fileinput

# Convert vsphere_virtual_machine to vsphere_cohesity_migrate_virtual_machine in
# terraform state. After this change in state, the vsphere_virtual_machine
# will be managed by vsphere_cohesity_migrate_virtual_machine resource.
with fileinput.FileInput('terraform.tfstate', inplace=True) as file:
    for line in file:
        print(line.replace('vsphere_virtual_machine',
            'vsphere_cohesity_migrate_virtual_machine'), end='')
