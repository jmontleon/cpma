package transform

import (
	"github.com/BurntSushi/toml"
	"github.com/fusor/cpma/pkg/ocp4"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func (c RegistriesTransform) Run(content []byte) (TransformOutput, error) {
	logrus.Info("RegistriesTransform::Run")

	var containers Containers
	var manifests ocp4.Manifests

	if _, err := toml.Decode(string(content), &containers); err != nil {
		// handle error
	}

	const (
		apiVersion = "config.openshift.io/v1"
		kind       = "Image"
		name       = "cluster"
		annokey    = "release.openshift.io/create-only"
		annoval    = "true"
	)

	var imageCR ImageCR
	imageCR.APIVersion = apiVersion
	imageCR.Kind = kind
	imageCR.Metadata.Name = name
	imageCR.Metadata.Annotations = make(map[string]string)
	imageCR.Metadata.Annotations[annokey] = annoval
	imageCR.Spec.RegistrySources.BlockedRegistries = containers.Registries["block"].List
	imageCR.Spec.RegistrySources.InsecureRegistries = containers.Registries["insecure"].List

	imageCRYAML, err := yaml.Marshal(&imageCR)
	if err != nil {
		HandleError(err)
	}

	manifest := ocp4.Manifest{Name: "100_CPMA-cluster-config-registries.yaml", CRD: imageCRYAML}
	manifests = append(manifests, manifest)

	return ManifestTransformOutput{
		Config:    *c.Config,
		Manifests: manifests,
	}, nil
}

func (c RegistriesTransform) Extract() []byte {
	logrus.Info("RegistriesTransform::Extract")
	return c.Config.Fetch(c.Config.RegistriesConfigFile)
}

func (c RegistriesTransform) Validate() error {
	return nil // Simulate fine
}