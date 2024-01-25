package adapter

import (
	"github.com/rednexela1941/eip/pkg/cip"
	"github.com/rednexela1941/eip/pkg/epath"
	"github.com/rednexela1941/eip/pkg/identity"
	"github.com/rednexela1941/eip/pkg/mr"
)

func ValidateElectronicKey(
	id *identity.Identity,
	ekey *epath.ElectronicKey,
) (valid bool, extendedStatus []cip.AdditionalStatus) {
	valid = true

	addStatus := func(s cip.AdditionalStatus) {
		valid = false
		extendedStatus = mr.AddAdditionalStatusToArray(
			extendedStatus,
			s,
		)
	}

	if !ekey.ValidFormat() {
		addStatus(cip.InvalidConnectionPath)
		return
	}

	cb := ekey.GetCompatibilityBit()

	canIgnore := func(isZero bool) bool {
		if isZero {
			return !cb
		}
		return false
	}

	if ekey.VendorID != id.VendorID && ekey.VendorID != 0 {
		addStatus(cip.VendorIDMismatch)
	}

	if ekey.ProductCode != id.ProductCode && ekey.ProductCode != 0 {
		addStatus(cip.VendorIDMismatch)
	}

	if ekey.DeviceType != id.DeviceType && ekey.DeviceType != 0 {
		addStatus(cip.DeviceTypeMismatch)
	}

	major := ekey.GetMajorRevision()
	if major != id.Revision.Major && !canIgnore(major == 0) {
		addStatus(cip.RevisionMismatch)
	}

	minor := ekey.GetMinorRevision()
	if minor != id.Revision.Minor {
		if !cb || minor > id.Revision.Minor || minor == 0 {
			if !canIgnore(minor == 0) {
				addStatus(cip.RevisionMismatch)
			}
		}
	}

	if ekey.Format == epath.Format5 && ekey.SerialNumber != id.SerialNumber {
		addStatus(cip.SerialNumberMismatch)
	}

	return
}
