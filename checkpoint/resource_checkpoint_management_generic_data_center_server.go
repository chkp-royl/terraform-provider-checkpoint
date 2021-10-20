package checkpoint

import (
	"fmt"
	"log"

	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceManagementGenericDataCenterServer() *schema.Resource {
	return &schema.Resource{
		Create: createManagementGenericDataCenterServer,
		Read:   readManagementGenericDataCenterServer,
		Update: updateManagementGenericDataCenterServer,
		Delete: deleteManagementGenericDataCenterServer,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Object name. Must be unique in the domain.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL of the JSON feed (e.g. https://example.com/file.json).",
			},
			"interval": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Update interval of the feed in seconds.",
			},
			"custom_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When set to false, The admin is not using Key and Value for a Custom Header in order to connect to the feed server.\n\nWhen set to true, The admin is using Key and Value for a Custom Header in order to connect to the feed server.",
				Default:     false,
			},
			"custom_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key for the Custom Header, relevant and required only when custom_header set to true.",
			},
			"custom_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Value for the Custom Header, relevant and required only when custom_header set to true.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Collection of tag identifiers.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"color": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Color of the object. Should be one of existing colors.",
				Default:     "black",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comments string.",
			},
			"details_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The level of detail for some of the fields in the response can vary from showing only the UID value of the object to a fully detailed representation of the object.",
				Default:     "standard",
			},
			"ignore_warnings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Apply changes ignoring warnings. By Setting this parameter to 'true' test connection failure will be ignored.",
				Default:     false,
			},
			"ignore_errors": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Apply changes ignoring errors. You won't be able to publish such a changes. If ignore-warnings flag was omitted - warnings will also be ignored.",
				Default:     false,
			},
		},
	}
}

func createManagementGenericDataCenterServer(d *schema.ResourceData, m interface{}) error {
	client := m.(*checkpoint.ApiClient)

	genericDataCenterServer := make(map[string]interface{})

	if v, ok := d.GetOk("name"); ok {
		genericDataCenterServer["name"] = v.(string)
	}

	genericDataCenterServer["type"] = "generic"

	if v, ok := d.GetOk("url"); ok {
		genericDataCenterServer["url"] = v.(string)
	}

	if v, ok := d.GetOk("interval"); ok {
		genericDataCenterServer["interval"] = v.(string)
	}

	if v, ok := d.GetOk("custom_header"); ok {
		genericDataCenterServer["custom-header"] = v.(bool)
	}

	if v, ok := d.GetOk("custom_key"); ok {
		genericDataCenterServer["custom-key"] = v.(string)
	}

	if v, ok := d.GetOk("custom_value"); ok {
		genericDataCenterServer["custom-value"] = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		genericDataCenterServer["tags"] = v.(*schema.Set).List()
	}

	if v, ok := d.GetOk("color"); ok {
		genericDataCenterServer["color"] = v.(string)
	}

	if v, ok := d.GetOk("comments"); ok {
		genericDataCenterServer["comments"] = v.(string)
	}

	if v, ok := d.GetOk("details_level"); ok {
		genericDataCenterServer["details-level"] = v.(string)
	}

	if v, ok := d.GetOkExists("ignore_warnings"); ok {
		genericDataCenterServer["ignore-warnings"] = v.(bool)
	}

	if v, ok := d.GetOkExists("ignore_errors"); ok {
		genericDataCenterServer["ignore-errors"] = v.(bool)
	}

	log.Println("Create genericDataCenterServer - Map = ", genericDataCenterServer)

	addGenericDataCenterServerRes, err := client.ApiCall("add-data-center-server", genericDataCenterServer, client.GetSessionID(), true, false)
	if err != nil || !addGenericDataCenterServerRes.Success {
		if addGenericDataCenterServerRes.ErrorMsg != "" {
			return fmt.Errorf(addGenericDataCenterServerRes.ErrorMsg)
		}
		return fmt.Errorf(err.Error())
	}
	payload := map[string]interface{}{
		"name": genericDataCenterServer["name"],
	}
	showGenericDataCenterServerRes, err := client.ApiCall("show-data-center-server", payload, client.GetSessionID(), true, false)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if !showGenericDataCenterServerRes.Success {
		if objectNotFound(showGenericDataCenterServerRes.GetData()["code"].(string)) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(showGenericDataCenterServerRes.ErrorMsg)
	}
	d.SetId(showGenericDataCenterServerRes.GetData()["uid"].(string))
	return readManagementGenericDataCenterServer(d, m)
}

