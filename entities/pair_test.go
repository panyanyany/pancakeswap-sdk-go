package entities

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/panyanyany/pancakeswap-sdk-go/constants"
	"github.com/panyanyany/pancakeswap-sdk-go/utils"
)

func TestGetAddress(t *testing.T) {
	var tests = []struct {
		Input  [2]string
		Output string
	}{
		{
			// WBNB,USDC
			[2]string{"0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"},
			"0x06118051f808d91277E6AB68E16450f98D066001",
		},
		{
			// WBNB,USDC
			// cover cache
			[2]string{"0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"},
			"0x06118051f808d91277E6AB68E16450f98D066001",
		},
		{
			// WBTC,DAI
			[2]string{"0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c", "0x1af3f329e8be154074d8769d1ffa4ee058b1dbc3"},
			"0xcB85A1505ec138C2b14b12bb1BFd870667Cc45Ac",
		},
		{
			// WBTC,AAVE
			// cover cache
			[2]string{"0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c", "0xfb6115445bff7b52feb98650c87f44907e58f802"},
			"0x05B1a7F60F23AFE449EF9989526Af81B74F7Dd05",
		},
	}
	for i, test := range tests {
		output := _PairAddressCache.GetAddress(common.HexToAddress(test.Input[0]), common.HexToAddress(test.Input[1]))
		if output.String() != utils.ValidateAndParseAddress(test.Output).String() {
			t.Errorf("test #%d: failed to match when it should (%s != %s)", i, output, test.Output)
		}
	}
}

