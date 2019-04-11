// Copyright 2018 Cohesity Inc.
//
// Author: Prakash Vaghela (prakash.vaghela@cohesity.com)

package vsphere

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func resourceCohesityDatastoreDestroy() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the datastore to be destroyed.",
		},
	}

	return &schema.Resource{
		Create:        resourceCohesityDatastoreCreate,
		Read:          resourceCohesityDatastoreRead,
		Update:        resourceCohesityDatastoreUpdate,
		Delete:        resourceCohesityDatastoreDelete,
		CustomizeDiff: resourceCohesityDatastoreCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: resourceCohesityDatastoreImport,
		},
		SchemaVersion: 3,
		MigrateState:  resourceCohesityDatastoreMigrateState,
		Schema:        s,
	}
}

func resourceCohesityDatastoreCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCohesityDatastoreRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCohesityDatastoreUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCohesityDatastoreDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCohesityDatastoreCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	return nil
}

func resourceCohesityDatastoreImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}

func resourceCohesityDatastoreMigrateState(version int, os *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	return nil, nil
}
