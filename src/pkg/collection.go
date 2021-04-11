package pkg

import (
	"fmt"
	"splitter/conf"
)

type PackageCollection struct {
	RootPackage *Package
	Packages    map[string]*Package
	Conf        *conf.Config
}

func (c *PackageCollection) Add(p *Package) {
	c.Packages[p.Composer.Items.Name] = p
}

func FromConfig(conf *conf.Config) (*PackageCollection, error) {
	root, err := loadRootPackage(conf)
	if err != nil {
		return nil, err
	}
	collection := &PackageCollection{
		Packages:    make(map[string]*Package, 0),
		RootPackage: root,
		Conf:        conf,
	}

	for _, pkg := range conf.Packages.Items {
		p, err := loadPackage(pkg, conf)
		if err != nil {
			return nil, fmt.Errorf("pkg %s cannot be loaded: %s", pkg.Url, err)
		}
		collection.Add(p)
	}

	return collection, nil
}
