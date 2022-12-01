package datagen

import (
	"github.com/flipkart-incubator/diligent/pkg/keygen"
	"github.com/flipkart-incubator/diligent/pkg/strtr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataSet(t *testing.T) {
	dgSpec := NewSpec(1000, 1024, false)
	kg := keygen.NewLeveledKeyGen(dgSpec.KeyGenSpec)
	dg := NewDataGen(dgSpec)
	assert.Equal(t, kg.NumKeys(), dg.NumRecords())
	assert.Equal(t, kg.Key(0).String(), dg.Key(0))
}

func TestRecs(t *testing.T) {
	dgSpec := NewSpec(1000, 1024, false)
	dg := NewDataGen(dgSpec)
	r0 := dg.Record(0)
	r1 := dg.Record(1)

	kg := keygen.NewLeveledKeyGen(dgSpec.KeyGenSpec)
	utr := strtr.NewTr(dgSpec.UniqTrSpec)
	str := strtr.NewTr(dgSpec.SmallGrpTrSpec)
	ltr := strtr.NewTr(dgSpec.LargeGrpTrSpec)

	assert.Equal(t, kg.Key(0).String(), r0.Pk)
	assert.Equal(t, utr.Apply(kg.Key(0).String()), r0.Uniq)
	assert.Equal(t, len(kg.Key(0).String()), len(r0.Uniq))
	assert.Equal(t, str.Apply(kg.Key(0).Prefix(1)), r0.SmallGrp)
	assert.Equal(t, ltr.Apply(kg.Key(0).Prefix(0)), r0.LargeGrp)
	assert.Equal(t, dgSpec.FixedValue, r0.FixedValue)
	assert.Equal(t, int64(1), r0.SeqNum)
	assert.NotEqual(t, 0, r0.TimeStamp)
	assert.NotEqual(t, "", r0.Payload)

	assert.Equal(t, kg.Key(1).String(), r1.Pk)
	assert.Equal(t, utr.Apply(kg.Key(1).String()), r1.Uniq)
	assert.Equal(t, len(kg.Key(1).String()), len(r1.Uniq))
	assert.Equal(t, str.Apply(kg.Key(1).Prefix(1)), r1.SmallGrp)
	assert.Equal(t, ltr.Apply(kg.Key(1).Prefix(0)), r1.LargeGrp)
	assert.Equal(t, dgSpec.FixedValue, r1.FixedValue)
	assert.Equal(t, int64(2), r1.SeqNum)
	assert.NotEqual(t, 0, r1.TimeStamp)
	assert.NotEqual(t, "", r1.Payload)
}