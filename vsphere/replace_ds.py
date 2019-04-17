#!/usr/bin/env python3
import fileinput

with fileinput.FileInput('terraform.tfstate', inplace=True) as file:
        for line in file:
                    print(line.replace('vsphere_nas_datastore', 'vsphere_cohesity_datastore_destroy'), end='')
