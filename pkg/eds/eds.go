// Package eds implements helpers to generate EtherNet/IP .eds files for an adapter device.
package eds

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/url"
	"sort"
	"text/template"
	"time"

	"github.com/rednexela1941/eip/pkg/adapter"
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/identity"
	"github.com/rednexela1941/eip/pkg/param"
)

func NoOP() int { return 0 }

type (
	FileSection struct {
		DescText         string
		CreateTime       time.Time
		ModificationTime time.Time
		Revision         identity.Revision
		HomeURL          *url.URL
	}

	DeviceSection struct {
		VendorName  string
		ProductType string
	}

	Header struct {
		FileSection
		DeviceSection
	}

	FileInitFunc func(header *Header)

	// data passed to the template.
	edsData struct {
		*Header
		*adapter.Adapter
		RPI    *param.AssemblyParam
		Config *param.AssemblyParam
	}

	edsWriter struct {
		data     *edsData
		template *template.Template
	}
)

//go:embed template.eds.tmpl
var templateStr string

var funcs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
}

/***
Special EDS Classes

- Assembly: Volume 1: 7-3.6.8
- Connection Manager: Volume 1: 7-3.6.10
- Port: Volume 1: 7-3.6.12

***/

func (self *edsData) AssemblyClass() *adapter.Class {
	a, ok := self.Adapter.Classes[cip.AssemblyClassCode]
	if !ok {
		log.Fatal("no assembly class found")
	}
	return a
}

func (self *edsData) AssemblyList() []*adapter.AssemblyInstance {
	return self.Adapter.AssemblyInstances
}

func (self *edsData) ParamList() []*param.AssemblyParam {
	params := make([]*param.AssemblyParam, 0)

	if self.RPI == nil || self.Config == nil {
		log.Fatal("didn't init()")
	}
	params = append(params, self.RPI)
	params = append(params, self.Config)

	// TODO: Add RPI and Config types.
	for _, inst := range self.Adapter.AssemblyInstances {
		params = append(params, inst.Parameters...)
	}
	// add param indices
	for i, p := range params {
		p.Index = i + 1
	}

	return params
}

func (self *edsData) PortClass() *adapter.Class {
	// See Volume 1: 7-3.6.12
	c, ok := self.Adapter.Classes[cip.PortClassCode]
	if !ok {
		return nil
	}
	return c
}

func (self *edsData) ConnectionManagerClass() *adapter.Class {
	c, ok := self.Adapter.Classes[cip.ConnectionManagerClassCode]
	if !ok {
		return nil
	}
	return c
}

func (self *edsData) StandardClasses() []*adapter.Class {
	classes := make([]*adapter.Class, 0)

	sortFn := func(i, j int) bool {
		a := classes[i]
		b := classes[j]
		return a.ClassCode < b.ClassCode
	}

	for _, c := range self.Adapter.Classes {
		switch c.ClassCode {
		case cip.AssemblyClassCode, cip.ConnectionManagerClassCode, cip.PortClassCode:
			continue
		default:
			classes = append(classes, c)
		}
	}

	sort.Slice(classes, sortFn)
	return classes
}

func (self *edsData) init() {
	rpi := param.NewUDINTParam("RPI", nil).SetHelpString(
		"Requested Packet Interval",
	).SetUnitsString(
		"Microseconds",
	).SetMinString(
		fmt.Sprintf("%d", 10*1000), // 10 millis
	).SetMaxString(
		fmt.Sprintf("%d", 1000*1000), // 1 second
	).SetDefaultValueString(
		fmt.Sprintf("%d", 100*1000), // 100 millis
	)

	config := param.NewBYTEParam("Config Data", nil).SetHelpString(
		"Config Data",
	)
	self.RPI = rpi
	self.Config = config

	self.ParamList() // to init the Index numbers
}

var defaultCreateDate = time.Date(2016, 11, 8, 10, 15, 0, 0, time.UTC)

func NewEDSWriter(init FileInitFunc, a *adapter.Adapter) (io.WriterTo, error) {
	url, _ := url.Parse("https://github.com/rednexela1941/eip")

	header := &Header{
		FileSection: FileSection{
			DescText:         fmt.Sprintf("%s EtherNet/IP EDS File", a.Identity.ProductName),
			CreateTime:       defaultCreateDate,
			ModificationTime: time.Now(),
			Revision:         a.Identity.Revision,
			HomeURL:          url,
		},
		DeviceSection: DeviceSection{
			VendorName:  "Default",
			ProductType: "Generic Device",
		},
	}

	if init != nil {
		init(header)
	}

	data := &edsData{
		Header:  header,
		Adapter: a,
	}

	tmpl, err := template.New("edsfile").Funcs(funcs).Parse(templateStr)
	if err != nil {
		return nil, err
	}

	data.init()

	return &edsWriter{
		data:     data,
		template: tmpl,
	}, nil
}

func (self *edsWriter) WriteTo(w io.Writer) (int64, error) {
	err := self.template.Execute(w, self.data)
	return 0, err
}
