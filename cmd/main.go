package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	terminate "github.com/glifio/go-termination-calculator"
	"github.com/spf13/cobra"
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "go-termination-calculator [options]",
	Short: "Calculate the termination costs for a sector",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		epoch, err := cmd.Flags().GetUint64("epoch")
		if err != nil {
			log.Fatal(err)
		}

		ss, err := cmd.Flags().GetUint64("sector-size")
		if err != nil {
			log.Fatal(err)
		}
		var sectorSize uint64
		switch ss {
		case 32:
			sectorSize = 32 * 1073741824
		case 64:
			sectorSize = 64 * 1073741824
		default:
			log.Fatal("Unknown sector size")
		}

		qapPos, err := cmd.Flags().GetString("qap-position")
		if err != nil {
			log.Fatal(err)
		}
		qapPosition, ok := new(big.Int).SetString(qapPos, 10)
		if !ok {
			log.Fatal("Bad qap-position")
		}

		qapVel, err := cmd.Flags().GetString("qap-velocity")
		if err != nil {
			log.Fatal(err)
		}
		qapVelocity, ok := new(big.Int).SetString(qapVel, 10)
		if !ok {
			log.Fatal("Bad qap-velocity")
		}

		rewardPos, err := cmd.Flags().GetString("reward-position")
		if err != nil {
			log.Fatal(err)
		}
		rewardPosition, ok := new(big.Int).SetString(rewardPos, 10)
		if !ok {
			log.Fatal("Bad reward-position")
		}

		rewardVel, err := cmd.Flags().GetString("reward-velocity")
		if err != nil {
			log.Fatal(err)
		}
		rewardVelocity, ok := new(big.Int).SetString(rewardVel, 10)
		if !ok {
			log.Fatal("Bad reward-velocity")
		}

		activation, err := cmd.Flags().GetUint64("activation")
		if err != nil {
			log.Fatal(err)
		}

		expiration, err := cmd.Flags().GetUint64("expiration")
		if err != nil {
			log.Fatal(err)
		}

		dealWght, err := cmd.Flags().GetString("deal-weight")
		if err != nil {
			log.Fatal(err)
		}
		dealWeight, ok := new(big.Int).SetString(dealWght, 10)
		if !ok {
			log.Fatal("Bad deal-weight")
		}

		verifiedDealWght, err := cmd.Flags().GetString("verified-deal-weight")
		if err != nil {
			log.Fatal(err)
		}
		verifiedDealWeight, ok := new(big.Int).SetString(verifiedDealWght, 10)
		if !ok {
			log.Fatal("Bad verified-deal-weight")
		}

		expectedDayRew, err := cmd.Flags().GetString("expected-day-reward")
		if err != nil {
			log.Fatal(err)
		}
		expectedDayReward, ok := new(big.Int).SetString(expectedDayRew, 10)
		if !ok {
			log.Fatal("Bad expected-day-reward")
		}

		expectedStoragePl, err := cmd.Flags().GetString("expected-storage-pledge")
		if err != nil {
			log.Fatal(err)
		}
		expectedStoragePledge, ok := new(big.Int).SetString(expectedStoragePl, 10)
		if !ok {
			log.Fatal("Bad expected-storage-pledge")
		}

		powerBaseEpoch, err := cmd.Flags().GetUint64("power-base-epoch")
		if err != nil {
			log.Fatal(err)
		}

		replacedDayRew, err := cmd.Flags().GetString("replaced-day-reward")
		if err != nil {
			log.Fatal(err)
		}
		replacedDayReward, ok := new(big.Int).SetString(replacedDayRew, 10)
		if !ok {
			log.Fatal("Bad replaced-day-reward")
		}

		burn, err := terminate.TerminateSector(
			epoch,
			sectorSize,
			qapPosition,
			qapVelocity,
			rewardPosition,
			rewardVelocity,
			activation,
			expiration,
			dealWeight,
			verifiedDealWeight,
			expectedDayReward,
			expectedStoragePledge,
			powerBaseEpoch,
			replacedDayReward,
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", burn)
	},
}

func init() {
	rootCmd.Flags().Uint64("epoch", 0, "")
	rootCmd.Flags().Uint64("sector-size", 32, "")
	rootCmd.Flags().String("qap-position", "", "")
	rootCmd.Flags().String("qap-velocity", "", "")
	rootCmd.Flags().String("reward-position", "", "")
	rootCmd.Flags().String("reward-velocity", "", "")
	rootCmd.Flags().Uint64("activation", 0, "")
	rootCmd.Flags().Uint64("expiration", 0, "")
	rootCmd.Flags().String("deal-weight", "", "")
	rootCmd.Flags().String("verified-deal-weight", "", "")
	rootCmd.Flags().String("expected-day-reward", "", "")
	rootCmd.Flags().String("expected-storage-pledge", "", "")
	rootCmd.Flags().Uint64("power-base-epoch", 0, "")
	rootCmd.Flags().String("replaced-day-reward", "", "")
}
