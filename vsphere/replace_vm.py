#!/usr/bin/env python3
import fileinput

with fileinput.FileInput('terraform.tfstate', inplace=True) as file:
        for line in file:
                    print(line.replace('vsphere_virtual_machine', 'vsphere_cohesity_migrate_virtual_machine'), end='')
