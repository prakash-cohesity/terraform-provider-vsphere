#!/usr/bin/env python3
import fileinput

# Convert vsphere_nas_datastore to vsphere_cohesity_datastore_destroy in
# terraform state. After this change in state, the vsphere_nas_datastore
# will be managed by vsphere_cohesity_datastore_destroy
with fileinput.FileInput('terraform.tfstate', inplace=True) as file:
    for line in file:
        print(line.replace('vsphere_nas_datastore',
            'vsphere_cohesity_datastore_destroy'), end='')
