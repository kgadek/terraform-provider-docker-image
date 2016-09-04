package dockerImage

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLocalDockerImage() *schema.Resource {
	return &schema.Resource{
		Create: dataSourceLocalDockerImageCreate,
		Read:   dataSourceLocalDockerImageRead,
		Exists: dataSourceLocalDockerImageExists,
		Delete: dataSourceLocalDockerImageDelete,

		Schema: map[string]*schema.Schema{
			"dockerfile_path": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Absolute path to the Dockerfile to build.",
			},

			"tag": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A short human readable identifier for this image.",
			},
		},
	}
}

func dataSourceLocalDockerImageCreate(d *schema.ResourceData, meta interface{}) error {
	pathToDockerfile := d.Get("dockerfile_path").(string)
	tag := d.Get("tag").(string)

	hash, err := dockerExec(meta.(Config).DockerExecutable).buildContainer(pathToDockerfile, tag)
	if err != nil {
		return err
	}

	d.SetId(hash)
	return nil
}

// following template data provider convention of doing real work in Exists.
// not sure how appropriate it is for a docker image.
func dataSourceLocalDockerImageRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func dataSourceLocalDockerImageExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	pathToDockerfile := d.Get("dockerfile_path").(string)
	tag := d.Get("tag").(string)

	hash, err := dockerExec(meta.(Config).DockerExecutable).buildContainer(pathToDockerfile, tag)
	if err != nil {
		return false, err
	}

	return hash == d.Id(), nil
}

func dataSourceLocalDockerImageDelete(d *schema.ResourceData, meta interface{}) error {
	return dockerExec(meta.(Config).DockerExecutable).deleteContainer(d.Id())
}