package v1alpha1

import (
	"io/ioutil"
	logger "keep/pkg/util/loggerv1.0.1"
	"sigs.k8s.io/yaml"
)

func (c *DeviceManagerConfig) Parse(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("Failed to read configfile ", filename, err)
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		logger.Error("Failed to unmarshal configfile ", filename, err)
		return err
	}
	return nil
}