// nolint funlen
func TestPair(t *testing.T) {
	USDC, _ := NewToken(constants.Mainnet, common.HexToAddress("0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d"), 18, "USDC", "USD Coin")
	DAI, _ := NewToken(constants.Mainnet, common.HexToAddress("0x1af3f329e8be154074d8769d1ffa4ee058b1dbc3"), 18, "DAI", "DAI Stablecoin")
	tokenAmountUSDC, _ := NewTokenAmount(USDC, constants.B100)
	tokenAmountDAI, _ := NewTokenAmount(DAI, constants.B100)
	tokenAmountUSDC101, _ := NewTokenAmount(USDC, big.NewInt(101))
	tokenAmountDAI101, _ := NewTokenAmount(DAI, big.NewInt(101))

	// cannot be used for tokens on different chains
	{
		tokenAmountB, _ := NewTokenAmount(WETH[constants.Testnet], constants.B100)
		_, output := NewPair(tokenAmountUSDC, tokenAmountB)
		expect := ErrDiffChainID
		if expect != output {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
	}

	// returns the correct address
	{
		output := _PairAddressCache.GetAddress(DAI.Address, USDC.Address)
		expect := "0xadBba1EF326A33FDB754f14e62A96D5278b942Bd"
		if output.String() != expect {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
	}

	{
		pairA, _ := NewPair(tokenAmountUSDC, tokenAmountDAI)
		pairB, _ := NewPair(tokenAmountDAI, tokenAmountUSDC)
		expect := DAI
		// always is the token that sorts before
		output := pairA.Token0()
		if !expect.Equals(output) {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
		output = pairB.Token0()
		if !expect.Equals(output) {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}

		expect = USDC
		// always is the token that sorts after
		output = pairA.Token1()
		if !expect.Equals(output) {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
		output = pairB.Token1()
		if !expect.Equals(output) {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
	}

	{
		pairA, _ := NewPair(tokenAmountUSDC, tokenAmountDAI101)
		pairB, _ := NewPair(tokenAmountDAI101, tokenAmountUSDC)
		expect := tokenAmountDAI101
		// always comes from the token that sorts before
		output := pairA.Reserve0()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output = pairB.Reserve0()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		expect = tokenAmountUSDC
		// always comes from the token that sorts after
		output = pairA.Reserve1()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output = pairB.Reserve1()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
	}

	{
		pairA, _ := NewPair(tokenAmountUSDC101, tokenAmountDAI)
		pairB, _ := NewPair(tokenAmountDAI, tokenAmountUSDC101)
		expect := NewPrice(DAI.Currency, USDC.Currency, constants.B100, big.NewInt(101))
		// returns price of token0 in terms of token1
		output := pairA.Token0Price()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output = pairB.Token0Price()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		expect = NewPrice(USDC.Currency, DAI.Currency, big.NewInt(101), constants.B100)
		// returns price of token1 in terms of token0
		output = pairA.Token1Price()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output = pairB.Token1Price()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
	}

	{
		pair, _ := NewPair(tokenAmountUSDC101, tokenAmountDAI)
		// returns price of token in terms of other token
		expect := pair.Token0Price()
		output, _ := pair.PriceOf(tokenAmountDAI.Token)
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		expect = pair.Token1Price()
		output, _ = pair.PriceOf(tokenAmountUSDC101.Token)
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		{
			// throws if invalid token
			expect := ErrDiffToken
			_, output := pair.PriceOf(WETH[constants.Mainnet])
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
	}

	{
		pairA, _ := NewPair(tokenAmountUSDC, tokenAmountDAI101)
		pairB, _ := NewPair(tokenAmountDAI101, tokenAmountUSDC)
		expect := tokenAmountUSDC
		// returns reserves of the given token
		output, _ := pairA.ReserveOf(USDC)
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output, _ = pairB.ReserveOf(USDC)
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		expect = tokenAmountUSDC
		// always comes from the token that sorts after
		output = pairA.Reserve1()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}
		output = pairB.Reserve1()
		if !expect.Fraction.EqualTo(output.Fraction) {
			t.Errorf("expect[%+v], but got[%+v]", expect.Fraction, output.Fraction)
		}

		{
			// throws if not in the pair
			expect := ErrDiffToken
			_, output := pairB.ReserveOf(WETH[constants.Mainnet])
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}
	}

	{
		pairA, _ := NewPair(tokenAmountUSDC, tokenAmountDAI)
		pairB, _ := NewPair(tokenAmountDAI, tokenAmountUSDC)
		expect := constants.Mainnet
		// returns the token0 chainId
		output := pairA.ChainID()
		if expect != output {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}
		output = pairB.ChainID()
		if expect != output {
			t.Errorf("expect[%+v], but got[%+v]", expect, output)
		}

		{
			expect := true
			// involvesToken
			output := pairA.InvolvesToken(USDC)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
			output = pairA.InvolvesToken(DAI)
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
			expect = false
			output = pairA.InvolvesToken(WETH[constants.Mainnet])
			if expect != output {
				t.Errorf("expect[%+v], but got[%+v]", expect, output)
			}
		}

		{
			tokenA, _ := NewToken(constants.Testnet, common.HexToAddress("0x0000000000000000000000000000000000000001"), 18, "", "")
			tokenB, _ := NewToken(constants.Testnet, common.HexToAddress("0x0000000000000000000000000000000000000002"), 18, "", "")
			tokenAmountA, _ := NewTokenAmount(tokenA, big.NewInt(0))
			tokenAmountB, _ := NewTokenAmount(tokenB, big.NewInt(0))
			pair, _ := NewPair(tokenAmountA, tokenAmountB)
			{
				tokenAmount, _ := NewTokenAmount(pair.LiquidityToken, big.NewInt(0))
				tokenAmountA, _ := NewTokenAmount(tokenA, big.NewInt(1000))
				tokenAmountB, _ := NewTokenAmount(tokenB, big.NewInt(1000))
				// getLiquidityMinted:0
				expect := ErrInsufficientInputAmount
				_, output := pair.GetLiquidityMinted(tokenAmount, tokenAmountA, tokenAmountB)
				if expect != output {
					t.Errorf("expect[%+v], but got[%+v]", expect, output)
				}

				tokenAmountA, _ = NewTokenAmount(tokenA, big.NewInt(1000000))
				tokenAmountB, _ = NewTokenAmount(tokenB, big.NewInt(1))
				_, output = pair.GetLiquidityMinted(tokenAmount, tokenAmountA, tokenAmountB)
				if expect != output {
					t.Errorf("expect[%+v], but got[%+v]", expect, output)
				}

				tokenAmountA, _ = NewTokenAmount(tokenA, big.NewInt(1001))
				tokenAmountB, _ = NewTokenAmount(tokenB, big.NewInt(1001))
				{
					expect := "1"
					liquidity, _ := pair.GetLiquidityMinted(tokenAmount, tokenAmountA, tokenAmountB)
					output := liquidity.Raw().String()
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
			}

			// getLiquidityMinted:!0
			tokenAmountA, _ = NewTokenAmount(tokenA, big.NewInt(10000))
			tokenAmountB, _ = NewTokenAmount(tokenB, big.NewInt(10000))
			pair, _ = NewPair(tokenAmountA, tokenAmountB)
			{
				tokenAmount, _ := NewTokenAmount(pair.LiquidityToken, big.NewInt(10000))
				tokenAmountA, _ = NewTokenAmount(tokenA, big.NewInt(2000))
				tokenAmountB, _ = NewTokenAmount(tokenB, big.NewInt(2000))
				expect := "2000"
				liquidity, _ := pair.GetLiquidityMinted(tokenAmount, tokenAmountA, tokenAmountB)
				output := liquidity.Raw().String()
				if expect != output {
					t.Errorf("expect[%+v], but got[%+v]", expect, output)
				}
			}

			// getLiquidityValue:!feeOn
			tokenAmountA, _ = NewTokenAmount(tokenA, big.NewInt(1000))
			tokenAmountB, _ = NewTokenAmount(tokenB, big.NewInt(1000))
			pair, _ = NewPair(tokenAmountA, tokenAmountB)
			tokenAmount, _ := NewTokenAmount(pair.LiquidityToken, big.NewInt(1000))
			tokenAmount500, _ := NewTokenAmount(pair.LiquidityToken, big.NewInt(500))
			{
				liquidityValue, _ := pair.GetLiquidityValue(tokenA, tokenAmount, tokenAmount, false, nil)
				{
					expect := true
					output := liquidityValue.Token.Equals(tokenA)
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
				{
					expect := "1000"
					output := liquidityValue.Raw().String()
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}

				liquidityValue, _ = pair.GetLiquidityValue(tokenA, tokenAmount, tokenAmount500, false, nil)
				// 500
				{
					expect := true
					output := liquidityValue.Token.Equals(tokenA)
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
				{
					expect := "500"
					output := liquidityValue.Raw().String()
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}

				liquidityValue, _ = pair.GetLiquidityValue(tokenB, tokenAmount, tokenAmount, false, nil)
				// tokenB
				{
					expect := true
					output := liquidityValue.Token.Equals(tokenB)
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
				{
					expect := "1000"
					output := liquidityValue.Raw().String()
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
			}

			// getLiquidityValue:feeOn
			{
				liquidityValue, _ := pair.GetLiquidityValue(tokenA, tokenAmount500, tokenAmount500, true, big.NewInt(500*500))
				{
					expect := true
					output := liquidityValue.Token.Equals(tokenA)
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
				{
					expect := "917" // ceiling(1000 - (500 * (1 / 6)))
					output := liquidityValue.Raw().String()
					if expect != output {
						t.Errorf("expect[%+v], but got[%+v]", expect, output)
					}
				}
			}
		}
	}
}
