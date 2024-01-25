package cm

import (
	"github.com/rednexela1941/eip/pkg/epath"
)

// Volume 1: Table 3-5.16
func (self *SharedForwardOpenRequest) GetConfigPath() (*epath.ApplicationPath, bool) {
	path := self.ConnectionPath
	aps := path.ApplicationPaths
	o2tType, t2oType := self.OtoTParameters.Type, self.TtoOParameters.Type

	if len(aps) == 0 {
		return nil, false
	}

	if o2tType == ConnectionTypeNull && t2oType == ConnectionTypeNull {
		if path.HasDataSegment() {
			return nil, false
		}
		return &aps[0], true
	}

	if t2oType == ConnectionTypeNull {
		if path.HasDataSegment() {
			return &aps[0], true
		}
		return nil, true
	}

	if o2tType == ConnectionTypeNull {
		if path.HasDataSegment() {
			// First path configurtion, if only one, it is for consumption as well.
			return &aps[0], true
		}
		return nil, true
	}

	// neither are null.
	return &aps[0], true
}

// Volume 1: Table 3-5.16
func (self *SharedForwardOpenRequest) GetProductionPath() (*epath.ApplicationPath, bool) {
	path := self.ConnectionPath
	aps := path.ApplicationPaths
	o2tType, t2oType := self.OtoTParameters.Type, self.TtoOParameters.Type

	if len(aps) == 0 || len(aps) > 3 {
		return nil, false
	}

	if o2tType == ConnectionTypeNull && t2oType == ConnectionTypeNull {
		return nil, true
	}

	if t2oType == ConnectionTypeNull {
		return nil, true
	}

	if o2tType == ConnectionTypeNull {
		if path.HasDataSegment() {
			// First path configurtion, if only one, it is for consumption as well.
			if len(aps) == 1 {
				return &aps[0], true
			} else {
				return &aps[1], true
			}
		}
		return &aps[0], true
	}

	if len(aps) == 1 {
		return &aps[0], true
	}
	if len(aps) == 2 {
		return &aps[1], true
	}
	return &aps[2], true
}

// Volume 1: Table 3-5.16
func (self *SharedForwardOpenRequest) GetConsumptionPath() (*epath.ApplicationPath, bool) {
	path := self.ConnectionPath
	aps := path.ApplicationPaths
	o2tType, t2oType := self.OtoTParameters.Type, self.TtoOParameters.Type

	if len(aps) == 0 || len(aps) > 3 {
		return nil, false
	}

	if o2tType.IsNull() && t2oType.IsNull() {
		if path.HasDataSegment() {
			return nil, true
		}
		return nil, true
	}

	if t2oType.IsNull() {
		if path.HasDataSegment() {
			if len(aps) == 1 {
				return &aps[0], true
			}
			return &aps[1], true
		}
		if len(aps) != 1 {
			return nil, false
		}
		return &aps[0], true
	}

	if o2tType.IsNull() {
		return nil, true
	}

	if len(aps) == 1 {
		return &aps[0], true
	}

	return &aps[1], true
}
