// Copyright 2018 Cohesity Inc.
//
// Author: Prakash Vaghela (prakash.vaghela@cohesity.com)

package vsphere

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/customattribute"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/structure"
)

func resourceCohesityDatastoreDestroy() *schema.Resource {
	s := map[string]*schema.Schema{
		"ds_id": {
			Type:        schema.TypeString,
			Description: "Datastore ID to be deleted.",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the datastore.",
			Required:    true,
		},
		"host_system_ids": {
			Type:        schema.TypeSet,
			Description: "The managed object IDs of the hosts to mount the datastore on.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			Required:    true,
		},
		"folder": {
			Type:          schema.TypeString,
			Description:   "The path to the datastore folder to put the datastore in.",
			Optional:      true,
			ConflictsWith: []string{"datastore_cluster_id"},
			StateFunc:     folder.NormalizePath,
		},
		"datastore_cluster_id": {
			Type:          schema.TypeString,
			Description:   "The managed object ID of the datastore cluster to place the datastore in.",
			Optional:      true,
			ConflictsWith: []string{"folder"},
		},
	}
	structure.MergeSchema(s, schemaHostNasVolumeSpec())
	structure.MergeSchema(s, schemaDatastoreSummary())

	// Add tags schema
	s[vSphereTagAttributeKey] = tagsSchema()
	// Add custom attribute schema
	s[customattribute.ConfigKey] = customattribute.ConfigSchema()

	return &schema.Resource{
		Create:        resourceCohesityDatastoreCreate,
		Read:          resourceCohesityDatastoreNoOp,
		Update:        resourceCohesityDatastoreNoOp,
		Delete:        resourceCohesityDatastoreNoOp,
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
	id := d.Get("ds_id").(string)

	d.SetId(id)
	// Unmount the NAS datastore
	return resourceVSphereNasDatastoreDelete(d, meta)
}

func resourceCohesityDatastoreNoOp(d *schema.ResourceData, meta interface{}) error {
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
