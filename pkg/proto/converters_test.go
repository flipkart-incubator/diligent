package proto

import (
	"github.com/flipkart-incubator/diligent/pkg/charset"
	"github.com/flipkart-incubator/diligent/pkg/datagen"
	"github.com/flipkart-incubator/diligent/pkg/intgen"
	"github.com/flipkart-incubator/diligent/pkg/keygen"
	"github.com/flipkart-incubator/diligent/pkg/strtr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataSpecConversion(t *testing.T) {
	dataSpec1 := datagen.NewSpec(10000, 1024)
	dsProto := DataSpecToProto(dataSpec1)
	dataSpec2 := DataSpecFromProto(dsProto)
	assert.Equal(t, dataSpec1, dataSpec2)
}

func TestKeyGenSpecConversion(t *testing.T) {
	kgs1 := keygen.NewRandomLeveledKeyGenSpec([]int{30, 20, 10}, 5)
	kgsProto := KeyGenSpecToProto(kgs1)
	kgs2 := KeyGenSpecFromProto(kgsProto)
	assert.Equal(t, kgs1, kgs2)
}

func TestTrSpecConversion(t *testing.T) {
	tr1 := strtr.NewRandomTrSpec(charset.AlphaUp)
	trsProto := TrSpecToProto(tr1)
	tr2 := TrSpecFromProto(trsProto)
	assert.Equal(t, tr1, tr2)
}

func TestRangeConversion(t *testing.T) {
	range1 := intgen.NewRange(0, 1771)
	rangeProto := RangeToProto(range1)
	range2 := RangeFromProto(rangeProto)
	assert.Equal(t, range1, range2)
}
