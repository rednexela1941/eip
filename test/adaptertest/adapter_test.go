package adaptertest

import (
	"testing"

	"github.com/rednexela1941/eip/pkg/encap"
	"github.com/rednexela1941/eip/pkg/network"
	"github.com/rednexela1941/eip/test/conformance"
	data "github.com/rednexela1941/eip/test/data_test"
	"github.com/rednexela1941/eip/test/helpers"
	"github.com/kr/pretty"
)

func NoOp(i interface{}) { return }

func TestIdentityObject(t *testing.T) {
	a := MakeTestAdapter()
	info := network.NewDummyInfo()
	// conn := network_test.NewDummyConn(false, false, [4]byte{192, 168, 1, 100})
	iter, err := conformance.NewIteratorFilter("Identity Object")
	if err != nil {
		t.Error(err)
		return
	}
	defer iter.Close()

	for i := 0; i < 10000; i++ {
		item := iter.Next()
		if item == nil {
			return
		}

		req, err := encap.NewRequest(item.SendData)
		if err != nil {
			t.Error(err)
			return
		}
		res, err := a.Handle(info, req)
		if err != nil {
			t.Error(err)
			return
		}

		data, err := res.Encode()
		if err != nil {
			t.Error(err)
			break
		}

		if !helpers.Compare(data, item.RecvData) {
			pretty.Println("REQUEST", req)
			pretty.Println("REPLY", res)
			t.Errorf("mismatch on line number %d (%s.%s)", item.LineNumber, item.ObjectName, item.TestName)
			helpers.CompareLog(data, item.RecvData)
			break
		}

		NoOp(res)

	}
}

func TestAdapterBasic(t *testing.T) {
	return
	info := network.NewDummyInfo()

	a := MakeTestAdapter()
	for i, tp := range data.ConnectionManagerTestPairs {
		if i > 10 {
			break
		}
		// fmt.Println("CMData test line:", tp.LineNumber())

		testData := tp.Send()
		p, err := encap.NewRequest(testData)
		if err != nil {
			t.Error(err)
			t.Errorf("failed on recv item %d of %d\n", i, len(data.ConnectionManagerTestPairs))
			break
		}
		r, err := a.Handle(info, p)
		if err != nil {
			t.Error(err)
			t.Errorf("failed on recv item %d of %d\n", i, len(data.ConnectionManagerTestPairs))
			break
		}

		// NoOp(data)

		NoOp(r) // silence compiler
		// pretty.Println("response", r)

	}
}
