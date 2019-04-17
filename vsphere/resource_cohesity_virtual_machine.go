// Copyright 2018 Cohesity Inc.
//
// Author: Prakash Vaghela (prakash.vaghela@cohesity.com)

package vsphere

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/virtualmachine"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/virtualdevice"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/vmworkflow"
	"github.com/vmware/govmomi/vim25/types"
)

func resourceCohesityMigrateVirtualMachine() *schema.Resource {
	s := map[string]*schema.Schema{
		"vm_uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "UUID of the virtual machine to vMotion to the new datastore.",
		},
		"resource_pool_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of a resource pool to put the virtual machine in.",
		},
		"datastore_id": {
			Type:          schema.TypeString,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"datastore_cluster_id"},
			Description:   "The ID of the virtual machine's datastore. The virtual machine configuration is placed here, along with any virtual disks that are created without datastores.",
		},
		"datastore_cluster_id": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"datastore_id"},
			Description:   "The ID of a datastore cluster to put the virtual machine in.",
		},
		"folder": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the folder to locate the virtual machine in.",
			StateFunc:   folder.NormalizePath,
		},
		"host_system_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The ID of an optional host system to pin the virtual machine to.",
		},
		"wait_for_guest_ip_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "The amount of time, in minutes, to wait for an available IP address on this virtual machine. A value less than 1 disables the waiter.",
		},
		"wait_for_guest_net_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     5,
			Description: "The amount of time, in minutes, to wait for an available IP address on this virtual machine. A value less than 1 disables the waiter.",
		},
		"wait_for_guest_net_routable": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Controls whether or not the guest network waiter waits for a routable address. When false, the waiter does not wait for a default gateway, nor are IP addresses checked against any discovered default gateways as part of its success criteria.",
		},
		"ignored_guest_ips": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of IP addresses to ignore while waiting for an IP",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"shutdown_wait_timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      3,
			Description:  "The amount of time, in minutes, to wait for shutdown when making necessary updates to the virtual machine.",
			ValidateFunc: validation.IntBetween(1, 10),
		},
		"migrate_wait_timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      30,
			Description:  "The amount of time, in minutes, to wait for a vMotion operation to complete before failing.",
			ValidateFunc: validation.IntAtLeast(10),
		},
		"force_power_off": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Set to true to force power-off a virtual machine if a graceful guest shutdown failed for a necessary operation.",
		},
		"scsi_controller_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1,
			Description:  "The number of SCSI controllers that Terraform manages on this virtual machine. This directly affects the amount of disks you can add to the virtual machine and the maximum disk unit number. Note that lowering this value does not remove controllers.",
			ValidateFunc: validation.IntBetween(1, 4),
		},
		"scsi_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      virtualdevice.SubresourceControllerTypeParaVirtual,
			Description:  "The type of SCSI bus this virtual machine will have. Can be one of lsilogic, lsilogic-sas or pvscsi.",
			ValidateFunc: validation.StringInSlice(virtualdevice.SCSIBusTypeAllowedValues, false),
		},
		"scsi_bus_sharing": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      string(types.VirtualSCSISharingNoSharing),
			Description:  "Mode for sharing the SCSI bus. The modes are physicalSharing, virtualSharing, and noSharing.",
			ValidateFunc: validation.StringInSlice(virtualdevice.SCSIBusSharingAllowedValues, false),
		},
		// NOTE: disk is only optional so that we can flag it as computed and use
		// it in ResourceDiff. We validate this field in ResourceDiff to enforce it
		// having a minimum count of 1 for now - but may support diskless VMs
		// later.
		"disk": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "A specification for a virtual disk device on this virtual machine.",
			MaxItems:    60,
			Elem:        &schema.Resource{Schema: virtualdevice.DiskSubresourceSchema()},
		},
		"network_interface": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "A specification for a virtual NIC on this virtual machine.",
			MaxItems:    10,
			Elem:        &schema.Resource{Schema: virtualdevice.NetworkInterfaceSubresourceSchema()},
		},
		"cdrom": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A specification for a CDROM device on this virtual machine.",
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: virtualdevice.CdromSubresourceSchema()},
		},
		"clone": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A specification for cloning a virtual machine from template.",
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: vmworkflow.VirtualMachineCloneSchema()},
		},
		"reboot_required": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Value internal to Terraform used to determine if a configuration set change requires a reboot.",
		},
		"vmware_tools_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The state of VMware tools in the guest. This will determine the proper course of action for some device operations.",
		},
		"vmx_path": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The path of the virtual machine's configuration file in the VM's datastore.",
		},
		"imported": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "A flag internal to Terraform that indicates that this resource was either imported or came from a earlier major version of this resource. Reset after the first post-import or post-upgrade apply.",
		},
		"moid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The machine object ID from VMWare",
		},
		vSphereTagAttributeKey:    tagsSchema(),
		customattribute.ConfigKey: customattribute.ConfigSchema(),
	}
	structure.MergeSchema(s, schemaVirtualMachineConfigSpec())
	structure.MergeSchema(s, schemaVirtualMachineGuestInfo())

	return &schema.Resource{
		Create:        resourceCohesityMigrateVirtualMachineCreate,
		Read:          resourceCohesityMigrateVirtualMachineNoOp,
		Update:        resourceCohesityMigrateVirtualMachineNoOp,
		Delete:        resourceVSphereVirtualMachineDelete,
		CustomizeDiff: resourceCohesityMigrateVirtualMachineCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: resourceCohesityMigrateVirtualMachineImport,
		},
		SchemaVersion: 3,
		MigrateState:  resourceVSphereVirtualMachineMigrateState,
		Schema:        s,
	}
}

func resourceCohesityMigrateVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning migrate", resourceCohesityMigrateVirtualMachineIDString(d))
	client := meta.(*VSphereClient).vimClient
	id := d.Get("vm_uuid").(string)

	d.SetId(id)
	_, err := virtualmachine.FromUUID(client, id)
	if err != nil {
		return fmt.Errorf("cannot locate virtual machine with UUID %q: %s", id, err)
	}

	// Now that any pending changes have been done (namely, any disks that don't
	// need to be migrated have been deleted), proceed with vMotion if we have
	// one pending.
	if err := resourceVSphereVirtualMachineUpdateLocation(d, meta); err != nil {
		return fmt.Errorf("error running VM migration: %s", err)
	}

	// All done with migration.
	log.Printf("[DEBUG] %s: Migrate complete", resourceVSphereVirtualMachineIDString(d))

	// Getting error cannot find disk device: invalid ID "" in the read below
	//return resourceVSphereVirtualMachineRead(d, meta)
	return nil
}

func resourceCohesityMigrateVirtualMachineNoOp(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCohesityMigrateVirtualMachineCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	return nil
}

func resourceCohesityMigrateVirtualMachineImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}

// resourceCohesityMigrateVirtualMachineIDString prints a friendly string for the
// vsphere_virtual_machine resource.
func resourceCohesityMigrateVirtualMachineIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, "vsphere_cohesity_migrate_virtual_machine")
}
