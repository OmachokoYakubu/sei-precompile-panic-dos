package pointer_test

import (
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/sei-protocol/sei-chain/precompiles/pointer"
	sdk "github.com/sei-protocol/sei-chain/sei-cosmos/types"
	testkeeper "github.com/sei-protocol/sei-chain/testutil/keeper"
	"github.com/sei-protocol/sei-chain/x/evm/state"
	"github.com/sei-protocol/sei-chain/x/evm/types"
	"github.com/stretchr/testify/require"
)

type mockBadWasmdKeeper struct {
	response []byte
}

func (m mockBadWasmdKeeper) QuerySmartSafe(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error) {
	return m.response, nil
}

func (m mockBadWasmdKeeper) QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte) ([]byte, error) {
	return m.response, nil
}

func TestAddCW20PanicReflection(t *testing.T) {
	testApp := testkeeper.EVMTestApp
	keepers := testApp.GetPrecompileKeepers()

	p, err := pointer.NewPrecompile(keepers)
	require.Nil(t, err)

	// Use reflection to access the private executor and its private wasmdKeeper field
	// p is a *pcommon.DynamicGasPrecompile
	// Its Executor field is an interface, which holds a *pointer.PrecompileExecutor
	
	executorValue := reflect.ValueOf(p).Elem().FieldByName("executor").Elem()
	// executorValue is now a *pointer.PrecompileExecutor
	
	// We need to use unsafe to set the private wasmdKeeper field
	wasmdKeeperField := executorValue.Elem().FieldByName("wasmdKeeper")
	
	badJsonResponse := []byte(`{"name": 123, "symbol": "BAD"}`)
	mock := mockBadWasmdKeeper{response: badJsonResponse}
	
	// Set the private field using unsafe
	ptr := unsafe.Pointer(wasmdKeeperField.UnsafeAddr())
	*(*interface{})(ptr) = interface{}(mock)

	ctx := testApp.GetContextForDeliverTx([]byte{}).WithBlockTime(time.Now())
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeterWithMultiplier(ctx))
	_, caller := testkeeper.MockAddressPair()
	suppliedGas := uint64(10000000)
	cfg := types.DefaultChainConfig().EthereumConfig(testApp.EvmKeeper.ChainID(ctx))

	m, exists := p.ABI.Methods["addCW20Pointer"]
	require.True(t, exists)
	methodID := m.ID
	
	cwAddr := "sei17pxq7vdvjnl6v9unu6w90yxgvg734ztx9mzlpv"
	args, err := m.Inputs.Pack(cwAddr)
	require.Nil(t, err)

	statedb := state.NewDBImpl(ctx, &testApp.EvmKeeper, false)
	blockCtx, _ := testApp.EvmKeeper.GetVMBlockContext(ctx, core.GasPool(suppliedGas))
	evmInstance := vm.NewEVM(*blockCtx, statedb, cfg, vm.Config{}, testApp.EvmKeeper.CustomPrecompiles(ctx))

	// This call should panic
	defer func() {
		if r := recover(); r != nil {
			t.Logf("SUCCESS: Recovered from expected panic: %v", r)
		} else {
			t.Errorf("FAIL: Expected panic but did not get one")
		}
	}()

	p.RunAndCalculateGas(evmInstance, caller, caller, append(methodID, args...), suppliedGas, nil, nil, false, false)
}