func readManagementGenericDataCenterServer(d *schema.ResourceData, m interface{}) error {
	client := m.(*checkpoint.ApiClient)
	payload := map[string]interface{}{
		"uid": d.Id(),
	}
	showGenericDataCenterServerRes, err := client.ApiCall("show-data-center-server", payload, client.GetSessionID(), true, false)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if !showGenericDataCenterServerRes.Success {
		if objectNotFound(showGenericDataCenterServerRes.GetData()["code"].(string)) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(showGenericDataCenterServerRes.ErrorMsg)
	}
	genericDataCenterServer := showGenericDataCenterServerRes.GetData()

	if v := genericDataCenterServer["name"]; v != nil {
		_ = d.Set("name", v)
	}

	if v := genericDataCenterServer["url"]; v != nil {
		_ = d.Set("url", v)
	}

	if v := genericDataCenterServer["interval"]; v != nil {
		_ = d.Set("interval", v)
	}

	if v := genericDataCenterServer["custom_header"]; v != nil {
		_ = d.Set("custom_header", v.(bool))
	}
	if v := genericDataCenterServer["custom_key"]; v != nil {
		_ = d.Set("custom_key", v)
	}

	if v := genericDataCenterServer["custom_value"]; v != nil {
		_ = d.Set("custom_value", v)
	}

	if genericDataCenterServer["tags"] != nil {
		tagsJson, ok := genericDataCenterServer["tags"].([]interface{})
		if ok {
			tagsIds := make([]string, 0)
			if len(tagsJson) > 0 {
				for _, tags := range tagsJson {
					tags := tags.(map[string]interface{})
					tagsIds = append(tagsIds, tags["name"].(string))
				}
			}
			_ = d.Set("tags", tagsIds)
		}
	} else {
		_ = d.Set("tags", nil)
	}

	if v := genericDataCenterServer["color"]; v != nil {
		_ = d.Set("color", v)
	}

	if v := genericDataCenterServer["comments"]; v != nil {
		_ = d.Set("comments", v)
	}

	if v := genericDataCenterServer["details-level"]; v != nil {
		_ = d.Set("details_level", v)
	}

	if v := genericDataCenterServer["ignore-warnings"]; v != nil {
		_ = d.Set("ignore_warnings", v)
	}

	if v := genericDataCenterServer["ignore-errors"]; v != nil {
		_ = d.Set("ignore_errors", v)
	}

	return nil

}

func updateManagementGenericDataCenterServer(d *schema.ResourceData, m interface{}) error {

	client := m.(*checkpoint.ApiClient)
	genericDataCenterServer := make(map[string]interface{})

	if ok := d.HasChange("name"); ok {
		oldName, newName := d.GetChange("name")
		genericDataCenterServer["name"] = oldName
		genericDataCenterServer["new-name"] = newName
	} else {
		genericDataCenterServer["name"] = d.Get("name")
	}

	if ok := d.HasChange("url"); ok {
		genericDataCenterServer["url"] = d.Get("url")
	}

	if ok := d.HasChange("interval"); ok {
		genericDataCenterServer["interval"] = d.Get("interval")
	}

	if ok := d.HasChange("custom_header"); ok {
		genericDataCenterServer["custom_header"] = d.Get("custom_header").(bool)
	}

	if ok := d.HasChange("custom_key"); ok {
		genericDataCenterServer["custom_key"] = d.Get("custom_key")
	}

	if ok := d.HasChange("custom_value"); ok {
		genericDataCenterServer["custom_value"] = d.Get("custom_value")
	}

	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			genericDataCenterServer["tags"] = v.(*schema.Set).List()
		} else {
			oldTags, _ := d.GetChange("tags")
			genericDataCenterServer["tags"] = map[string]interface{}{"remove": oldTags.(*schema.Set).List()}
		}
	}

	if ok := d.HasChange("color"); ok {
		genericDataCenterServer["color"] = d.Get("color")
	}

	if ok := d.HasChange("comments"); ok {
		genericDataCenterServer["comments"] = d.Get("comments")
	}

	if ok := d.HasChange("details_level"); ok {
		genericDataCenterServer["details-level"] = d.Get("details_level")
	}

	if v, ok := d.GetOkExists("ignore_warnings"); ok {
		genericDataCenterServer["ignore-warnings"] = v.(bool)
	}

	if v, ok := d.GetOkExists("ignore_errors"); ok {
		genericDataCenterServer["ignore-errors"] = v.(bool)
	}

	log.Println("Update genericDataCenterServer - Map = ", genericDataCenterServer)

	updateGenericDataCenterServerRes, err := client.ApiCall("set-data-center-server", genericDataCenterServer, client.GetSessionID(), true, false)
	if err != nil || !updateGenericDataCenterServerRes.Success {
		if updateGenericDataCenterServerRes.ErrorMsg != "" {
			return fmt.Errorf(updateGenericDataCenterServerRes.ErrorMsg)
		}
		return fmt.Errorf(err.Error())
	}

	return readManagementGenericDataCenterServer(d, m)
}

func deleteManagementGenericDataCenterServer(d *schema.ResourceData, m interface{}) error {

	client := m.(*checkpoint.ApiClient)

	genericDataCenterServerPayload := map[string]interface{}{
		"uid": d.Id(),
	}

	log.Println("Delete genericDataCenterServer")

	deleteGenericDataCenterServerRes, err := client.ApiCall("delete-data-center-server", genericDataCenterServerPayload, client.GetSessionID(), true, false)
	if err != nil || !deleteGenericDataCenterServerRes.Success {
		if deleteGenericDataCenterServerRes.ErrorMsg != "" {
			return fmt.Errorf(deleteGenericDataCenterServerRes.ErrorMsg)
		}
		return fmt.Errorf(err.Error())
	}
	d.SetId("")

	return nil
}
