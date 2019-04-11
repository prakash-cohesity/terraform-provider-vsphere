// Copyright 2018 Cohesity Inc.
//
// Author: Prakash Vaghela (prakash.vaghela@cohesity.com)

package vsphere

import (
	"log"
	"reflect"
	"unsafe"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/structure"
)

func resourceCohesityVirtualMachine() *schema.Resource {
	s := map[string]*schema.Schema{
		"nas_datastore_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the virtual machine's NAS datastore. The virtual machine will be created from here and then moved to the target datastore. It is expected that all the virtual machine configuration and disk files are available in the corresponding folder in this datastore.",
		},
		"local_datastore_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the virtual machine local datastore. The virtual machine will be moved to this datastore on comletion of creation from NAS datastore.",
		},
		"vsphere_vm": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Specification of vsphere_virtual_machine that is being created and migrated by this resource.",
			MaxItems:    1,
			Elem:        resourceVSphereVirtualMachine(),
		},
	}

	return &schema.Resource{
		Create:        resourceCohesityVirtualMachineCreate,
		Read:          resourceCohesityVirtualMachineRead,
		Update:        resourceCohesityVirtualMachineUpdate,
		Delete:        resourceCohesityVirtualMachineDelete,
		CustomizeDiff: resourceCohesityVirtualMachineCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: resourceCohesityVirtualMachineImport,
		},
		SchemaVersion: 3,
		MigrateState:  resourceCohesityVirtualMachineMigrateState,
		Schema:        s,
	}
}

func resourceCohesityVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceCohesityVirtualMachineIDString(d))
	v := reflect.ValueOf(*d)
	y := v.FieldByName("schema")
	iter := y.MapRange()
	for iter.Next() {
		if iter.Key().String() == "vsphere_vm" {
			val := iter.Value().Elem()
			sch := reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()
			res := sch.Interface().(schema.Schema).Elem.(*schema.Resource)
			log.Printf("[DEBUG] %#v %T: elem", res, res)
			_, err := res.Apply(d.State(), nil, meta)
			log.Printf("[DEBUG] %s: Create complete", resourceCohesityVirtualMachineIDString(d))
			return err
		}
	}

	// All done!
	log.Printf("[DEBUG] %s: Create complete", resourceCohesityVirtualMachineIDString(d))
	return nil
}

func resourceCohesityVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Reading state of virtual machine", resourceCohesityVirtualMachineIDString(d))
	// Add code here
	log.Printf("[DEBUG] %s: Read complete", resourceCohesityVirtualMachineIDString(d))
	return nil
}

func resourceCohesityVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Performing update", resourceCohesityVirtualMachineIDString(d))
	// Add code here
	log.Printf("[DEBUG] %s: Update complete", resourceCohesityVirtualMachineIDString(d))
	return nil
}

func resourceCohesityVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Performing delete", resourceCohesityVirtualMachineIDString(d))
	// Add code here
	log.Printf("[DEBUG] %s: Delete complete", resourceCohesityVirtualMachineIDString(d))
	return nil
}

func resourceCohesityVirtualMachineCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	log.Printf("[DEBUG] %s: Performing diff customization and validation", resourceCohesityVirtualMachineIDString(d))
	// Add code here
	log.Printf("[DEBUG] %s: Diff customization and validation complete", resourceCohesityVirtualMachineIDString(d))
	return nil
}

func resourceCohesityVirtualMachineImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Performing import complete", resourceCohesityVirtualMachineIDString(d))
	// Add code here
	log.Printf("[DEBUG] %s: Import complete, resource is ready for read", resourceCohesityVirtualMachineIDString(d))
	return nil, nil
}

// resourceCohesityVirtualMachineIDString prints a friendly string for the
// vsphere_virtual_machine resource.
func resourceCohesityVirtualMachineIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, "vsphere_virtual_machine")
}

// resourceCohesityVirtualMachineMigrateState is the master state migration function for
// the vsphere_virtual_machine resource.
func resourceCohesityVirtualMachineMigrateState(version int, os *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	return nil, nil
}
